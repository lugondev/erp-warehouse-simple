package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
	"gorm.io/gorm"
)

type InventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) GetByID(ctx context.Context, id string) (*entity.Inventory, error) {
	var inventory entity.Inventory
	if err := r.db.WithContext(ctx).First(&inventory, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &inventory, nil
}

func (r *InventoryRepository) GetByProductAndWarehouse(ctx context.Context, productID, warehouseID string) (*entity.Inventory, error) {
	var inventory entity.Inventory
	if err := r.db.WithContext(ctx).First(&inventory, "product_id = ? AND warehouse_id = ?", productID, warehouseID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &inventory, nil
}

func (r *InventoryRepository) List(ctx context.Context, filter *entity.InventoryFilter) ([]entity.Inventory, error) {
	var inventories []entity.Inventory
	query := r.db.WithContext(ctx)

	if filter != nil {
		if filter.WarehouseID != "" {
			query = query.Where("warehouse_id = ?", filter.WarehouseID)
		}
		if filter.ProductID != "" {
			query = query.Where("product_id = ?", filter.ProductID)
		}
		if filter.BatchNumber != "" {
			query = query.Where("batch_number = ?", filter.BatchNumber)
		}
		if filter.LotNumber != "" {
			query = query.Where("lot_number = ?", filter.LotNumber)
		}
	}

	if err := query.Find(&inventories).Error; err != nil {
		return nil, err
	}
	return inventories, nil
}

func (r *InventoryRepository) UpdateQuantity(ctx context.Context, id string, quantity float64) error {
	return r.db.WithContext(ctx).
		Model(&entity.Inventory{}).
		Where("id = ?", id).
		UpdateColumn("quantity", quantity).
		Error
}

func (r *InventoryRepository) CreateStockEntry(ctx context.Context, entry *entity.StockEntry) error {
	if entry.ID == "" {
		entry.ID = uuid.New().String()
	}
	return r.db.WithContext(ctx).Create(entry).Error
}

func (r *InventoryRepository) CreateInventoryHistory(ctx context.Context, history *entity.InventoryHistory) error {
	if history.ID == "" {
		history.ID = uuid.New().String()
	}
	return r.db.WithContext(ctx).Create(history).Error
}

func (r *InventoryRepository) CreateOrUpdateInventory(ctx context.Context, inventory *entity.Inventory) error {
	tx := r.db.WithContext(ctx)

	var existing entity.Inventory
	err := tx.First(&existing, "product_id = ? AND warehouse_id = ?", inventory.ProductID, inventory.WarehouseID).Error

	if err == gorm.ErrRecordNotFound {
		if inventory.ID == "" {
			inventory.ID = uuid.New().String()
		}
		return tx.Create(inventory).Error
	}
	if err != nil {
		return err
	}

	// Update existing inventory
	existing.Quantity = inventory.Quantity
	existing.BinLocation = inventory.BinLocation
	existing.ShelfNumber = inventory.ShelfNumber
	existing.ZoneCode = inventory.ZoneCode
	existing.BatchNumber = inventory.BatchNumber
	existing.LotNumber = inventory.LotNumber
	existing.ManufactureDate = inventory.ManufactureDate
	existing.ExpiryDate = inventory.ExpiryDate

	return tx.Save(&existing).Error
}

func (r *InventoryRepository) ProcessStockEntry(ctx context.Context, entry *entity.StockEntry, createdBy string) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create stock entry
	entry.CreatedBy = createdBy
	if err := r.createStockEntryTx(ctx, tx, entry); err != nil {
		tx.Rollback()
		return err
	}

	// Get or create inventory
	inventory, err := r.getOrCreateInventoryTx(ctx, tx, entry.ProductID, entry.WarehouseID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Calculate new quantity
	newQty := inventory.Quantity
	if entry.Type == "IN" {
		newQty += entry.Quantity
	} else if entry.Type == "OUT" {
		newQty -= entry.Quantity
		if newQty < 0 {
			tx.Rollback()
			return ErrInsufficientStock
		}
	}

	// Create history record
	history := &entity.InventoryHistory{
		InventoryID: inventory.ID,
		Type:        entry.Type,
		Quantity:    entry.Quantity,
		PreviousQty: inventory.Quantity,
		NewQty:      newQty,
		Reference:   entry.ID,
		Note:        entry.Note,
		CreatedBy:   createdBy,
	}
	if err := r.createInventoryHistoryTx(ctx, tx, history); err != nil {
		tx.Rollback()
		return err
	}

	// Update inventory quantity
	if err := r.updateInventoryQuantityTx(ctx, tx, inventory.ID, newQty); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *InventoryRepository) createStockEntryTx(ctx context.Context, tx *gorm.DB, entry *entity.StockEntry) error {
	if entry.ID == "" {
		entry.ID = uuid.New().String()
	}
	return tx.WithContext(ctx).Create(entry).Error
}

func (r *InventoryRepository) getOrCreateInventoryTx(ctx context.Context, tx *gorm.DB, productID, warehouseID string) (*entity.Inventory, error) {
	var inventory entity.Inventory
	err := tx.WithContext(ctx).First(&inventory, "product_id = ? AND warehouse_id = ?", productID, warehouseID).Error

	if err == gorm.ErrRecordNotFound {
		inventory = entity.Inventory{
			ID:          uuid.New().String(),
			ProductID:   productID,
			WarehouseID: warehouseID,
			Quantity:    0,
		}
		if err := tx.WithContext(ctx).Create(&inventory).Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return &inventory, nil
}

func (r *InventoryRepository) createInventoryHistoryTx(ctx context.Context, tx *gorm.DB, history *entity.InventoryHistory) error {
	if history.ID == "" {
		history.ID = uuid.New().String()
	}
	return tx.WithContext(ctx).Create(history).Error
}

func (r *InventoryRepository) updateInventoryQuantityTx(ctx context.Context, tx *gorm.DB, id string, quantity float64) error {
	return tx.WithContext(ctx).
		Model(&entity.Inventory{}).
		Where("id = ?", id).
		UpdateColumn("quantity", quantity).
		Error
}
