package usecase

import (
	"context"

	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
	"github.com/lugondev/erp-warehouse-simple/internal/infrastructure/repository"
)

type StocksUseCase struct {
	repo      *repository.StocksRepository
	storeRepo *repository.StoreRepository
}

func NewStocksUseCase(repo *repository.StocksRepository, storeRepo *repository.StoreRepository) *StocksUseCase {
	return &StocksUseCase{
		repo:      repo,
		storeRepo: storeRepo,
	}
}

func (u *StocksUseCase) GetStock(ctx context.Context, id string) (*entity.Stock, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *StocksUseCase) ListStocks(ctx context.Context, filter *entity.StockFilter) ([]entity.Stock, error) {
	return u.repo.List(ctx, filter)
}

func (u *StocksUseCase) ProcessStockEntry(ctx context.Context, entry *entity.StockEntry, userID string) error {
	// Validate store exists and is active
	store, err := u.storeRepo.GetByID(ctx, entry.StoreID)
	if err != nil {
		return err
	}
	if store.Status != entity.StoreStatusActive {
		return repository.ErrInvalidData
	}

	// Process stock entry with transaction
	return u.repo.ProcessStockEntry(ctx, entry, userID)
}

func (u *StocksUseCase) CheckStock(ctx context.Context, skuID string, storeID string) (*entity.Stock, error) {
	stock, err := u.repo.GetBySKUAndStore(ctx, skuID, storeID)
	if err != nil {
		if err == repository.ErrRecordNotFound {
			// Return zero quantity stock if not found
			return &entity.Stock{
				SKUID:    skuID,
				StoreID:  storeID,
				Quantity: 0,
			}, nil
		}
		return nil, err
	}
	return stock, nil
}

func (u *StocksUseCase) UpdateStockLocation(ctx context.Context, id string, binLocation, shelfNumber, zoneCode string) error {
	stock, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	stock.BinLocation = binLocation
	stock.ShelfNumber = shelfNumber
	stock.ZoneCode = zoneCode

	return u.repo.CreateOrUpdateStock(ctx, stock)
}

func (u *StocksUseCase) BatchStockEntry(ctx context.Context, entries []entity.StockEntry, userID string) error {
	// Process each entry in a separate transaction
	for _, entry := range entries {
		if err := u.ProcessStockEntry(ctx, &entry, userID); err != nil {
			return err
		}
	}
	return nil
}

func (u *StocksUseCase) GetStockHistory(ctx context.Context, stockID string) ([]entity.StockHistory, error) {
	// This would require adding a new repository method
	// For now, return empty slice
	return []entity.StockHistory{}, nil
}

func (u *StocksUseCase) ValidateStockEntry(entry *entity.StockEntry) error {
	if entry.Quantity <= 0 {
		return repository.ErrInvalidData
	}

	if entry.Type != "IN" && entry.Type != "OUT" {
		return repository.ErrInvalidData
	}

	if entry.StoreID == "" || entry.SKUID == "" {
		return repository.ErrInvalidData
	}

	return nil
}
