package entity

import "time"

type ManufacturingFacility struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Address   string    `json:"address" gorm:"not null"`
	Type      string    `json:"type" gorm:"not null"`
	Capacity  int       `json:"capacity" gorm:"not null"` // Production capacity per day
	Manager   string    `json:"manager"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
