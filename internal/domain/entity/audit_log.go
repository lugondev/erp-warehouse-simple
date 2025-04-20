package entity

import (
	"time"
)

type ActionType string

const (
	ActionCreate ActionType = "create"
	ActionRead   ActionType = "read"
	ActionUpdate ActionType = "update"
	ActionDelete ActionType = "delete"
	ActionLogin  ActionType = "login"
	ActionLogout ActionType = "logout"
)

type AuditLog struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	UserID    uint       `json:"user_id"`
	User      *User      `json:"user" gorm:"foreignKey:UserID"`
	Action    ActionType `json:"action" gorm:"type:varchar(20)"`
	Resource  string     `json:"resource" gorm:"type:varchar(50)"`
	Detail    string     `json:"detail" gorm:"type:text"`
	IP        string     `json:"ip" gorm:"type:varchar(45)"`
	UserAgent string     `json:"user_agent" gorm:"type:text"`
	CreatedAt time.Time  `json:"created_at"`
}

type AuditLogRepository interface {
	Create(log *AuditLog) error
	FindByUserID(userID uint, limit, offset int) ([]AuditLog, error)
	FindByAction(action ActionType, limit, offset int) ([]AuditLog, error)
	FindByDateRange(start, end time.Time, limit, offset int) ([]AuditLog, error)
	List(limit, offset int) ([]AuditLog, error)
	Count(filter map[string]interface{}) (int64, error)
}
