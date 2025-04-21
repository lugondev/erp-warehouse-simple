package entity

import (
	"database/sql/driver"
	"time"

	"github.com/lib/pq"
)

// Additional permissions specific to roles and modules
const (
	ModuleIntegrate Permission = "module:integrate"
)

// GormPermissionSlice is a helper type for GORM operations
type GormPermissionSlice []Permission

// Value implements the driver.Valuer interface
func (p GormPermissionSlice) Value() (driver.Value, error) {
	if len(p) == 0 {
		return "{}", nil
	}
	var values []string
	for _, perm := range p {
		values = append(values, string(perm))
	}
	return pq.Array(values), nil
}

// Scan implements the sql.Scanner interface
func (p *GormPermissionSlice) Scan(value interface{}) error {
	if value == nil {
		*p = make([]Permission, 0)
		return nil
	}
	var array []string
	if err := pq.Array(&array).Scan(value); err != nil {
		return err
	}
	*p = make([]Permission, len(array))
	for i, v := range array {
		(*p)[i] = Permission(v)
	}
	return nil
}

type Role struct {
	ID          uint                `json:"id" gorm:"primaryKey"`
	Name        string              `json:"name" gorm:"unique;not null"`
	Permissions GormPermissionSlice `json:"permissions" gorm:"type:text[]"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

type RoleRepository interface {
	Create(role *Role) error
	FindByID(id uint) (*Role, error)
	FindByName(name string) (*Role, error)
	Update(role *Role) error
	Delete(id uint) error
	List() ([]Role, error)
}
