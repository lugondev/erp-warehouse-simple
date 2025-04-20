package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// CustomerType represents the type of customer
type CustomerType string

const (
	CustomerTypeIndividual CustomerType = "INDIVIDUAL"
	CustomerTypeCorporate  CustomerType = "CORPORATE"
	CustomerTypeReseller   CustomerType = "RESELLER"
	CustomerTypeWholesaler CustomerType = "WHOLESALER"
)

// CustomerLoyaltyTier represents the loyalty tier of a customer
type CustomerLoyaltyTier string

const (
	CustomerLoyaltyTierStandard CustomerLoyaltyTier = "STANDARD"
	CustomerLoyaltyTierSilver   CustomerLoyaltyTier = "SILVER"
	CustomerLoyaltyTierGold     CustomerLoyaltyTier = "GOLD"
	CustomerLoyaltyTierPlatinum CustomerLoyaltyTier = "PLATINUM"
)

// AddressType represents the type of address
type AddressType string

const (
	AddressTypeBilling  AddressType = "BILLING"
	AddressTypeShipping AddressType = "SHIPPING"
	AddressTypeBoth     AddressType = "BOTH"
)

// CustomerAddress represents a customer's address
type CustomerAddress struct {
	ID         uint        `json:"id" gorm:"primaryKey"`
	CustomerID uint        `json:"customer_id" gorm:"not null"`
	Type       AddressType `json:"type" gorm:"not null;default:'BOTH'"`
	Street     string      `json:"street" gorm:"not null"`
	City       string      `json:"city" gorm:"not null"`
	State      string      `json:"state"`
	PostalCode string      `json:"postal_code"`
	Country    string      `json:"country" gorm:"not null"`
	IsDefault  bool        `json:"is_default" gorm:"default:false"`
	CreatedAt  time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
	Customer   *Customer   `json:"-" gorm:"foreignKey:CustomerID"`
}

// CustomerContact represents contact information for a customer
type CustomerContact struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Position    string `json:"position"`
}

// Scan implements the sql.Scanner interface for CustomerContacts
func (cc *CustomerContacts) Scan(value interface{}) error {
	if value == nil {
		*cc = make(CustomerContacts, 0)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan CustomerContacts: value is not []byte")
	}

	if err := json.Unmarshal(bytes, cc); err != nil {
		return err
	}
	return nil
}

// Value implements the driver.Valuer interface for CustomerContacts
func (cc CustomerContacts) Value() (driver.Value, error) {
	if cc == nil {
		return nil, nil
	}
	return json.Marshal(cc)
}

// CustomerContacts is a slice of CustomerContact
type CustomerContacts []CustomerContact

// Customer represents a customer in the system
type Customer struct {
	ID            uint                `json:"id" gorm:"primaryKey"`
	Code          string              `json:"code" gorm:"uniqueIndex;not null"`
	Name          string              `json:"name" gorm:"not null"`
	Type          CustomerType        `json:"type" gorm:"not null;default:'INDIVIDUAL'"`
	Email         string              `json:"email" gorm:"uniqueIndex"`
	PhoneNumber   string              `json:"phone_number"`
	TaxID         string              `json:"tax_id"`
	Contacts      CustomerContacts    `json:"contacts" gorm:"type:jsonb"`
	CreditLimit   float64             `json:"credit_limit" gorm:"type:decimal(15,2);default:0"`
	CurrentDebt   float64             `json:"current_debt" gorm:"type:decimal(15,2);default:0"`
	LoyaltyTier   CustomerLoyaltyTier `json:"loyalty_tier" gorm:"not null;default:'STANDARD'"`
	LoyaltyPoints int                 `json:"loyalty_points" gorm:"default:0"`
	Notes         string              `json:"notes" gorm:"type:text"`
	CreatedAt     time.Time           `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time           `json:"updated_at" gorm:"autoUpdateTime"`
	Addresses     []CustomerAddress   `json:"addresses,omitempty" gorm:"foreignKey:CustomerID"`
	SalesOrders   []SalesOrder        `json:"sales_orders,omitempty" gorm:"foreignKey:CustomerID"`
}

// CustomerOrderHistory represents a summary of a customer's order history
type CustomerOrderHistory struct {
	TotalOrders       int       `json:"total_orders"`
	TotalSpent        float64   `json:"total_spent"`
	FirstOrderDate    time.Time `json:"first_order_date"`
	LastOrderDate     time.Time `json:"last_order_date"`
	AverageOrderValue float64   `json:"average_order_value"`
	FrequentItems     []string  `json:"frequent_items"`
}

// CustomerDebt represents a customer's debt information
type CustomerDebt struct {
	TotalDebt         float64   `json:"total_debt"`
	OverdueDebt       float64   `json:"overdue_debt"`
	UpcomingPayments  float64   `json:"upcoming_payments"`
	LastPaymentDate   time.Time `json:"last_payment_date"`
	LastPaymentAmount float64   `json:"last_payment_amount"`
}

// CustomerFilter represents filters for searching customers
type CustomerFilter struct {
	Code        string               `json:"code,omitempty"`
	Name        string               `json:"name,omitempty"`
	Type        *CustomerType        `json:"type,omitempty"`
	Email       string               `json:"email,omitempty"`
	PhoneNumber string               `json:"phone_number,omitempty"`
	LoyaltyTier *CustomerLoyaltyTier `json:"loyalty_tier,omitempty"`
	City        string               `json:"city,omitempty"`
	Country     string               `json:"country,omitempty"`
}

// CustomerRepository defines the interface for customer data access
type CustomerRepository interface {
	Create(customer *Customer) error
	FindByID(id uint) (*Customer, error)
	FindByCode(code string) (*Customer, error)
	FindByEmail(email string) (*Customer, error)
	Update(customer *Customer) error
	Delete(id uint) error
	List(filter CustomerFilter) ([]Customer, error)

	// Address methods
	CreateAddress(address *CustomerAddress) error
	UpdateAddress(address *CustomerAddress) error
	DeleteAddress(id uint) error
	FindAddressesByCustomerID(customerID uint) ([]CustomerAddress, error)

	// Order history methods
	GetOrderHistory(customerID uint) (*CustomerOrderHistory, error)

	// Debt management methods
	GetCustomerDebt(customerID uint) (*CustomerDebt, error)
	UpdateCustomerDebt(customerID uint, amount float64) error

	// Loyalty methods
	UpdateLoyaltyPoints(customerID uint, points int) error
	UpdateLoyaltyTier(customerID uint, tier CustomerLoyaltyTier) error
}
