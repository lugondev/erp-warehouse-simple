package repository

import (
	"context"
	"time"

	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
	"gorm.io/gorm"
)

// ReportRepository handles database operations for reports and analytics
type ReportRepository struct {
	db *gorm.DB
}

// NewReportRepository creates a new report repository
func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

// CreateReport creates a new report
func (r *ReportRepository) CreateReport(ctx context.Context, report *entity.Report) error {
	return r.db.WithContext(ctx).Create(report).Error
}

// GetReportByID retrieves a report by ID
func (r *ReportRepository) GetReportByID(ctx context.Context, id string) (*entity.Report, error) {
	var report entity.Report
	if err := r.db.WithContext(ctx).First(&report, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &report, nil
}

// UpdateReport updates a report
func (r *ReportRepository) UpdateReport(ctx context.Context, report *entity.Report) error {
	result := r.db.WithContext(ctx).Save(report)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

// DeleteReport deletes a report
func (r *ReportRepository) DeleteReport(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&entity.Report{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

// ListReports lists reports based on filter criteria
func (r *ReportRepository) ListReports(ctx context.Context, filter *entity.ReportFilter) ([]entity.Report, int64, error) {
	var reports []entity.Report
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Report{})

	// Apply filters
	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.StartDate != nil {
		query = query.Where("start_date >= ?", filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("end_date <= ?", filter.EndDate)
	}
	if filter.CreatedBy != nil {
		query = query.Where("created_by = ?", *filter.CreatedBy)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get reports
	if err := query.Order("created_at DESC").Find(&reports).Error; err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}

// CreateReportSchedule creates a new report schedule
func (r *ReportRepository) CreateReportSchedule(ctx context.Context, schedule *entity.ReportSchedule) error {
	return r.db.WithContext(ctx).Create(schedule).Error
}

// GetReportScheduleByID retrieves a report schedule by ID
func (r *ReportRepository) GetReportScheduleByID(ctx context.Context, id string) (*entity.ReportSchedule, error) {
	var schedule entity.ReportSchedule
	if err := r.db.WithContext(ctx).First(&schedule, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &schedule, nil
}

// UpdateReportSchedule updates a report schedule
func (r *ReportRepository) UpdateReportSchedule(ctx context.Context, schedule *entity.ReportSchedule) error {
	result := r.db.WithContext(ctx).Save(schedule)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

// DeleteReportSchedule deletes a report schedule
func (r *ReportRepository) DeleteReportSchedule(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&entity.ReportSchedule{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

// ListReportSchedules lists report schedules based on filter criteria
func (r *ReportRepository) ListReportSchedules(ctx context.Context, filter *entity.ReportScheduleFilter) ([]entity.ReportSchedule, int64, error) {
	var schedules []entity.ReportSchedule
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.ReportSchedule{})

	// Apply filters
	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.Type != "" {
		query = query.Where("report_type = ?", filter.Type)
	}
	if filter.Frequency != nil {
		query = query.Where("frequency = ?", *filter.Frequency)
	}
	if filter.Active != nil {
		query = query.Where("active = ?", *filter.Active)
	}
	if filter.CreatedBy != nil {
		query = query.Where("created_by = ?", *filter.CreatedBy)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get schedules
	if err := query.Order("created_at DESC").Find(&schedules).Error; err != nil {
		return nil, 0, err
	}

	return schedules, total, nil
}

// GetDueSchedules gets report schedules that are due to run
func (r *ReportRepository) GetDueSchedules(ctx context.Context) ([]entity.ReportSchedule, error) {
	var schedules []entity.ReportSchedule
	now := time.Now()

	if err := r.db.WithContext(ctx).
		Where("active = ? AND (next_run_at IS NULL OR next_run_at <= ?)", true, now).
		Find(&schedules).Error; err != nil {
		return nil, err
	}

	return schedules, nil
}

// UpdateScheduleNextRun updates the next run time for a schedule
func (r *ReportRepository) UpdateScheduleNextRun(ctx context.Context, id string, lastRun, nextRun time.Time) error {
	result := r.db.WithContext(ctx).Model(&entity.ReportSchedule{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_run_at": lastRun,
			"next_run_at": nextRun,
			"updated_at":  time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

// GetInventoryValueReport generates an inventory value report
func (r *ReportRepository) GetInventoryValueReport(ctx context.Context, warehouseID string, asOfDate time.Time) ([]entity.InventoryValueReport, error) {
	var report []entity.InventoryValueReport

	query := `
		SELECT 
			i.product_id,
			it.name AS product_name,
			it.sku,
			i.warehouse_id,
			w.name AS warehouse_name,
			i.quantity,
			it.unit_of_measure,
			COALESCE(po.unit_price, it.price) AS unit_cost,
			i.quantity * COALESCE(po.unit_price, it.price) AS total_value
		FROM 
			inventories i
		JOIN 
			items it ON i.product_id = it.id
		JOIN 
			warehouses w ON i.warehouse_id = w.id
		LEFT JOIN (
			SELECT 
				pri.item_id,
				AVG(pri.unit_price) AS unit_price
			FROM 
				purchase_receipts pr
			JOIN 
				purchase_receipt_items pri ON pr.id = pri.purchase_receipt_id
			WHERE 
				pr.receipt_date <= ?
			GROUP BY 
				pri.item_id
		) po ON i.product_id = po.item_id
		WHERE 
			i.quantity > 0
	`

	args := []interface{}{asOfDate}
	if warehouseID != "" {
		query += " AND i.warehouse_id = ?"
		args = append(args, warehouseID)
	}

	query += " ORDER BY it.name"

	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&report).Error; err != nil {
		return nil, err
	}

	return report, nil
}

// GetInventoryAgeReport generates an inventory age report
func (r *ReportRepository) GetInventoryAgeReport(ctx context.Context, warehouseID string, asOfDate time.Time) ([]entity.InventoryAgeReport, error) {
	var report []entity.InventoryAgeReport

	query := `
		SELECT 
			i.product_id,
			it.name AS product_name,
			it.sku,
			i.warehouse_id,
			w.name AS warehouse_name,
			i.quantity,
			it.unit_of_measure,
			COALESCE(se.created_at, i.created_at) AS receipt_date,
			EXTRACT(DAY FROM ? - COALESCE(se.created_at, i.created_at))::int AS days_in_inventory,
			i.quantity * COALESCE(po.unit_price, it.price) AS value
		FROM 
			inventories i
		JOIN 
			items it ON i.product_id = it.id
		JOIN 
			warehouses w ON i.warehouse_id = w.id
		LEFT JOIN (
			SELECT 
				product_id,
				warehouse_id,
				MIN(created_at) AS created_at
			FROM 
				stock_entries
			WHERE 
				type = 'IN'
			GROUP BY 
				product_id, warehouse_id
		) se ON i.product_id = se.product_id AND i.warehouse_id = se.warehouse_id
		LEFT JOIN (
			SELECT 
				pri.item_id,
				AVG(pri.unit_price) AS unit_price
			FROM 
				purchase_receipts pr
			JOIN 
				purchase_receipt_items pri ON pr.id = pri.purchase_receipt_id
			WHERE 
				pr.receipt_date <= ?
			GROUP BY 
				pri.item_id
		) po ON i.product_id = po.item_id
		WHERE 
			i.quantity > 0
	`

	args := []interface{}{asOfDate, asOfDate}
	if warehouseID != "" {
		query += " AND i.warehouse_id = ?"
		args = append(args, warehouseID)
	}

	query += " ORDER BY days_in_inventory DESC"

	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&report).Error; err != nil {
		return nil, err
	}

	return report, nil
}

// GetProductSalesReport generates a product sales report
func (r *ReportRepository) GetProductSalesReport(ctx context.Context, startDate, endDate time.Time) ([]entity.ProductSalesReport, error) {
	var report []entity.ProductSalesReport

	query := `
		SELECT 
			soi.item_id AS product_id,
			it.name AS product_name,
			it.sku,
			it.category,
			SUM(soi.quantity) AS quantity_sold,
			it.unit_of_measure,
			SUM(soi.total_price) AS revenue,
			SUM(soi.quantity * COALESCE(po.unit_price, it.price)) AS cost,
			SUM(soi.total_price) - SUM(soi.quantity * COALESCE(po.unit_price, it.price)) AS profit,
			CASE 
				WHEN SUM(soi.total_price) > 0 THEN 
					(SUM(soi.total_price) - SUM(soi.quantity * COALESCE(po.unit_price, it.price))) / SUM(soi.total_price) * 100
				ELSE 0
			END AS profit_margin
		FROM 
			sales_orders so
		JOIN 
			sales_order_items soi ON so.id = soi.sales_order_id
		JOIN 
			items it ON soi.item_id = it.id
		LEFT JOIN (
			SELECT 
				pri.item_id,
				AVG(pri.unit_price) AS unit_price
			FROM 
				purchase_receipts pr
			JOIN 
				purchase_receipt_items pri ON pr.id = pri.purchase_receipt_id
			GROUP BY 
				pri.item_id
		) po ON soi.item_id = po.item_id
		WHERE 
			so.order_date BETWEEN ? AND ?
			AND so.status NOT IN ('CANCELLED', 'DRAFT')
		GROUP BY 
			soi.item_id, it.name, it.sku, it.category, it.unit_of_measure
		ORDER BY 
			profit DESC
	`

	if err := r.db.WithContext(ctx).Raw(query, startDate, endDate).Scan(&report).Error; err != nil {
		return nil, err
	}

	return report, nil
}

// GetCustomerSalesReport generates a customer sales report
func (r *ReportRepository) GetCustomerSalesReport(ctx context.Context, startDate, endDate time.Time) ([]entity.CustomerSalesReport, error) {
	var report []entity.CustomerSalesReport

	query := `
		SELECT 
			so.customer_id,
			c.name AS customer_name,
			COUNT(DISTINCT so.id) AS order_count,
			SUM(so.grand_total) AS total_revenue,
			SUM(
				(SELECT COALESCE(SUM(soi.quantity * COALESCE(po.unit_price, it.price)), 0)
				FROM sales_order_items soi
				JOIN items it ON soi.item_id = it.id
				LEFT JOIN (
					SELECT 
						pri.item_id,
						AVG(pri.unit_price) AS unit_price
					FROM 
						purchase_receipts pr
					JOIN 
						purchase_receipt_items pri ON pr.id = pri.purchase_receipt_id
					GROUP BY 
						pri.item_id
				) po ON soi.item_id = po.item_id
				WHERE soi.sales_order_id = so.id)
			) AS total_cost,
			SUM(so.grand_total) - SUM(
				(SELECT COALESCE(SUM(soi.quantity * COALESCE(po.unit_price, it.price)), 0)
				FROM sales_order_items soi
				JOIN items it ON soi.item_id = it.id
				LEFT JOIN (
					SELECT 
						pri.item_id,
						AVG(pri.unit_price) AS unit_price
					FROM 
						purchase_receipts pr
					JOIN 
						purchase_receipt_items pri ON pr.id = pri.purchase_receipt_id
					GROUP BY 
						pri.item_id
				) po ON soi.item_id = po.item_id
				WHERE soi.sales_order_id = so.id)
			) AS total_profit,
			CASE 
				WHEN SUM(so.grand_total) > 0 THEN 
					(SUM(so.grand_total) - SUM(
						(SELECT COALESCE(SUM(soi.quantity * COALESCE(po.unit_price, it.price)), 0)
						FROM sales_order_items soi
						JOIN items it ON soi.item_id = it.id
						LEFT JOIN (
							SELECT 
								pri.item_id,
								AVG(pri.unit_price) AS unit_price
							FROM 
								purchase_receipts pr
							JOIN 
								purchase_receipt_items pri ON pr.id = pri.purchase_receipt_id
							GROUP BY 
								pri.item_id
						) po ON soi.item_id = po.item_id
						WHERE soi.sales_order_id = so.id)
					)) / SUM(so.grand_total) * 100
				ELSE 0
			END AS profit_margin
		FROM 
			sales_orders so
		JOIN 
			customers c ON so.customer_id = c.id
		WHERE 
			so.order_date BETWEEN ? AND ?
			AND so.status NOT IN ('CANCELLED', 'DRAFT')
		GROUP BY 
			so.customer_id, c.name
		ORDER BY 
			total_revenue DESC
	`

	if err := r.db.WithContext(ctx).Raw(query, startDate, endDate).Scan(&report).Error; err != nil {
		return nil, err
	}

	return report, nil
}

// GetSupplierPurchaseReport generates a supplier purchase report
func (r *ReportRepository) GetSupplierPurchaseReport(ctx context.Context, startDate, endDate time.Time) ([]entity.SupplierPurchaseReport, error) {
	var report []entity.SupplierPurchaseReport

	query := `
		SELECT 
			po.supplier_id,
			s.name AS supplier_name,
			COUNT(DISTINCT po.id) AS order_count,
			SUM(po.grand_total) AS total_cost
		FROM 
			purchase_orders po
		JOIN 
			suppliers s ON po.supplier_id = s.id
		WHERE 
			po.order_date BETWEEN ? AND ?
			AND po.status NOT IN ('CANCELLED', 'DRAFT')
		GROUP BY 
			po.supplier_id, s.name
		ORDER BY 
			total_cost DESC
	`

	if err := r.db.WithContext(ctx).Raw(query, startDate, endDate).Scan(&report).Error; err != nil {
		return nil, err
	}

	return report, nil
}

// GetProfitAndLossReport generates a profit and loss report
func (r *ReportRepository) GetProfitAndLossReport(ctx context.Context, startDate, endDate time.Time) (*entity.ProfitAndLossReport, error) {
	var report entity.ProfitAndLossReport
	report.StartDate = startDate
	report.EndDate = endDate

	// Get revenue
	revenueQuery := `
		SELECT COALESCE(SUM(grand_total), 0) AS revenue
		FROM sales_orders
		WHERE order_date BETWEEN ? AND ?
		AND status NOT IN ('CANCELLED', 'DRAFT')
	`
	if err := r.db.WithContext(ctx).Raw(revenueQuery, startDate, endDate).Scan(&report.Revenue).Error; err != nil {
		return nil, err
	}

	// Get cost of goods sold
	cogQuery := `
		SELECT COALESCE(SUM(
			(SELECT COALESCE(SUM(soi.quantity * COALESCE(po.unit_price, it.price)), 0)
			FROM sales_order_items soi
			JOIN items it ON soi.item_id = it.id
			LEFT JOIN (
				SELECT 
					pri.item_id,
					AVG(pri.unit_price) AS unit_price
				FROM 
					purchase_receipts pr
				JOIN 
					purchase_receipt_items pri ON pr.id = pri.purchase_receipt_id
				GROUP BY 
					pri.item_id
			) po ON soi.item_id = po.item_id
			WHERE soi.sales_order_id = so.id)
		), 0) AS cost_of_goods
		FROM sales_orders so
		WHERE so.order_date BETWEEN ? AND ?
		AND so.status NOT IN ('CANCELLED', 'DRAFT')
	`
	if err := r.db.WithContext(ctx).Raw(cogQuery, startDate, endDate).Scan(&report.CostOfGoods).Error; err != nil {
		return nil, err
	}

	// Get expenses (from purchase orders not related to inventory)
	expensesQuery := `
		SELECT COALESCE(SUM(grand_total), 0) AS expenses
		FROM purchase_orders
		WHERE order_date BETWEEN ? AND ?
		AND status NOT IN ('CANCELLED', 'DRAFT')
		AND id NOT IN (
			SELECT DISTINCT purchase_order_id 
			FROM purchase_receipts
		)
	`
	if err := r.db.WithContext(ctx).Raw(expensesQuery, startDate, endDate).Scan(&report.Expenses).Error; err != nil {
		return nil, err
	}

	// Calculate gross profit
	report.GrossProfit = report.Revenue - report.CostOfGoods

	// Calculate net profit
	report.NetProfit = report.GrossProfit - report.Expenses

	// Calculate profit margin
	if report.Revenue > 0 {
		report.ProfitMargin = (report.NetProfit / report.Revenue) * 100
	}

	return &report, nil
}

// GetDashboardMetrics generates dashboard metrics
func (r *ReportRepository) GetDashboardMetrics(ctx context.Context, period string) (*entity.DashboardMetrics, error) {
	var metrics entity.DashboardMetrics
	var startDate time.Time

	// Determine start date based on period
	now := time.Now()
	switch period {
	case "day":
		startDate = now.AddDate(0, 0, -1)
	case "week":
		startDate = now.AddDate(0, 0, -7)
	case "month":
		startDate = now.AddDate(0, -1, 0)
	case "quarter":
		startDate = now.AddDate(0, -3, 0)
	case "year":
		startDate = now.AddDate(-1, 0, 0)
	default:
		startDate = now.AddDate(0, -1, 0) // Default to month
	}

	// Get revenue
	revenueQuery := `
		SELECT COALESCE(SUM(grand_total), 0) AS total_revenue
		FROM sales_orders
		WHERE order_date BETWEEN ? AND ?
		AND status NOT IN ('CANCELLED', 'DRAFT')
	`
	if err := r.db.WithContext(ctx).Raw(revenueQuery, startDate, now).Scan(&metrics.TotalRevenue).Error; err != nil {
		return nil, err
	}

	// Get cost
	costQuery := `
		SELECT COALESCE(SUM(grand_total), 0) AS total_cost
		FROM purchase_orders
		WHERE order_date BETWEEN ? AND ?
		AND status NOT IN ('CANCELLED', 'DRAFT')
	`
	if err := r.db.WithContext(ctx).Raw(costQuery, startDate, now).Scan(&metrics.TotalCost).Error; err != nil {
		return nil, err
	}

	// Calculate gross profit
	metrics.GrossProfit = metrics.TotalRevenue - metrics.TotalCost

	// Calculate profit margin
	if metrics.TotalRevenue > 0 {
		metrics.ProfitMargin = (metrics.GrossProfit / metrics.TotalRevenue) * 100
	}

	// Get inventory value
	inventoryValueQuery := `
		SELECT COALESCE(SUM(i.quantity * COALESCE(po.unit_price, it.price)), 0) AS inventory_value,
		       COUNT(DISTINCT i.id) AS inventory_count
		FROM inventories i
		JOIN items it ON i.product_id = it.id
		LEFT JOIN (
			SELECT 
				pri.item_id,
				AVG(pri.unit_price) AS unit_price
			FROM 
				purchase_receipts pr
			JOIN 
				purchase_receipt_items pri ON pr.id = pri.purchase_receipt_id
			GROUP BY 
				pri.item_id
		) po ON i.product_id = po.item_id
		WHERE i.quantity > 0
	`
	var inventoryData struct {
		InventoryValue float64 `gorm:"column:inventory_value"`
		InventoryCount int     `gorm:"column:inventory_count"`
	}
	if err := r.db.WithContext(ctx).Raw(inventoryValueQuery).Scan(&inventoryData).Error; err != nil {
		return nil, err
	}
	metrics.InventoryValue = inventoryData.InventoryValue
	metrics.InventoryCount = inventoryData.InventoryCount

	// Get order counts
	orderCountQuery := `
		SELECT 
			COUNT(CASE WHEN status IN ('DRAFT', 'CONFIRMED', 'PROCESSING') THEN 1 END) AS pending_orders,
			COUNT(CASE WHEN status = 'COMPLETED' THEN 1 END) AS completed_orders
		FROM sales_orders
		WHERE order_date BETWEEN ? AND ?
	`
	var orderCounts struct {
		PendingOrders   int `gorm:"column:pending_orders"`
		CompletedOrders int `gorm:"column:completed_orders"`
	}
	if err := r.db.WithContext(ctx).Raw(orderCountQuery, startDate, now).Scan(&orderCounts).Error; err != nil {
		return nil, err
	}
	metrics.PendingOrders = orderCounts.PendingOrders
	metrics.CompletedOrders = orderCounts.CompletedOrders

	// Get pending purchase orders
	poCountQuery := `
		SELECT COUNT(*) AS pending_purchase_orders
		FROM purchase_orders
		WHERE status IN ('DRAFT', 'SUBMITTED', 'APPROVED', 'SENT')
		AND order_date BETWEEN ? AND ?
	`
	if err := r.db.WithContext(ctx).Raw(poCountQuery, startDate, now).Scan(&metrics.PendingPurchaseOrders).Error; err != nil {
		return nil, err
	}

	// Get top selling products
	topProductsQuery := `
		SELECT 
			soi.item_id AS product_id,
			it.name AS product_name,
			SUM(soi.quantity) AS quantity,
			SUM(soi.total_price) AS revenue
		FROM 
			sales_orders so
		JOIN 
			sales_order_items soi ON so.id = soi.sales_order_id
		JOIN 
			items it ON soi.item_id = it.id
		WHERE 
			so.order_date BETWEEN ? AND ?
			AND so.status NOT IN ('CANCELLED', 'DRAFT')
		GROUP BY 
			soi.item_id, it.name
		ORDER BY 
			revenue DESC
		LIMIT 5
	`
	if err := r.db.WithContext(ctx).Raw(topProductsQuery, startDate, now).Scan(&metrics.TopSellingProducts).Error; err != nil {
		return nil, err
	}

	// Get revenue by month
	revenueByMonthQuery := `
		SELECT 
			TO_CHAR(order_date, 'YYYY-MM') AS month,
			SUM(grand_total) AS revenue
		FROM 
			sales_orders
		WHERE 
			order_date BETWEEN ? AND ?
			AND status NOT IN ('CANCELLED', 'DRAFT')
		GROUP BY 
			TO_CHAR(order_date, 'YYYY-MM')
		ORDER BY 
			month
	`
	var monthlyRevenue []struct {
		Month   string  `gorm:"column:month"`
		Revenue float64 `gorm:"column:revenue"`
	}
	if err := r.db.WithContext(ctx).Raw(revenueByMonthQuery, startDate.AddDate(0, -11, 0), now).Scan(&monthlyRevenue).Error; err != nil {
		return nil, err
	}

	// Convert to map
	metrics.RevenueByMonth = make(map[string]float64)
	for _, mr := range monthlyRevenue {
		metrics.RevenueByMonth[mr.Month] = mr.Revenue
	}

	return &metrics, nil
}
