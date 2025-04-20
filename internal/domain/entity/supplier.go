package entity

import (
	"errors"
	"time"
)

var (
	// ErrInvalidRating is returned when a rating score is not between 0 and 5
	ErrInvalidRating = errors.New("rating score must be between 0 and 5")
)

type SupplierType string

const (
	Manufacturer SupplierType = "manufacturer"
	Wholesaler   SupplierType = "wholesaler"
	Distributor  SupplierType = "distributor"
	Retailer     SupplierType = "retailer"
)

type Supplier struct {
	ID           uint         `json:"id" gorm:"primaryKey"`
	Code         string       `json:"code" gorm:"unique;not null"`
	Name         string       `json:"name" gorm:"not null"`
	Type         SupplierType `json:"type"`
	ContactInfo  ContactInfo  `json:"contact_info" gorm:"embedded"`
	PaymentTerms PaymentTerms `json:"payment_terms" gorm:"embedded"`
	Rating       float64      `json:"rating" gorm:"default:0"`
	Products     []Product    `json:"products" gorm:"many2many:supplier_products;"`
	Contracts    []Contract   `json:"contracts" gorm:"foreignKey:SupplierID"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

type ContactInfo struct {
	Address string `json:"address"`
	Country string `json:"country"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Website string `json:"website"`
	TaxID   string `json:"tax_id"`
}

type PaymentTerms struct {
	PaymentMethod string `json:"payment_method"`
	PaymentDays   int    `json:"payment_days"`
	Currency      string `json:"currency"`
}

type Product struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	Code        string  `json:"code" gorm:"unique;not null"`
	Name        string  `json:"name" gorm:"not null"`
	Description string  `json:"description"`
	UnitPrice   float64 `json:"unit_price"`
	Currency    string  `json:"currency"`
}

type Contract struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	SupplierID uint      `json:"supplier_id"`
	ContractNo string    `json:"contract_no" gorm:"unique;not null"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	Terms      string    `json:"terms" gorm:"type:text"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type SupplierRating struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	SupplierID uint      `json:"supplier_id"`
	Score      float64   `json:"score"`
	Category   string    `json:"category"`
	Comment    string    `json:"comment"`
	RatedBy    uint      `json:"rated_by"`
	CreatedAt  time.Time `json:"created_at"`
}
