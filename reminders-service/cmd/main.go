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

	"github.com/emzola/numer/reminders-service/config"
	"github.com/emzola/numer/reminders-service/internal/client"
	"github.com/emzola/numer/reminders-service/internal/grpcutil"
	"github.com/emzola/numer/reminders-service/internal/scheduler"
	"github.com/emzola/numer/reminders-service/pkg/discovery"
	consul "github.com/emzola/numer/reminders-service/pkg/discovery/consul"
	pb "github.com/emzola/numer/reminders-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const serviceName = "reminders-service"

func main() {
	// Initialize logger
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	// Load configuration
	var cfg config.Params
	flag.StringVar(&cfg.GRPCServerAddress, "server-address", os.Getenv("GRPC_SERVER_ADDRESS"), "GRPC server address")

	ctx, cancel := context.WithCancel(context.Background())

	// Service discovery
	registry, err := consul.NewRegistry("consul:8500")
	if err != nil {
		logger.Error("failed to create new consul-based service registry instance", slog.Any("error", err))
	}
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("host.docker.internal%v", cfg.GRPCServerAddress)); err != nil {
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

	// Set up connection to the Notification service
	notificationConn, err := grpcutil.ServiceConnection(ctx, "notification-service", registry)
	if err != nil {
		logger.Error("could not connect to notification service", slog.Any("error", err))
	}
	defer notificationConn.Close()

	// Set up connection to the User service
	userConn, err := grpcutil.ServiceConnection(ctx, "user-service", registry)
	if err != nil {
		logger.Error("could not connect to user service", slog.Any("error", err))
	}
	defer userConn.Close()

	// Initialize and start reminders scheduler
	scheduler := scheduler.NewScheduler(
		client.NewInvoiceClient(invoiceConn),
		client.NewNotificationClient(invoiceConn),
		client.NewUserClient(userConn),
	)
	go func() {
		scheduler.Start()
	}()

	// Initialize server
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterRemindersServiceServer(grpcServer, scheduler)

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

	logger.Info("reminders service running", slog.String("port", cfg.GRPCServerAddress))
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}

	wg.Wait()
}
