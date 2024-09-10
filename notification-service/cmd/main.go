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

	"github.com/emzola/numer/notification-service/config"
	"github.com/emzola/numer/notification-service/internal/email"
	"github.com/emzola/numer/notification-service/internal/handler"
	"github.com/emzola/numer/notification-service/pkg/discovery"
	consul "github.com/emzola/numer/notification-service/pkg/discovery/consul"
	pb "github.com/emzola/numer/notification-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const serviceName = "notification-service"

func main() {
	// Initialize logger
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	// Load configuration
	var cfg config.Params
	flag.StringVar(&cfg.GRPCServerAddress, "server-address", os.Getenv("GRPC_SERVER_ADDRESS"), "GRPC server address")
	flag.StringVar(&cfg.SMTPHost, "smtp-host", os.Getenv("SMTP_HOST"), "SMTP host")
	flag.StringVar(&cfg.SMTPPort, "smtp-port", os.Getenv("SMTP_PORT"), "SMTP port")
	flag.StringVar(&cfg.SMTPUsername, "smtp-username", os.Getenv("SMTP_USERNAME"), "SMTP username")
	flag.StringVar(&cfg.SMTPPassword, "smtp-password", os.Getenv("SMTP_PASSWORD"), "SMTP password")

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

	// SMTP settings for sending emails
	emailSender := email.NewEmailSender(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUsername,
		cfg.SMTPPassword,
	)
	handler := handler.NewNotificationHandler(emailSender)

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterNotificationServiceServer(grpcServer, handler)

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

	logger.Info("notification service running", slog.String("port", cfg.GRPCServerAddress))
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}

	wg.Wait()
}
