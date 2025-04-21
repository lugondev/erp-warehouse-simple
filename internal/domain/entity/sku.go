package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// SKU represents a stock keeping unit in the system
type SKU struct {
	ID             string    `json:"id" gorm:"primaryKey;type:uuid"`
	SKUCode        string    `json:"sku_code" gorm:"uniqueIndex;not null"`
	Name           string    `json:"name" gorm:"not null"`
	Description    string    `json:"description"`
	UnitOfMeasure  string    `json:"unit_of_measure" gorm:"not null"`
	Price          float64   `json:"price" gorm:"default:0"`
	Category       string    `json:"category"`
	TechnicalSpecs JSONMap   `json:"technical_specs" gorm:"type:jsonb"`
	ManufacturerID *uint     `json:"manufacturer_id"`
	VendorID       *uint     `json:"vendor_id"`
	ImageURL       string    `json:"image_url"`
	Status         SKUStatus `json:"status" gorm:"default:'ACTIVE'"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Manufacturer   *Vendor   `json:"manufacturer,omitempty" gorm:"foreignKey:ManufacturerID"`
	Vendor         *Vendor   `json:"vendor,omitempty" gorm:"foreignKey:VendorID"`
}

// SKUStatus represents the status of a SKU
type SKUStatus string

const (
	SKUStatusActive   SKUStatus = "ACTIVE"
	SKUStatusInactive SKUStatus = "INACTIVE"
	SKUStatusArchived SKUStatus = "ARCHIVED"
)

// JSONMap is a helper type for JSON fields
type JSONMap map[string]interface{}

// Scan implements the sql.Scanner interface for JSONMap
func (m *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*m = make(JSONMap)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan JSONMap: value is not []byte")
	}

	if err := json.Unmarshal(bytes, m); err != nil {
		return err
	}

	return nil
}

// Value implements the driver.Valuer interface for JSONMap
func (m JSONMap) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

// SKUCategory represents a product category
type SKUCategory struct {
	ID          string        `json:"id" gorm:"primaryKey;type:uuid"`
	Name        string        `json:"name" gorm:"not null;unique"`
	Description string        `json:"description"`
	ParentID    *string       `json:"parent_id"`
	CreatedAt   time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
	Parent      *SKUCategory  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children    []SKUCategory `json:"children,omitempty" gorm:"foreignKey:ParentID"`
}

// SKUFilter represents filters for searching SKUs
type SKUFilter struct {
	SKUCode        string     `json:"sku_code,omitempty"`
	Name           string     `json:"name,omitempty"`
	Category       string     `json:"category,omitempty"`
	ManufacturerID *uint      `json:"manufacturer_id,omitempty"`
	VendorID       *uint      `json:"vendor_id,omitempty"`
	Status         *SKUStatus `json:"status,omitempty"`
	MinPrice       *float64   `json:"min_price,omitempty"`
	MaxPrice       *float64   `json:"max_price,omitempty"`
}
