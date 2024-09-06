package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"database/sql"

	userpb "github.com/emzola/numer/userservice/genproto"
	"github.com/emzola/numer/userservice/internal/handler"
	"github.com/emzola/numer/userservice/internal/repository"
	"github.com/emzola/numer/userservice/internal/service"
	"github.com/emzola/numer/userservice/pkg/discovery"
	consul "github.com/emzola/numer/userservice/pkg/discovery/consul"
	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
)

const serviceName = "userservice"

func main() {
	// Initialize logger
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	// Open and parse yaml configuration file
	configFilePath := filepath.Join("userservice", "configs", "base.yaml")
	f, err := os.Open(configFilePath)
	if err != nil {
		logger.Error("failed to open configuration file", slog.Any("error", err))
	}
	defer f.Close()
	var cfg config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		logger.Error("failed to parse configuration", slog.Any("error", err))
	}
	port := cfg.API.Port

	ctx, cancel := context.WithCancel(context.Background())

	// Service discovery
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				logger.Error("failed to report healthy state", slog.Any("error", err))
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)

	// Establish database connection
	connStr := "postgres://" + os.Getenv("USER_DB_USER") + ":" + os.Getenv("USER_DB_PASSWORD") + "@" + os.Getenv("USER_DB_HOST") + os.Getenv("USER_DB_PORT") + "/" + os.Getenv("USER_DB_NAME") + "?sslmode=disable"
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		panic("failed to connect to database")
	}
	defer db.Close()
	logger.Info("database connection established")

	// Initialize repository, service and server
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)
	server := handler.NewUserServiceServer(svc)

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	userpb.RegisterUserServiceServer(grpcServer, server)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		logger.Error("failed to listen", slog.Any("error", err))
	}

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		s := <-stop
		cancel()
		logger.Info("shutting down gracefully", slog.String("signal", s.String()))
		grpcServer.GracefulStop()
		logger.Info("server stopped")
	}()

	logger.Info("user service running", slog.Int("port", port))
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}

	wg.Wait()
}
