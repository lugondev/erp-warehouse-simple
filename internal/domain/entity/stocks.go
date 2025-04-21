package entity

import (
	"time"
)

// Stock represents stock of an item in a store
type Stock struct {
	ID              string    `json:"id" gorm:"primaryKey;type:uuid"`
	SKUID           string    `json:"sku_id" gorm:"not null"`
	StoreID         string    `json:"store_id" gorm:"not null"`
	Quantity        float64   `json:"quantity" gorm:"not null;default:0"`
	BinLocation     string    `json:"bin_location"`
	ShelfNumber     string    `json:"shelf_number"`
	ZoneCode        string    `json:"zone_code"`
	BatchNumber     string    `json:"batch_number"`
	LotNumber       string    `json:"lot_number"`
	ManufactureDate time.Time `json:"manufacture_date"`
	ExpiryDate      time.Time `json:"expiry_date"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	SKU             *SKU      `json:"sku,omitempty" gorm:"foreignKey:SKUID"`
	Store           *Store    `json:"store,omitempty" gorm:"foreignKey:StoreID"`
}

// StockEntry represents a stock movement entry
type StockEntry struct {
	ID              string    `json:"id" gorm:"primaryKey;type:uuid"`
	SKUID           string    `json:"sku_id" gorm:"not null"`
	StoreID         string    `json:"store_id" gorm:"not null"`
	Type            string    `json:"type" gorm:"not null"` // IN, OUT
	Quantity        float64   `json:"quantity" gorm:"not null"`
	BatchNumber     string    `json:"batch_number"`
	LotNumber       string    `json:"lot_number"`
	ManufactureDate time.Time `json:"manufacture_date"`
	ExpiryDate      time.Time `json:"expiry_date"`
	Reference       string    `json:"reference"`
	Note            string    `json:"note"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	CreatedBy       string    `json:"created_by" gorm:"not null"`
	SKU             *SKU      `json:"sku,omitempty" gorm:"foreignKey:SKUID"`
	Store           *Store    `json:"store,omitempty" gorm:"foreignKey:StoreID"`
}

// StockHistory represents a history of stock changes
type StockHistory struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid"`
	StockID     string    `json:"stock_id" gorm:"not null"`
	Type        string    `json:"type" gorm:"not null"` // IN, OUT, ADJUST
	Quantity    float64   `json:"quantity" gorm:"not null"`
	PreviousQty float64   `json:"previous_qty" gorm:"not null"`
	NewQty      float64   `json:"new_qty" gorm:"not null"`
	Reference   string    `json:"reference"` // Reference to a StockEntry
	Note        string    `json:"note"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	CreatedBy   string    `json:"created_by" gorm:"not null"`
	Stock       *Stock    `json:"stock,omitempty" gorm:"foreignKey:StockID"`
}

// StockFilter represents filters for searching stocks
type StockFilter struct {
	SKUID          string    `json:"sku_id,omitempty"`
	StoreID        string    `json:"store_id,omitempty"`
	MinQuantity    float64   `json:"min_quantity,omitempty"`
	MaxQuantity    float64   `json:"max_quantity,omitempty"`
	BatchNumber    string    `json:"batch_number,omitempty"`
	LotNumber      string    `json:"lot_number,omitempty"`
	ZoneCode       string    `json:"zone_code,omitempty"`
	BinLocation    string    `json:"bin_location,omitempty"`
	ShelfNumber    string    `json:"shelf_number,omitempty"`
	ExpiryDateFrom time.Time `json:"expiry_date_from,omitempty"`
	ExpiryDateTo   time.Time `json:"expiry_date_to,omitempty"`
}
