package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// FinanceInvoiceType represents the type of invoice (sales or purchase)
type FinanceInvoiceType string

const (
	FinanceSalesInvoice    FinanceInvoiceType = "SALES"
	FinancePurchaseInvoice FinanceInvoiceType = "PURCHASE"
)

// FinanceInvoiceStatus represents the status of a finance invoice
type FinanceInvoiceStatus string

const (
	FinanceInvoiceDraft         FinanceInvoiceStatus = "DRAFT"
	FinanceInvoicePending       FinanceInvoiceStatus = "PENDING"
	FinanceInvoiceApproved      FinanceInvoiceStatus = "APPROVED"
	FinanceInvoicePaid          FinanceInvoiceStatus = "PAID"
	FinanceInvoicePartiallyPaid FinanceInvoiceStatus = "PARTIALLY_PAID"
	FinanceInvoiceCancelled     FinanceInvoiceStatus = "CANCELLED"
	FinanceInvoiceOverdue       FinanceInvoiceStatus = "OVERDUE"
)

// FinanceInvoiceItem represents a line item in a finance invoice
type FinanceInvoiceItem struct {
	ID          int64     `json:"id" db:"id"`
	InvoiceID   int64     `json:"invoice_id" db:"invoice_id"`
	ProductID   int64     `json:"product_id" db:"product_id"`
	ProductName string    `json:"product_name" db:"product_name"`
	Quantity    float64   `json:"quantity" db:"quantity"`
	UnitPrice   float64   `json:"unit_price" db:"unit_price"`
	TaxRate     float64   `json:"tax_rate" db:"tax_rate"`
	TaxAmount   float64   `json:"tax_amount" db:"tax_amount"`
	Subtotal    float64   `json:"subtotal" db:"subtotal"`
	Total       float64   `json:"total" db:"total"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// FinanceInvoiceItems is a slice of FinanceInvoiceItem
type FinanceInvoiceItems []FinanceInvoiceItem

// Scan implements the sql.Scanner interface for FinanceInvoiceItems
func (fii *FinanceInvoiceItems) Scan(value interface{}) error {
	if value == nil {
		*fii = make(FinanceInvoiceItems, 0)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan FinanceInvoiceItems: value is not []byte")
	}

	if err := json.Unmarshal(bytes, fii); err != nil {
		return err
	}
	return nil
}

// Value implements the driver.Valuer interface for FinanceInvoiceItems
func (fii FinanceInvoiceItems) Value() (driver.Value, error) {
	if fii == nil {
		return nil, nil
	}
	return json.Marshal(fii)
}

// FinanceInvoice represents a finance invoice entity
type FinanceInvoice struct {
	ID             int64                `json:"id" db:"id"`
	InvoiceNumber  string               `json:"invoice_number" db:"invoice_number"`
	Type           FinanceInvoiceType   `json:"type" db:"type"`
	ReferenceID    string               `json:"reference_id" db:"reference_id"`
	EntityID       int64                `json:"entity_id" db:"entity_id"`
	EntityType     string               `json:"entity_type" db:"entity_type"` // "CUSTOMER" or "SUPPLIER"
	EntityName     string               `json:"entity_name" db:"entity_name"`
	IssueDate      time.Time            `json:"issue_date" db:"issue_date"`
	DueDate        time.Time            `json:"due_date" db:"due_date"`
	Items          FinanceInvoiceItems  `json:"items" db:"items"`
	Subtotal       float64              `json:"subtotal" db:"subtotal"`
	TaxTotal       float64              `json:"tax_total" db:"tax_total"`
	DiscountAmount float64              `json:"discount_amount" db:"discount_amount"`
	Total          float64              `json:"total" db:"total"`
	AmountPaid     float64              `json:"amount_paid" db:"amount_paid"`
	AmountDue      float64              `json:"amount_due" db:"amount_due"`
	Status         FinanceInvoiceStatus `json:"status" db:"status"`
	Notes          string               `json:"notes" db:"notes"`
	CreatedBy      int64                `json:"created_by" db:"created_by"`
	CreatedAt      time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at" db:"updated_at"`
}

// FinanceInvoiceFilter represents filters for querying finance invoices
type FinanceInvoiceFilter struct {
	InvoiceNumber string               `json:"invoice_number,omitempty"`
	Type          FinanceInvoiceType   `json:"type,omitempty"`
	ReferenceID   string               `json:"reference_id,omitempty"`
	EntityID      int64                `json:"entity_id,omitempty"`
	EntityType    string               `json:"entity_type,omitempty"`
	Status        FinanceInvoiceStatus `json:"status,omitempty"`
	StartDate     *time.Time           `json:"start_date,omitempty"`
	EndDate       *time.Time           `json:"end_date,omitempty"`
	Page          int                  `json:"page,omitempty"`
	PageSize      int                  `json:"page_size,omitempty"`
}

// CreateFinanceInvoiceRequest represents the request to create a new finance invoice
type CreateFinanceInvoiceRequest struct {
	Type           FinanceInvoiceType  `json:"type" binding:"required,oneof=SALES PURCHASE"`
	ReferenceID    string              `json:"reference_id"`
	EntityID       int64               `json:"entity_id" binding:"required"`
	EntityType     string              `json:"entity_type" binding:"required,oneof=CUSTOMER SUPPLIER"`
	IssueDate      time.Time           `json:"issue_date" binding:"required"`
	DueDate        time.Time           `json:"due_date" binding:"required"`
	Items          FinanceInvoiceItems `json:"items" binding:"required,dive"`
	DiscountAmount float64             `json:"discount_amount"`
	Notes          string              `json:"notes"`
}

// UpdateFinanceInvoiceRequest represents the request to update a finance invoice
type UpdateFinanceInvoiceRequest struct {
	ReferenceID    string               `json:"reference_id"`
	DueDate        time.Time            `json:"due_date"`
	Items          FinanceInvoiceItems  `json:"items"`
	DiscountAmount float64              `json:"discount_amount"`
	Notes          string               `json:"notes"`
	Status         FinanceInvoiceStatus `json:"status"`
}

// FinanceInvoiceResponse represents the response for finance invoice operations
type FinanceInvoiceResponse struct {
	Invoice *FinanceInvoice `json:"invoice"`
	Error   string          `json:"error,omitempty"`
}

// FinanceInvoiceListResponse represents the response for listing finance invoices
type FinanceInvoiceListResponse struct {
	Invoices []FinanceInvoice `json:"invoices"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
	Error    string           `json:"error,omitempty"`
}

// FinanceReport represents a financial report
type FinanceReport struct {
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	TotalRevenue float64   `json:"total_revenue"`
	TotalCost    float64   `json:"total_cost"`
	GrossProfit  float64   `json:"gross_profit"`
	TotalTax     float64   `json:"total_tax"`
	NetProfit    float64   `json:"net_profit"`
}

// FinanceReportRequest represents a request for a financial report
type FinanceReportRequest struct {
	StartDate  time.Time `json:"start_date" binding:"required"`
	EndDate    time.Time `json:"end_date" binding:"required"`
	ReportType string    `json:"report_type" binding:"required,oneof=REVENUE EXPENSE PROFIT_LOSS TAX"`
}

// FinanceReportResponse represents the response for a financial report
type FinanceReportResponse struct {
	Report *FinanceReport `json:"report"`
	Error  string         `json:"error,omitempty"`
}
