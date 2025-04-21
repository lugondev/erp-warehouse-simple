package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
	"gorm.io/gorm"
)

type StoreRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) *StoreRepository {
	return &StoreRepository{db: db}
}

// CreateStore creates a new store
func (r *StoreRepository) CreateStore(ctx context.Context, store *entity.Store) error {
	if store.ID == "" {
		store.ID = uuid.New().String()
	}
	return r.db.WithContext(ctx).Create(store).Error
}

// GetByID retrieves a store by ID
func (r *StoreRepository) GetByID(ctx context.Context, id string) (*entity.Store, error) {
	var store entity.Store
	if err := r.db.WithContext(ctx).First(&store, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &store, nil
}

// GetByCode retrieves a store by code
func (r *StoreRepository) GetByCode(ctx context.Context, code string) (*entity.Store, error) {
	var store entity.Store
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&store).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &store, nil
}

// Update updates a store
func (r *StoreRepository) Update(ctx context.Context, store *entity.Store) error {
	return r.db.WithContext(ctx).Save(store).Error
}

// Delete deletes a store
func (r *StoreRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.Store{}, "id = ?", id).Error
}

// List lists stores with optional filtering
func (r *StoreRepository) List(ctx context.Context, filter *entity.StoreFilter) ([]entity.Store, error) {
	var stores []entity.Store
	query := r.db.WithContext(ctx).Model(&entity.Store{})

	// Apply filters if provided
	if filter != nil {
		if filter.Name != "" {
			query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
		}
		if filter.Type != nil {
			query = query.Where("type = ?", filter.Type)
		}
		if filter.Status != nil {
			query = query.Where("status = ?", filter.Status)
		}
		if filter.ManagerID != nil {
			query = query.Where("manager_id = ?", filter.ManagerID)
		}
	}

	// Execute the query
	if err := query.Find(&stores).Error; err != nil {
		return nil, err
	}

	return stores, nil
}

// GetStoreStocks gets all stocks in a store
func (r *StoreRepository) GetStoreStocks(ctx context.Context, storeID string) ([]entity.Stock, error) {
	var stocks []entity.Stock
	if err := r.db.WithContext(ctx).
		Where("store_id = ?", storeID).
		Find(&stocks).Error; err != nil {
		return nil, err
	}
	return stocks, nil
}

// GetStoreStockValue calculates the total value of all stocks in a store
func (r *StoreRepository) GetStoreStockValue(ctx context.Context, storeID string) (float64, error) {
	type Result struct {
		TotalValue float64
	}
	var result Result

	// This query assumes you have a way to get the value of each stock item
	// by joining with the SKU table to get the price
	err := r.db.WithContext(ctx).
		Table("stocks").
		Select("COALESCE(SUM(stocks.quantity * skus.price), 0) as total_value").
		Joins("JOIN skus ON stocks.sku_id = skus.id").
		Where("stocks.store_id = ?", storeID).
		Scan(&result).Error

	if err != nil {
		return 0, err
	}

	return result.TotalValue, nil
}

// GetStoresWithLowStock finds stores with stock levels below threshold for specified SKUs
func (r *StoreRepository) GetStoresWithLowStock(ctx context.Context, skuIDs []string, threshold float64) (map[string][]entity.Stock, error) {
	var stocks []entity.Stock
	err := r.db.WithContext(ctx).
		Where("sku_id IN ? AND quantity <= ?", skuIDs, threshold).
		Find(&stocks).Error

	if err != nil {
		return nil, err
	}

	// Group results by store ID
	storeStocks := make(map[string][]entity.Stock)
	for _, stock := range stocks {
		storeStocks[stock.StoreID] = append(storeStocks[stock.StoreID], stock)
	}

	return storeStocks, nil
}

// AssignManager assigns a manager to a store
func (r *StoreRepository) AssignManager(ctx context.Context, storeID string, managerID uint) error {
	return r.db.WithContext(ctx).
		Model(&entity.Store{}).
		Where("id = ?", storeID).
		Update("manager_id", managerID).Error
}

// UpdateStatus updates a store's status
func (r *StoreRepository) UpdateStatus(ctx context.Context, storeID string, status entity.StoreStatus) error {
	return r.db.WithContext(ctx).
		Model(&entity.Store{}).
		Where("id = ?", storeID).
		Update("status", status).Error
}
