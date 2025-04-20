package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lugondev/erp-warehouse-simple/internal/infrastructure/config"
	"github.com/lugondev/erp-warehouse-simple/internal/infrastructure/gateway"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create API Gateway
	apiGateway, err := gateway.NewGateway(cfg)
	if err != nil {
		log.Fatalf("Failed to create API Gateway: %v", err)
	}

	// Start API Gateway in a goroutine
	go func() {
		if err := apiGateway.Start(); err != nil {
			log.Fatalf("Failed to start API Gateway: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down API Gateway...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Gracefully shutdown the server
	if err := apiGateway.Stop(ctx); err != nil {
		log.Fatalf("API Gateway forced to shutdown: %v", err)
	}

	log.Println("API Gateway exited properly")
}
