package entity

import "time"

type Inventory struct {
	ID              string    `json:"id" gorm:"primaryKey"`
	ProductID       string    `json:"product_id" gorm:"not null"`
	WarehouseID     string    `json:"warehouse_id" gorm:"not null"`
	Quantity        float64   `json:"quantity" gorm:"not null"`
	BinLocation     string    `json:"bin_location"`
	ShelfNumber     string    `json:"shelf_number"`
	ZoneCode        string    `json:"zone_code"`
	BatchNumber     string    `json:"batch_number"`
	LotNumber       string    `json:"lot_number"`
	ManufactureDate time.Time `json:"manufacture_date"`
	ExpiryDate      time.Time `json:"expiry_date"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type StockEntry struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	WarehouseID string    `json:"warehouse_id" gorm:"not null"`
	ProductID   string    `json:"product_id" gorm:"not null"`
	Type        string    `json:"type" gorm:"not null"` // IN/OUT
	Quantity    float64   `json:"quantity" gorm:"not null"`
	BatchNumber string    `json:"batch_number"`
	LotNumber   string    `json:"lot_number"`
	Reference   string    `json:"reference"` // PO/SO number
	Note        string    `json:"note"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	CreatedBy   string    `json:"created_by" gorm:"not null"`
}

type InventoryHistory struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	InventoryID string    `json:"inventory_id" gorm:"not null"`
	Type        string    `json:"type" gorm:"not null"` // IN/OUT/ADJUST
	Quantity    float64   `json:"quantity" gorm:"not null"`
	PreviousQty float64   `json:"previous_qty"`
	NewQty      float64   `json:"new_qty"`
	Reference   string    `json:"reference"` // Stock Entry ID
	Note        string    `json:"note"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	CreatedBy   string    `json:"created_by" gorm:"not null"`
}

type InventoryFilter struct {
	WarehouseID string `json:"warehouse_id,omitempty"`
	ProductID   string `json:"product_id,omitempty"`
	BatchNumber string `json:"batch_number,omitempty"`
	LotNumber   string `json:"lot_number,omitempty"`
}
