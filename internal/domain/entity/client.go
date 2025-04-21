package entity

import (
	"encoding/json"
	"time"
)

// Client represents a client in the system
type Client struct {
	ID            uint              `json:"id" gorm:"primaryKey"`
	Code          string            `json:"code" gorm:"unique;not null"`
	Name          string            `json:"name" gorm:"not null"`
	Type          string            `json:"type" gorm:"not null;default:'INDIVIDUAL'"`
	Email         string            `json:"email" gorm:"unique"`
	PhoneNumber   string            `json:"phone_number"`
	TaxID         string            `json:"tax_id"`
	Contacts      json.RawMessage   `json:"contacts" gorm:"type:jsonb"`
	CreditLimit   float64           `json:"credit_limit" gorm:"type:decimal(15,2);default:0"`
	CurrentDebt   float64           `json:"current_debt" gorm:"type:decimal(15,2);default:0"`
	LoyaltyTier   ClientLoyaltyTier `json:"loyalty_tier" gorm:"not null;default:'STANDARD'"`
	LoyaltyPoints int               `json:"loyalty_points" gorm:"default:0"`
	Notes         string            `json:"notes" gorm:"type:text"`
	CreatedAt     time.Time         `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time         `json:"updated_at" gorm:"autoUpdateTime"`
	Addresses     []ClientAddress   `json:"addresses,omitempty" gorm:"foreignKey:ClientID"`
	Orders        []SalesOrder      `json:"orders,omitempty" gorm:"foreignKey:ClientID"`
}

// ClientLoyaltyTier represents the loyalty tier of a client
type ClientLoyaltyTier string

const (
	ClientLoyaltyTierStandard ClientLoyaltyTier = "STANDARD"
	ClientLoyaltyTierSilver   ClientLoyaltyTier = "SILVER"
	ClientLoyaltyTierGold     ClientLoyaltyTier = "GOLD"
	ClientLoyaltyTierPlatinum ClientLoyaltyTier = "PLATINUM"
)

// Client types
const (
	ClientTypeIndividual  = "INDIVIDUAL"
	ClientTypeCorporate   = "CORPORATE"
	ClientTypeGovernment  = "GOVERNMENT"
	ClientTypeDistributor = "DISTRIBUTOR"
	ClientTypeReseller    = "RESELLER"
)

// ClientAddress represents an address associated with a client
type ClientAddress struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	ClientID   uint      `json:"client_id" gorm:"not null"`
	Type       string    `json:"type" gorm:"not null;default:'BOTH'"` // SHIPPING, BILLING, BOTH
	Street     string    `json:"street" gorm:"not null"`
	City       string    `json:"city" gorm:"not null"`
	State      string    `json:"state"`
	PostalCode string    `json:"postal_code"`
	Country    string    `json:"country" gorm:"not null"`
	IsDefault  bool      `json:"is_default" gorm:"default:false"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Client     *Client   `json:"-" gorm:"foreignKey:ClientID"`
}

// ClientFilter represents filters for searching clients
type ClientFilter struct {
	Code        string             `json:"code,omitempty"`
	Name        string             `json:"name,omitempty"`
	Type        *string            `json:"type,omitempty"`
	Email       string             `json:"email,omitempty"`
	PhoneNumber string             `json:"phone_number,omitempty"`
	LoyaltyTier *ClientLoyaltyTier `json:"loyalty_tier,omitempty"`
	City        string             `json:"city,omitempty"`
	Country     string             `json:"country,omitempty"`
}

// ClientOrderHistory represents a client's order history summary
type ClientOrderHistory struct {
	TotalOrders       int       `json:"total_orders"`
	TotalSpent        float64   `json:"total_spent"`
	AverageOrderValue float64   `json:"average_order_value"`
	FirstOrderDate    time.Time `json:"first_order_date"`
	LastOrderDate     time.Time `json:"last_order_date"`
	FrequentItems     []string  `json:"frequent_items"`
}

// ClientDebt represents a client's debt information
type ClientDebt struct {
	TotalDebt         float64   `json:"total_debt"`
	OverdueDebt       float64   `json:"overdue_debt"`
	UpcomingPayments  float64   `json:"upcoming_payments"`
	LastPaymentAmount float64   `json:"last_payment_amount"`
	LastPaymentDate   time.Time `json:"last_payment_date"`
}

// ClientRepository defines the interface for client data access
type ClientRepository interface {
	Create(client *Client) error
	FindByID(id uint) (*Client, error)
	FindByCode(code string) (*Client, error)
	FindByEmail(email string) (*Client, error)
	Update(client *Client) error
	Delete(id uint) error
	List(filter ClientFilter) ([]Client, error)

	CreateAddress(address *ClientAddress) error
	UpdateAddress(address *ClientAddress) error
	DeleteAddress(id uint) error
	FindAddressesByClientID(clientID uint) ([]ClientAddress, error)

	GetOrderHistory(clientID uint) (*ClientOrderHistory, error)
	GetClientDebt(clientID uint) (*ClientDebt, error)
	UpdateClientDebt(clientID uint, amount float64) error
	UpdateLoyaltyPoints(clientID uint, points int) error
	UpdateLoyaltyTier(clientID uint, tier ClientLoyaltyTier) error
}
