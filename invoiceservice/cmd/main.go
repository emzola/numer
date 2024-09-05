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

	"database/sql"

	invoicepb "github.com/emzola/numer/invoiceservice/genproto"
	"github.com/emzola/numer/invoiceservice/internal/handler"
	"github.com/emzola/numer/invoiceservice/internal/repository"
	"github.com/emzola/numer/invoiceservice/internal/service"
	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
)

func main() {
	// Initialize logger
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	// Open and parse yaml configuration file
	configFilePath := filepath.Join("invoiceservice", "configs", "base.yaml")
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

	_, cancel := context.WithCancel(context.Background())

	// Establish database connection
	connStr := "postgres://" + os.Getenv("INVOICE_DB_USER") + ":" + os.Getenv("INVOICE_DB_PASSWORD") + "@invoice-db:" + os.Getenv("INVOICE_DB_PORT") + "/" + os.Getenv("INVOICE_DB_NAME") + "?sslmode=disable"
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		panic("failed to connect to database")
	}
	defer db.Close()
	logger.Info("database connection established")

	// Initialize repository, service and server
	repo := repository.NewInvoiceRepository(db)
	svc := service.NewInvoiceService(repo)
	server := handler.NewInvoiceServiceServer(svc)

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	invoicepb.RegisterInvoiceServiceServer(grpcServer, server)

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

	logger.Info("invoice service running", slog.Int("port", port))
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}

	wg.Wait()
}
