package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/emzola/numer/gateway-service/config"
	"github.com/emzola/numer/gateway-service/internal/handler"
	"github.com/emzola/numer/gateway-service/pkg/discovery"
	consul "github.com/emzola/numer/gateway-service/pkg/discovery/consul"
)

const serviceName = "gateway-service"

func main() {
	// Initialize logger
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	// Load configuration
	var cfg config.Params
	flag.StringVar(&cfg.Port, "server-port", os.Getenv("PORT"), "HTTP Port")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Service discovery
	registry, err := consul.NewRegistry("consul:8500")
	if err != nil {
		logger.Error("failed to create new consul-based service registry instance", slog.Any("error", err))
	}
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("host.docker.internal%v", cfg.Port)); err != nil {
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

	// Initialize handler and server
	handler := handler.NewGatewayHandler(registry)
	srv := &http.Server{
		Addr:         cfg.Port,
		Handler:      handler.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start the server in a goroutine so it doesn't block
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("ListenAndServe(): ", slog.Any("error", err))
		}
	}()
	logger.Info("gateway service running", slog.String("port", cfg.Port))

	// Create a channel to listen for interrupt signals
	// Block until a signla is received
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	<-stopChan
	logger.Info("shutdown signal received, gracefully shutting down...")

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server shutdown failed", slog.Any("error", err))
	}
	logger.Info("server gracefully stopped")

	wg.Wait()
}
