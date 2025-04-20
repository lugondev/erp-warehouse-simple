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

// ReportHandlers handles report-related HTTP requests
type ReportHandlers struct {
	reportUseCase *usecase.ReportUseCase
}

// NewReportHandlers creates a new report handlers instance
func NewReportHandlers(reportUseCase *usecase.ReportUseCase) *ReportHandlers {
	return &ReportHandlers{
		reportUseCase: reportUseCase,
	}
}

// RegisterRoutes registers report-related routes
func (h *ReportHandlers) RegisterRoutes(router *gin.RouterGroup) {
	reportRouter := router.Group("/reports")
	{
		// Report management
		reportRouter.POST("", middleware.PermissionMiddleware(entity.ReportCreate), h.CreateReport)
		reportRouter.GET("", middleware.PermissionMiddleware(entity.ReportRead), h.ListReports)
		reportRouter.GET("/:id", middleware.PermissionMiddleware(entity.ReportRead), h.GetReport)
		reportRouter.DELETE("/:id", middleware.PermissionMiddleware(entity.ReportDelete), h.DeleteReport)
		reportRouter.POST("/:id/export", middleware.PermissionMiddleware(entity.ReportExport), h.ExportReport)

		// Report schedule management
		reportRouter.POST("/schedules", middleware.PermissionMiddleware(entity.ReportScheduleCreate), h.CreateReportSchedule)
		reportRouter.GET("/schedules", middleware.PermissionMiddleware(entity.ReportScheduleRead), h.ListReportSchedules)
		reportRouter.GET("/schedules/:id", middleware.PermissionMiddleware(entity.ReportScheduleRead), h.GetReportSchedule)
		reportRouter.PUT("/schedules/:id", middleware.PermissionMiddleware(entity.ReportScheduleUpdate), h.UpdateReportSchedule)
		reportRouter.DELETE("/schedules/:id", middleware.PermissionMiddleware(entity.ReportScheduleDelete), h.DeleteReportSchedule)

		// Inventory reports
		reportRouter.GET("/inventory/value", middleware.PermissionMiddleware(entity.ReportRead), h.GetInventoryValueReport)
		reportRouter.GET("/inventory/age", middleware.PermissionMiddleware(entity.ReportRead), h.GetInventoryAgeReport)

		// Sales reports
		reportRouter.GET("/sales/products", middleware.PermissionMiddleware(entity.ReportRead), h.GetProductSalesReport)
		reportRouter.GET("/sales/customers", middleware.PermissionMiddleware(entity.ReportRead), h.GetCustomerSalesReport)

		// Purchase reports
		reportRouter.GET("/purchases/suppliers", middleware.PermissionMiddleware(entity.ReportRead), h.GetSupplierPurchaseReport)

		// Financial reports
		reportRouter.GET("/financial/profit-loss", middleware.PermissionMiddleware(entity.ReportRead), h.GetProfitAndLossReport)

		// Dashboard metrics
		reportRouter.GET("/dashboard/metrics", middleware.PermissionMiddleware(entity.ReportRead), h.GetDashboardMetrics)
	}
}

// CreateReport handles the creation of a new report
// @Summary Create a new report
// @Description Create a new report
// @Tags Reports
// @Accept json
// @Produce json
// @Param report body entity.CreateReportRequest true "Report details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/reports [post]
func (h *ReportHandlers) CreateReport(c *gin.Context) {
	var req entity.CreateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr := auth.GetUserIDFromContext(c)
	userID, _ := strconv.ParseUint(userIDStr, 10, 32)

	report, err := h.reportUseCase.CreateReport(c.Request.Context(), &req, uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"report": report})
}

// GetReport handles the retrieval of a report by ID
// @Summary Get a report by ID
// @Description Get a report by its ID
// @Tags Reports
// @Produce json
// @Param id path string true "Report ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/reports/{id} [get]
func (h *ReportHandlers) GetReport(c *gin.Context) {
	id := c.Param("id")

	report, err := h.reportUseCase.GetReportByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"report": report})
}

// ListReports handles the listing of reports based on filter criteria
// @Summary List reports
// @Description List reports based on filter criteria
// @Tags Reports
// @Produce json
// @Param name query string false "Report name"
// @Param type query string false "Report type"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param status query string false "Report status"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/reports [get]
func (h *ReportHandlers) ListReports(c *gin.Context) {
	filter := &entity.ReportFilter{
		Name: c.Query("name"),
		Type: entity.ReportType(c.Query("type")),
	}

	if status := c.Query("status"); status != "" {
		reportStatus := entity.ReportStatus(status)
		filter.Status = &reportStatus
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

	if createdByStr := c.Query("created_by"); createdByStr != "" {
		if createdBy, err := strconv.ParseUint(createdByStr, 10, 32); err == nil {
			createdByUint := uint(createdBy)
			filter.CreatedBy = &createdByUint
		}
	}

	reports, total, err := h.reportUseCase.ListReports(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reports": reports,
		"total":   total,
	})
}

// DeleteReport handles the deletion of a report
// @Summary Delete a report
// @Description Delete a report by its ID
// @Tags Reports
// @Produce json
// @Param id path string true "Report ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/reports/{id} [delete]
func (h *ReportHandlers) DeleteReport(c *gin.Context) {
	id := c.Param("id")

	if err := h.reportUseCase.DeleteReport(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Report deleted successfully"})
}

// ExportReport handles the export of a report to a specific format
// @Summary Export a report
// @Description Export a report to a specific format (CSV, Excel, PDF)
// @Tags Reports
// @Produce json
// @Param id path string true "Report ID"
// @Param format query string true "Export format (CSV, EXCEL, PDF)"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/reports/{id}/export [post]
func (h *ReportHandlers) ExportReport(c *gin.Context) {
	id := c.Param("id")
	formatStr := c.Query("format")

	if formatStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format is required"})
		return
	}

	format := entity.ReportFormat(formatStr)
	if format != entity.ReportFormatCSV && format != entity.ReportFormatExcel && format != entity.ReportFormatPDF && format != entity.ReportFormatJSON {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format. Supported formats: CSV, EXCEL, PDF, JSON"})
		return
	}

	fileURL, err := h.reportUseCase.ExportReport(c.Request.Context(), id, format)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Report exported successfully",
		"file_url": fileURL,
	})
}

// CreateReportSchedule handles the creation of a new report schedule
// @Summary Create a new report schedule
// @Description Create a new report schedule
// @Tags Reports
// @Accept json
// @Produce json
// @Param schedule body entity.CreateReportScheduleRequest true "Report schedule details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/reports/schedules [post]
func (h *ReportHandlers) CreateReportSchedule(c *gin.Context) {
	var req entity.CreateReportScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr := auth.GetUserIDFromContext(c)
	userID, _ := strconv.ParseUint(userIDStr, 10, 32)

	schedule, err := h.reportUseCase.CreateReportSchedule(c.Request.Context(), &req, uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"schedule": schedule})
}

// GetReportSchedule handles the retrieval of a report schedule by ID
// @Summary Get a report schedule by ID
// @Description Get a report schedule by its ID
// @Tags Reports
// @Produce json
// @Param id path string true "Report Schedule ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/reports/schedules/{id} [get]
func (h *ReportHandlers) GetReportSchedule(c *gin.Context) {
	id := c.Param("id")

	schedule, err := h.reportUseCase.GetReportScheduleByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"schedule": schedule})
}

// UpdateReportSchedule handles the update of a report schedule
// @Summary Update a report schedule
// @Description Update a report schedule
// @Tags Reports
// @Accept json
// @Produce json
// @Param id path string true "Report Schedule ID"
// @Param schedule body entity.UpdateReportScheduleRequest true "Report schedule details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/reports/schedules/{id} [put]
func (h *ReportHandlers) UpdateReportSchedule(c *gin.Context) {
	id := c.Param("id")

	var req entity.UpdateReportScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schedule, err := h.reportUseCase.UpdateReportSchedule(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"schedule": schedule})
}

// ListReportSchedules handles the listing of report schedules based on filter criteria
// @Summary List report schedules
// @Description List report schedules based on filter criteria
// @Tags Reports
// @Produce json
// @Param name query string false "Schedule name"
// @Param type query string false "Report type"
// @Param frequency query string false "Schedule frequency"
// @Param active query bool false "Active status"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/reports/schedules [get]
func (h *ReportHandlers) ListReportSchedules(c *gin.Context) {
	filter := &entity.ReportScheduleFilter{
		Name: c.Query("name"),
		Type: entity.ReportType(c.Query("type")),
	}

	if frequencyStr := c.Query("frequency"); frequencyStr != "" {
		frequency := entity.ReportScheduleFrequency(frequencyStr)
		filter.Frequency = &frequency
	}

	if activeStr := c.Query("active"); activeStr != "" {
		active, err := strconv.ParseBool(activeStr)
		if err == nil {
			filter.Active = &active
		}
	}

	if createdByStr := c.Query("created_by"); createdByStr != "" {
		if createdBy, err := strconv.ParseUint(createdByStr, 10, 32); err == nil {
			createdByUint := uint(createdBy)
			filter.CreatedBy = &createdByUint
		}
	}

	schedules, total, err := h.reportUseCase.ListReportSchedules(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"schedules": schedules,
		"total":     total,
	})
}

// DeleteReportSchedule handles the deletion of a report schedule
// @Summary Delete a report schedule
// @Description Delete a report schedule by its ID
// @Tags Reports
// @Produce json
// @Param id path string true "Report Schedule ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/reports/schedules/{id} [delete]
func (h *ReportHandlers) DeleteReportSchedule(c *gin.Context) {
	id := c.Param("id")

	if err := h.reportUseCase.DeleteReportSchedule(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Report schedule deleted successfully"})
}

// GetInventoryValueReport handles the retrieval of an inventory value report
// @Summary Get inventory value report
// @Description Get inventory value report
// @Tags Reports
// @Produce json
// @Param warehouse_id query string false "Warehouse ID"
// @Param as_of_date query string false "As of date (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/reports/inventory/value [get]
func (h *ReportHandlers) GetInventoryValueReport(c *gin.Context) {
	warehouseID := c.Query("warehouse_id")

	var asOfDate time.Time
	if asOfDateStr := c.Query("as_of_date"); asOfDateStr != "" {
		if date, err := time.Parse("2006-01-02", asOfDateStr); err == nil {
			asOfDate = date
		}
	}

	report, err := h.reportUseCase.GetInventoryValueReport(c.Request.Context(), warehouseID, asOfDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"report": report})
}

// GetInventoryAgeReport handles the retrieval of an inventory age report
// @Summary Get inventory age report
// @Description Get inventory age report
// @Tags Reports
// @Produce json
// @Param warehouse_id query string false "Warehouse ID"
// @Param as_of_date query string false "As of date (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/reports/inventory/age [get]
func (h *ReportHandlers) GetInventoryAgeReport(c *gin.Context) {
	warehouseID := c.Query("warehouse_id")

	var asOfDate time.Time
	if asOfDateStr := c.Query("as_of_date"); asOfDateStr != "" {
		if date, err := time.Parse("2006-01-02", asOfDateStr); err == nil {
			asOfDate = date
		}
	}

	report, err := h.reportUseCase.GetInventoryAgeReport(c.Request.Context(), warehouseID, asOfDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"report": report})
}

// GetProductSalesReport handles the retrieval of a product sales report
// @Summary Get product sales report
// @Description Get product sales report
// @Tags Reports
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/reports/sales/products [get]
func (h *ReportHandlers) GetProductSalesReport(c *gin.Context) {
	var startDate, endDate time.Time

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if date, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = date
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if date, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = date
		}
	}

	report, err := h.reportUseCase.GetProductSalesReport(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"report": report})
}

// GetCustomerSalesReport handles the retrieval of a customer sales report
// @Summary Get customer sales report
// @Description Get customer sales report
// @Tags Reports
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/reports/sales/customers [get]
func (h *ReportHandlers) GetCustomerSalesReport(c *gin.Context) {
	var startDate, endDate time.Time

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if date, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = date
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if date, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = date
		}
	}

	report, err := h.reportUseCase.GetCustomerSalesReport(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"report": report})
}

// GetSupplierPurchaseReport handles the retrieval of a supplier purchase report
// @Summary Get supplier purchase report
// @Description Get supplier purchase report
// @Tags Reports
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/reports/purchases/suppliers [get]
func (h *ReportHandlers) GetSupplierPurchaseReport(c *gin.Context) {
	var startDate, endDate time.Time

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if date, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = date
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if date, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = date
		}
	}

	report, err := h.reportUseCase.GetSupplierPurchaseReport(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"report": report})
}

// GetProfitAndLossReport handles the retrieval of a profit and loss report
// @Summary Get profit and loss report
// @Description Get profit and loss report
// @Tags Reports
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/reports/financial/profit-loss [get]
func (h *ReportHandlers) GetProfitAndLossReport(c *gin.Context) {
	var startDate, endDate time.Time

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if date, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = date
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if date, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = date
		}
	}

	report, err := h.reportUseCase.GetProfitAndLossReport(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"report": report})
}

// GetDashboardMetrics handles the retrieval of dashboard metrics
// @Summary Get dashboard metrics
// @Description Get dashboard metrics
// @Tags Reports
// @Produce json
// @Param period query string false "Period (day, week, month, quarter, year)"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/reports/dashboard/metrics [get]
func (h *ReportHandlers) GetDashboardMetrics(c *gin.Context) {
	period := c.DefaultQuery("period", "month")

	metrics, err := h.reportUseCase.GetDashboardMetrics(c.Request.Context(), period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"metrics": metrics})
}
