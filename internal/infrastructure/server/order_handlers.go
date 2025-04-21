package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lugondev/erp-warehouse-simple/internal/application/usecase"
	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
	"github.com/lugondev/erp-warehouse-simple/internal/infrastructure/auth"
)

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string `json:"message" example:"Operation successful"`
}

type OrderHandlers struct {
	orderUseCase *usecase.OrderUseCase
}

func NewOrderHandlers(orderUseCase *usecase.OrderUseCase) *OrderHandlers {
	return &OrderHandlers{
		orderUseCase: orderUseCase,
	}
}

// CreateSalesOrderRequest represents the request to create a sales order
type CreateSalesOrderRequest struct {
	CustomerID      uint                    `json:"customer_id" binding:"required"`
	Items           []entity.SalesOrderItem `json:"items" binding:"required,dive"`
	ShippingAddress string                  `json:"shipping_address"`
	BillingAddress  string                  `json:"billing_address"`
	PaymentMethod   entity.PaymentMethod    `json:"payment_method"`
	Notes           string                  `json:"notes"`
	WarehouseID     string                  `json:"warehouse_id" binding:"required"`
}

// CreateSalesOrder creates a new sales order
// @Summary Create a new sales order
// @Description Create a new sales order with items
// @Tags orders
// @Accept json
// @Produce json
// @Param order body CreateSalesOrderRequest true "Sales Order"
// @Success 201 {object} entity.SalesOrder
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders [post]
func (h *OrderHandlers) CreateSalesOrder(c *gin.Context) {
	var req CreateSalesOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Get user ID from context
	userID := auth.GetUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	// Create sales order entity
	order := &entity.SalesOrder{
		CustomerID:      req.CustomerID,
		OrderDate:       time.Now(),
		Items:           req.Items,
		ShippingAddress: req.ShippingAddress,
		BillingAddress:  req.BillingAddress,
		PaymentMethod:   req.PaymentMethod,
		Notes:           req.Notes,
		Status:          entity.SalesOrderStatusDraft,
		PaymentStatus:   entity.PaymentStatusPending,
	}

	// Create the order
	if err := h.orderUseCase.CreateSalesOrder(c.Request.Context(), order, req.WarehouseID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// GetSalesOrder gets a sales order by ID
// @Summary Get a sales order
// @Description Get a sales order by ID
// @Tags orders
// @Produce json
// @Param id path string true "Sales Order ID"
// @Success 200 {object} entity.SalesOrder
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/{id} [get]
func (h *OrderHandlers) GetSalesOrder(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id is required"})
		return
	}

	order, err := h.orderUseCase.GetSalesOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// SalesOrderFilter represents the filter for listing sales orders
type SalesOrderFilter struct {
	OrderNumber   string    `form:"order_number"`
	CustomerID    *uint     `form:"customer_id"`
	Status        string    `form:"status"`
	PaymentStatus string    `form:"payment_status"`
	StartDate     time.Time `form:"start_date" time_format:"2006-01-02"`
	EndDate       time.Time `form:"end_date" time_format:"2006-01-02"`
	ItemID        string    `form:"item_id"`
}

// ListSalesOrders lists sales orders with optional filtering
// @Summary List sales orders
// @Description List sales orders with optional filtering
// @Tags orders
// @Produce json
// @Param order_number query string false "Order Number"
// @Param customer_id query integer false "Customer ID"
// @Param status query string false "Order Status"
// @Param payment_status query string false "Payment Status"
// @Param start_date query string false "Start Date (YYYY-MM-DD)"
// @Param end_date query string false "End Date (YYYY-MM-DD)"
// @Param item_id query string false "Item ID"
// @Success 200 {array} entity.SalesOrder
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders [get]
func (h *OrderHandlers) ListSalesOrders(c *gin.Context) {
	var filter SalesOrderFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Convert filter to entity filter
	entityFilter := &entity.SalesOrderFilter{
		OrderNumber: filter.OrderNumber,
		CustomerID:  filter.CustomerID,
		ItemID:      filter.ItemID,
	}

	// Convert string status to entity status if provided
	if filter.Status != "" {
		status := entity.SalesOrderStatus(filter.Status)
		entityFilter.Status = &status
	}

	// Convert string payment status to entity payment status if provided
	if filter.PaymentStatus != "" {
		paymentStatus := entity.PaymentStatus(filter.PaymentStatus)
		entityFilter.PaymentStatus = &paymentStatus
	}

	// Set date filters if provided
	if !filter.StartDate.IsZero() {
		entityFilter.StartDate = &filter.StartDate
	}
	if !filter.EndDate.IsZero() {
		entityFilter.EndDate = &filter.EndDate
	}

	orders, err := h.orderUseCase.ListSalesOrders(c.Request.Context(), entityFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// ConfirmSalesOrder confirms a sales order
// @Summary Confirm a sales order
// @Description Change a sales order status from draft to confirmed
// @Tags orders
// @Produce json
// @Param id path string true "Sales Order ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/{id}/confirm [post]
func (h *OrderHandlers) ConfirmSalesOrder(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id is required"})
		return
	}

	// Get user ID from context
	userID := auth.GetUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	if err := h.orderUseCase.ConfirmSalesOrder(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sales order confirmed successfully"})
}

// CancelSalesOrder cancels a sales order
// @Summary Cancel a sales order
// @Description Cancel a sales order
// @Tags orders
// @Produce json
// @Param id path string true "Sales Order ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/{id}/cancel [post]
func (h *OrderHandlers) CancelSalesOrder(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id is required"})
		return
	}

	if err := h.orderUseCase.CancelSalesOrder(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sales order cancelled successfully"})
}

// CompleteSalesOrder completes a sales order
// @Summary Complete a sales order
// @Description Mark a sales order as completed
// @Tags orders
// @Produce json
// @Param id path string true "Sales Order ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/{id}/complete [post]
func (h *OrderHandlers) CompleteSalesOrder(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id is required"})
		return
	}

	if err := h.orderUseCase.CompleteSalesOrder(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sales order completed successfully"})
}

// CreateDeliveryOrderRequest represents the request to create a delivery order
type CreateDeliveryOrderRequest struct {
	DeliveryDate    time.Time                  `json:"delivery_date" binding:"required"`
	Items           []entity.DeliveryOrderItem `json:"items" binding:"required,dive"`
	ShippingAddress string                     `json:"shipping_address"`
	TrackingNumber  string                     `json:"tracking_number"`
	ShippingMethod  string                     `json:"shipping_method"`
	WarehouseID     string                     `json:"warehouse_id" binding:"required"`
	Notes           string                     `json:"notes"`
}

// CreateDeliveryOrder creates a delivery order for a sales order
// @Summary Create a delivery order
// @Description Create a delivery order for a sales order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Sales Order ID"
// @Param delivery body CreateDeliveryOrderRequest true "Delivery Order"
// @Success 201 {object} entity.DeliveryOrder
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/{id}/deliveries [post]
func (h *OrderHandlers) CreateDeliveryOrder(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "sales order id is required"})
		return
	}

	var req CreateDeliveryOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Get user ID from context
	userID := auth.GetUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	// Create delivery order entity
	delivery := &entity.DeliveryOrder{
		SalesOrderID:    id,
		DeliveryDate:    req.DeliveryDate,
		Items:           req.Items,
		ShippingAddress: req.ShippingAddress,
		TrackingNumber:  req.TrackingNumber,
		ShippingMethod:  req.ShippingMethod,
		WarehouseID:     req.WarehouseID,
		Notes:           req.Notes,
		Status:          entity.DeliveryOrderStatusPending,
	}

	// Create the delivery order
	if err := h.orderUseCase.CreateDeliveryOrder(c.Request.Context(), delivery, userID); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, delivery)
}

// GetDeliveryOrder gets a delivery order by ID
// @Summary Get a delivery order
// @Description Get a delivery order by ID
// @Tags orders
// @Produce json
// @Param id path string true "Delivery Order ID"
// @Success 200 {object} entity.DeliveryOrder
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/deliveries/{id} [get]
func (h *OrderHandlers) GetDeliveryOrder(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id is required"})
		return
	}

	delivery, err := h.orderUseCase.GetDeliveryOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, delivery)
}

// DeliveryOrderFilter represents the filter for listing delivery orders
type DeliveryOrderFilter struct {
	DeliveryNumber string    `form:"delivery_number"`
	SalesOrderID   string    `form:"sales_order_id"`
	Status         string    `form:"status"`
	StartDate      time.Time `form:"start_date" time_format:"2006-01-02"`
	EndDate        time.Time `form:"end_date" time_format:"2006-01-02"`
	WarehouseID    string    `form:"warehouse_id"`
}

// ListDeliveryOrders lists delivery orders with optional filtering
// @Summary List delivery orders
// @Description List delivery orders with optional filtering
// @Tags orders
// @Produce json
// @Param delivery_number query string false "Delivery Number"
// @Param sales_order_id query string false "Sales Order ID"
// @Param status query string false "Delivery Status"
// @Param start_date query string false "Start Date (YYYY-MM-DD)"
// @Param end_date query string false "End Date (YYYY-MM-DD)"
// @Param warehouse_id query string false "Warehouse ID"
// @Success 200 {array} entity.DeliveryOrder
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/deliveries [get]
func (h *OrderHandlers) ListDeliveryOrders(c *gin.Context) {
	var filter DeliveryOrderFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Convert filter to entity filter
	entityFilter := &entity.DeliveryOrderFilter{
		DeliveryNumber: filter.DeliveryNumber,
		SalesOrderID:   filter.SalesOrderID,
		WarehouseID:    filter.WarehouseID,
	}

	// Convert string status to entity status if provided
	if filter.Status != "" {
		status := entity.DeliveryOrderStatus(filter.Status)
		entityFilter.Status = &status
	}

	// Set date filters if provided
	if !filter.StartDate.IsZero() {
		entityFilter.StartDate = &filter.StartDate
	}
	if !filter.EndDate.IsZero() {
		entityFilter.EndDate = &filter.EndDate
	}

	deliveries, err := h.orderUseCase.ListDeliveryOrders(c.Request.Context(), entityFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, deliveries)
}

// PrepareDelivery updates a delivery order status to preparing
// @Summary Prepare a delivery
// @Description Update a delivery order status to preparing
// @Tags orders
// @Produce json
// @Param id path string true "Delivery Order ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/deliveries/{id}/prepare [post]
func (h *OrderHandlers) PrepareDelivery(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id is required"})
		return
	}

	if err := h.orderUseCase.PrepareDelivery(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Delivery prepared successfully"})
}

// ShipDelivery processes a delivery by updating inventory and changing status
// @Summary Ship a delivery
// @Description Process a delivery by updating inventory and changing status
// @Tags orders
// @Produce json
// @Param id path string true "Delivery Order ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/deliveries/{id}/ship [post]
func (h *OrderHandlers) ShipDelivery(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id is required"})
		return
	}

	// Get user ID from context
	userID := auth.GetUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	if err := h.orderUseCase.ShipDelivery(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Delivery shipped successfully"})
}

// CompleteDelivery marks a delivery as delivered
// @Summary Complete a delivery
// @Description Mark a delivery as delivered
// @Tags orders
// @Produce json
// @Param id path string true "Delivery Order ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/deliveries/{id}/complete [post]
func (h *OrderHandlers) CompleteDelivery(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id is required"})
		return
	}

	if err := h.orderUseCase.CompleteDelivery(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Delivery completed successfully"})
}

// CreateInvoiceRequest represents the request to create an invoice
type CreateInvoiceRequest struct {
	DueDate     time.Time `json:"due_date" binding:"required"`
	Amount      float64   `json:"amount"`
	TaxAmount   float64   `json:"tax_amount"`
	TotalAmount float64   `json:"total_amount"`
	Notes       string    `json:"notes"`
}

// CreateInvoice creates an invoice for a sales order
// @Summary Create an invoice
// @Description Create an invoice for a sales order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Sales Order ID"
// @Param invoice body CreateInvoiceRequest true "Invoice"
// @Success 201 {object} entity.Invoice
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/{id}/invoices [post]
func (h *OrderHandlers) CreateInvoice(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "sales order id is required"})
		return
	}

	var req CreateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Get user ID from context
	userID := auth.GetUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	// Create invoice entity
	invoice := &entity.Invoice{
		SalesOrderID: id,
		IssueDate:    time.Now(),
		DueDate:      req.DueDate,
		Amount:       req.Amount,
		TaxAmount:    req.TaxAmount,
		TotalAmount:  req.TotalAmount,
		Notes:        req.Notes,
		Status:       entity.InvoiceStatusDraft,
	}

	// Create the invoice
	if err := h.orderUseCase.CreateInvoice(c.Request.Context(), invoice, userID); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, invoice)
}

// GetInvoice gets an invoice by ID
// @Summary Get an invoice
// @Description Get an invoice by ID
// @Tags orders
// @Produce json
// @Param id path string true "Invoice ID"
// @Success 200 {object} entity.Invoice
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/invoices/{id} [get]
func (h *OrderHandlers) GetInvoice(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id is required"})
		return
	}

	invoice, err := h.orderUseCase.GetInvoice(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, invoice)
}

// InvoiceFilter represents the filter for listing invoices
type InvoiceFilter struct {
	InvoiceNumber string    `form:"invoice_number"`
	SalesOrderID  string    `form:"sales_order_id"`
	Status        string    `form:"status"`
	StartDate     time.Time `form:"start_date" time_format:"2006-01-02"`
	EndDate       time.Time `form:"end_date" time_format:"2006-01-02"`
}

// ListInvoices lists invoices with optional filtering
// @Summary List invoices
// @Description List invoices with optional filtering
// @Tags orders
// @Produce json
// @Param invoice_number query string false "Invoice Number"
// @Param sales_order_id query string false "Sales Order ID"
// @Param status query string false "Invoice Status"
// @Param start_date query string false "Start Date (YYYY-MM-DD)"
// @Param end_date query string false "End Date (YYYY-MM-DD)"
// @Success 200 {array} entity.Invoice
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/invoices [get]
func (h *OrderHandlers) ListInvoices(c *gin.Context) {
	var filter InvoiceFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Convert filter to entity filter
	entityFilter := &entity.InvoiceFilter{
		InvoiceNumber: filter.InvoiceNumber,
		SalesOrderID:  filter.SalesOrderID,
	}

	// Convert string status to entity status if provided
	if filter.Status != "" {
		status := entity.InvoiceStatus(filter.Status)
		entityFilter.Status = &status
	}

	// Set date filters if provided
	if !filter.StartDate.IsZero() {
		entityFilter.StartDate = &filter.StartDate
	}
	if !filter.EndDate.IsZero() {
		entityFilter.EndDate = &filter.EndDate
	}

	invoices, err := h.orderUseCase.ListInvoices(c.Request.Context(), entityFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, invoices)
}

// IssueInvoice changes an invoice from draft to issued status
// @Summary Issue an invoice
// @Description Change an invoice from draft to issued status
// @Tags orders
// @Produce json
// @Param id path string true "Invoice ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/invoices/{id}/issue [post]
func (h *OrderHandlers) IssueInvoice(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id is required"})
		return
	}

	if err := h.orderUseCase.IssueInvoice(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invoice issued successfully"})
}

// PayInvoice marks an invoice as paid
// @Summary Pay an invoice
// @Description Mark an invoice as paid
// @Tags orders
// @Produce json
// @Param id path string true "Invoice ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/invoices/{id}/pay [post]
func (h *OrderHandlers) PayInvoice(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id is required"})
		return
	}

	if err := h.orderUseCase.PayInvoice(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invoice paid successfully"})
}
