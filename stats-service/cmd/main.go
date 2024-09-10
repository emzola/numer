package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/emzola/numer/stats-service/config"
	"github.com/emzola/numer/stats-service/internal/client"
	"github.com/emzola/numer/stats-service/internal/grpcutil"
	"github.com/emzola/numer/stats-service/internal/handler"
	"github.com/emzola/numer/stats-service/internal/service"
	"github.com/emzola/numer/stats-service/pkg/discovery"
	consul "github.com/emzola/numer/stats-service/pkg/discovery/consul"
	pb "github.com/emzola/numer/stats-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const serviceName = "stats-service"

func main() {
	// Initialize logger
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	// Load configuration
	var cfg config.Params
	flag.StringVar(&cfg.GRPCServerAddress, "server-address", os.Getenv("GRPC_SERVER_ADDRESS"), "GRPC server address")
	flag.StringVar(&cfg.InvoiceGRPCServerAddress, "invoice-server-address", os.Getenv("INVOICE_GRPC_SERVER_ADDRESS"), "Invoice Service server address")

	ctx, cancel := context.WithCancel(context.Background())

	// Service discovery
	registry, err := consul.NewRegistry("consul:8500")
	if err != nil {
		logger.Error("failed to create new consul-based service registry instance", slog.Any("error", err))
	}
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost%v", cfg.GRPCServerAddress)); err != nil {
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

	// Set up connection to the Invoice service
	invoiceConn, err := grpcutil.ServiceConnection(ctx, "invoice-service", registry)
	if err != nil {
		logger.Error("could not connect to invoice service", slog.Any("error", err))
	}
	defer invoiceConn.Close()

	// Initialize client, service and server
	invoiceStatsClient := client.NewStatsClient(invoiceConn)
	svc := service.NewStatsService(invoiceStatsClient)
	handler := handler.NewStatsHandler(svc)

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterStatsServiceServer(grpcServer, handler)

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

	logger.Info("stats service running", slog.String("port", cfg.GRPCServerAddress))
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}

	wg.Wait()
}
