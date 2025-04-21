package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Purchase Request Status
type PurchaseRequestStatus string

const (
	PurchaseRequestStatusDraft     PurchaseRequestStatus = "DRAFT"
	PurchaseRequestStatusSubmitted PurchaseRequestStatus = "SUBMITTED"
	PurchaseRequestStatusApproved  PurchaseRequestStatus = "APPROVED"
	PurchaseRequestStatusRejected  PurchaseRequestStatus = "REJECTED"
	PurchaseRequestStatusCancelled PurchaseRequestStatus = "CANCELLED"
	PurchaseRequestStatusOrdered   PurchaseRequestStatus = "ORDERED"
)

// Purchase Order Status
type PurchaseOrderStatus string

const (
	PurchaseOrderStatusDraft     PurchaseOrderStatus = "DRAFT"
	PurchaseOrderStatusSubmitted PurchaseOrderStatus = "SUBMITTED"
	PurchaseOrderStatusApproved  PurchaseOrderStatus = "APPROVED"
	PurchaseOrderStatusSent      PurchaseOrderStatus = "SENT"
	PurchaseOrderStatusConfirmed PurchaseOrderStatus = "CONFIRMED"
	PurchaseOrderStatusPartial   PurchaseOrderStatus = "PARTIALLY_RECEIVED"
	PurchaseOrderStatusReceived  PurchaseOrderStatus = "RECEIVED"
	PurchaseOrderStatusCancelled PurchaseOrderStatus = "CANCELLED"
	PurchaseOrderStatusClosed    PurchaseOrderStatus = "CLOSED"
)

// Payment Status
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusPartial   PaymentStatus = "PARTIAL"
	PaymentStatusPaid      PaymentStatus = "PAID"
	PaymentStatusOverdue   PaymentStatus = "OVERDUE"
	PaymentStatusCancelled PaymentStatus = "CANCELLED"
)

// Purchase Request Item represents an item in a purchase request
type PurchaseRequestItem struct {
	SKUID       string  `json:"sku_id" gorm:"not null"`
	Quantity    float64 `json:"quantity" gorm:"not null"`
	Description string  `json:"description"`
	SKU         *SKU    `json:"sku,omitempty" gorm:"foreignKey:SKUID"`
}

// Scan implements the sql.Scanner interface for PurchaseRequestItems
func (pri *PurchaseRequestItems) Scan(value interface{}) error {
	if value == nil {
		*pri = make(PurchaseRequestItems, 0)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan PurchaseRequestItems: value is not []byte")
	}

	if err := json.Unmarshal(bytes, pri); err != nil {
		return err
	}
	return nil
}

// Value implements the driver.Valuer interface for PurchaseRequestItems
func (pri PurchaseRequestItems) Value() (driver.Value, error) {
	if pri == nil {
		return nil, nil
	}
	return json.Marshal(pri)
}

// PurchaseRequestItems is a slice of PurchaseRequestItem
type PurchaseRequestItems []PurchaseRequestItem

// PurchaseRequest represents a request to purchase items
type PurchaseRequest struct {
	ID              string                `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	RequestNumber   string                `json:"request_number" gorm:"uniqueIndex;not null"`
	RequesterID     uint                  `json:"requester_id" gorm:"not null"`
	RequestDate     time.Time             `json:"request_date" gorm:"not null"`
	RequiredDate    time.Time             `json:"required_date"`
	Items           PurchaseRequestItems  `json:"items" gorm:"type:jsonb;not null"`
	Reason          string                `json:"reason" gorm:"type:text"`
	Status          PurchaseRequestStatus `json:"status" gorm:"not null;default:'DRAFT'"`
	ApproverID      *uint                 `json:"approver_id"`
	ApprovalDate    *time.Time            `json:"approval_date"`
	ApprovalNotes   string                `json:"approval_notes" gorm:"type:text"`
	DepartmentID    *uint                 `json:"department_id"`
	TotalEstimated  float64               `json:"total_estimated" gorm:"type:decimal(15,2)"`
	CurrencyCode    string                `json:"currency_code" gorm:"default:'USD'"`
	AttachmentURLs  []string              `json:"attachment_urls" gorm:"type:text[]"`
	CreatedAt       time.Time             `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time             `json:"updated_at" gorm:"autoUpdateTime"`
	Requester       *User                 `json:"requester,omitempty" gorm:"foreignKey:RequesterID"`
	Approver        *User                 `json:"approver,omitempty" gorm:"foreignKey:ApproverID"`
	PurchaseOrderID *string               `json:"purchase_order_id" gorm:"type:uuid"`
	PurchaseOrder   *PurchaseOrder        `json:"purchase_order,omitempty" gorm:"foreignKey:PurchaseOrderID"`
}

// PurchaseOrderItem represents an item in a purchase order
type PurchaseOrderItem struct {
	SKUID       string  `json:"sku_id" gorm:"not null"`
	Quantity    float64 `json:"quantity" gorm:"not null"`
	UnitPrice   float64 `json:"unit_price" gorm:"type:decimal(15,2);not null"`
	TaxRate     float64 `json:"tax_rate" gorm:"type:decimal(5,2);default:0"`
	TaxAmount   float64 `json:"tax_amount" gorm:"type:decimal(15,2);default:0"`
	Discount    float64 `json:"discount" gorm:"type:decimal(15,2);default:0"`
	TotalPrice  float64 `json:"total_price" gorm:"type:decimal(15,2);not null"`
	Description string  `json:"description"`
	SKU         *SKU    `json:"sku,omitempty" gorm:"foreignKey:SKUID"`
}

// Scan implements the sql.Scanner interface for PurchaseOrderItems
func (poi *PurchaseOrderItems) Scan(value interface{}) error {
	if value == nil {
		*poi = make(PurchaseOrderItems, 0)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan PurchaseOrderItems: value is not []byte")
	}

	if err := json.Unmarshal(bytes, poi); err != nil {
		return err
	}
	return nil
}

// Value implements the driver.Valuer interface for PurchaseOrderItems
func (poi PurchaseOrderItems) Value() (driver.Value, error) {
	if poi == nil {
		return nil, nil
	}
	return json.Marshal(poi)
}

// PurchaseOrderItems is a slice of PurchaseOrderItem
type PurchaseOrderItems []PurchaseOrderItem

// PurchaseOrder represents an order to a supplier
type PurchaseOrder struct {
	ID               string              `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OrderNumber      string              `json:"order_number" gorm:"uniqueIndex;not null"`
	VendorID         uint                `json:"vendor_id" gorm:"not null"`
	OrderDate        time.Time           `json:"order_date" gorm:"not null"`
	ExpectedDate     time.Time           `json:"expected_date"`
	Items            PurchaseOrderItems  `json:"items" gorm:"type:jsonb;not null"`
	SubTotal         float64             `json:"sub_total" gorm:"type:decimal(15,2);not null"`
	TaxTotal         float64             `json:"tax_total" gorm:"type:decimal(15,2);default:0"`
	DiscountTotal    float64             `json:"discount_total" gorm:"type:decimal(15,2);default:0"`
	GrandTotal       float64             `json:"grand_total" gorm:"type:decimal(15,2);not null"`
	CurrencyCode     string              `json:"currency_code" gorm:"default:'USD'"`
	PaymentTerms     string              `json:"payment_terms"`
	Status           PurchaseOrderStatus `json:"status" gorm:"not null;default:'DRAFT'"`
	PaymentStatus    PaymentStatus       `json:"payment_status" gorm:"not null;default:'PENDING'"`
	ShippingAddress  string              `json:"shipping_address" gorm:"type:text"`
	ShippingMethod   string              `json:"shipping_method"`
	Notes            string              `json:"notes" gorm:"type:text"`
	AttachmentURLs   []string            `json:"attachment_urls" gorm:"type:text[]"`
	CreatedByID      uint                `json:"created_by_id" gorm:"not null"`
	ApprovedByID     *uint               `json:"approved_by_id"`
	ApprovalDate     *time.Time          `json:"approval_date"`
	CreatedAt        time.Time           `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time           `json:"updated_at" gorm:"autoUpdateTime"`
	Vendor           *Vendor             `json:"vendor,omitempty" gorm:"foreignKey:VendorID"`
	CreatedBy        *User               `json:"created_by,omitempty" gorm:"foreignKey:CreatedByID"`
	ApprovedBy       *User               `json:"approved_by,omitempty" gorm:"foreignKey:ApprovedByID"`
	PurchaseRequests []PurchaseRequest   `json:"purchase_requests,omitempty" gorm:"foreignKey:PurchaseOrderID"`
}

// PurchaseReceiptItem represents an item in a purchase receipt
type PurchaseReceiptItem struct {
	SKUID            string  `json:"sku_id" gorm:"not null"`
	OrderedQuantity  float64 `json:"ordered_quantity" gorm:"not null"`
	ReceivedQuantity float64 `json:"received_quantity" gorm:"not null"`
	RejectedQuantity float64 `json:"rejected_quantity" gorm:"default:0"`
	UnitPrice        float64 `json:"unit_price" gorm:"type:decimal(15,2);not null"`
	TotalPrice       float64 `json:"total_price" gorm:"type:decimal(15,2);not null"`
	Notes            string  `json:"notes"`
	SKU              *SKU    `json:"sku,omitempty" gorm:"foreignKey:SKUID"`
}

// Scan implements the sql.Scanner interface for PurchaseReceiptItems
func (pri *PurchaseReceiptItems) Scan(value interface{}) error {
	if value == nil {
		*pri = make(PurchaseReceiptItems, 0)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan PurchaseReceiptItems: value is not []byte")
	}

	if err := json.Unmarshal(bytes, pri); err != nil {
		return err
	}
	return nil
}

// Value implements the driver.Valuer interface for PurchaseReceiptItems
func (pri PurchaseReceiptItems) Value() (driver.Value, error) {
	if pri == nil {
		return nil, nil
	}
	return json.Marshal(pri)
}

// PurchaseReceiptItems is a slice of PurchaseReceiptItem
type PurchaseReceiptItems []PurchaseReceiptItem

// PurchaseReceipt represents a receipt of goods from a purchase order
type PurchaseReceipt struct {
	ID              string               `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ReceiptNumber   string               `json:"receipt_number" gorm:"uniqueIndex;not null"`
	PurchaseOrderID string               `json:"purchase_order_id" gorm:"type:uuid;not null"`
	ReceiptDate     time.Time            `json:"receipt_date" gorm:"not null"`
	Items           PurchaseReceiptItems `json:"items" gorm:"type:jsonb;not null"`
	StoreID         string               `json:"store_id" gorm:"not null"`
	ReceivedByID    uint                 `json:"received_by_id" gorm:"not null"`
	Notes           string               `json:"notes" gorm:"type:text"`
	AttachmentURLs  []string             `json:"attachment_urls" gorm:"type:text[]"`
	CreatedAt       time.Time            `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time            `json:"updated_at" gorm:"autoUpdateTime"`
	PurchaseOrder   *PurchaseOrder       `json:"purchase_order,omitempty" gorm:"foreignKey:PurchaseOrderID"`
	ReceivedBy      *User                `json:"received_by,omitempty" gorm:"foreignKey:ReceivedByID"`
}

// PurchasePayment represents a payment for a purchase order
type PurchasePayment struct {
	ID              string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	PaymentNumber   string         `json:"payment_number" gorm:"uniqueIndex;not null"`
	PurchaseOrderID string         `json:"purchase_order_id" gorm:"type:uuid;not null"`
	PaymentDate     time.Time      `json:"payment_date" gorm:"not null"`
	Amount          float64        `json:"amount" gorm:"type:decimal(15,2);not null"`
	PaymentMethod   string         `json:"payment_method" gorm:"not null"`
	ReferenceNumber string         `json:"reference_number"`
	Notes           string         `json:"notes" gorm:"type:text"`
	CreatedByID     uint           `json:"created_by_id" gorm:"not null"`
	CreatedAt       time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	PurchaseOrder   *PurchaseOrder `json:"purchase_order,omitempty" gorm:"foreignKey:PurchaseOrderID"`
	CreatedBy       *User          `json:"created_by,omitempty" gorm:"foreignKey:CreatedByID"`
}

// PurchaseRequestFilter represents filters for searching purchase requests
type PurchaseRequestFilter struct {
	RequestNumber string                 `json:"request_number,omitempty"`
	RequesterID   *uint                  `json:"requester_id,omitempty"`
	Status        *PurchaseRequestStatus `json:"status,omitempty"`
	StartDate     *time.Time             `json:"start_date,omitempty"`
	EndDate       *time.Time             `json:"end_date,omitempty"`
	SKUID         string                 `json:"sku_id,omitempty"`
}

// PurchaseOrderFilter represents filters for searching purchase orders
type PurchaseOrderFilter struct {
	OrderNumber   string               `json:"order_number,omitempty"`
	VendorID      *uint                `json:"vendor_id,omitempty"`
	Status        *PurchaseOrderStatus `json:"status,omitempty"`
	PaymentStatus *PaymentStatus       `json:"payment_status,omitempty"`
	StartDate     *time.Time           `json:"start_date,omitempty"`
	EndDate       *time.Time           `json:"end_date,omitempty"`
	SKUID         string               `json:"sku_id,omitempty"`
}
