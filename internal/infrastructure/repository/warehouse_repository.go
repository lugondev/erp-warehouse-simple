package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
	"gorm.io/gorm"
)

type WarehouseRepository struct {
	db *gorm.DB
}

func NewWarehouseRepository(db *gorm.DB) *WarehouseRepository {
	return &WarehouseRepository{db: db}
}

func (r *WarehouseRepository) Create(ctx context.Context, warehouse *entity.Warehouse) error {
	if warehouse.ID == "" {
		warehouse.ID = uuid.New().String()
	}
	return r.db.WithContext(ctx).Create(warehouse).Error
}

func (r *WarehouseRepository) Update(ctx context.Context, warehouse *entity.Warehouse) error {
	return r.db.WithContext(ctx).Save(warehouse).Error
}

func (r *WarehouseRepository) GetByID(ctx context.Context, id string) (*entity.Warehouse, error) {
	var warehouse entity.Warehouse
	if err := r.db.WithContext(ctx).First(&warehouse, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &warehouse, nil
}

func (r *WarehouseRepository) List(ctx context.Context, filter *entity.WarehouseFilter) ([]entity.Warehouse, error) {
	var warehouses []entity.Warehouse
	query := r.db.WithContext(ctx)

	if filter != nil {
		if filter.Type != nil {
			query = query.Where("type = ?", *filter.Type)
		}
		if filter.Status != nil {
			query = query.Where("status = ?", *filter.Status)
		}
	}

	if err := query.Find(&warehouses).Error; err != nil {
		return nil, err
	}
	return warehouses, nil
}

func (r *WarehouseRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&entity.Warehouse{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
