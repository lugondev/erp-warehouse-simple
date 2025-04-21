package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
	"gorm.io/gorm"
)

type StocksRepository struct {
	db *gorm.DB
}

func NewStocksRepository(db *gorm.DB) *StocksRepository {
	return &StocksRepository{db: db}
}

// GetByID retrieves a stock record by ID
func (r *StocksRepository) GetByID(ctx context.Context, id string) (*entity.Stock, error) {
	var stock entity.Stock
	if err := r.db.WithContext(ctx).First(&stock, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &stock, nil
}

// GetBySKUAndStore retrieves a stock record by SKU ID and store ID
func (r *StocksRepository) GetBySKUAndStore(ctx context.Context, skuID, storeID string) (*entity.Stock, error) {
	var stock entity.Stock
	if err := r.db.WithContext(ctx).
		Where("sku_id = ? AND store_id = ?", skuID, storeID).
		First(&stock).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &stock, nil
}

// List retrieves stocks with filtering
func (r *StocksRepository) List(ctx context.Context, filter *entity.StockFilter) ([]entity.Stock, error) {
	var stocks []entity.Stock
	query := r.db.WithContext(ctx).
		Model(&entity.Stock{})

	// Apply filters if provided
	if filter != nil {
		if filter.SKUID != "" {
			query = query.Where("sku_id = ?", filter.SKUID)
		}
		if filter.StoreID != "" {
			query = query.Where("store_id = ?", filter.StoreID)
		}
		if filter.MinQuantity > 0 {
			query = query.Where("quantity >= ?", filter.MinQuantity)
		}
		if filter.MaxQuantity > 0 {
			query = query.Where("quantity <= ?", filter.MaxQuantity)
		}
		if filter.BatchNumber != "" {
			query = query.Where("batch_number = ?", filter.BatchNumber)
		}
		if filter.LotNumber != "" {
			query = query.Where("lot_number = ?", filter.LotNumber)
		}
		if filter.ZoneCode != "" {
			query = query.Where("zone_code = ?", filter.ZoneCode)
		}
		if filter.BinLocation != "" {
			query = query.Where("bin_location = ?", filter.BinLocation)
		}
		if filter.ShelfNumber != "" {
			query = query.Where("shelf_number = ?", filter.ShelfNumber)
		}
		if filter.ExpiryDateFrom.IsZero() == false {
			query = query.Where("expiry_date >= ?", filter.ExpiryDateFrom)
		}
		if filter.ExpiryDateTo.IsZero() == false {
			query = query.Where("expiry_date <= ?", filter.ExpiryDateTo)
		}
	}

	if err := query.Find(&stocks).Error; err != nil {
		return nil, err
	}

	return stocks, nil
}

// CreateStockHistory creates a new stock history record
func (r *StocksRepository) CreateStockHistory(ctx context.Context, history *entity.StockHistory) error {
	if history.ID == "" {
		history.ID = uuid.New().String()
	}
	return r.db.WithContext(ctx).Create(history).Error
}

// CreateOrUpdateStock creates or updates a stock record
func (r *StocksRepository) CreateOrUpdateStock(ctx context.Context, stock *entity.Stock) error {
	var existing entity.Stock
	result := r.db.WithContext(ctx).
		Where("sku_id = ? AND store_id = ?", stock.SKUID, stock.StoreID).
		First(&existing)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Create new stock record
			if stock.ID == "" {
				stock.ID = uuid.New().String()
			}
			return r.db.WithContext(ctx).Create(stock).Error
		}
		return result.Error
	}

	// Update existing stock record
	stock.ID = existing.ID
	return r.db.WithContext(ctx).Save(stock).Error
}

// ProcessStockEntry handles a stock entry and updates inventory accordingly
func (r *StocksRepository) ProcessStockEntry(ctx context.Context, entry *entity.StockEntry, userID string) error {
	// Generate ID if not provided
	if entry.ID == "" {
		entry.ID = uuid.New().String()
	}

	// Start transaction
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get or create inventory record
		stock, err := r.getOrCreateStockTx(ctx, tx, entry.SKUID, entry.StoreID)
		if err != nil {
			return err
		}

		// Calculate new quantity
		previousQty := stock.Quantity
		var newQty float64

		switch entry.Type {
		case "IN":
			newQty = previousQty + entry.Quantity
		case "OUT":
			newQty = previousQty - entry.Quantity
			if newQty < 0 {
				return ErrInsufficientStock
			}
		default:
			return ErrInvalidData
		}

		// Create stock entry
		if err := tx.Create(entry).Error; err != nil {
			return err
		}

		// Update stock record
		stock.Quantity = newQty
		// Update other fields if provided in the entry
		if entry.BatchNumber != "" {
			stock.BatchNumber = entry.BatchNumber
		}
		if entry.LotNumber != "" {
			stock.LotNumber = entry.LotNumber
		}
		if !entry.ManufactureDate.IsZero() {
			stock.ManufactureDate = entry.ManufactureDate
		}
		if !entry.ExpiryDate.IsZero() {
			stock.ExpiryDate = entry.ExpiryDate
		}

		if err := tx.Save(stock).Error; err != nil {
			return err
		}

		// Create stock history record
		history := &entity.StockHistory{
			StockID:     stock.ID,
			Type:        entry.Type,
			Quantity:    entry.Quantity,
			PreviousQty: previousQty,
			NewQty:      newQty,
			Reference:   entry.ID,
			Note:        entry.Note,
			CreatedBy:   userID,
		}

		return r.createStockHistoryTx(ctx, tx, history)
	})
}

// getOrCreateStockTx gets or creates a stock record within a transaction
func (r *StocksRepository) getOrCreateStockTx(ctx context.Context, tx *gorm.DB, skuID, storeID string) (*entity.Stock, error) {
	var stock entity.Stock
	result := tx.WithContext(ctx).
		Where("sku_id = ? AND store_id = ?", skuID, storeID).
		First(&stock)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Create new stock record
			stock = entity.Stock{
				ID:       uuid.New().String(),
				SKUID:    skuID,
				StoreID:  storeID,
				Quantity: 0,
			}
			if err := tx.Create(&stock).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, result.Error
		}
	}

	return &stock, nil
}

// createStockHistoryTx creates a stock history record within a transaction
func (r *StocksRepository) createStockHistoryTx(ctx context.Context, tx *gorm.DB, history *entity.StockHistory) error {
	if history.ID == "" {
		history.ID = uuid.New().String()
	}
	return tx.WithContext(ctx).Create(history).Error
}

// GetStockLevels retrieves current stock levels for all SKUs in a store
func (r *StocksRepository) GetStockLevels(ctx context.Context, storeID string) ([]entity.Stock, error) {
	var stocks []entity.Stock
	if err := r.db.WithContext(ctx).
		Where("store_id = ?", storeID).
		Find(&stocks).Error; err != nil {
		return nil, err
	}
	return stocks, nil
}

// GetSKUsWithLowStock retrieves SKUs with stock level below threshold
func (r *StocksRepository) GetSKUsWithLowStock(ctx context.Context, threshold float64) ([]entity.Stock, error) {
	var stocks []entity.Stock
	if err := r.db.WithContext(ctx).
		Model(&entity.Stock{}).
		Where("quantity <= ?", threshold).
		Find(&stocks).Error; err != nil {
		return nil, err
	}
	return stocks, nil
}

// GetStockMovements retrieves stock movements for a specific SKU
func (r *StocksRepository) GetStockMovements(ctx context.Context, skuID string, fromDate, toDate string) ([]entity.StockHistory, error) {
	var histories []entity.StockHistory
	query := r.db.WithContext(ctx).
		Joins("JOIN stocks ON stock_history.stock_id = stocks.id").
		Where("stocks.sku_id = ?", skuID)

	if fromDate != "" {
		query = query.Where("stock_history.created_at >= ?", fromDate)
	}
	if toDate != "" {
		query = query.Where("stock_history.created_at <= ?", toDate)
	}

	if err := query.Order("stock_history.created_at DESC").
		Find(&histories).Error; err != nil {
		return nil, err
	}

	return histories, nil
}

// GetStockHistoryByStockID retrieves stock history for a specific stock
func (r *StocksRepository) GetStockHistoryByStockID(ctx context.Context, stockID string) ([]entity.StockHistory, error) {
	var histories []entity.StockHistory
	if err := r.db.WithContext(ctx).
		Where("stock_id = ?", stockID).
		Order("created_at DESC").
		Find(&histories).Error; err != nil {
		return nil, err
	}
	return histories, nil
}

// AdjustStock adjusts a stock level directly (e.g., after physical count)
func (r *StocksRepository) AdjustStock(ctx context.Context, stockID string, newQuantity float64, note string, userID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get current stock record
		var stock entity.Stock
		if err := tx.First(&stock, "id = ?", stockID).Error; err != nil {
			return err
		}

		previousQty := stock.Quantity
		diff := newQuantity - previousQty

		// Update stock quantity
		if err := tx.Model(&entity.Stock{}).
			Where("id = ?", stockID).
			Update("quantity", newQuantity).Error; err != nil {
			return err
		}

		// Create history record for the adjustment
		history := &entity.StockHistory{
			ID:          uuid.New().String(),
			StockID:     stockID,
			Type:        "ADJUST",
			Quantity:    diff,
			PreviousQty: previousQty,
			NewQty:      newQuantity,
			Note:        note,
			CreatedBy:   userID,
		}

		return tx.Create(history).Error
	})
}
