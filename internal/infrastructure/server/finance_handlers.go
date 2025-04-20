package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lugondev/erp-warehouse-simple/internal/application/usecase"
	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
	"github.com/lugondev/erp-warehouse-simple/internal/infrastructure/auth"
	"github.com/lugondev/erp-warehouse-simple/internal/infrastructure/server/middleware"
)

// FinanceHandlers handles finance-related HTTP requests
type FinanceHandlers struct {
	financeUseCase *usecase.FinanceUseCase
}

// NewFinanceHandlers creates a new finance handlers instance
func NewFinanceHandlers(financeUseCase *usecase.FinanceUseCase) *FinanceHandlers {
	return &FinanceHandlers{
		financeUseCase: financeUseCase,
	}
}

// RegisterRoutes registers finance-related routes
func (h *FinanceHandlers) RegisterRoutes(router *gin.RouterGroup) {
	financeRouter := router.Group("/finance")
	{
		// Invoice routes
		financeRouter.POST("/invoices", middleware.PermissionMiddleware(entity.FinanceInvoiceCreate), h.CreateInvoice)
		financeRouter.GET("/invoices", middleware.PermissionMiddleware(entity.FinanceInvoiceRead), h.ListInvoices)
		financeRouter.GET("/invoices/:id", middleware.PermissionMiddleware(entity.FinanceInvoiceRead), h.GetInvoice)
		financeRouter.PUT("/invoices/:id", middleware.PermissionMiddleware(entity.FinanceInvoiceUpdate), h.UpdateInvoice)
		financeRouter.PATCH("/invoices/:id/status", middleware.PermissionMiddleware(entity.FinanceInvoiceUpdate), h.UpdateInvoiceStatus)
		financeRouter.POST("/invoices/:id/cancel", middleware.PermissionMiddleware(entity.FinanceInvoiceUpdate), h.CancelInvoice)

		// Payment routes
		financeRouter.POST("/payments", middleware.PermissionMiddleware(entity.FinancePaymentCreate), h.CreatePayment)
		financeRouter.GET("/payments", middleware.PermissionMiddleware(entity.FinancePaymentRead), h.ListPayments)
		financeRouter.GET("/payments/:id", middleware.PermissionMiddleware(entity.FinancePaymentRead), h.GetPayment)
		financeRouter.PUT("/payments/:id", middleware.PermissionMiddleware(entity.FinancePaymentUpdate), h.UpdatePayment)
		financeRouter.POST("/payments/:id/confirm", middleware.PermissionMiddleware(entity.FinancePaymentProcess), h.ConfirmPayment)
		financeRouter.POST("/payments/:id/cancel", middleware.PermissionMiddleware(entity.FinancePaymentProcess), h.CancelPayment)
		financeRouter.POST("/payments/:id/refund", middleware.PermissionMiddleware(entity.FinancePaymentProcess), h.RefundPayment)

		// Report routes
		financeRouter.GET("/reports/accounts-receivable", middleware.PermissionMiddleware(entity.FinanceReportRead), h.GetAccountsReceivable)
		financeRouter.GET("/reports/accounts-payable", middleware.PermissionMiddleware(entity.FinanceReportRead), h.GetAccountsPayable)
		financeRouter.GET("/reports/finance", middleware.PermissionMiddleware(entity.FinanceReportRead), h.GetFinanceReport)
	}
}

// CreateInvoice handles the creation of a new invoice
// @Summary Create a new invoice
// @Description Create a new finance invoice
// @Tags Finance
// @Accept json
// @Produce json
// @Param invoice body entity.CreateFinanceInvoiceRequest true "Invoice details"
// @Success 201 {object} entity.FinanceInvoiceResponse
// @Failure 400 {object} entity.FinanceInvoiceResponse
// @Failure 500 {object} entity.FinanceInvoiceResponse
// @Router /api/finance/invoices [post]
func (h *FinanceHandlers) CreateInvoice(c *gin.Context) {
	var req entity.CreateFinanceInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr := auth.GetUserIDFromContext(c)
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)
	invoice, err := h.financeUseCase.CreateInvoice(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"invoice": invoice})
}

// GetInvoice handles the retrieval of an invoice by ID
// @Summary Get an invoice by ID
// @Description Get a finance invoice by its ID
// @Tags Finance
// @Produce json
// @Param id path int true "Invoice ID"
// @Success 200 {object} entity.FinanceInvoiceResponse
// @Failure 404 {object} entity.FinanceInvoiceResponse
// @Failure 500 {object} entity.FinanceInvoiceResponse
// @Router /api/finance/invoices/{id} [get]
func (h *FinanceHandlers) GetInvoice(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID"})
		return
	}

	invoice, err := h.financeUseCase.GetInvoiceByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"invoice": invoice})
}

// UpdateInvoice handles the update of an invoice
// @Summary Update an invoice
// @Description Update a finance invoice
// @Tags Finance
// @Accept json
// @Produce json
// @Param id path int true "Invoice ID"
// @Param invoice body entity.UpdateFinanceInvoiceRequest true "Invoice details"
// @Success 200 {object} entity.FinanceInvoiceResponse
// @Failure 400 {object} entity.FinanceInvoiceResponse
// @Failure 404 {object} entity.FinanceInvoiceResponse
// @Failure 500 {object} entity.FinanceInvoiceResponse
// @Router /api/finance/invoices/{id} [put]
func (h *FinanceHandlers) UpdateInvoice(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID"})
		return
	}

	var req entity.UpdateFinanceInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	invoice, err := h.financeUseCase.UpdateInvoice(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"invoice": invoice})
}

// UpdateInvoiceStatus handles the update of an invoice status
// @Summary Update an invoice status
// @Description Update a finance invoice status
// @Tags Finance
// @Accept json
// @Produce json
// @Param id path int true "Invoice ID"
// @Param status body map[string]string true "Status"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/finance/invoices/{id}/status [patch]
func (h *FinanceHandlers) UpdateInvoiceStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status := entity.FinanceInvoiceStatus(req.Status)
	if err := h.financeUseCase.UpdateInvoiceStatus(c.Request.Context(), id, status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invoice status updated successfully"})
}

// CancelInvoice handles the cancellation of an invoice
// @Summary Cancel an invoice
// @Description Cancel a finance invoice
// @Tags Finance
// @Produce json
// @Param id path int true "Invoice ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/finance/invoices/{id}/cancel [post]
func (h *FinanceHandlers) CancelInvoice(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID"})
		return
	}

	if err := h.financeUseCase.CancelInvoice(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invoice cancelled successfully"})
}

// ListInvoices handles the listing of invoices based on filter criteria
// @Summary List invoices
// @Description List finance invoices based on filter criteria
// @Tags Finance
// @Produce json
// @Param invoice_number query string false "Invoice number"
// @Param type query string false "Invoice type (SALES/PURCHASE)"
// @Param entity_id query int false "Entity ID"
// @Param entity_type query string false "Entity type (CUSTOMER/SUPPLIER)"
// @Param status query string false "Invoice status"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} entity.FinanceInvoiceListResponse
// @Failure 400 {object} entity.FinanceInvoiceListResponse
// @Failure 500 {object} entity.FinanceInvoiceListResponse
// @Router /api/finance/invoices [get]
func (h *FinanceHandlers) ListInvoices(c *gin.Context) {
	filter := &entity.FinanceInvoiceFilter{
		InvoiceNumber: c.Query("invoice_number"),
		Type:          entity.FinanceInvoiceType(c.Query("type")),
		EntityType:    c.Query("entity_type"),
		Status:        entity.FinanceInvoiceStatus(c.Query("status")),
	}

	if entityID, err := strconv.ParseInt(c.Query("entity_id"), 10, 64); err == nil {
		filter.EntityID = entityID
	}

	if page, err := strconv.Atoi(c.Query("page")); err == nil {
		filter.Page = page
	}

	if pageSize, err := strconv.Atoi(c.Query("page_size")); err == nil {
		filter.PageSize = pageSize
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

	invoices, total, err := h.financeUseCase.ListInvoices(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"invoices":  invoices,
		"total":     total,
		"page":      filter.Page,
		"page_size": filter.PageSize,
	})
}

// CreatePayment handles the creation of a new payment
// @Summary Create a new payment
// @Description Create a new finance payment
// @Tags Finance
// @Accept json
// @Produce json
// @Param payment body entity.CreateFinancePaymentRequest true "Payment details"
// @Success 201 {object} entity.FinancePaymentResponse
// @Failure 400 {object} entity.FinancePaymentResponse
// @Failure 500 {object} entity.FinancePaymentResponse
// @Router /api/finance/payments [post]
func (h *FinanceHandlers) CreatePayment(c *gin.Context) {
	var req entity.CreateFinancePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr := auth.GetUserIDFromContext(c)
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)
	payment, err := h.financeUseCase.CreatePayment(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"payment": payment})
}

// GetPayment handles the retrieval of a payment by ID
// @Summary Get a payment by ID
// @Description Get a finance payment by its ID
// @Tags Finance
// @Produce json
// @Param id path int true "Payment ID"
// @Success 200 {object} entity.FinancePaymentResponse
// @Failure 404 {object} entity.FinancePaymentResponse
// @Failure 500 {object} entity.FinancePaymentResponse
// @Router /api/finance/payments/{id} [get]
func (h *FinanceHandlers) GetPayment(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	payment, err := h.financeUseCase.GetPaymentByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payment": payment})
}

// UpdatePayment handles the update of a payment
// @Summary Update a payment
// @Description Update a finance payment
// @Tags Finance
// @Accept json
// @Produce json
// @Param id path int true "Payment ID"
// @Param payment body entity.UpdateFinancePaymentRequest true "Payment details"
// @Success 200 {object} entity.FinancePaymentResponse
// @Failure 400 {object} entity.FinancePaymentResponse
// @Failure 404 {object} entity.FinancePaymentResponse
// @Failure 500 {object} entity.FinancePaymentResponse
// @Router /api/finance/payments/{id} [put]
func (h *FinanceHandlers) UpdatePayment(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	var req entity.UpdateFinancePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.financeUseCase.UpdatePayment(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payment": payment})
}

// ConfirmPayment handles the confirmation of a payment
// @Summary Confirm a payment
// @Description Confirm a finance payment
// @Tags Finance
// @Produce json
// @Param id path int true "Payment ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/finance/payments/{id}/confirm [post]
func (h *FinanceHandlers) ConfirmPayment(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	if err := h.financeUseCase.ConfirmPayment(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment confirmed successfully"})
}

// CancelPayment handles the cancellation of a payment
// @Summary Cancel a payment
// @Description Cancel a finance payment
// @Tags Finance
// @Produce json
// @Param id path int true "Payment ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/finance/payments/{id}/cancel [post]
func (h *FinanceHandlers) CancelPayment(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	if err := h.financeUseCase.CancelPayment(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment cancelled successfully"})
}

// RefundPayment handles the refund of a payment
// @Summary Refund a payment
// @Description Refund a finance payment
// @Tags Finance
// @Produce json
// @Param id path int true "Payment ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/finance/payments/{id}/refund [post]
func (h *FinanceHandlers) RefundPayment(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	if err := h.financeUseCase.RefundPayment(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment refunded successfully"})
}

// ListPayments handles the listing of payments based on filter criteria
// @Summary List payments
// @Description List finance payments based on filter criteria
// @Tags Finance
// @Produce json
// @Param payment_number query string false "Payment number"
// @Param invoice_id query int false "Invoice ID"
// @Param invoice_number query string false "Invoice number"
// @Param entity_id query int false "Entity ID"
// @Param entity_type query string false "Entity type (CUSTOMER/SUPPLIER)"
// @Param status query string false "Payment status"
// @Param payment_method query string false "Payment method"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} entity.FinancePaymentListResponse
// @Failure 400 {object} entity.FinancePaymentListResponse
// @Failure 500 {object} entity.FinancePaymentListResponse
// @Router /api/finance/payments [get]
func (h *FinanceHandlers) ListPayments(c *gin.Context) {
	filter := &entity.FinancePaymentFilter{
		PaymentNumber: c.Query("payment_number"),
		InvoiceNumber: c.Query("invoice_number"),
		EntityType:    c.Query("entity_type"),
		Status:        entity.FinancePaymentStatus(c.Query("status")),
		PaymentMethod: entity.FinancePaymentMethod(c.Query("payment_method")),
	}

	if invoiceID, err := strconv.ParseInt(c.Query("invoice_id"), 10, 64); err == nil {
		filter.InvoiceID = invoiceID
	}

	if entityID, err := strconv.ParseInt(c.Query("entity_id"), 10, 64); err == nil {
		filter.EntityID = entityID
	}

	if page, err := strconv.Atoi(c.Query("page")); err == nil {
		filter.Page = page
	}

	if pageSize, err := strconv.Atoi(c.Query("page_size")); err == nil {
		filter.PageSize = pageSize
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

	payments, total, err := h.financeUseCase.ListPayments(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payments":  payments,
		"total":     total,
		"page":      filter.Page,
		"page_size": filter.PageSize,
	})
}

// GetAccountsReceivable handles the retrieval of accounts receivable data
// @Summary Get accounts receivable
// @Description Get accounts receivable data
// @Tags Finance
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {array} entity.FinanceAccountsReceivable
// @Failure 500 {object} map[string]string
// @Router /api/finance/reports/accounts-receivable [get]
func (h *FinanceHandlers) GetAccountsReceivable(c *gin.Context) {
	var startDate, endDate *time.Time

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if date, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &date
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if date, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = &date
		}
	}

	receivables, err := h.financeUseCase.GetAccountsReceivable(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"receivables": receivables})
}

// GetAccountsPayable handles the retrieval of accounts payable data
// @Summary Get accounts payable
// @Description Get accounts payable data
// @Tags Finance
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {array} entity.FinanceAccountsPayable
// @Failure 500 {object} map[string]string
// @Router /api/finance/reports/accounts-payable [get]
func (h *FinanceHandlers) GetAccountsPayable(c *gin.Context) {
	var startDate, endDate *time.Time

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if date, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &date
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if date, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = &date
		}
	}

	payables, err := h.financeUseCase.GetAccountsPayable(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payables": payables})
}

// GetFinanceReport handles the retrieval of a finance report
// @Summary Get finance report
// @Description Get a finance report for the specified period
// @Tags Finance
// @Produce json
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} entity.FinanceReport
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/finance/reports/finance [get]
func (h *FinanceHandlers) GetFinanceReport(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start date and end date are required"})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
		return
	}

	report, err := h.financeUseCase.GetFinanceReport(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"report": report})
}
