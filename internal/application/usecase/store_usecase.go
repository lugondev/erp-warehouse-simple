package usecase

import (
	"context"

	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
	"github.com/lugondev/erp-warehouse-simple/internal/infrastructure/repository"
)

type StoreUseCase struct {
	repo *repository.StoreRepository
}

func NewStoreUseCase(repo *repository.StoreRepository) *StoreUseCase {
	return &StoreUseCase{repo: repo}
}

func (u *StoreUseCase) CreateStore(ctx context.Context, store *entity.Store) error {
	if store.Status == "" {
		store.Status = entity.StoreStatusActive
	}
	return u.repo.CreateStore(ctx, store)
}

func (u *StoreUseCase) UpdateStore(ctx context.Context, store *entity.Store) error {
	existing, err := u.repo.GetByID(ctx, store.ID)
	if err != nil {
		return err
	}

	// Update fields
	existing.Name = store.Name
	existing.Address = store.Address
	existing.Type = store.Type
	existing.ManagerID = store.ManagerID
	existing.Contact = store.Contact
	existing.Status = store.Status

	return u.repo.Update(ctx, existing)
}

func (u *StoreUseCase) GetStore(ctx context.Context, id string) (*entity.Store, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *StoreUseCase) ListStores(ctx context.Context, filter *entity.StoreFilter) ([]entity.Store, error) {
	return u.repo.List(ctx, filter)
}

func (u *StoreUseCase) DeleteStore(ctx context.Context, id string) error {
	// Could add additional checks here (e.g., ensure store is empty)
	return u.repo.Delete(ctx, id)
}

// Additional business logic methods

func (u *StoreUseCase) GetStoreStocks(ctx context.Context, storeID string) ([]entity.Stock, error) {
	return u.repo.GetStoreStocks(ctx, storeID)
}

func (u *StoreUseCase) GetStoreStockValue(ctx context.Context, storeID string) (float64, error) {
	return u.repo.GetStoreStockValue(ctx, storeID)
}

func (u *StoreUseCase) GetStoresWithLowStock(ctx context.Context, skuIDs []string, threshold float64) (map[string][]entity.Stock, error) {
	return u.repo.GetStoresWithLowStock(ctx, skuIDs, threshold)
}

func (u *StoreUseCase) AssignManager(ctx context.Context, storeID string, managerID uint) error {
	return u.repo.AssignManager(ctx, storeID, managerID)
}

func (u *StoreUseCase) UpdateStatus(ctx context.Context, storeID string, status entity.StoreStatus) error {
	return u.repo.UpdateStatus(ctx, storeID, status)
}
