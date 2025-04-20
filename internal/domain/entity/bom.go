package entity

import "time"

// BillOfMaterial represents a bill of materials for a product
type BillOfMaterial struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ProductID uint      `json:"product_id" gorm:"not null"`
	Name      string    `json:"name" gorm:"not null"`
	Version   string    `json:"version" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BOMItem represents an item in the bill of materials
type BOMItem struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	BOMID          uint      `json:"bom_id" gorm:"not null"`
	MaterialID     uint      `json:"material_id" gorm:"not null"` // References inventory item
	QuantityNeeded float64   `json:"quantity_needed" gorm:"not null"`
	UnitOfMeasure  string    `json:"unit_of_measure" gorm:"not null"`
	Notes          string    `json:"notes"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// MRPCalculation represents material requirements planning calculation
type MRPCalculation struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	ProductionID  uint      `json:"production_id" gorm:"not null"`
	MaterialID    uint      `json:"material_id" gorm:"not null"`
	RequiredQty   float64   `json:"required_qty" gorm:"not null"`
	AvailableQty  float64   `json:"available_qty"`
	ShortageQty   float64   `json:"shortage_qty"`
	UnitOfMeasure string    `json:"unit_of_measure"`
	CalculatedAt  time.Time `json:"calculated_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
