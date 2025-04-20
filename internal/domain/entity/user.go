package entity

import (
	"time"
)

type UserStatus string

const (
	StatusActive   UserStatus = "active"
	StatusInactive UserStatus = "inactive"
	StatusLocked   UserStatus = "locked"
)

type User struct {
	ID                 uint       `json:"id" gorm:"primaryKey"`
	Username           string     `json:"username" gorm:"unique;not null"`
	Email              string     `json:"email" gorm:"unique;not null"`
	Password           string     `json:"-" gorm:"not null"`
	RoleID             uint       `json:"role_id" gorm:"not null"`
	Role               *Role      `json:"role" gorm:"foreignKey:RoleID"`
	Status             UserStatus `json:"status" gorm:"type:varchar(20);default:'active'"`
	LastLogin          *time.Time `json:"last_login,omitempty"`
	RefreshToken       string     `json:"-" gorm:"type:text"`
	RefreshTokenExpiry time.Time  `json:"-"`
	PasswordResetToken string     `json:"-" gorm:"type:varchar(100)"`
	ResetTokenExpiry   time.Time  `json:"-"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type UserRepository interface {
	Create(user *User) error
	FindByID(id uint) (*User, error)
	FindByEmail(email string) (*User, error)
	FindByUsername(username string) (*User, error)
	Update(user *User) error
	Delete(id uint) error
	List() ([]User, error)
	UpdateRefreshToken(userID uint, token string, expiry time.Time) error
	FindByRefreshToken(token string) (*User, error)
	UpdatePasswordResetToken(userID uint, token string, expiry time.Time) error
	FindByPasswordResetToken(token string) (*User, error)
	UpdatePassword(userID uint, hashedPassword string) error
	UpdateLastLogin(userID uint) error
}

// HasPermission checks if the user has a specific permission
func (u *User) HasPermission(permission Permission) bool {
	if u.Role == nil {
		return false
	}

	for _, p := range u.Role.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// IsActive checks if the user account is active
func (u *User) IsActive() bool {
	return u.Status == StatusActive
}

// IsLocked checks if the user account is locked
func (u *User) IsLocked() bool {
	return u.Status == StatusLocked
}
