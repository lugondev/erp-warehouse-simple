package repository

import (
	"time"

	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
	"gorm.io/gorm"
)

type AuditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

func (r *AuditLogRepository) Create(log *entity.AuditLog) error {
	return r.db.Create(log).Error
}

func (r *AuditLogRepository) FindByUserID(userID uint, limit, offset int) ([]entity.AuditLog, error) {
	var logs []entity.AuditLog
	err := r.db.Preload("User").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	return logs, err
}

func (r *AuditLogRepository) FindByAction(action entity.ActionType, limit, offset int) ([]entity.AuditLog, error) {
	var logs []entity.AuditLog
	err := r.db.Preload("User").
		Where("action = ?", action).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	return logs, err
}

func (r *AuditLogRepository) FindByDateRange(start, end time.Time, limit, offset int) ([]entity.AuditLog, error) {
	var logs []entity.AuditLog
	err := r.db.Preload("User").
		Where("created_at BETWEEN ? AND ?", start, end).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	return logs, err
}

func (r *AuditLogRepository) List(limit, offset int) ([]entity.AuditLog, error) {
	var logs []entity.AuditLog
	err := r.db.Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	return logs, err
}

func (r *AuditLogRepository) Count(filter map[string]interface{}) (int64, error) {
	var count int64
	query := r.db.Model(&entity.AuditLog{})

	for key, value := range filter {
		query = query.Where(key+" = ?", value)
	}

	err := query.Count(&count).Error
	return count, err
}
