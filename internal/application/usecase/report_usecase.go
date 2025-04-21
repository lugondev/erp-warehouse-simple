package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
	"github.com/lugondev/erp-warehouse-simple/internal/infrastructure/repository"
)

// ReportUseCase handles business logic for reports and analytics
type ReportUseCase struct {
	reportRepo   *repository.ReportRepository
	stocksRepo   *repository.StocksRepository
	orderRepo    *repository.OrderRepository
	purchaseRepo *repository.PurchaseRepository
	skuRepo      *repository.SKURepository
}

// NewReportUseCase creates a new report use case
func NewReportUseCase(
	reportRepo *repository.ReportRepository,
	stocksRepo *repository.StocksRepository,
	orderRepo *repository.OrderRepository,
	purchaseRepo *repository.PurchaseRepository,
	skuRepo *repository.SKURepository,
) *ReportUseCase {
	return &ReportUseCase{
		reportRepo:   reportRepo,
		stocksRepo:   stocksRepo,
		orderRepo:    orderRepo,
		purchaseRepo: purchaseRepo,
		skuRepo:      skuRepo,
	}
}

// CreateReport creates a new report
func (u *ReportUseCase) CreateReport(ctx context.Context, req *entity.CreateReportRequest, userID uint) (*entity.Report, error) {
	report := &entity.Report{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Parameters:  req.Parameters,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Format:      req.Format,
		Status:      entity.ReportStatusPending,
		CreatedBy:   userID,
	}

	if err := u.reportRepo.CreateReport(ctx, report); err != nil {
		return nil, fmt.Errorf("error creating report: %w", err)
	}

	// Generate report data based on type
	if err := u.generateReport(ctx, report); err != nil {
		// Update report status to failed
		report.Status = entity.ReportStatusFailed
		_ = u.reportRepo.UpdateReport(ctx, report)
		return nil, fmt.Errorf("error generating report: %w", err)
	}

	return report, nil
}

// GetReportByID retrieves a report by ID
func (u *ReportUseCase) GetReportByID(ctx context.Context, id string) (*entity.Report, error) {
	report, err := u.reportRepo.GetReportByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting report: %w", err)
	}
	return report, nil
}

// ListReports lists reports based on filter criteria
func (u *ReportUseCase) ListReports(ctx context.Context, filter *entity.ReportFilter) ([]entity.Report, int64, error) {
	reports, total, err := u.reportRepo.ListReports(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("error listing reports: %w", err)
	}
	return reports, total, nil
}

// DeleteReport deletes a report
func (u *ReportUseCase) DeleteReport(ctx context.Context, id string) error {
	if err := u.reportRepo.DeleteReport(ctx, id); err != nil {
		return fmt.Errorf("error deleting report: %w", err)
	}
	return nil
}

// CreateReportSchedule creates a new report schedule
func (u *ReportUseCase) CreateReportSchedule(ctx context.Context, req *entity.CreateReportScheduleRequest, userID uint) (*entity.ReportSchedule, error) {
	// Calculate next run time based on frequency
	nextRun := u.calculateNextRunTime(time.Now(), req.Frequency)

	schedule := &entity.ReportSchedule{
		Name:        req.Name,
		Description: req.Description,
		ReportType:  req.ReportType,
		Parameters:  req.Parameters,
		Frequency:   req.Frequency,
		Format:      req.Format,
		Active:      true,
		Recipients:  req.Recipients,
		CreatedBy:   userID,
		NextRunAt:   &nextRun,
	}

	if err := u.reportRepo.CreateReportSchedule(ctx, schedule); err != nil {
		return nil, fmt.Errorf("error creating report schedule: %w", err)
	}

	return schedule, nil
}

// GetReportScheduleByID retrieves a report schedule by ID
func (u *ReportUseCase) GetReportScheduleByID(ctx context.Context, id string) (*entity.ReportSchedule, error) {
	schedule, err := u.reportRepo.GetReportScheduleByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting report schedule: %w", err)
	}
	return schedule, nil
}

// UpdateReportSchedule updates a report schedule
func (u *ReportUseCase) UpdateReportSchedule(ctx context.Context, id string, req *entity.UpdateReportScheduleRequest) (*entity.ReportSchedule, error) {
	schedule, err := u.reportRepo.GetReportScheduleByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting report schedule: %w", err)
	}

	// Update fields if provided
	if req.Name != "" {
		schedule.Name = req.Name
	}
	if req.Description != "" {
		schedule.Description = req.Description
	}
	if req.Parameters != nil {
		schedule.Parameters = req.Parameters
	}
	if req.Frequency != "" {
		schedule.Frequency = req.Frequency
		// Recalculate next run time if frequency changed
		nextRun := u.calculateNextRunTime(time.Now(), req.Frequency)
		schedule.NextRunAt = &nextRun
	}
	if req.Format != "" {
		schedule.Format = req.Format
	}
	if req.Active != nil {
		schedule.Active = *req.Active
	}
	if req.Recipients != nil {
		schedule.Recipients = req.Recipients
	}

	if err := u.reportRepo.UpdateReportSchedule(ctx, schedule); err != nil {
		return nil, fmt.Errorf("error updating report schedule: %w", err)
	}

	return schedule, nil
}

// ListReportSchedules lists report schedules based on filter criteria
func (u *ReportUseCase) ListReportSchedules(ctx context.Context, filter *entity.ReportScheduleFilter) ([]entity.ReportSchedule, int64, error) {
	schedules, total, err := u.reportRepo.ListReportSchedules(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("error listing report schedules: %w", err)
	}
	return schedules, total, nil
}

// DeleteReportSchedule deletes a report schedule
func (u *ReportUseCase) DeleteReportSchedule(ctx context.Context, id string) error {
	if err := u.reportRepo.DeleteReportSchedule(ctx, id); err != nil {
		return fmt.Errorf("error deleting report schedule: %w", err)
	}
	return nil
}

// RunScheduledReports runs all due scheduled reports
func (u *ReportUseCase) RunScheduledReports(ctx context.Context) error {
	schedules, err := u.reportRepo.GetDueSchedules(ctx)
	if err != nil {
		return fmt.Errorf("error getting due schedules: %w", err)
	}

	for _, schedule := range schedules {
		// Create report from schedule
		report := &entity.Report{
			Name:        schedule.Name,
			Description: schedule.Description,
			Type:        schedule.ReportType,
			Parameters:  schedule.Parameters,
			StartDate:   time.Now().AddDate(0, -1, 0), // Default to last month
			EndDate:     time.Now(),
			Format:      schedule.Format,
			Status:      entity.ReportStatusPending,
			CreatedBy:   schedule.CreatedBy,
		}

		// Adjust date range based on frequency
		switch schedule.Frequency {
		case entity.ReportScheduleDaily:
			report.StartDate = time.Now().AddDate(0, 0, -1)
		case entity.ReportScheduleWeekly:
			report.StartDate = time.Now().AddDate(0, 0, -7)
		case entity.ReportScheduleMonthly:
			report.StartDate = time.Now().AddDate(0, -1, 0)
		case entity.ReportScheduleQuarterly:
			report.StartDate = time.Now().AddDate(0, -3, 0)
		case entity.ReportScheduleYearly:
			report.StartDate = time.Now().AddDate(-1, 0, 0)
		}

		// Create the report
		if err := u.reportRepo.CreateReport(ctx, report); err != nil {
			continue // Skip to next schedule if this one fails
		}

		// Generate report data
		if err := u.generateReport(ctx, report); err != nil {
			// Update report status to failed
			report.Status = entity.ReportStatusFailed
			_ = u.reportRepo.UpdateReport(ctx, report)
			continue
		}

		// TODO: Send email with report to recipients

		// Update schedule's last run and next run times
		now := time.Now()
		nextRun := u.calculateNextRunTime(now, schedule.Frequency)
		if err := u.reportRepo.UpdateScheduleNextRun(ctx, schedule.ID, now, nextRun); err != nil {
			continue
		}
	}

	return nil
}

// GetInventoryValueReport generates an inventory value report
func (u *ReportUseCase) GetInventoryValueReport(ctx context.Context, warehouseID string, asOfDate time.Time) ([]entity.InventoryValueReport, error) {
	if asOfDate.IsZero() {
		asOfDate = time.Now()
	}

	report, err := u.reportRepo.GetInventoryValueReport(ctx, warehouseID, asOfDate)
	if err != nil {
		return nil, fmt.Errorf("error generating inventory value report: %w", err)
	}

	return report, nil
}

// GetInventoryAgeReport generates an inventory age report
func (u *ReportUseCase) GetInventoryAgeReport(ctx context.Context, warehouseID string, asOfDate time.Time) ([]entity.InventoryAgeReport, error) {
	if asOfDate.IsZero() {
		asOfDate = time.Now()
	}

	report, err := u.reportRepo.GetInventoryAgeReport(ctx, warehouseID, asOfDate)
	if err != nil {
		return nil, fmt.Errorf("error generating inventory age report: %w", err)
	}

	return report, nil
}

// GetProductSalesReport generates a product sales report
func (u *ReportUseCase) GetProductSalesReport(ctx context.Context, startDate, endDate time.Time) ([]entity.ProductSalesReport, error) {
	if startDate.IsZero() {
		startDate = time.Now().AddDate(0, -1, 0) // Default to last month
	}
	if endDate.IsZero() {
		endDate = time.Now()
	}

	report, err := u.reportRepo.GetProductSalesReport(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error generating product sales report: %w", err)
	}

	return report, nil
}

// GetCustomerSalesReport generates a customer sales report
func (u *ReportUseCase) GetCustomerSalesReport(ctx context.Context, startDate, endDate time.Time) ([]entity.CustomerSalesReport, error) {
	if startDate.IsZero() {
		startDate = time.Now().AddDate(0, -1, 0) // Default to last month
	}
	if endDate.IsZero() {
		endDate = time.Now()
	}

	report, err := u.reportRepo.GetCustomerSalesReport(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error generating customer sales report: %w", err)
	}

	return report, nil
}

// GetSupplierPurchaseReport generates a supplier purchase report
func (u *ReportUseCase) GetSupplierPurchaseReport(ctx context.Context, startDate, endDate time.Time) ([]entity.SupplierPurchaseReport, error) {
	if startDate.IsZero() {
		startDate = time.Now().AddDate(0, -1, 0) // Default to last month
	}
	if endDate.IsZero() {
		endDate = time.Now()
	}

	report, err := u.reportRepo.GetSupplierPurchaseReport(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error generating supplier purchase report: %w", err)
	}

	return report, nil
}

// GetProfitAndLossReport generates a profit and loss report
func (u *ReportUseCase) GetProfitAndLossReport(ctx context.Context, startDate, endDate time.Time) (*entity.ProfitAndLossReport, error) {
	if startDate.IsZero() {
		startDate = time.Now().AddDate(0, -1, 0) // Default to last month
	}
	if endDate.IsZero() {
		endDate = time.Now()
	}

	report, err := u.reportRepo.GetProfitAndLossReport(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error generating profit and loss report: %w", err)
	}

	return report, nil
}

// GetDashboardMetrics generates dashboard metrics
func (u *ReportUseCase) GetDashboardMetrics(ctx context.Context, period string) (*entity.DashboardMetrics, error) {
	metrics, err := u.reportRepo.GetDashboardMetrics(ctx, period)
	if err != nil {
		return nil, fmt.Errorf("error generating dashboard metrics: %w", err)
	}

	return metrics, nil
}

// ExportReport exports a report to the specified format
func (u *ReportUseCase) ExportReport(ctx context.Context, reportID string, format entity.ReportFormat) (string, error) {
	report, err := u.reportRepo.GetReportByID(ctx, reportID)
	if err != nil {
		return "", fmt.Errorf("error getting report: %w", err)
	}

	// TODO: Implement export functionality for different formats
	// This would generate the file and return the file URL

	// For now, just update the report with a dummy file URL
	fileURL := fmt.Sprintf("/reports/%s.%s", report.ID, string(format))
	report.FileURL = fileURL
	report.Format = format

	if err := u.reportRepo.UpdateReport(ctx, report); err != nil {
		return "", fmt.Errorf("error updating report: %w", err)
	}

	return fileURL, nil
}

// Helper functions

// generateReport generates the report data based on the report type
func (u *ReportUseCase) generateReport(ctx context.Context, report *entity.Report) error {
	var err error

	switch report.Type {
	case entity.ReportTypeInventory:
		// Check if it's an inventory value or age report
		reportSubtype, ok := report.Parameters["subtype"]
		if !ok {
			reportSubtype = "value" // Default to value report
		}

		warehouseID, _ := report.Parameters["warehouse_id"].(string)

		if reportSubtype == "age" {
			_, err = u.GetInventoryAgeReport(ctx, warehouseID, report.EndDate)
		} else {
			_, err = u.GetInventoryValueReport(ctx, warehouseID, report.EndDate)
		}

	case entity.ReportTypeSales:
		// Check if it's a product or customer sales report
		reportSubtype, ok := report.Parameters["subtype"]
		if !ok {
			reportSubtype = "product" // Default to product report
		}

		if reportSubtype == "customer" {
			_, err = u.GetCustomerSalesReport(ctx, report.StartDate, report.EndDate)
		} else {
			_, err = u.GetProductSalesReport(ctx, report.StartDate, report.EndDate)
		}

	case entity.ReportTypePurchase:
		_, err = u.GetSupplierPurchaseReport(ctx, report.StartDate, report.EndDate)

	case entity.ReportTypeProfitAndLoss:
		_, err = u.GetProfitAndLossReport(ctx, report.StartDate, report.EndDate)

	case entity.ReportTypeFinancial:
		// Financial reports are handled by the finance use case
		// This is just a placeholder
		err = nil

	case entity.ReportTypeCustom:
		// Custom reports would be implemented based on parameters
		err = nil

	default:
		err = fmt.Errorf("unsupported report type: %s", report.Type)
	}

	if err != nil {
		return err
	}

	// Update report status to completed
	report.Status = entity.ReportStatusCompleted
	return u.reportRepo.UpdateReport(ctx, report)
}

// calculateNextRunTime calculates the next run time based on frequency
func (u *ReportUseCase) calculateNextRunTime(from time.Time, frequency entity.ReportScheduleFrequency) time.Time {
	switch frequency {
	case entity.ReportScheduleDaily:
		return from.AddDate(0, 0, 1)
	case entity.ReportScheduleWeekly:
		return from.AddDate(0, 0, 7)
	case entity.ReportScheduleMonthly:
		return from.AddDate(0, 1, 0)
	case entity.ReportScheduleQuarterly:
		return from.AddDate(0, 3, 0)
	case entity.ReportScheduleYearly:
		return from.AddDate(1, 0, 0)
	default:
		return from.AddDate(0, 1, 0) // Default to monthly
	}
}
