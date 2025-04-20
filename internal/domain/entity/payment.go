package entity

import (
	"time"
)

// FinancePaymentMethod represents the method of payment
type FinancePaymentMethod string

const (
	FinancePaymentMethodCash          FinancePaymentMethod = "CASH"
	FinancePaymentMethodBankTransfer  FinancePaymentMethod = "BANK_TRANSFER"
	FinancePaymentMethodCreditCard    FinancePaymentMethod = "CREDIT_CARD"
	FinancePaymentMethodCheck         FinancePaymentMethod = "CHECK"
	FinancePaymentMethodDigitalWallet FinancePaymentMethod = "DIGITAL_WALLET"
	FinancePaymentMethodOther         FinancePaymentMethod = "OTHER"
)

// FinancePaymentStatus represents the status of a payment
type FinancePaymentStatus string

const (
	FinancePaymentPending   FinancePaymentStatus = "PENDING"
	FinancePaymentCompleted FinancePaymentStatus = "COMPLETED"
	FinancePaymentFailed    FinancePaymentStatus = "FAILED"
	FinancePaymentCancelled FinancePaymentStatus = "CANCELLED"
	FinancePaymentRefunded  FinancePaymentStatus = "REFUNDED"
)

// FinancePayment represents a payment entity
type FinancePayment struct {
	ID              int64                `json:"id" db:"id"`
	PaymentNumber   string               `json:"payment_number" db:"payment_number"`
	InvoiceID       int64                `json:"invoice_id" db:"invoice_id"`
	InvoiceNumber   string               `json:"invoice_number" db:"invoice_number"`
	EntityID        int64                `json:"entity_id" db:"entity_id"`
	EntityType      string               `json:"entity_type" db:"entity_type"` // "CUSTOMER" or "SUPPLIER"
	EntityName      string               `json:"entity_name" db:"entity_name"`
	PaymentDate     time.Time            `json:"payment_date" db:"payment_date"`
	PaymentMethod   FinancePaymentMethod `json:"payment_method" db:"payment_method"`
	Amount          float64              `json:"amount" db:"amount"`
	Status          FinancePaymentStatus `json:"status" db:"status"`
	Notes           string               `json:"notes" db:"notes"`
	ReferenceNumber string               `json:"reference_number" db:"reference_number"`
	CreatedBy       int64                `json:"created_by" db:"created_by"`
	CreatedAt       time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at" db:"updated_at"`
}

// FinancePaymentFilter represents filters for querying payments
type FinancePaymentFilter struct {
	PaymentNumber string               `json:"payment_number,omitempty"`
	InvoiceID     int64                `json:"invoice_id,omitempty"`
	InvoiceNumber string               `json:"invoice_number,omitempty"`
	EntityID      int64                `json:"entity_id,omitempty"`
	EntityType    string               `json:"entity_type,omitempty"`
	Status        FinancePaymentStatus `json:"status,omitempty"`
	PaymentMethod FinancePaymentMethod `json:"payment_method,omitempty"`
	StartDate     *time.Time           `json:"start_date,omitempty"`
	EndDate       *time.Time           `json:"end_date,omitempty"`
	Page          int                  `json:"page,omitempty"`
	PageSize      int                  `json:"page_size,omitempty"`
}

// CreateFinancePaymentRequest represents the request to create a new payment
type CreateFinancePaymentRequest struct {
	InvoiceID       int64                `json:"invoice_id" binding:"required"`
	PaymentDate     time.Time            `json:"payment_date" binding:"required"`
	PaymentMethod   FinancePaymentMethod `json:"payment_method" binding:"required"`
	Amount          float64              `json:"amount" binding:"required,gt=0"`
	ReferenceNumber string               `json:"reference_number"`
	Notes           string               `json:"notes"`
}

// UpdateFinancePaymentRequest represents the request to update a payment
type UpdateFinancePaymentRequest struct {
	PaymentMethod   FinancePaymentMethod `json:"payment_method"`
	Amount          float64              `json:"amount" binding:"gt=0"`
	Status          FinancePaymentStatus `json:"status"`
	ReferenceNumber string               `json:"reference_number"`
	Notes           string               `json:"notes"`
}

// FinancePaymentResponse represents the response for payment operations
type FinancePaymentResponse struct {
	Payment *FinancePayment `json:"payment"`
	Error   string          `json:"error,omitempty"`
}

// FinancePaymentListResponse represents the response for listing payments
type FinancePaymentListResponse struct {
	Payments []FinancePayment `json:"payments"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
	Error    string           `json:"error,omitempty"`
}

// FinanceAccountsReceivable represents an accounts receivable record
type FinanceAccountsReceivable struct {
	EntityID        int64     `json:"entity_id" db:"entity_id"`
	EntityName      string    `json:"entity_name" db:"entity_name"`
	InvoiceID       int64     `json:"invoice_id" db:"invoice_id"`
	InvoiceNumber   string    `json:"invoice_number" db:"invoice_number"`
	InvoiceDate     time.Time `json:"invoice_date" db:"invoice_date"`
	DueDate         time.Time `json:"due_date" db:"due_date"`
	TotalAmount     float64   `json:"total_amount" db:"total_amount"`
	AmountPaid      float64   `json:"amount_paid" db:"amount_paid"`
	AmountDue       float64   `json:"amount_due" db:"amount_due"`
	DaysOverdue     int       `json:"days_overdue" db:"days_overdue"`
	Status          string    `json:"status" db:"status"`
	LastPaymentDate time.Time `json:"last_payment_date" db:"last_payment_date"`
}

// FinanceAccountsPayable represents an accounts payable record
type FinanceAccountsPayable struct {
	EntityID        int64     `json:"entity_id" db:"entity_id"`
	EntityName      string    `json:"entity_name" db:"entity_name"`
	InvoiceID       int64     `json:"invoice_id" db:"invoice_id"`
	InvoiceNumber   string    `json:"invoice_number" db:"invoice_number"`
	InvoiceDate     time.Time `json:"invoice_date" db:"invoice_date"`
	DueDate         time.Time `json:"due_date" db:"due_date"`
	TotalAmount     float64   `json:"total_amount" db:"total_amount"`
	AmountPaid      float64   `json:"amount_paid" db:"amount_paid"`
	AmountDue       float64   `json:"amount_due" db:"amount_due"`
	DaysOverdue     int       `json:"days_overdue" db:"days_overdue"`
	Status          string    `json:"status" db:"status"`
	LastPaymentDate time.Time `json:"last_payment_date" db:"last_payment_date"`
}
