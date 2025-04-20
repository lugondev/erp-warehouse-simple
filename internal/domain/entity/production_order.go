package entity

import "time"

type ProductionOrderStatus string

const (
	OrderStatusPending   ProductionOrderStatus = "pending"
	OrderStatusInProcess ProductionOrderStatus = "in_process"
	OrderStatusCompleted ProductionOrderStatus = "completed"
	OrderStatusCancelled ProductionOrderStatus = "cancelled"
)

type ProductionOrder struct {
	ID           uint                  `json:"id" gorm:"primaryKey"`
	ProductID    uint                  `json:"product_id" gorm:"not null"`
	Quantity     int                   `json:"quantity" gorm:"not null"`
	StartDate    time.Time             `json:"start_date"`
	Deadline     time.Time             `json:"deadline" gorm:"not null"`
	Status       ProductionOrderStatus `json:"status" gorm:"not null;default:'pending'"`
	FacilityID   uint                  `json:"facility_id" gorm:"not null"`
	CompletedQty int                   `json:"completed_qty" gorm:"default:0"`
	DefectQty    int                   `json:"defect_qty" gorm:"default:0"`
	Notes        string                `json:"notes"`
	CreatedAt    time.Time             `json:"created_at"`
	UpdatedAt    time.Time             `json:"updated_at"`
}
