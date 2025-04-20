package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// ReportType represents the type of report
type ReportType string

const (
	ReportTypeInventory     ReportType = "INVENTORY"
	ReportTypePurchase      ReportType = "PURCHASE"
	ReportTypeSales         ReportType = "SALES"
	ReportTypeProfitAndLoss ReportType = "PROFIT_LOSS"
	ReportTypeFinancial     ReportType = "FINANCIAL"
	ReportTypeCustom        ReportType = "CUSTOM"
)

// ReportFormat represents the format of a report export
type ReportFormat string

const (
	ReportFormatCSV   ReportFormat = "CSV"
	ReportFormatExcel ReportFormat = "EXCEL"
	ReportFormatPDF   ReportFormat = "PDF"
	ReportFormatJSON  ReportFormat = "JSON"
)

// ReportScheduleFrequency represents how often a report is scheduled
type ReportScheduleFrequency string

const (
	ReportScheduleDaily     ReportScheduleFrequency = "DAILY"
	ReportScheduleWeekly    ReportScheduleFrequency = "WEEKLY"
	ReportScheduleMonthly   ReportScheduleFrequency = "MONTHLY"
	ReportScheduleQuarterly ReportScheduleFrequency = "QUARTERLY"
	ReportScheduleYearly    ReportScheduleFrequency = "YEARLY"
)

// ReportStatus represents the status of a report
type ReportStatus string

const (
	ReportStatusPending   ReportStatus = "PENDING"
	ReportStatusCompleted ReportStatus = "COMPLETED"
	ReportStatusFailed    ReportStatus = "FAILED"
)

// ReportParameters represents the parameters for a report
type ReportParameters map[string]interface{}

// Scan implements the sql.Scanner interface for ReportParameters
func (rp *ReportParameters) Scan(value interface{}) error {
	if value == nil {
		*rp = make(ReportParameters)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan ReportParameters: value is not []byte")
	}

	if err := json.Unmarshal(bytes, rp); err != nil {
		return err
	}
	return nil
}

// Value implements the driver.Valuer interface for ReportParameters
func (rp ReportParameters) Value() (driver.Value, error) {
	if rp == nil {
		return nil, nil
	}
	return json.Marshal(rp)
}

// Report represents a generated report
type Report struct {
	ID          string           `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string           `json:"name" gorm:"not null"`
	Description string           `json:"description"`
	Type        ReportType       `json:"type" gorm:"not null"`
	Parameters  ReportParameters `json:"parameters" gorm:"type:jsonb"`
	StartDate   time.Time        `json:"start_date"`
	EndDate     time.Time        `json:"end_date"`
	CreatedBy   uint             `json:"created_by" gorm:"not null"`
	CreatedAt   time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
	FileURL     string           `json:"file_url"`
	Format      ReportFormat     `json:"format"`
	Status      ReportStatus     `json:"status" gorm:"not null;default:'PENDING'"`
}

// ReportSchedule represents a scheduled report
type ReportSchedule struct {
	ID          string                  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string                  `json:"name" gorm:"not null"`
	Description string                  `json:"description"`
	ReportType  ReportType              `json:"report_type" gorm:"not null"`
	Parameters  ReportParameters        `json:"parameters" gorm:"type:jsonb"`
	Frequency   ReportScheduleFrequency `json:"frequency" gorm:"not null"`
	Format      ReportFormat            `json:"format" gorm:"not null"`
	Active      bool                    `json:"active" gorm:"not null;default:true"`
	Recipients  []string                `json:"recipients" gorm:"type:text[]"`
	CreatedBy   uint                    `json:"created_by" gorm:"not null"`
	CreatedAt   time.Time               `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time               `json:"updated_at" gorm:"autoUpdateTime"`
	LastRunAt   *time.Time              `json:"last_run_at"`
	NextRunAt   *time.Time              `json:"next_run_at"`
}

// InventoryAgeReport represents an inventory age report item
type InventoryAgeReport struct {
	ProductID       string    `json:"product_id"`
	ProductName     string    `json:"product_name"`
	SKU             string    `json:"sku"`
	WarehouseID     string    `json:"warehouse_id"`
	WarehouseName   string    `json:"warehouse_name"`
	Quantity        float64   `json:"quantity"`
	UnitOfMeasure   string    `json:"unit_of_measure"`
	ReceiptDate     time.Time `json:"receipt_date"`
	DaysInInventory int       `json:"days_in_inventory"`
	Value           float64   `json:"value"`
}

// InventoryValueReport represents an inventory value report item
type InventoryValueReport struct {
	ProductID     string  `json:"product_id"`
	ProductName   string  `json:"product_name"`
	SKU           string  `json:"sku"`
	WarehouseID   string  `json:"warehouse_id"`
	WarehouseName string  `json:"warehouse_name"`
	Quantity      float64 `json:"quantity"`
	UnitOfMeasure string  `json:"unit_of_measure"`
	UnitCost      float64 `json:"unit_cost"`
	TotalValue    float64 `json:"total_value"`
}

// ProductSalesReport represents a product sales report item
type ProductSalesReport struct {
	ProductID     string  `json:"product_id"`
	ProductName   string  `json:"product_name"`
	SKU           string  `json:"sku"`
	Category      string  `json:"category"`
	QuantitySold  float64 `json:"quantity_sold"`
	UnitOfMeasure string  `json:"unit_of_measure"`
	Revenue       float64 `json:"revenue"`
	Cost          float64 `json:"cost"`
	Profit        float64 `json:"profit"`
	ProfitMargin  float64 `json:"profit_margin"`
}

// CustomerSalesReport represents a customer sales report item
type CustomerSalesReport struct {
	CustomerID   uint    `json:"customer_id"`
	CustomerName string  `json:"customer_name"`
	OrderCount   int     `json:"order_count"`
	TotalRevenue float64 `json:"total_revenue"`
	TotalCost    float64 `json:"total_cost"`
	TotalProfit  float64 `json:"total_profit"`
	ProfitMargin float64 `json:"profit_margin"`
}

// SupplierPurchaseReport represents a supplier purchase report item
type SupplierPurchaseReport struct {
	SupplierID   uint    `json:"supplier_id"`
	SupplierName string  `json:"supplier_name"`
	OrderCount   int     `json:"order_count"`
	TotalCost    float64 `json:"total_cost"`
}

// ProfitAndLossReport represents a profit and loss report
type ProfitAndLossReport struct {
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	Revenue      float64   `json:"revenue"`
	CostOfGoods  float64   `json:"cost_of_goods"`
	GrossProfit  float64   `json:"gross_profit"`
	Expenses     float64   `json:"expenses"`
	NetProfit    float64   `json:"net_profit"`
	ProfitMargin float64   `json:"profit_margin"`
}

// DashboardMetrics represents key metrics for the dashboard
type DashboardMetrics struct {
	TotalRevenue          float64 `json:"total_revenue"`
	TotalCost             float64 `json:"total_cost"`
	GrossProfit           float64 `json:"gross_profit"`
	ProfitMargin          float64 `json:"profit_margin"`
	InventoryValue        float64 `json:"inventory_value"`
	InventoryCount        int     `json:"inventory_count"`
	PendingOrders         int     `json:"pending_orders"`
	CompletedOrders       int     `json:"completed_orders"`
	PendingPurchaseOrders int     `json:"pending_purchase_orders"`
	TopSellingProducts    []struct {
		ProductID   string  `json:"product_id"`
		ProductName string  `json:"product_name"`
		Quantity    float64 `json:"quantity"`
		Revenue     float64 `json:"revenue"`
	} `json:"top_selling_products"`
	RevenueByMonth map[string]float64 `json:"revenue_by_month"`
}

// ReportFilter represents filters for searching reports
type ReportFilter struct {
	Name      string        `json:"name,omitempty"`
	Type      ReportType    `json:"type,omitempty"`
	StartDate *time.Time    `json:"start_date,omitempty"`
	EndDate   *time.Time    `json:"end_date,omitempty"`
	CreatedBy *uint         `json:"created_by,omitempty"`
	Status    *ReportStatus `json:"status,omitempty"`
}

// ReportScheduleFilter represents filters for searching report schedules
type ReportScheduleFilter struct {
	Name      string                   `json:"name,omitempty"`
	Type      ReportType               `json:"type,omitempty"`
	Frequency *ReportScheduleFrequency `json:"frequency,omitempty"`
	Active    *bool                    `json:"active,omitempty"`
	CreatedBy *uint                    `json:"created_by,omitempty"`
}

// CreateReportRequest represents a request to create a new report
type CreateReportRequest struct {
	Name        string           `json:"name" binding:"required"`
	Description string           `json:"description"`
	Type        ReportType       `json:"type" binding:"required"`
	Parameters  ReportParameters `json:"parameters"`
	StartDate   time.Time        `json:"start_date" binding:"required"`
	EndDate     time.Time        `json:"end_date" binding:"required"`
	Format      ReportFormat     `json:"format" binding:"required"`
}

// CreateReportScheduleRequest represents a request to create a new report schedule
type CreateReportScheduleRequest struct {
	Name        string                  `json:"name" binding:"required"`
	Description string                  `json:"description"`
	ReportType  ReportType              `json:"report_type" binding:"required"`
	Parameters  ReportParameters        `json:"parameters"`
	Frequency   ReportScheduleFrequency `json:"frequency" binding:"required"`
	Format      ReportFormat            `json:"format" binding:"required"`
	Recipients  []string                `json:"recipients" binding:"required"`
}

// UpdateReportScheduleRequest represents a request to update a report schedule
type UpdateReportScheduleRequest struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Parameters  ReportParameters        `json:"parameters"`
	Frequency   ReportScheduleFrequency `json:"frequency"`
	Format      ReportFormat            `json:"format"`
	Active      *bool                   `json:"active"`
	Recipients  []string                `json:"recipients"`
}
