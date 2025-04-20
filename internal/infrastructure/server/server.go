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
	warehouseUC     *usecase.WarehouseUseCase
	inventoryUC     *usecase.InventoryUseCase
	supplierUC      *usecase.SupplierUseCase
	manufacturingUC *usecase.ManufacturingUseCase
	itemUC          *usecase.ItemUseCase
	purchaseUC      *usecase.PurchaseUseCase
	orderUC         *usecase.OrderUseCase
	customerUC      usecase.CustomerUseCase
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
	warehouseRepo := repository.NewWarehouseRepository(db)
	inventoryRepo := repository.NewInventoryRepository(db)
	supplierRepo := repository.NewSupplierRepository(db)
	manufacturingRepo := repository.NewManufacturingRepository(db)
	itemRepo := repository.NewItemRepository(db)
	purchaseRepo := repository.NewPurchaseRepository(db)
	orderRepo := repository.NewOrderRepository(db, inventoryRepo)
	customerRepo := repository.NewCustomerRepository(db)
	financeRepo := repository.NewFinanceRepository(db)
	reportRepo := repository.NewReportRepository(db)

	// Initialize use cases
	userUC := usecase.NewUserUseCase(userRepo)
	roleUC := usecase.NewRoleUseCase(roleRepo)
	warehouseUC := usecase.NewWarehouseUseCase(warehouseRepo)
	inventoryUC := usecase.NewInventoryUseCase(inventoryRepo, warehouseRepo)
	supplierUC := usecase.NewSupplierUseCase(supplierRepo)
	manufacturingUC := usecase.NewManufacturingUseCase(manufacturingRepo, inventoryRepo)
	itemUC := usecase.NewItemUseCase(itemRepo)
	purchaseUC := usecase.NewPurchaseUseCase(purchaseRepo, inventoryRepo, supplierRepo, itemRepo)
	orderUC := usecase.NewOrderUseCase(orderRepo, inventoryRepo)
	customerUC := usecase.NewCustomerUseCase(customerRepo)
	financeUC := usecase.NewFinanceUseCase(financeRepo)
	reportUC := usecase.NewReportUseCase(reportRepo, inventoryRepo, orderRepo, purchaseRepo, itemRepo)

	// Initialize services
	jwtService := auth.NewJWTService(cfg.JWT.AccessSecret, cfg.JWT.RefreshSecret)
	auditService := service.NewAuditService(auditRepo)

	// Initialize server
	server := &Server{
		config:          cfg,
		router:          gin.Default(),
		userUC:          userUC,
		roleUC:          roleUC,
		warehouseUC:     warehouseUC,
		inventoryUC:     inventoryUC,
		supplierUC:      supplierUC,
		manufacturingUC: manufacturingUC,
		itemUC:          itemUC,
		purchaseUC:      purchaseUC,
		orderUC:         orderUC,
		customerUC:      customerUC,
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
		warehouseHandler := NewWarehouseHandler(s.warehouseUC, s.inventoryUC)
		inventoryHandler := NewInventoryHandler(s.inventoryUC)
		supplierHandler := NewSupplierHandler(s.supplierUC)
		manufacturingHandler := NewManufacturingHandler(s.manufacturingUC)
		itemHandler := NewItemHandler(s.itemUC)
		purchaseHandler := NewPurchaseHandler(s.purchaseUC)

		// Warehouse routes
		warehouses := protected.Group("/warehouses")
		{
			warehouses.POST("", middleware.PermissionMiddleware(entity.WarehouseCreate), warehouseHandler.CreateWarehouse)
			warehouses.GET("", middleware.PermissionMiddleware(entity.WarehouseRead), warehouseHandler.ListWarehouses)
			warehouses.GET("/:id", middleware.PermissionMiddleware(entity.WarehouseRead), warehouseHandler.GetWarehouse)
			warehouses.PUT("/:id", middleware.PermissionMiddleware(entity.WarehouseUpdate), warehouseHandler.UpdateWarehouse)
			warehouses.DELETE("/:id", middleware.PermissionMiddleware(entity.WarehouseDelete), warehouseHandler.DeleteWarehouse)
		}

		// Inventory routes
		inventory := protected.Group("/inventory")
		{
			inventory.GET("", middleware.PermissionMiddleware(entity.InventoryRead), inventoryHandler.ListInventory)
			inventory.GET("/check-stock", middleware.PermissionMiddleware(entity.InventoryRead), inventoryHandler.CheckStock)
			inventory.POST("/stock-entries", middleware.PermissionMiddleware(entity.StockEntryCreate), inventoryHandler.ProcessStockEntry)
			inventory.POST("/batch-stock-entries", middleware.PermissionMiddleware(entity.StockEntryCreate), inventoryHandler.BatchStockEntry)
			inventory.PUT("/:id/location", middleware.PermissionMiddleware(entity.InventoryUpdate), inventoryHandler.UpdateStockLocation)
			inventory.GET("/:id/history", middleware.PermissionMiddleware(entity.StockEntryRead), inventoryHandler.GetInventoryHistory)
		}

		// Supplier routes
		suppliers := protected.Group("/suppliers")
		{
			suppliers.POST("", middleware.PermissionMiddleware(entity.SupplierCreate), supplierHandler.CreateSupplier)
			suppliers.GET("", middleware.PermissionMiddleware(entity.SupplierRead), supplierHandler.ListSuppliers)
			suppliers.GET("/:id", middleware.PermissionMiddleware(entity.SupplierRead), supplierHandler.GetSupplier)
			suppliers.PUT("/:id", middleware.PermissionMiddleware(entity.SupplierUpdate), supplierHandler.UpdateSupplier)
			suppliers.DELETE("/:id", middleware.PermissionMiddleware(entity.SupplierDelete), supplierHandler.DeleteSupplier)

			// Product management
			suppliers.POST("/products", middleware.PermissionMiddleware(entity.ProductCreate), supplierHandler.CreateProduct)
			suppliers.POST("/:id/products/:productId", middleware.PermissionMiddleware(entity.ProductCreate), supplierHandler.AddProductToSupplier)
			suppliers.DELETE("/:id/products/:productId", middleware.PermissionMiddleware(entity.ProductDelete), supplierHandler.RemoveProductFromSupplier)

			// Contract management
			suppliers.POST("/:id/contracts", middleware.PermissionMiddleware(entity.ContractCreate), supplierHandler.CreateContract)
			suppliers.PUT("/contracts/:contractId", middleware.PermissionMiddleware(entity.ContractUpdate), supplierHandler.UpdateContract)
			suppliers.GET("/contracts/:contractId", middleware.PermissionMiddleware(entity.ContractRead), supplierHandler.GetContract)

			// Rating management
			suppliers.POST("/:id/ratings", middleware.PermissionMiddleware(entity.RatingCreate), supplierHandler.AddRating)
			suppliers.GET("/:id/ratings", middleware.PermissionMiddleware(entity.RatingRead), supplierHandler.GetRatings)
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

		// Item routes
		items := protected.Group("/items")
		{
			items.POST("", middleware.PermissionMiddleware(entity.ProductCreate), itemHandler.CreateItem)
			items.GET("", middleware.PermissionMiddleware(entity.ProductRead), itemHandler.ListItems)
			items.GET("/search", middleware.PermissionMiddleware(entity.ProductRead), itemHandler.SearchItems)
			items.GET("/:id", middleware.PermissionMiddleware(entity.ProductRead), itemHandler.GetItem)
			items.GET("/sku/:sku", middleware.PermissionMiddleware(entity.ProductRead), itemHandler.GetItemBySKU)
			items.PUT("/:id", middleware.PermissionMiddleware(entity.ProductUpdate), itemHandler.UpdateItem)
			items.DELETE("/:id", middleware.PermissionMiddleware(entity.ProductDelete), itemHandler.DeleteItem)
			items.POST("/bulk", middleware.PermissionMiddleware(entity.ProductCreate), itemHandler.BulkCreateItems)
			items.PUT("/bulk", middleware.PermissionMiddleware(entity.ProductUpdate), itemHandler.BulkUpdateItems)
		}

		// Item category routes
		itemCategories := protected.Group("/item-categories")
		{
			itemCategories.POST("", middleware.PermissionMiddleware(entity.ProductCreate), itemHandler.CreateItemCategory)
			itemCategories.GET("", middleware.PermissionMiddleware(entity.ProductRead), itemHandler.ListItemCategories)
			itemCategories.GET("/tree", middleware.PermissionMiddleware(entity.ProductRead), itemHandler.GetItemCategoriesTree)
			itemCategories.GET("/:id", middleware.PermissionMiddleware(entity.ProductRead), itemHandler.GetItemCategory)
			itemCategories.PUT("/:id", middleware.PermissionMiddleware(entity.ProductUpdate), itemHandler.UpdateItemCategory)
			itemCategories.DELETE("/:id", middleware.PermissionMiddleware(entity.ProductDelete), itemHandler.DeleteItemCategory)
			itemCategories.GET("/:id/items", middleware.PermissionMiddleware(entity.ProductRead), itemHandler.GetItemsByCategory)
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

		// Customer routes
		customerHandler := NewCustomerHandler(s.customerUC)
		customerHandler.RegisterRoutes(protected)

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
