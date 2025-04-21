package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// SalesOrderStatus represents the status of a sales order
type SalesOrderStatus string

const (
	SalesOrderStatusDraft      SalesOrderStatus = "DRAFT"
	SalesOrderStatusConfirmed  SalesOrderStatus = "CONFIRMED"
	SalesOrderStatusProcessing SalesOrderStatus = "PROCESSING"
	SalesOrderStatusShipped    SalesOrderStatus = "SHIPPED"
	SalesOrderStatusDelivered  SalesOrderStatus = "DELIVERED"
	SalesOrderStatusCompleted  SalesOrderStatus = "COMPLETED"
	SalesOrderStatusCancelled  SalesOrderStatus = "CANCELLED"
)

// DeliveryOrderStatus represents the status of a delivery order
type DeliveryOrderStatus string

const (
	DeliveryOrderStatusPending   DeliveryOrderStatus = "PENDING"
	DeliveryOrderStatusPreparing DeliveryOrderStatus = "PREPARING"
	DeliveryOrderStatusInTransit DeliveryOrderStatus = "IN_TRANSIT"
	DeliveryOrderStatusDelivered DeliveryOrderStatus = "DELIVERED"
	DeliveryOrderStatusCancelled DeliveryOrderStatus = "CANCELLED"
	DeliveryOrderStatusReturned  DeliveryOrderStatus = "RETURNED"
)

// PaymentMethod represents the payment method for a sales order
type PaymentMethod string

const (
	PaymentMethodCash          PaymentMethod = "CASH"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodBankTransfer  PaymentMethod = "BANK_TRANSFER"
	PaymentMethodDigitalWallet PaymentMethod = "DIGITAL_WALLET"
)

// InvoiceStatus represents the status of an invoice
type InvoiceStatus string

const (
	InvoiceStatusDraft     InvoiceStatus = "DRAFT"
	InvoiceStatusIssued    InvoiceStatus = "ISSUED"
	InvoiceStatusPaid      InvoiceStatus = "PAID"
	InvoiceStatusPartial   InvoiceStatus = "PARTIAL"
	InvoiceStatusOverdue   InvoiceStatus = "OVERDUE"
	InvoiceStatusCancelled InvoiceStatus = "CANCELLED"
)

// SalesOrderItem represents an item in a sales order
type SalesOrderItem struct {
	SKUID       string  `json:"sku_id" gorm:"not null"`
	Quantity    float64 `json:"quantity" gorm:"not null"`
	UnitPrice   float64 `json:"unit_price" gorm:"type:decimal(15,2);not null"`
	Discount    float64 `json:"discount" gorm:"type:decimal(15,2);default:0"`
	TaxRate     float64 `json:"tax_rate" gorm:"type:decimal(5,2);default:0"`
	TaxAmount   float64 `json:"tax_amount" gorm:"type:decimal(15,2);default:0"`
	TotalPrice  float64 `json:"total_price" gorm:"type:decimal(15,2);not null"`
	Description string  `json:"description"`
	SKU         *SKU    `json:"sku,omitempty" gorm:"foreignKey:SKUID"`
}

// Scan implements the sql.Scanner interface for SalesOrderItems
func (soi *SalesOrderItems) Scan(value interface{}) error {
	if value == nil {
		*soi = make(SalesOrderItems, 0)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan SalesOrderItems: value is not []byte")
	}

	if err := json.Unmarshal(bytes, soi); err != nil {
		return err
	}
	return nil
}

// Value implements the driver.Valuer interface for SalesOrderItems
func (soi SalesOrderItems) Value() (driver.Value, error) {
	if soi == nil {
		return nil, nil
	}
	return json.Marshal(soi)
}

// SalesOrderItems is a slice of SalesOrderItem
type SalesOrderItems []SalesOrderItem

// SalesOrder represents a customer order
type SalesOrder struct {
	ID              string           `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OrderNumber     string           `json:"order_number" gorm:"uniqueIndex;not null"`
	ClientID        uint             `json:"client_id" gorm:"not null"`
	OrderDate       time.Time        `json:"order_date" gorm:"not null"`
	Items           SalesOrderItems  `json:"items" gorm:"type:jsonb;not null"`
	SubTotal        float64          `json:"sub_total" gorm:"type:decimal(15,2);not null"`
	TaxTotal        float64          `json:"tax_total" gorm:"type:decimal(15,2);default:0"`
	DiscountTotal   float64          `json:"discount_total" gorm:"type:decimal(15,2);default:0"`
	GrandTotal      float64          `json:"grand_total" gorm:"type:decimal(15,2);not null"`
	Status          SalesOrderStatus `json:"status" gorm:"not null;default:'DRAFT'"`
	PaymentMethod   PaymentMethod    `json:"payment_method"`
	PaymentStatus   PaymentStatus    `json:"payment_status" gorm:"not null;default:'PENDING'"`
	ShippingAddress string           `json:"shipping_address" gorm:"type:text"`
	BillingAddress  string           `json:"billing_address" gorm:"type:text"`
	Notes           string           `json:"notes" gorm:"type:text"`
	CreatedByID     uint             `json:"created_by_id" gorm:"not null"`
	CreatedAt       time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
	Client          *User            `json:"client,omitempty" gorm:"foreignKey:ClientID"` // Using User as Client for now
	CreatedBy       *User            `json:"created_by,omitempty" gorm:"foreignKey:CreatedByID"`
	DeliveryOrders  []DeliveryOrder  `json:"delivery_orders,omitempty" gorm:"foreignKey:SalesOrderID"`
	Invoices        []Invoice        `json:"invoices,omitempty" gorm:"foreignKey:SalesOrderID"`
}

// DeliveryOrderItem represents an item in a delivery order
type DeliveryOrderItem struct {
	SKUID             string  `json:"sku_id" gorm:"not null"`
	OrderedQuantity   float64 `json:"ordered_quantity" gorm:"not null"`
	ShippedQuantity   float64 `json:"shipped_quantity" gorm:"not null"`
	RemainingQuantity float64 `json:"remaining_quantity" gorm:"default:0"`
	Notes             string  `json:"notes"`
	SKU               *SKU    `json:"sku,omitempty" gorm:"foreignKey:SKUID"`
}

// Scan implements the sql.Scanner interface for DeliveryOrderItems
func (doi *DeliveryOrderItems) Scan(value interface{}) error {
	if value == nil {
		*doi = make(DeliveryOrderItems, 0)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan DeliveryOrderItems: value is not []byte")
	}

	if err := json.Unmarshal(bytes, doi); err != nil {
		return err
	}
	return nil
}

// Value implements the driver.Valuer interface for DeliveryOrderItems
func (doi DeliveryOrderItems) Value() (driver.Value, error) {
	if doi == nil {
		return nil, nil
	}
	return json.Marshal(doi)
}

// DeliveryOrderItems is a slice of DeliveryOrderItem
type DeliveryOrderItems []DeliveryOrderItem

// DeliveryOrder represents a delivery of goods from a sales order
type DeliveryOrder struct {
	ID              string              `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	DeliveryNumber  string              `json:"delivery_number" gorm:"uniqueIndex;not null"`
	SalesOrderID    string              `json:"sales_order_id" gorm:"type:uuid;not null"`
	DeliveryDate    time.Time           `json:"delivery_date" gorm:"not null"`
	Items           DeliveryOrderItems  `json:"items" gorm:"type:jsonb;not null"`
	ShippingAddress string              `json:"shipping_address" gorm:"type:text;not null"`
	Status          DeliveryOrderStatus `json:"status" gorm:"not null;default:'PENDING'"`
	TrackingNumber  string              `json:"tracking_number"`
	ShippingMethod  string              `json:"shipping_method"`
	StoreID         string              `json:"store_id" gorm:"not null"`
	Notes           string              `json:"notes" gorm:"type:text"`
	CreatedByID     uint                `json:"created_by_id" gorm:"not null"`
	CreatedAt       time.Time           `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time           `json:"updated_at" gorm:"autoUpdateTime"`
	SalesOrder      *SalesOrder         `json:"sales_order,omitempty" gorm:"foreignKey:SalesOrderID"`
	CreatedBy       *User               `json:"created_by,omitempty" gorm:"foreignKey:CreatedByID"`
}

// Invoice represents an invoice for a sales order
type Invoice struct {
	ID            string        `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	InvoiceNumber string        `json:"invoice_number" gorm:"uniqueIndex;not null"`
	SalesOrderID  string        `json:"sales_order_id" gorm:"type:uuid;not null"`
	IssueDate     time.Time     `json:"issue_date" gorm:"not null"`
	DueDate       time.Time     `json:"due_date" gorm:"not null"`
	Amount        float64       `json:"amount" gorm:"type:decimal(15,2);not null"`
	TaxAmount     float64       `json:"tax_amount" gorm:"type:decimal(15,2);default:0"`
	TotalAmount   float64       `json:"total_amount" gorm:"type:decimal(15,2);not null"`
	Status        InvoiceStatus `json:"status" gorm:"not null;default:'DRAFT'"`
	Notes         string        `json:"notes" gorm:"type:text"`
	CreatedByID   uint          `json:"created_by_id" gorm:"not null"`
	CreatedAt     time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
	SalesOrder    *SalesOrder   `json:"sales_order,omitempty" gorm:"foreignKey:SalesOrderID"`
	CreatedBy     *User         `json:"created_by,omitempty" gorm:"foreignKey:CreatedByID"`
}

// SalesOrderFilter represents filters for searching sales orders
type SalesOrderFilter struct {
	OrderNumber   string            `json:"order_number,omitempty"`
	ClientID      *uint             `json:"client_id,omitempty"`
	Status        *SalesOrderStatus `json:"status,omitempty"`
	PaymentStatus *PaymentStatus    `json:"payment_status,omitempty"`
	StartDate     *time.Time        `json:"start_date,omitempty"`
	EndDate       *time.Time        `json:"end_date,omitempty"`
	SKUID         string            `json:"sku_id,omitempty"`
}

// DeliveryOrderFilter represents filters for searching delivery orders
type DeliveryOrderFilter struct {
	DeliveryNumber string               `json:"delivery_number,omitempty"`
	SalesOrderID   string               `json:"sales_order_id,omitempty"`
	Status         *DeliveryOrderStatus `json:"status,omitempty"`
	StartDate      *time.Time           `json:"start_date,omitempty"`
	EndDate        *time.Time           `json:"end_date,omitempty"`
	StoreID        string               `json:"store_id,omitempty"`
}

// InvoiceFilter represents filters for searching invoices
type InvoiceFilter struct {
	InvoiceNumber string         `json:"invoice_number,omitempty"`
	SalesOrderID  string         `json:"sales_order_id,omitempty"`
	Status        *InvoiceStatus `json:"status,omitempty"`
	StartDate     *time.Time     `json:"start_date,omitempty"`
	EndDate       *time.Time     `json:"end_date,omitempty"`
}
