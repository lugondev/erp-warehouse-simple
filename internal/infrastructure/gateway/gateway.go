package gateway

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lugondev/erp-warehouse-simple/internal/infrastructure/auth"
	"github.com/lugondev/erp-warehouse-simple/internal/infrastructure/config"
	"github.com/lugondev/erp-warehouse-simple/internal/infrastructure/gateway/middleware"
	"github.com/lugondev/erp-warehouse-simple/internal/infrastructure/gateway/proxy"
	"github.com/lugondev/erp-warehouse-simple/internal/infrastructure/gateway/websocket"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Gateway represents the API Gateway
type Gateway struct {
	config     *config.APIGatewayConfig
	router     *gin.Engine
	proxy      *proxy.ServiceProxy
	jwtService *auth.JWTService
	server     *http.Server
	wsHub      *websocket.Hub
}

// NewGateway creates a new API Gateway
func NewGateway(cfg *config.Config) (*Gateway, error) {
	if !cfg.APIGateway.Enabled {
		return nil, fmt.Errorf("API Gateway is disabled in configuration")
	}

	// Set Gin mode
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	router := gin.New()

	// Initialize JWT service
	jwtService := auth.NewJWTService(cfg.JWT.AccessSecret, cfg.JWT.RefreshSecret)

	// Initialize service proxy
	serviceProxy := proxy.NewServiceProxy(cfg.APIGateway.Services)

	// Initialize WebSocket hub
	wsHub := websocket.NewHub()
	go wsHub.Run()

	// Create gateway
	gateway := &Gateway{
		config:     &cfg.APIGateway,
		router:     router,
		proxy:      serviceProxy,
		jwtService: jwtService,
		wsHub:      wsHub,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%s", cfg.APIGateway.Port),
			Handler: router,
		},
	}

	// Setup routes and middleware
	gateway.setupMiddleware()
	gateway.setupRoutes()

	return gateway, nil
}

// setupMiddleware configures middleware for the API Gateway
func (g *Gateway) setupMiddleware() {
	// Recovery middleware
	g.router.Use(gin.Recovery())

	// Logger middleware
	if g.config.Logging {
		g.router.Use(middleware.Logger())
	}

	// CORS middleware
	g.router.Use(middleware.CORS())

	// Rate limiting middleware
	g.router.Use(middleware.RateLimit(g.config.RateLimit.RequestsPerSecond, g.config.RateLimit.Burst))

	// Tracing middleware
	if g.config.Tracing {
		g.router.Use(middleware.Tracing())
	}
}

// setupRoutes configures routes for the API Gateway
func (g *Gateway) setupRoutes() {
	// Health check
	g.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Swagger documentation
	g.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// WebSocket endpoint
	g.router.GET("/ws", func(c *gin.Context) {
		websocket.ServeWs(g.wsHub, c.Writer, c.Request)
	})

	// API routes
	api := g.router.Group("/api")
	{
		// Public routes
		public := api.Group("/v1")
		{
			// Auth routes - direct to auth service
			auth := public.Group("/auth")
			{
				auth.POST("/register", g.proxy.ProxyRequest("auth", "/api/v1/auth/register"))
				auth.POST("/login", g.proxy.ProxyRequest("auth", "/api/v1/auth/login"))
				auth.POST("/refresh-token", g.proxy.ProxyRequest("auth", "/api/v1/auth/refresh-token"))
				auth.POST("/forgot-password", g.proxy.ProxyRequest("auth", "/api/v1/auth/forgot-password"))
				auth.POST("/reset-password", g.proxy.ProxyRequest("auth", "/api/v1/auth/reset-password"))
			}
		}

		// Protected routes
		protected := api.Group("/v1")
		protected.Use(middleware.Auth(g.jwtService))
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("", g.proxy.ProxyRequest("user", "/api/v1/users"))
				users.GET("/:id", g.proxy.ProxyRequest("user", "/api/v1/users/:id"))
				users.PUT("/:id", g.proxy.ProxyRequest("user", "/api/v1/users/:id"))
				users.DELETE("/:id", g.proxy.ProxyRequest("user", "/api/v1/users/:id"))
				users.POST("/logout", g.proxy.ProxyRequest("user", "/api/v1/users/logout"))
			}

			// Role routes
			roles := protected.Group("/roles")
			{
				roles.POST("", g.proxy.ProxyRequest("user", "/api/v1/roles"))
				roles.GET("", g.proxy.ProxyRequest("user", "/api/v1/roles"))
				roles.GET("/:id", g.proxy.ProxyRequest("user", "/api/v1/roles/:id"))
				roles.PUT("/:id", g.proxy.ProxyRequest("user", "/api/v1/roles/:id"))
				roles.DELETE("/:id", g.proxy.ProxyRequest("user", "/api/v1/roles/:id"))
			}

			// Audit log routes
			audit := protected.Group("/audit")
			{
				audit.GET("/logs", g.proxy.ProxyRequest("audit", "/api/v1/audit/logs"))
				audit.GET("/logs/user/:id", g.proxy.ProxyRequest("audit", "/api/v1/audit/logs/user/:id"))
			}

			// Store routes
			stores := protected.Group("/stores")
			{
				stores.POST("", g.proxy.ProxyRequest("store", "/api/v1/stores"))
				stores.GET("", g.proxy.ProxyRequest("store", "/api/v1/stores"))
				stores.GET("/:id", g.proxy.ProxyRequest("store", "/api/v1/stores/:id"))
				stores.PUT("/:id", g.proxy.ProxyRequest("store", "/api/v1/stores/:id"))
				stores.DELETE("/:id", g.proxy.ProxyRequest("store", "/api/v1/stores/:id"))
				stores.POST("/:id/manager", g.proxy.ProxyRequest("store", "/api/v1/stores/:id/manager"))
				stores.PUT("/:id/status", g.proxy.ProxyRequest("store", "/api/v1/stores/:id/status"))
				stores.GET("/:id/stock-value", g.proxy.ProxyRequest("store", "/api/v1/stores/:id/stock-value"))
				stores.GET("/:id/stocks", g.proxy.ProxyRequest("store", "/api/v1/stores/:id/stocks"))
			}

			// Stock routes
			stocks := protected.Group("/stocks")
			{
				stocks.GET("", g.proxy.ProxyRequest("stock", "/api/v1/stocks"))
				stocks.GET("/check-stock", g.proxy.ProxyRequest("stock", "/api/v1/stocks/check-stock"))
				stocks.POST("/stock-entries", g.proxy.ProxyRequest("stock", "/api/v1/stocks/stock-entries"))
				stocks.POST("/batch-stock-entries", g.proxy.ProxyRequest("stock", "/api/v1/stocks/batch-stock-entries"))
				stocks.PUT("/:id/location", g.proxy.ProxyRequest("stock", "/api/v1/stocks/:id/location"))
				stocks.GET("/:id/history", g.proxy.ProxyRequest("stock", "/api/v1/stocks/:id/history"))
			}

			// Vendors routes
			vendors := protected.Group("/vendors")
			{
				vendors.POST("", g.proxy.ProxyRequest("vendor", "/api/v1/vendors"))
				vendors.GET("", g.proxy.ProxyRequest("vendor", "/api/v1/vendors"))
				vendors.GET("/:id", g.proxy.ProxyRequest("vendor", "/api/v1/vendors/:id"))
				vendors.PUT("/:id", g.proxy.ProxyRequest("vendor", "/api/v1/vendors/:id"))
				vendors.DELETE("/:id", g.proxy.ProxyRequest("vendor", "/api/v1/vendors/:id"))
				vendors.POST("/products", g.proxy.ProxyRequest("vendor", "/api/v1/vendors/products"))
				vendors.POST("/:id/products/:productId", g.proxy.ProxyRequest("vendor", "/api/v1/vendors/:id/products/:productId"))
				vendors.DELETE("/:id/products/:productId", g.proxy.ProxyRequest("vendor", "/api/v1/vendors/:id/products/:productId"))
				vendors.POST("/:id/contracts", g.proxy.ProxyRequest("vendor", "/api/v1/vendors/:id/contracts"))
				vendors.PUT("/contracts/:contractId", g.proxy.ProxyRequest("vendor", "/api/v1/vendors/contracts/:contractId"))
				vendors.GET("/contracts/:contractId", g.proxy.ProxyRequest("vendor", "/api/v1/vendors/contracts/:contractId"))
				vendors.POST("/:id/ratings", g.proxy.ProxyRequest("vendor", "/api/v1/vendors/:id/ratings"))
				vendors.GET("/:id/ratings", g.proxy.ProxyRequest("vendor", "/api/v1/vendors/:id/ratings"))
			}

			// Manufacturing routes
			manufacturing := protected.Group("/manufacturing")
			{
				manufacturing.POST("/facilities", g.proxy.ProxyRequest("manufacturing", "/api/v1/manufacturing/facilities"))
				manufacturing.GET("/facilities", g.proxy.ProxyRequest("manufacturing", "/api/v1/manufacturing/facilities"))
				manufacturing.GET("/facilities/:id", g.proxy.ProxyRequest("manufacturing", "/api/v1/manufacturing/facilities/:id"))
				manufacturing.POST("/orders", g.proxy.ProxyRequest("manufacturing", "/api/v1/manufacturing/orders"))
				manufacturing.POST("/orders/:id/start", g.proxy.ProxyRequest("manufacturing", "/api/v1/manufacturing/orders/:id/start"))
				manufacturing.PUT("/orders/:id/progress", g.proxy.ProxyRequest("manufacturing", "/api/v1/manufacturing/orders/:id/progress"))
				manufacturing.POST("/bom", g.proxy.ProxyRequest("manufacturing", "/api/v1/manufacturing/bom"))
			}

			// Purchase routes
			purchases := protected.Group("/purchases")
			{
				purchases.POST("/orders", g.proxy.ProxyRequest("purchase", "/api/v1/purchases/orders"))
				purchases.GET("/orders", g.proxy.ProxyRequest("purchase", "/api/v1/purchases/orders"))
				purchases.GET("/orders/:id", g.proxy.ProxyRequest("purchase", "/api/v1/purchases/orders/:id"))
				purchases.PUT("/orders/:id", g.proxy.ProxyRequest("purchase", "/api/v1/purchases/orders/:id"))
				purchases.POST("/orders/:id/approve", g.proxy.ProxyRequest("purchase", "/api/v1/purchases/orders/:id/approve"))
				purchases.POST("/orders/:id/reject", g.proxy.ProxyRequest("purchase", "/api/v1/purchases/orders/:id/reject"))
				purchases.POST("/orders/:id/receive", g.proxy.ProxyRequest("purchase", "/api/v1/purchases/orders/:id/receive"))
			}

			// Order routes
			orders := protected.Group("/orders")
			{
				orders.POST("", g.proxy.ProxyRequest("order", "/api/v1/orders"))
				orders.GET("", g.proxy.ProxyRequest("order", "/api/v1/orders"))
				orders.GET("/:id", g.proxy.ProxyRequest("order", "/api/v1/orders/:id"))
				orders.POST("/:id/confirm", g.proxy.ProxyRequest("order", "/api/v1/orders/:id/confirm"))
				orders.POST("/:id/cancel", g.proxy.ProxyRequest("order", "/api/v1/orders/:id/cancel"))
				orders.POST("/:id/complete", g.proxy.ProxyRequest("order", "/api/v1/orders/:id/complete"))
				orders.POST("/:id/deliveries", g.proxy.ProxyRequest("order", "/api/v1/orders/:id/deliveries"))
				orders.GET("/deliveries", g.proxy.ProxyRequest("order", "/api/v1/orders/deliveries"))
				orders.GET("/deliveries/:id", g.proxy.ProxyRequest("order", "/api/v1/orders/deliveries/:id"))
				orders.POST("/deliveries/:id/prepare", g.proxy.ProxyRequest("order", "/api/v1/orders/deliveries/:id/prepare"))
				orders.POST("/deliveries/:id/ship", g.proxy.ProxyRequest("order", "/api/v1/orders/deliveries/:id/ship"))
				orders.POST("/deliveries/:id/complete", g.proxy.ProxyRequest("order", "/api/v1/orders/deliveries/:id/complete"))
				orders.POST("/:id/invoices", g.proxy.ProxyRequest("order", "/api/v1/orders/:id/invoices"))
				orders.GET("/invoices", g.proxy.ProxyRequest("order", "/api/v1/orders/invoices"))
				orders.GET("/invoices/:id", g.proxy.ProxyRequest("order", "/api/v1/orders/invoices/:id"))
				orders.POST("/invoices/:id/issue", g.proxy.ProxyRequest("order", "/api/v1/orders/invoices/:id/issue"))
				orders.POST("/invoices/:id/pay", g.proxy.ProxyRequest("order", "/api/v1/orders/invoices/:id/pay"))
			}

			// Customer routes
			clients := protected.Group("/clients")
			{
				clients.POST("", g.proxy.ProxyRequest("client", "/api/v1/clients"))
				clients.GET("", g.proxy.ProxyRequest("client", "/api/v1/clients"))
				clients.GET("/:id", g.proxy.ProxyRequest("client", "/api/v1/clients/:id"))
				clients.PUT("/:id", g.proxy.ProxyRequest("client", "/api/v1/clients/:id"))
				clients.DELETE("/:id", g.proxy.ProxyRequest("client", "/api/v1/clients/:id"))
			}

			// Finance routes
			finance := protected.Group("/finance")
			{
				finance.GET("/invoices", g.proxy.ProxyRequest("finance", "/api/v1/finance/invoices"))
				finance.GET("/payments", g.proxy.ProxyRequest("finance", "/api/v1/finance/payments"))
				finance.POST("/payments", g.proxy.ProxyRequest("finance", "/api/v1/finance/payments"))
				finance.GET("/reports/revenue", g.proxy.ProxyRequest("finance", "/api/v1/finance/reports/revenue"))
				finance.GET("/reports/expenses", g.proxy.ProxyRequest("finance", "/api/v1/finance/reports/expenses"))
			}

			// Report routes
			reports := protected.Group("/reports")
			{
				reports.GET("/inventory", g.proxy.ProxyRequest("report", "/api/v1/reports/inventory"))
				reports.GET("/sales", g.proxy.ProxyRequest("report", "/api/v1/reports/sales"))
				reports.GET("/purchases", g.proxy.ProxyRequest("report", "/api/v1/reports/purchases"))
				reports.GET("/manufacturing", g.proxy.ProxyRequest("report", "/api/v1/reports/manufacturing"))
				reports.GET("/custom", g.proxy.ProxyRequest("report", "/api/v1/reports/custom"))
			}

			// SKU routes
			skus := protected.Group("/skus")
			{
				skus.POST("", g.proxy.ProxyRequest("sku", "/api/v1/skus"))
				skus.GET("", g.proxy.ProxyRequest("sku", "/api/v1/skus"))
				skus.GET("/search", g.proxy.ProxyRequest("sku", "/api/v1/skus/search"))
				skus.GET("/:id", g.proxy.ProxyRequest("sku", "/api/v1/skus/:id"))
				skus.GET("/sku/:sku", g.proxy.ProxyRequest("sku", "/api/v1/skus/sku/:sku"))
				skus.PUT("/:id", g.proxy.ProxyRequest("sku", "/api/v1/skus/:id"))
				skus.DELETE("/:id", g.proxy.ProxyRequest("sku", "/api/v1/skus/:id"))
				skus.POST("/bulk", g.proxy.ProxyRequest("sku", "/api/v1/skus/bulk"))
				skus.PUT("/bulk", g.proxy.ProxyRequest("sku", "/api/v1/skus/bulk"))
			}

			// SKU category routes
			skuCategories := protected.Group("/sku-categories")
			{
				skuCategories.POST("", g.proxy.ProxyRequest("sku", "/api/v1/sku-categories"))
				skuCategories.GET("", g.proxy.ProxyRequest("sku", "/api/v1/sku-categories"))
				skuCategories.GET("/tree", g.proxy.ProxyRequest("sku", "/api/v1/sku-categories/tree"))
				skuCategories.GET("/:id", g.proxy.ProxyRequest("sku", "/api/v1/sku-categories/:id"))
				skuCategories.PUT("/:id", g.proxy.ProxyRequest("sku", "/api/v1/sku-categories/:id"))
				skuCategories.DELETE("/:id", g.proxy.ProxyRequest("sku", "/api/v1/sku-categories/:id"))
				skuCategories.GET("/:id/items", g.proxy.ProxyRequest("sku", "/api/v1/sku-categories/:id/skus"))
			}
		}
	}
}

// Start starts the API Gateway
func (g *Gateway) Start() error {
	log.Printf("API Gateway starting on port %s", g.config.Port)
	return g.server.ListenAndServe()
}

// Stop stops the API Gateway
func (g *Gateway) Stop(ctx context.Context) error {
	log.Println("Shutting down API Gateway...")
	return g.server.Shutdown(ctx)
}

// ServiceStatus represents the status of a service
type ServiceStatus struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Status  string `json:"status"`
	Latency int64  `json:"latency_ms"`
}

// CheckServices checks the status of all services
func (g *Gateway) CheckServices() []ServiceStatus {
	var statuses []ServiceStatus
	var wg sync.WaitGroup
	var mu sync.Mutex

	for name, svc := range g.config.Services {
		wg.Add(1)
		go func(name string, svc config.ServiceConfig) {
			defer wg.Add(-1)

			status := "down"
			var latency int64 = -1

			if svc.URL != "" && svc.HealthCheck != "" {
				healthURL := fmt.Sprintf("%s%s", svc.URL, svc.HealthCheck)
				start := time.Now()
				resp, err := http.Get(healthURL)
				if err == nil {
					latency = time.Since(start).Milliseconds()
					if resp.StatusCode >= 200 && resp.StatusCode < 300 {
						status = "up"
					}
					resp.Body.Close()
				}
			}

			mu.Lock()
			statuses = append(statuses, ServiceStatus{
				Name:    name,
				URL:     svc.URL,
				Status:  status,
				Latency: latency,
			})
			mu.Unlock()
		}(name, svc)
	}

	wg.Wait()
	return statuses
}
