package usecase

import (
	"context"

	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
	"github.com/lugondev/erp-warehouse-simple/internal/infrastructure/repository"
)

type InventoryUseCase struct {
	repo          *repository.InventoryRepository
	warehouseRepo *repository.WarehouseRepository
}

func NewInventoryUseCase(repo *repository.InventoryRepository, warehouseRepo *repository.WarehouseRepository) *InventoryUseCase {
	return &InventoryUseCase{
		repo:          repo,
		warehouseRepo: warehouseRepo,
	}
}

func (u *InventoryUseCase) GetInventory(ctx context.Context, id string) (*entity.Inventory, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *InventoryUseCase) ListInventory(ctx context.Context, filter *entity.InventoryFilter) ([]entity.Inventory, error) {
	return u.repo.List(ctx, filter)
}

func (u *InventoryUseCase) ProcessStockEntry(ctx context.Context, entry *entity.StockEntry, userID string) error {
	// Validate warehouse exists and is active
	warehouse, err := u.warehouseRepo.GetByID(ctx, entry.WarehouseID)
	if err != nil {
		return err
	}
	if warehouse.Status != entity.WarehouseStatusActive {
		return repository.ErrInvalidData
	}

	// Process stock entry with transaction
	return u.repo.ProcessStockEntry(ctx, entry, userID)
}

func (u *InventoryUseCase) CheckStock(ctx context.Context, productID string, warehouseID string) (*entity.Inventory, error) {
	inventory, err := u.repo.GetByProductAndWarehouse(ctx, productID, warehouseID)
	if err != nil {
		if err == repository.ErrRecordNotFound {
			// Return zero quantity inventory if not found
			return &entity.Inventory{
				ProductID:   productID,
				WarehouseID: warehouseID,
				Quantity:    0,
			}, nil
		}
		return nil, err
	}
	return inventory, nil
}

func (u *InventoryUseCase) UpdateStockLocation(ctx context.Context, id string, binLocation, shelfNumber, zoneCode string) error {
	inventory, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	inventory.BinLocation = binLocation
	inventory.ShelfNumber = shelfNumber
	inventory.ZoneCode = zoneCode

	return u.repo.CreateOrUpdateInventory(ctx, inventory)
}

func (u *InventoryUseCase) BatchStockEntry(ctx context.Context, entries []entity.StockEntry, userID string) error {
	// Process each entry in a separate transaction
	for _, entry := range entries {
		if err := u.ProcessStockEntry(ctx, &entry, userID); err != nil {
			return err
		}
	}
	return nil
}

func (u *InventoryUseCase) GetInventoryHistory(ctx context.Context, inventoryID string) ([]entity.InventoryHistory, error) {
	// This would require adding a new repository method
	// For now, return empty slice
	return []entity.InventoryHistory{}, nil
}

func (u *InventoryUseCase) ValidateStockEntry(entry *entity.StockEntry) error {
	if entry.Quantity <= 0 {
		return repository.ErrInvalidData
	}

	if entry.Type != "IN" && entry.Type != "OUT" {
		return repository.ErrInvalidData
	}

	if entry.WarehouseID == "" || entry.ProductID == "" {
		return repository.ErrInvalidData
	}

	return nil
}
