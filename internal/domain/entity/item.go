package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// ItemStatus represents the status of an item
type ItemStatus string

const (
	ItemStatusActive   ItemStatus = "ACTIVE"
	ItemStatusInactive ItemStatus = "INACTIVE"
	ItemStatusArchived ItemStatus = "ARCHIVED"
)

// TechnicalSpecs represents the technical specifications of an item as a JSON object
type TechnicalSpecs map[string]interface{}

// Scan implements the sql.Scanner interface for TechnicalSpecs
func (ts *TechnicalSpecs) Scan(value interface{}) error {
	if value == nil {
		*ts = make(TechnicalSpecs)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan TechnicalSpecs: value is not []byte")
	}

	if err := json.Unmarshal(bytes, ts); err != nil {
		return err
	}
	return nil
}

// Value implements the driver.Valuer interface for TechnicalSpecs
func (ts TechnicalSpecs) Value() (driver.Value, error) {
	if ts == nil {
		return nil, nil
	}
	return json.Marshal(ts)
}

// Item represents a product/SKU in the system
type Item struct {
	ID             string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	SKU            string         `json:"sku" gorm:"uniqueIndex;not null"`
	Name           string         `json:"name" gorm:"not null"`
	Description    string         `json:"description" gorm:"type:text"`
	UnitOfMeasure  string         `json:"unit_of_measure" gorm:"not null"`
	Price          float64        `json:"price" gorm:"type:decimal(15,2)"`
	Category       string         `json:"category"`
	TechnicalSpecs TechnicalSpecs `json:"technical_specs" gorm:"type:jsonb"`
	ManufacturerID *uint          `json:"manufacturer_id" gorm:"index"`
	SupplierID     *uint          `json:"supplier_id" gorm:"index"`
	ImageURL       string         `json:"image_url"`
	Status         ItemStatus     `json:"status" gorm:"not null;default:'ACTIVE'"`
	CreatedAt      time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	Manufacturer   *Supplier      `json:"manufacturer,omitempty" gorm:"foreignKey:ManufacturerID"`
	Supplier       *Supplier      `json:"supplier,omitempty" gorm:"foreignKey:SupplierID"`
}

// ItemCategory represents a category for items
type ItemCategory struct {
	ID          string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string         `json:"name" gorm:"uniqueIndex;not null"`
	Description string         `json:"description" gorm:"type:text"`
	ParentID    *string        `json:"parent_id" gorm:"type:uuid"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	Parent      *ItemCategory  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children    []ItemCategory `json:"children,omitempty" gorm:"foreignKey:ParentID"`
}

// ItemFilter represents filters for searching items
type ItemFilter struct {
	SKU            string      `json:"sku,omitempty"`
	Name           string      `json:"name,omitempty"`
	Category       string      `json:"category,omitempty"`
	ManufacturerID *uint       `json:"manufacturer_id,omitempty"`
	SupplierID     *uint       `json:"supplier_id,omitempty"`
	Status         *ItemStatus `json:"status,omitempty"`
	MinPrice       *float64    `json:"min_price,omitempty"`
	MaxPrice       *float64    `json:"max_price,omitempty"`
}
