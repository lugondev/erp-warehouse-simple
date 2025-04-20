package entity

import "time"

type WarehouseType string

const (
	WarehouseTypeRaw      WarehouseType = "RAW"
	WarehouseTypeFinished WarehouseType = "FINISHED"
	WarehouseTypeGeneral  WarehouseType = "GENERAL"
)

type WarehouseStatus string

const (
	WarehouseStatusActive   WarehouseStatus = "ACTIVE"
	WarehouseStatusInactive WarehouseStatus = "INACTIVE"
)

type Warehouse struct {
	ID        string          `json:"id" gorm:"primaryKey"`
	Name      string          `json:"name" gorm:"not null"`
	Address   string          `json:"address"`
	Type      WarehouseType   `json:"type" gorm:"not null"`
	ManagerID string          `json:"manager_id" gorm:"not null"`
	Contact   string          `json:"contact"`
	Status    WarehouseStatus `json:"status" gorm:"not null"`
	CreatedAt time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
}

type WarehouseFilter struct {
	Type   *WarehouseType   `json:"type,omitempty"`
	Status *WarehouseStatus `json:"status,omitempty"`
}
