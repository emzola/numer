package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/emzola/numer/invoiceservice/config"
	"github.com/emzola/numer/invoiceservice/internal/handler"
	"github.com/emzola/numer/invoiceservice/internal/repository"
	"github.com/emzola/numer/invoiceservice/internal/service"
	"github.com/emzola/numer/invoiceservice/pkg/discovery"
	consul "github.com/emzola/numer/invoiceservice/pkg/discovery/consul"
	pb "github.com/emzola/numer/invoiceservice/proto"
	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const serviceName = "invoiceservice"

func main() {
	// Initialize logger
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	// Load configuration
	var cfg config.Params
	flag.StringVar(&cfg.GRPCServerAddress, "server-address", os.Getenv("GRPC_SERVER_ADDRESS"), "GRPC server address")
	flag.StringVar(&cfg.DatabaseURL, "database-url", os.Getenv("INVOICE_DB_URL"), "POSTGRESQL database URL")

	ctx, cancel := context.WithCancel(context.Background())

	// Service discovery
	registry, err := consul.NewRegistry("consul:8500")
	if err != nil {
		logger.Error("failed to create new consul-based service registry instance", slog.Any("error", err))
	}
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("invoice-service%v", cfg.GRPCServerAddress)); err != nil {
		logger.Error("failed to create a service record in the registry", slog.Any("error", err))
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

	// Connect to PostgreSQL database
	dbpool, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to connect to the database", slog.Any("error", err))
	}
	defer dbpool.Close()
	logger.Info("database connection established")

	// Initialize repository, service and server
	repo := repository.NewInvoiceRepository(dbpool)
	svc := service.NewInvoiceService(repo)
	handler := handler.NewInvoiceHandler(svc)

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterInvoiceServiceServer(grpcServer, handler)

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0%v", cfg.GRPCServerAddress))
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

	logger.Info("invoice service running", slog.String("port", cfg.GRPCServerAddress))
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}

	wg.Wait()
}
