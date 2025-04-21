package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lugondev/erp-warehouse-simple/internal/application/usecase"
	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
	"github.com/lugondev/erp-warehouse-simple/internal/infrastructure/auth"
	"github.com/lugondev/erp-warehouse-simple/internal/infrastructure/config"
	"github.com/lugondev/erp-warehouse-simple/internal/infrastructure/database"
	"github.com/lugondev/erp-warehouse-simple/internal/infrastructure/repository"
	"github.com/lugondev/erp-warehouse-simple/internal/infrastructure/server/middleware"
	"github.com/lugondev/erp-warehouse-simple/internal/infrastructure/service"
)

type Server struct {
	config          *config.Config
	router          *gin.Engine
	userUC          *usecase.UserUseCase
	roleUC          *usecase.RoleUseCase
	storeUC         *usecase.StoreUseCase
	stocksUC        *usecase.StocksUseCase
	vendorUC        *usecase.VendorUseCase
	manufacturingUC *usecase.ManufacturingUseCase
	skuUC           *usecase.SKUUseCase
	purchaseUC      *usecase.PurchaseUseCase
	orderUC         *usecase.OrderUseCase
	clientUC        usecase.ClientUseCase // Changed from *usecase.ClientUseCase
	financeUC       *usecase.FinanceUseCase
	reportUC        *usecase.ReportUseCase
	jwtService      *auth.JWTService
	auditService    *service.AuditService
}

// Router returns the gin engine
func (s *Server) Router() *gin.Engine {
	return s.router
}

func NewServer(cfg *config.Config) (*Server, error) {
	// Initialize database
	db, err := database.NewDatabase(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	auditRepo := repository.NewAuditLogRepository(db)
	storeRepo := repository.NewStoreRepository(db)
	stocksRepo := repository.NewStocksRepository(db)
	vendorRepo := repository.NewVendorRepository(db)
	manufacturingRepo := repository.NewManufacturingRepository(db)
	skuRepo := repository.NewSKURepository(db)
	purchaseRepo := repository.NewPurchaseRepository(db)
	orderRepo := repository.NewOrderRepository(db, stocksRepo)
	clientRepo := repository.NewClientRepository(db)
	financeRepo := repository.NewFinanceRepository(db)
	reportRepo := repository.NewReportRepository(db)

	// Initialize use cases
	userUC := usecase.NewUserUseCase(userRepo)
	roleUC := usecase.NewRoleUseCase(roleRepo)
	storeUC := usecase.NewStoreUseCase(storeRepo)
	stocksUC := usecase.NewStocksUseCase(stocksRepo, storeRepo)
	vendorUC := usecase.NewVendorUseCase(vendorRepo)
	manufacturingUC := usecase.NewManufacturingUseCase(manufacturingRepo, stocksRepo)
	skuUC := usecase.NewSKUUseCase(skuRepo)
	purchaseUC := usecase.NewPurchaseUseCase(purchaseRepo, stocksRepo, vendorRepo, skuRepo)
	orderUC := usecase.NewOrderUseCase(orderRepo, stocksRepo)
	clientUC := usecase.NewClientUseCase(clientRepo)
	financeUC := usecase.NewFinanceUseCase(financeRepo)
	reportUC := usecase.NewReportUseCase(reportRepo, stocksRepo, orderRepo, purchaseRepo, skuRepo)

	// Initialize services
	jwtService := auth.NewJWTService(cfg.JWT.AccessSecret, cfg.JWT.RefreshSecret)
	auditService := service.NewAuditService(auditRepo)

	// Initialize server
	server := &Server{
		config:          cfg,
		router:          gin.Default(),
		userUC:          userUC,
		roleUC:          roleUC,
		storeUC:         storeUC,
		stocksUC:        stocksUC,
		vendorUC:        vendorUC,
		manufacturingUC: manufacturingUC,
		skuUC:           skuUC,
		purchaseUC:      purchaseUC,
		orderUC:         orderUC,
		clientUC:        clientUC, // Using interface instead of pointer
		financeUC:       financeUC,
		reportUC:        reportUC,
		jwtService:      jwtService,
		auditService:    auditService,
	}

	// Setup routes
	server.setupRoutes()

	return server, nil
}

func (s *Server) setupRoutes() {
	// Apply audit logging middleware globally
	s.router.Use(service.CreateAuditLogMiddleware(s.auditService))

	// Health check
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Public routes
	public := s.router.Group("/api/v1")
	{
		auth := public.Group("/auth")
		{
			auth.POST("/register", s.handleRegister)
			auth.POST("/login", s.handleLogin)
			auth.POST("/refresh-token", s.handleRefreshToken)
			auth.POST("/forgot-password", s.handleForgotPassword)
			auth.POST("/reset-password", s.handleResetPassword)
		}
	}

	// Protected routes
	protected := s.router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(s.jwtService))
	{
		// User routes
		users := protected.Group("/users")
		{
			users.GET("", middleware.PermissionMiddleware(entity.UserRead), s.handleListUsers)
			users.GET("/:id", middleware.PermissionMiddleware(entity.UserRead), s.handleGetUser)
			users.PUT("/:id", middleware.PermissionMiddleware(entity.UserUpdate), s.handleUpdateUser)
			users.DELETE("/:id", middleware.PermissionMiddleware(entity.UserDelete), s.handleDeleteUser)
			users.POST("/logout", s.handleLogout)
		}

		// Role routes
		roles := protected.Group("/roles")
		{
			roles.POST("", middleware.PermissionMiddleware(entity.RoleCreate), s.handleCreateRole)
			roles.GET("", middleware.PermissionMiddleware(entity.RoleRead), s.handleListRoles)
			roles.GET("/:id", middleware.PermissionMiddleware(entity.RoleRead), s.handleGetRole)
			roles.PUT("/:id", middleware.PermissionMiddleware(entity.RoleUpdate), s.handleUpdateRole)
			roles.DELETE("/:id", middleware.PermissionMiddleware(entity.RoleDelete), s.handleDeleteRole)
		}

		// Audit log routes
		audit := protected.Group("/audit")
		{
			audit.GET("/logs", middleware.PermissionMiddleware(entity.AuditLogRead), s.handleListAuditLogs)
			audit.GET("/logs/user/:id", middleware.PermissionMiddleware(entity.AuditLogRead), s.handleUserAuditLogs)
		}

		// Initialize handlers
		storeHandler := NewStoreHandler(s.storeUC, s.stocksUC)
		stocksHandler := NewStocksHandler(s.stocksUC)
		vendorHandler := NewVendorHandler(s.vendorUC)
		manufacturingHandler := NewManufacturingHandler(s.manufacturingUC)
		skuHandler := NewSKUHandler(s.skuUC)
		purchaseHandler := NewPurchaseHandler(s.purchaseUC)

		// Store routes
		stores := protected.Group("/stores")
		{
			stores.POST("", middleware.PermissionMiddleware(entity.StoreCreate), storeHandler.CreateStore)
			stores.GET("", middleware.PermissionMiddleware(entity.StoreRead), storeHandler.ListStores)
			stores.GET("/:id", middleware.PermissionMiddleware(entity.StoreRead), storeHandler.GetStore)
			stores.PUT("/:id", middleware.PermissionMiddleware(entity.StoreUpdate), storeHandler.UpdateStore)
			stores.DELETE("/:id", middleware.PermissionMiddleware(entity.StoreDelete), storeHandler.DeleteStore)
		}

		// Stocks routes
		stocks := protected.Group("/stocks")
		{
			stocks.GET("", middleware.PermissionMiddleware(entity.StockRead), stocksHandler.ListStocks)
			stocks.GET("/check-stock", middleware.PermissionMiddleware(entity.StockRead), stocksHandler.CheckStock)
			stocks.POST("/stock-entries", middleware.PermissionMiddleware(entity.StockEntryCreate), stocksHandler.ProcessStockEntry)
			stocks.POST("/batch-stock-entries", middleware.PermissionMiddleware(entity.StockEntryCreate), stocksHandler.BatchStockEntry)
			stocks.PUT("/:id/location", middleware.PermissionMiddleware(entity.StockUpdate), stocksHandler.UpdateStockLocation)
			stocks.GET("/:id/history", middleware.PermissionMiddleware(entity.StockEntryRead), stocksHandler.GetStockHistory)
		}

		// Vendor routes
		vendors := protected.Group("/vendors")
		{
			vendors.POST("", middleware.PermissionMiddleware(entity.VendorCreate), vendorHandler.CreateVendor)
			vendors.GET("", middleware.PermissionMiddleware(entity.VendorRead), vendorHandler.ListVendors)
			vendors.GET("/:id", middleware.PermissionMiddleware(entity.VendorRead), vendorHandler.GetVendor)
			vendors.PUT("/:id", middleware.PermissionMiddleware(entity.VendorUpdate), vendorHandler.UpdateVendor)
			vendors.DELETE("/:id", middleware.PermissionMiddleware(entity.VendorDelete), vendorHandler.DeleteVendor)

			// Product management
			vendors.POST("/products", middleware.PermissionMiddleware(entity.ProductCreate), vendorHandler.CreateProduct)
			vendors.POST("/:id/products/:productId", middleware.PermissionMiddleware(entity.ProductCreate), vendorHandler.AddProductToVendor)
			vendors.DELETE("/:id/products/:productId", middleware.PermissionMiddleware(entity.ProductDelete), vendorHandler.RemoveProductFromVendor)

			// Contract management
			vendors.POST("/:id/contracts", middleware.PermissionMiddleware(entity.ContractCreate), vendorHandler.CreateContract)
			vendors.PUT("/contracts/:contractId", middleware.PermissionMiddleware(entity.ContractUpdate), vendorHandler.UpdateContract)
			vendors.GET("/contracts/:contractId", middleware.PermissionMiddleware(entity.ContractRead), vendorHandler.GetContract)

			// Rating management
			vendors.POST("/:id/ratings", middleware.PermissionMiddleware(entity.RatingCreate), vendorHandler.AddRating)
			vendors.GET("/:id/ratings", middleware.PermissionMiddleware(entity.RatingRead), vendorHandler.GetRatings)
		}

		// Manufacturing routes
		manufacturing := protected.Group("/manufacturing")
		{
			// Facility routes
			manufacturing.POST("/facilities", middleware.PermissionMiddleware(entity.ManufacturingFacilityCreate), manufacturingHandler.CreateFacility)
			manufacturing.GET("/facilities", middleware.PermissionMiddleware(entity.ManufacturingFacilityRead), manufacturingHandler.ListFacilities)
			manufacturing.GET("/facilities/:id", middleware.PermissionMiddleware(entity.ManufacturingFacilityRead), manufacturingHandler.GetFacility)

			// Production routes
			manufacturing.POST("/orders", middleware.PermissionMiddleware(entity.ProductionOrderCreate), manufacturingHandler.CreateProductionOrder)
			manufacturing.POST("/orders/:id/start", middleware.PermissionMiddleware(entity.ProductionOrderUpdate), manufacturingHandler.StartProduction)
			manufacturing.PUT("/orders/:id/progress", middleware.PermissionMiddleware(entity.ProductionOrderUpdate), manufacturingHandler.UpdateProductionProgress)

			// BOM routes
			manufacturing.POST("/bom", middleware.PermissionMiddleware(entity.BOMCreate), manufacturingHandler.CreateBOM)
		}

		// SKU routes
		skus := protected.Group("/skus")
		{
			skus.POST("", middleware.PermissionMiddleware(entity.ProductCreate), skuHandler.CreateSKU)
			skus.GET("", middleware.PermissionMiddleware(entity.ProductRead), skuHandler.ListSKUs)
			skus.GET("/search", middleware.PermissionMiddleware(entity.ProductRead), skuHandler.SearchSKUs)
			skus.GET("/:id", middleware.PermissionMiddleware(entity.ProductRead), skuHandler.GetSKU)
			skus.GET("/code/:code", middleware.PermissionMiddleware(entity.ProductRead), skuHandler.GetSKUByCode)
			skus.PUT("/:id", middleware.PermissionMiddleware(entity.ProductUpdate), skuHandler.UpdateSKU)
			skus.DELETE("/:id", middleware.PermissionMiddleware(entity.ProductDelete), skuHandler.DeleteSKU)
			skus.POST("/bulk", middleware.PermissionMiddleware(entity.ProductCreate), skuHandler.BulkCreateSKUs)
			skus.PUT("/bulk", middleware.PermissionMiddleware(entity.ProductUpdate), skuHandler.BulkUpdateSKUs)
		}

		// SKU category routes
		skuCategories := protected.Group("/sku-categories")
		{
			skuCategories.POST("", middleware.PermissionMiddleware(entity.ProductCreate), skuHandler.CreateSKUCategory)
			skuCategories.GET("", middleware.PermissionMiddleware(entity.ProductRead), skuHandler.ListSKUCategories)
			skuCategories.GET("/tree", middleware.PermissionMiddleware(entity.ProductRead), skuHandler.GetSKUCategoriesTree)
			skuCategories.GET("/:id", middleware.PermissionMiddleware(entity.ProductRead), skuHandler.GetSKUCategory)
			skuCategories.PUT("/:id", middleware.PermissionMiddleware(entity.ProductUpdate), skuHandler.UpdateSKUCategory)
			skuCategories.DELETE("/:id", middleware.PermissionMiddleware(entity.ProductDelete), skuHandler.DeleteSKUCategory)
			skuCategories.GET("/:id/skus", middleware.PermissionMiddleware(entity.ProductRead), skuHandler.GetSKUsByCategory)
		}

		// Purchase routes
		purchaseHandler.RegisterRoutes(s.router)

		// Order routes
		orderHandler := NewOrderHandlers(s.orderUC)
		orders := protected.Group("/orders")
		{
			orders.POST("", middleware.PermissionMiddleware(entity.SalesOrderCreate), orderHandler.CreateSalesOrder)
			orders.GET("", middleware.PermissionMiddleware(entity.SalesOrderRead), orderHandler.ListSalesOrders)
			orders.GET("/:id", middleware.PermissionMiddleware(entity.SalesOrderRead), orderHandler.GetSalesOrder)
			orders.POST("/:id/confirm", middleware.PermissionMiddleware(entity.SalesOrderConfirm), orderHandler.ConfirmSalesOrder)
			orders.POST("/:id/cancel", middleware.PermissionMiddleware(entity.SalesOrderCancel), orderHandler.CancelSalesOrder)
			orders.POST("/:id/complete", middleware.PermissionMiddleware(entity.SalesOrderUpdate), orderHandler.CompleteSalesOrder)

			// Delivery routes
			orders.POST("/:id/deliveries", middleware.PermissionMiddleware(entity.DeliveryOrderCreate), orderHandler.CreateDeliveryOrder)
			orders.GET("/deliveries", middleware.PermissionMiddleware(entity.DeliveryOrderRead), orderHandler.ListDeliveryOrders)
			orders.GET("/deliveries/:id", middleware.PermissionMiddleware(entity.DeliveryOrderRead), orderHandler.GetDeliveryOrder)
			orders.POST("/deliveries/:id/prepare", middleware.PermissionMiddleware(entity.DeliveryOrderProcess), orderHandler.PrepareDelivery)
			orders.POST("/deliveries/:id/ship", middleware.PermissionMiddleware(entity.DeliveryOrderProcess), orderHandler.ShipDelivery)
			orders.POST("/deliveries/:id/complete", middleware.PermissionMiddleware(entity.DeliveryOrderProcess), orderHandler.CompleteDelivery)

			// Invoice routes
			orders.POST("/:id/invoices", middleware.PermissionMiddleware(entity.InvoiceCreate), orderHandler.CreateInvoice)
			orders.GET("/invoices", middleware.PermissionMiddleware(entity.InvoiceRead), orderHandler.ListInvoices)
			orders.GET("/invoices/:id", middleware.PermissionMiddleware(entity.InvoiceRead), orderHandler.GetInvoice)
			orders.POST("/invoices/:id/issue", middleware.PermissionMiddleware(entity.InvoiceIssue), orderHandler.IssueInvoice)
			orders.POST("/invoices/:id/pay", middleware.PermissionMiddleware(entity.InvoicePay), orderHandler.PayInvoice)
		}

		// Client routes
		clientHandler := NewClientHandler(s.clientUC)
		clientHandler.RegisterRoutes(protected)

		// Finance routes
		financeHandler := NewFinanceHandlers(s.financeUC)
		financeHandler.RegisterRoutes(protected)

		// Report routes
		reportHandler := NewReportHandlers(s.reportUC)
		reportHandler.RegisterRoutes(protected)
	}
}

func (s *Server) Run() error {
	return s.router.Run(fmt.Sprintf(":%s", s.config.Server.Port))
}
