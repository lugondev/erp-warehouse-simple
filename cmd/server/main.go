package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lugondev/erp-warehouse-simple/internal/infrastructure/config"
	"github.com/lugondev/erp-warehouse-simple/internal/infrastructure/server"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/lugondev/erp-warehouse-simple/docs"
)

// @title ERP Warehouse API
// @version 1.0
// @description Simple ERP Warehouse Management System API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url https://github.com/lugondev/erp-warehouse-simple
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set Gin mode
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize server
	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Add Swagger documentation route
	srv.Router().GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := srv.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
