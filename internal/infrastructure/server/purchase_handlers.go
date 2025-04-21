package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lugondev/erp-warehouse-simple/internal/application/usecase"
	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
)

type PurchaseHandler struct {
	purchaseUseCase *usecase.PurchaseUseCase
}

func NewPurchaseHandler(purchaseUseCase *usecase.PurchaseUseCase) *PurchaseHandler {
	return &PurchaseHandler{
		purchaseUseCase: purchaseUseCase,
	}
}

// RegisterRoutes registers purchase routes
func (h *PurchaseHandler) RegisterRoutes(router *gin.Engine) {
	purchase := router.Group("/api/purchase")
	{
		// Purchase Request routes
		requests := purchase.Group("/requests")
		{
			requests.POST("", h.CreatePurchaseRequest)
			requests.GET("", h.ListPurchaseRequests)
			requests.GET("/:id", h.GetPurchaseRequest)
			requests.PUT("/:id", h.UpdatePurchaseRequest)
			requests.DELETE("/:id", h.DeletePurchaseRequest)
			requests.POST("/:id/submit", h.SubmitPurchaseRequest)
			requests.POST("/:id/approve", h.ApprovePurchaseRequest)
			requests.POST("/:id/reject", h.RejectPurchaseRequest)
			requests.POST("/:id/order", h.CreateOrderFromRequest)
		}

		// Purchase Order routes
		orders := purchase.Group("/orders")
		{
			orders.POST("", h.CreatePurchaseOrder)
			orders.GET("", h.ListPurchaseOrders)
			orders.GET("/:id", h.GetPurchaseOrder)
			orders.PUT("/:id", h.UpdatePurchaseOrder)
			orders.DELETE("/:id", h.DeletePurchaseOrder)
			orders.POST("/:id/submit", h.SubmitPurchaseOrder)
			orders.POST("/:id/approve", h.ApprovePurchaseOrder)
			orders.POST("/:id/send", h.SendPurchaseOrder)
			orders.POST("/:id/confirm", h.ConfirmPurchaseOrder)
			orders.POST("/:id/cancel", h.CancelPurchaseOrder)
			orders.POST("/:id/close", h.ClosePurchaseOrder)
			orders.GET("/:id/receipts", h.ListPurchaseReceiptsByOrder)
			orders.GET("/:id/payments", h.ListPurchasePaymentsByOrder)
			orders.GET("/:id/payment-summary", h.GetPurchaseOrderPaymentSummary)
		}

		// Purchase Receipt routes
		receipts := purchase.Group("/receipts")
		{
			receipts.POST("", h.CreatePurchaseReceipt)
			receipts.GET("/:id", h.GetPurchaseReceipt)
		}

		// Purchase Payment routes
		payments := purchase.Group("/payments")
		{
			payments.POST("", h.CreatePurchasePayment)
			payments.GET("/:id", h.GetPurchasePayment)
		}
	}
}

// Purchase Request Handlers

// @Summary Create a new purchase request
// @Description Create a new purchase request
// @Tags purchase-requests
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body entity.PurchaseRequest true "Purchase request details"
// @Success 201 {object} entity.PurchaseRequest
// @Failure 400 {object} ErrorResponse
// @Router /purchase/requests [post]
func (h *PurchaseHandler) CreatePurchaseRequest(c *gin.Context) {
	var request entity.PurchaseRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "user not authenticated"})
		return
	}
	request.RequesterID = userID.(uint)

	if err := h.purchaseUseCase.CreatePurchaseRequest(c.Request.Context(), &request); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, request)
}

// @Summary Get a purchase request by ID
// @Description Get a purchase request by ID
// @Tags purchase-requests
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Purchase Request ID"
// @Success 200 {object} entity.PurchaseRequest
// @Failure 404 {object} ErrorResponse
// @Router /purchase/requests/{id} [get]
func (h *PurchaseHandler) GetPurchaseRequest(c *gin.Context) {
	id := c.Param("id")

	request, err := h.purchaseUseCase.GetPurchaseRequest(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, request)
}

// @Summary Update a purchase request
// @Description Update a purchase request
// @Tags purchase-requests
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Purchase Request ID"
// @Param request body entity.PurchaseRequest true "Updated purchase request details"
// @Success 200 {object} entity.PurchaseRequest
// @Failure 400 {object} ErrorResponse
// @Router /purchase/requests/{id} [put]
func (h *PurchaseHandler) UpdatePurchaseRequest(c *gin.Context) {
	id := c.Param("id")

	var request entity.PurchaseRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	request.ID = id

	if err := h.purchaseUseCase.UpdatePurchaseRequest(c.Request.Context(), &request); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, request)
}

// @Summary Delete a purchase request
// @Description Delete a purchase request
// @Tags purchase-requests
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Purchase Request ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Router /purchase/requests/{id} [delete]
func (h *PurchaseHandler) DeletePurchaseRequest(c *gin.Context) {
	id := c.Param("id")

	if err := h.purchaseUseCase.DeletePurchaseRequest(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary List purchase requests
// @Description List purchase requests with filters
// @Tags purchase-requests
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request_number query string false "Request number"
// @Param requester_id query integer false "Requester ID"
// @Param status query string false "Status"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param item_id query string false "Item ID"
// @Param page query integer false "Page number"
// @Param page_size query integer false "Page size"
// @Success 200 {object} map[string]interface{}
// @Router /purchase/requests [get]
func (h *PurchaseHandler) ListPurchaseRequests(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	filter := &entity.PurchaseRequestFilter{
		RequestNumber: c.Query("request_number"),
		SKUID:         c.Query("item_id"),
	}

	if requesterID, err := strconv.ParseUint(c.Query("requester_id"), 10, 32); err == nil {
		requesterIDUint := uint(requesterID)
		filter.RequesterID = &requesterIDUint
	}

	if status := c.Query("status"); status != "" {
		requestStatus := entity.PurchaseRequestStatus(status)
		filter.Status = &requestStatus
	}

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			filter.StartDate = &startDate
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			filter.EndDate = &endDate
		}
	}

	requests, total, err := h.purchaseUseCase.ListPurchaseRequests(c.Request.Context(), filter, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"requests":  requests,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// @Summary Submit a purchase request for approval
// @Description Submit a purchase request for approval
// @Tags purchase-requests
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Purchase Request ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Router /purchase/requests/{id}/submit [post]
func (h *PurchaseHandler) SubmitPurchaseRequest(c *gin.Context) {
	id := c.Param("id")

	if err := h.purchaseUseCase.SubmitPurchaseRequest(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Purchase request submitted successfully"})
}

// @Summary Approve a purchase request
// @Description Approve a purchase request
// @Tags purchase-requests
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Purchase Request ID"
// @Param approval body map[string]interface{} true "Approval details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Router /purchase/requests/{id}/approve [post]
func (h *PurchaseHandler) ApprovePurchaseRequest(c *gin.Context) {
	id := c.Param("id")

	var data struct {
		Notes string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "user not authenticated"})
		return
	}

	if err := h.purchaseUseCase.ApprovePurchaseRequest(c.Request.Context(), id, userID.(uint), data.Notes); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Purchase request approved successfully"})
}

// @Summary Reject a purchase request
// @Description Reject a purchase request
// @Tags purchase-requests
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Purchase Request ID"
// @Param rejection body map[string]interface{} true "Rejection details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Router /purchase/requests/{id}/reject [post]
func (h *PurchaseHandler) RejectPurchaseRequest(c *gin.Context) {
	id := c.Param("id")

	var data struct {
		Notes string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "user not authenticated"})
		return
	}

	if err := h.purchaseUseCase.RejectPurchaseRequest(c.Request.Context(), id, userID.(uint), data.Notes); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Purchase request rejected successfully"})
}

// @Summary Create a purchase order from a request
// @Description Create a purchase order from a purchase request
// @Tags purchase-requests
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Purchase Request ID"
// @Param order body map[string]interface{} true "Order details"
// @Success 201 {object} entity.PurchaseOrder
// @Failure 400 {object} ErrorResponse
// @Router /purchase/requests/{id}/order [post]
func (h *PurchaseHandler) CreateOrderFromRequest(c *gin.Context) {
	id := c.Param("id")

	var data struct {
		SupplierID uint `json:"supplier_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "user not authenticated"})
		return
	}

	order, err := h.purchaseUseCase.CreatePurchaseOrderFromRequest(c.Request.Context(), id, data.SupplierID, userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// Purchase Order Handlers

// @Summary Create a new purchase order
// @Description Create a new purchase order
// @Tags purchase-orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param order body entity.PurchaseOrder true "Purchase order details"
// @Success 201 {object} entity.PurchaseOrder
// @Failure 400 {object} ErrorResponse
// @Router /purchase/orders [post]
func (h *PurchaseHandler) CreatePurchaseOrder(c *gin.Context) {
	var order entity.PurchaseOrder
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "user not authenticated"})
		return
	}
	order.CreatedByID = userID.(uint)

	if err := h.purchaseUseCase.CreatePurchaseOrder(c.Request.Context(), &order); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// @Summary Get a purchase order by ID
// @Description Get a purchase order by ID
// @Tags purchase-orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Purchase Order ID"
// @Success 200 {object} entity.PurchaseOrder
// @Failure 404 {object} ErrorResponse
// @Router /purchase/orders/{id} [get]
func (h *PurchaseHandler) GetPurchaseOrder(c *gin.Context) {
	id := c.Param("id")

	order, err := h.purchaseUseCase.GetPurchaseOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// @Summary Update a purchase order
// @Description Update a purchase order
// @Tags purchase-orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Purchase Order ID"
// @Param order body entity.PurchaseOrder true "Updated purchase order details"
// @Success 200 {object} entity.PurchaseOrder
// @Failure 400 {object} ErrorResponse
// @Router /purchase/orders/{id} [put]
func (h *PurchaseHandler) UpdatePurchaseOrder(c *gin.Context) {
	id := c.Param("id")

	var order entity.PurchaseOrder
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	order.ID = id

	if err := h.purchaseUseCase.UpdatePurchaseOrder(c.Request.Context(), &order); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// @Summary Delete a purchase order
// @Description Delete a purchase order
// @Tags purchase-orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Purchase Order ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Router /purchase/orders/{id} [delete]
func (h *PurchaseHandler) DeletePurchaseOrder(c *gin.Context) {
	id := c.Param("id")

	if err := h.purchaseUseCase.DeletePurchaseOrder(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary List purchase orders
// @Description List purchase orders with filters
// @Tags purchase-orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param order_number query string false "Order number"
// @Param supplier_id query integer false "Supplier ID"
// @Param status query string false "Status"
// @Param payment_status query string false "Payment status"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param item_id query string false "Item ID"
// @Param page query integer false "Page number"
// @Param page_size query integer false "Page size"
// @Success 200 {object} map[string]interface{}
// @Router /purchase/orders [get]
func (h *PurchaseHandler) ListPurchaseOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	filter := &entity.PurchaseOrderFilter{
		OrderNumber: c.Query("order_number"),
		SKUID:       c.Query("item_id"),
	}

	if vendorID, err := strconv.ParseUint(c.Query("supplier_id"), 10, 32); err == nil {
		vendorIDUint := uint(vendorID)
		filter.VendorID = &vendorIDUint
	}

	if status := c.Query("status"); status != "" {
		orderStatus := entity.PurchaseOrderStatus(status)
		filter.Status = &orderStatus
	}

	if paymentStatus := c.Query("payment_status"); paymentStatus != "" {
		pmtStatus := entity.PaymentStatus(paymentStatus)
		filter.PaymentStatus = &pmtStatus
	}

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			filter.StartDate = &startDate
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			filter.EndDate = &endDate
		}
	}

	orders, total, err := h.purchaseUseCase.ListPurchaseOrders(c.Request.Context(), filter, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders":    orders,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// @Summary Submit a purchase order for approval
// @Description Submit a purchase order for approval
// @Tags purchase-orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Purchase Order ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Router /purchase/orders/{id}/submit [post]
func (h *PurchaseHandler) SubmitPurchaseOrder(c *gin.Context) {
	id := c.Param("id")

	if err := h.purchaseUseCase.SubmitPurchaseOrder(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Purchase order submitted successfully"})
}

// @Summary Approve a purchase order
// @Description Approve a purchase order
// @Tags purchase-orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Purchase Order ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Router /purchase/orders/{id}/approve [post]
func (h *PurchaseHandler) ApprovePurchaseOrder(c *gin.Context) {
	id := c.Param("id")

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "user not authenticated"})
		return
	}

	if err := h.purchaseUseCase.ApprovePurchaseOrder(c.Request.Context(), id, userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Purchase order approved successfully"})
}

// @Summary Send a purchase order to supplier
// @Description Mark a purchase order as sent to supplier
// @Tags purchase-orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Purchase Order ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Router /purchase/orders/{id}/send [post]
func (h *PurchaseHandler) SendPurchaseOrder(c *gin.Context) {
	id := c.Param("id")

	if err := h.purchaseUseCase.SendPurchaseOrder(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Purchase order marked as sent successfully"})
}

// @Summary Confirm a purchase order
// @Description Mark a purchase order as confirmed by supplier
// @Tags purchase-orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Purchase Order ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Router /purchase/orders/{id}/confirm [post]
func (h *PurchaseHandler) ConfirmPurchaseOrder(c *gin.Context) {
	id := c.Param("id")

	if err := h.purchaseUseCase.ConfirmPurchaseOrder(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Purchase order confirmed successfully"})
}

// @Summary Cancel a purchase order
// @Description Cancel a purchase order
// @Tags purchase-orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Purchase Order ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Router /purchase/orders/{id}/cancel [post]
func (h *PurchaseHandler) CancelPurchaseOrder(c *gin.Context) {
	id := c.Param("id")

	if err := h.purchaseUseCase.CancelPurchaseOrder(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Purchase order cancelled successfully"})
}

// @Summary Close a purchase order
// @Description Close a purchase order
// @Tags purchase-orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Purchase Order ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Router /purchase/orders/{id}/close [post]
func (h *PurchaseHandler) ClosePurchaseOrder(c *gin.Context) {
	id := c.Param("id")

	if err := h.purchaseUseCase.ClosePurchaseOrder(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Purchase order closed successfully"})
}

// Purchase Receipt Handlers

// @Summary Create a new purchase receipt
// @Description Create a new purchase receipt
// @Tags purchase-receipts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param receipt body entity.PurchaseReceipt true "Purchase receipt details"
// @Success 201 {object} entity.PurchaseReceipt
// @Failure 400 {object} ErrorResponse
// @Router /purchase/receipts [post]
func (h *PurchaseHandler) CreatePurchaseReceipt(c *gin.Context) {
	var receipt entity.PurchaseReceipt
	if err := c.ShouldBindJSON(&receipt); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "user not authenticated"})
		return
	}
	receipt.ReceivedByID = userID.(uint)

	if err := h.purchaseUseCase.CreatePurchaseReceipt(c.Request.Context(), &receipt, userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, receipt)
}

// @Summary Get a purchase receipt by ID
// @Description Get a purchase receipt by ID
// @Tags purchase-receipts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Purchase Receipt ID"
// @Success 200 {object} entity.PurchaseReceipt
// @Failure 404 {object} ErrorResponse
// @Router /purchase/receipts/{id} [get]
func (h *PurchaseHandler) GetPurchaseReceipt(c *gin.Context) {
	id := c.Param("id")

	receipt, err := h.purchaseUseCase.GetPurchaseReceipt(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, receipt)
}

// @Summary List purchase receipts by order
// @Description List purchase receipts for a purchase order
// @Tags purchase-orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Purchase Order ID"
// @Success 200 {array} entity.PurchaseReceipt
// @Failure 400 {object} ErrorResponse
// @Router /purchase/orders/{id}/receipts [get]
func (h *PurchaseHandler) ListPurchaseReceiptsByOrder(c *gin.Context) {
	id := c.Param("id")

	receipts, err := h.purchaseUseCase.ListPurchaseReceiptsByOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, receipts)
}

// Purchase Payment Handlers

// @Summary Create a new purchase payment
// @Description Create a new purchase payment
// @Tags purchase-payments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payment body entity.PurchasePayment true "Purchase payment details"
// @Success 201 {object} entity.PurchasePayment
// @Failure 400 {object} ErrorResponse
// @Router /purchase/payments [post]
func (h *PurchaseHandler) CreatePurchasePayment(c *gin.Context) {
	var payment entity.PurchasePayment
	if err := c.ShouldBindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "user not authenticated"})
		return
	}
	payment.CreatedByID = userID.(uint)

	if err := h.purchaseUseCase.CreatePurchasePayment(c.Request.Context(), &payment); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, payment)
}

// @Summary Get a purchase payment by ID
// @Description Get a purchase payment by ID
// @Tags purchase-payments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Purchase Payment ID"
// @Success 200 {object} entity.PurchasePayment
// @Failure 404 {object} ErrorResponse
// @Router /purchase/payments/{id} [get]
func (h *PurchaseHandler) GetPurchasePayment(c *gin.Context) {
	id := c.Param("id")

	payment, err := h.purchaseUseCase.GetPurchasePayment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// @Summary List purchase payments by order
// @Description List purchase payments for a purchase order
// @Tags purchase-orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Purchase Order ID"
// @Success 200 {array} entity.PurchasePayment
// @Failure 400 {object} ErrorResponse
// @Router /purchase/orders/{id}/payments [get]
func (h *PurchaseHandler) ListPurchasePaymentsByOrder(c *gin.Context) {
	id := c.Param("id")

	payments, err := h.purchaseUseCase.ListPurchasePaymentsByOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}

// @Summary Get purchase order payment summary
// @Description Get payment summary for a purchase order
// @Tags purchase-orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Purchase Order ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Router /purchase/orders/{id}/payment-summary [get]
func (h *PurchaseHandler) GetPurchaseOrderPaymentSummary(c *gin.Context) {
	id := c.Param("id")

	summary, err := h.purchaseUseCase.GetPurchaseOrderPaymentSummary(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}
