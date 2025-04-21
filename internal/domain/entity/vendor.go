package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Error definitions
var (
	ErrInvalidRating = errors.New("rating must be between 0 and 5")
)

// Vendor represents a supplier of goods or services
type Vendor struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Code          string         `json:"code" gorm:"unique;not null"`
	Name          string         `json:"name" gorm:"not null"`
	Type          string         `json:"type"`
	Address       string         `json:"address"`
	Country       string         `json:"country"`
	Email         string         `json:"email"`
	Phone         string         `json:"phone"`
	Website       string         `json:"website"`
	TaxID         string         `json:"tax_id"`
	PaymentMethod string         `json:"payment_method"`
	PaymentDays   int            `json:"payment_days"`
	Currency      string         `json:"currency"`
	Rating        float64        `json:"rating" gorm:"type:decimal(3,2);default:0"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	Products      []Product      `json:"products,omitempty" gorm:"many2many:vendor_products"`
	Contracts     []Contract     `json:"contracts,omitempty" gorm:"foreignKey:VendorID"`
	VendorRatings []VendorRating `json:"vendor_ratings,omitempty" gorm:"foreignKey:VendorID"`
}

// Product represents a product supplied by a vendor
type Product struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Code        string    `json:"code" gorm:"unique;not null"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	UnitPrice   float64   `json:"unit_price" gorm:"type:decimal(15,2)"`
	Currency    string    `json:"currency"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Vendors     []Vendor  `json:"vendors,omitempty" gorm:"many2many:vendor_products"`
}

// Contract represents a contract with a vendor
type Contract struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	VendorID   uint      `json:"vendor_id" gorm:"not null"`
	ContractNo string    `json:"contract_no" gorm:"unique;not null"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	Terms      string    `json:"terms"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Vendor     *Vendor   `json:"vendor,omitempty" gorm:"foreignKey:VendorID"`
}

// VendorRating represents a rating for a vendor
type VendorRating struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	VendorID  uint      `json:"vendor_id" gorm:"not null"`
	Score     float64   `json:"score" gorm:"type:decimal(3,2);not null"`
	Category  string    `json:"category"`
	Comment   string    `json:"comment"`
	RatedByID uint      `json:"rated_by_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	Vendor    *Vendor   `json:"vendor,omitempty" gorm:"foreignKey:VendorID"`
	RatedBy   *User     `json:"rated_by,omitempty" gorm:"foreignKey:RatedByID"`
}

// VendorContact represents a contact person at a vendor
type VendorContact struct {
	Name     string `json:"name"`
	Position string `json:"position"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Primary  bool   `json:"primary"`
}

// VendorContacts is a slice of VendorContact
type VendorContacts []VendorContact

// Scan implements the sql.Scanner interface for VendorContacts
func (vc *VendorContacts) Scan(value interface{}) error {
	if value == nil {
		*vc = VendorContacts{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan VendorContacts: value is not []byte")
	}

	if err := json.Unmarshal(bytes, vc); err != nil {
		return err
	}

	return nil
}

// Value implements the driver.Valuer interface for VendorContacts
func (vc VendorContacts) Value() (driver.Value, error) {
	if vc == nil {
		return nil, nil
	}
	return json.Marshal(vc)
}

// VendorFilter represents filters for searching vendors
type VendorFilter struct {
	Code       string   `json:"code,omitempty"`
	Name       string   `json:"name,omitempty"`
	Type       string   `json:"type,omitempty"`
	Country    string   `json:"country,omitempty"`
	ProductIDs []uint   `json:"product_ids,omitempty"`
	MinRating  *float64 `json:"min_rating,omitempty"`
}

// VendorRepository defines the interface for vendor data access
type VendorRepository interface {
	Create(vendor *Vendor) error
	FindByID(id uint) (*Vendor, error)
	FindByCode(code string) (*Vendor, error)
	Update(vendor *Vendor) error
	Delete(id uint) error
	List(filter VendorFilter) ([]Vendor, error)

	CreateProduct(product *Product) error
	UpdateProduct(product *Product) error
	DeleteProduct(id uint) error
	FindProductByID(id uint) (*Product, error)
	FindProductByCode(code string) (*Product, error)
	ListProducts() ([]Product, error)

	AddProductToVendor(vendorID, productID uint) error
	RemoveProductFromVendor(vendorID, productID uint) error
	GetVendorProducts(vendorID uint) ([]Product, error)

	CreateContract(contract *Contract) error
	UpdateContract(contract *Contract) error
	DeleteContract(id uint) error
	FindContractByID(id uint) (*Contract, error)
	ListVendorContracts(vendorID uint) ([]Contract, error)

	CreateVendorRating(rating *VendorRating) error
	GetVendorRatings(vendorID uint) ([]VendorRating, error)
	GetVendorAverageRating(vendorID uint) (float64, error)
}
