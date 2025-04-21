package entity

import (
	"time"
)

// Store represents a physical or virtual location where goods are stored
type Store struct {
	ID        string      `json:"id" gorm:"primaryKey;type:uuid"`
	Name      string      `json:"name" gorm:"not null;unique"`
	Code      string      `json:"code" gorm:"unique"`
	Address   string      `json:"address"`
	Type      StoreType   `json:"type" gorm:"not null"`
	ManagerID uint        `json:"manager_id" gorm:"not null"`
	Contact   string      `json:"contact"`
	Status    StoreStatus `json:"status" gorm:"not null;default:'ACTIVE'"`
	CreatedAt time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
	Manager   *User       `json:"manager,omitempty" gorm:"foreignKey:ManagerID"`
	Stocks    []Stock     `json:"stocks,omitempty" gorm:"foreignKey:StoreID"`
}

// StoreType represents the type of store
type StoreType string

const (
	StoreTypeRaw      StoreType = "RAW"
	StoreTypeFinished StoreType = "FINISHED"
	StoreTypeGeneral  StoreType = "GENERAL"
)

// StoreStatus represents the status of a store
type StoreStatus string

const (
	StoreStatusActive   StoreStatus = "ACTIVE"
	StoreStatusInactive StoreStatus = "INACTIVE"
)

// StoreFilter represents filters for searching stores
type StoreFilter struct {
	Name      string       `json:"name,omitempty"`
	Type      *StoreType   `json:"type,omitempty"`
	Status    *StoreStatus `json:"status,omitempty"`
	ManagerID *uint        `json:"manager_id,omitempty"`
}

// StockTransfer represents a transfer of stock between stores
type StockTransfer struct {
	ID                 string     `json:"id" gorm:"primaryKey;type:uuid"`
	SKUID              string     `json:"sku_id" gorm:"not null"`
	SourceStoreID      string     `json:"source_store_id" gorm:"not null"`
	DestinationStoreID string     `json:"destination_store_id" gorm:"not null"`
	Quantity           float64    `json:"quantity" gorm:"not null"`
	Status             string     `json:"status" gorm:"not null;default:'PENDING'"` // PENDING, COMPLETED, CANCELLED
	RequestedByID      uint       `json:"requested_by_id" gorm:"not null"`
	ApprovedByID       *uint      `json:"approved_by_id"`
	CompletedByID      *uint      `json:"completed_by_id"`
	RequestedAt        time.Time  `json:"requested_at" gorm:"autoCreateTime"`
	ApprovedAt         *time.Time `json:"approved_at"`
	CompletedAt        *time.Time `json:"completed_at"`
	Notes              string     `json:"notes"`
	SourceStore        *Store     `json:"source_store,omitempty" gorm:"foreignKey:SourceStoreID"`
	DestinationStore   *Store     `json:"destination_store,omitempty" gorm:"foreignKey:DestinationStoreID"`
	SKU                *SKU       `json:"sku,omitempty" gorm:"foreignKey:SKUID"`
	RequestedBy        *User      `json:"requested_by,omitempty" gorm:"foreignKey:RequestedByID"`
	ApprovedBy         *User      `json:"approved_by,omitempty" gorm:"foreignKey:ApprovedByID"`
	CompletedBy        *User      `json:"completed_by,omitempty" gorm:"foreignKey:CompletedByID"`
}

// StoreReport represents a summary report for a store
type StoreReport struct {
	ID             string    `json:"id" gorm:"primaryKey;type:uuid"`
	StoreID        string    `json:"store_id" gorm:"not null"`
	ReportDate     time.Time `json:"report_date" gorm:"not null"`
	TotalSKUs      int       `json:"total_skus"`
	TotalQuantity  float64   `json:"total_quantity"`
	TotalValue     float64   `json:"total_value"`
	LowStockSKUs   int       `json:"low_stock_skus"`
	OutOfStockSKUs int       `json:"out_of_stock_skus"`
	ExpiringItems  int       `json:"expiring_items"`
	ExpiredItems   int       `json:"expired_items"`
	CreatedByID    uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	Store          *Store    `json:"store,omitempty" gorm:"foreignKey:StoreID"`
	CreatedBy      *User     `json:"created_by,omitempty" gorm:"foreignKey:CreatedByID"`
}
