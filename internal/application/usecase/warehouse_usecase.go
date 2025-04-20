package usecase

import (
	"context"

	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
	"github.com/lugondev/erp-warehouse-simple/internal/infrastructure/repository"
)

type WarehouseUseCase struct {
	repo *repository.WarehouseRepository
}

func NewWarehouseUseCase(repo *repository.WarehouseRepository) *WarehouseUseCase {
	return &WarehouseUseCase{repo: repo}
}

func (u *WarehouseUseCase) CreateWarehouse(ctx context.Context, warehouse *entity.Warehouse) error {
	if warehouse.Status == "" {
		warehouse.Status = entity.WarehouseStatusActive
	}
	return u.repo.Create(ctx, warehouse)
}

func (u *WarehouseUseCase) UpdateWarehouse(ctx context.Context, warehouse *entity.Warehouse) error {
	existing, err := u.repo.GetByID(ctx, warehouse.ID)
	if err != nil {
		return err
	}

	// Update fields
	existing.Name = warehouse.Name
	existing.Address = warehouse.Address
	existing.Type = warehouse.Type
	existing.ManagerID = warehouse.ManagerID
	existing.Contact = warehouse.Contact
	existing.Status = warehouse.Status

	return u.repo.Update(ctx, existing)
}

func (u *WarehouseUseCase) GetWarehouse(ctx context.Context, id string) (*entity.Warehouse, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *WarehouseUseCase) ListWarehouses(ctx context.Context, filter *entity.WarehouseFilter) ([]entity.Warehouse, error) {
	return u.repo.List(ctx, filter)
}

func (u *WarehouseUseCase) DeleteWarehouse(ctx context.Context, id string) error {
	// Could add additional checks here (e.g., ensure warehouse is empty)
	return u.repo.Delete(ctx, id)
}
