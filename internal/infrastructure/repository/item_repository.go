package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
	"gorm.io/gorm"
)

type ItemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

// CreateItem creates a new item
func (r *ItemRepository) CreateItem(ctx context.Context, item *entity.Item) error {
	if item.ID == "" {
		item.ID = uuid.New().String()
	}
	return r.db.WithContext(ctx).Create(item).Error
}

// UpdateItem updates an existing item
func (r *ItemRepository) UpdateItem(ctx context.Context, item *entity.Item) error {
	return r.db.WithContext(ctx).Save(item).Error
}

// GetItemByID retrieves an item by ID
func (r *ItemRepository) GetItemByID(ctx context.Context, id string) (*entity.Item, error) {
	var item entity.Item
	if err := r.db.WithContext(ctx).
		Preload("Manufacturer").
		Preload("Supplier").
		First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// GetItemBySKU retrieves an item by SKU
func (r *ItemRepository) GetItemBySKU(ctx context.Context, sku string) (*entity.Item, error) {
	var item entity.Item
	if err := r.db.WithContext(ctx).
		Preload("Manufacturer").
		Preload("Supplier").
		First(&item, "sku = ?", sku).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// DeleteItem deletes an item by ID
func (r *ItemRepository) DeleteItem(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.Item{}, "id = ?", id).Error
}

// ListItems retrieves items with filters
func (r *ItemRepository) ListItems(ctx context.Context, filter *entity.ItemFilter, page, pageSize int) ([]entity.Item, int64, error) {
	var items []entity.Item
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Item{})

	// Apply filters
	if filter != nil {
		if filter.SKU != "" {
			query = query.Where("sku LIKE ?", "%"+filter.SKU+"%")
		}
		if filter.Name != "" {
			query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
		}
		if filter.Category != "" {
			query = query.Where("category = ?", filter.Category)
		}
		if filter.ManufacturerID != nil {
			query = query.Where("manufacturer_id = ?", filter.ManufacturerID)
		}
		if filter.SupplierID != nil {
			query = query.Where("supplier_id = ?", filter.SupplierID)
		}
		if filter.Status != nil {
			query = query.Where("status = ?", filter.Status)
		}
		if filter.MinPrice != nil {
			query = query.Where("price >= ?", filter.MinPrice)
		}
		if filter.MaxPrice != nil {
			query = query.Where("price <= ?", filter.MaxPrice)
		}
	}

	// Count total items
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).
		Preload("Manufacturer").
		Preload("Supplier").
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// SearchItems searches for items based on a search term
func (r *ItemRepository) SearchItems(ctx context.Context, searchTerm string, page, pageSize int) ([]entity.Item, int64, error) {
	var items []entity.Item
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Item{}).
		Where("sku ILIKE ? OR name ILIKE ? OR description ILIKE ?",
			"%"+searchTerm+"%", "%"+searchTerm+"%", "%"+searchTerm+"%")

	// Count total items
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).
		Preload("Manufacturer").
		Preload("Supplier").
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// CreateItemCategory creates a new item category
func (r *ItemRepository) CreateItemCategory(ctx context.Context, category *entity.ItemCategory) error {
	if category.ID == "" {
		category.ID = uuid.New().String()
	}
	return r.db.WithContext(ctx).Create(category).Error
}

// UpdateItemCategory updates an existing item category
func (r *ItemRepository) UpdateItemCategory(ctx context.Context, category *entity.ItemCategory) error {
	return r.db.WithContext(ctx).Save(category).Error
}

// GetItemCategoryByID retrieves an item category by ID
func (r *ItemRepository) GetItemCategoryByID(ctx context.Context, id string) (*entity.ItemCategory, error) {
	var category entity.ItemCategory
	if err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		First(&category, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// DeleteItemCategory deletes an item category by ID
func (r *ItemRepository) DeleteItemCategory(ctx context.Context, id string) error {
	// Check if category has children
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.ItemCategory{}).
		Where("parent_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("cannot delete category with children")
	}

	// Check if category is used by items
	if err := r.db.WithContext(ctx).Model(&entity.Item{}).
		Where("category = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("cannot delete category used by items")
	}

	return r.db.WithContext(ctx).Delete(&entity.ItemCategory{}, "id = ?", id).Error
}

// ListItemCategories retrieves all item categories
func (r *ItemRepository) ListItemCategories(ctx context.Context) ([]entity.ItemCategory, error) {
	var categories []entity.ItemCategory
	if err := r.db.WithContext(ctx).
		Preload("Parent").
		Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// GetItemCategoriesTree retrieves item categories in a hierarchical structure
func (r *ItemRepository) GetItemCategoriesTree(ctx context.Context) ([]entity.ItemCategory, error) {
	var rootCategories []entity.ItemCategory
	if err := r.db.WithContext(ctx).
		Preload("Children.Children").
		Where("parent_id IS NULL").
		Find(&rootCategories).Error; err != nil {
		return nil, err
	}
	return rootCategories, nil
}

// GetItemsByCategory retrieves items by category
func (r *ItemRepository) GetItemsByCategory(ctx context.Context, categoryID string, page, pageSize int) ([]entity.Item, int64, error) {
	var items []entity.Item
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Item{}).
		Where("category = ?", categoryID)

	// Count total items
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).
		Preload("Manufacturer").
		Preload("Supplier").
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// BulkCreateItems creates multiple items in a single transaction
func (r *ItemRepository) BulkCreateItems(ctx context.Context, items []*entity.Item) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if item.ID == "" {
				item.ID = uuid.New().String()
			}
			if err := tx.Create(item).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// BulkUpdateItems updates multiple items in a single transaction
func (r *ItemRepository) BulkUpdateItems(ctx context.Context, items []*entity.Item) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if err := tx.Save(item).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetItemsByIDs retrieves items by their IDs
func (r *ItemRepository) GetItemsByIDs(ctx context.Context, ids []string) ([]entity.Item, error) {
	var items []entity.Item
	if err := r.db.WithContext(ctx).
		Preload("Manufacturer").
		Preload("Supplier").
		Where("id IN ?", ids).
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// GetItemsBySKUs retrieves items by their SKUs
func (r *ItemRepository) GetItemsBySKUs(ctx context.Context, skus []string) ([]entity.Item, error) {
	var items []entity.Item
	if err := r.db.WithContext(ctx).
		Preload("Manufacturer").
		Preload("Supplier").
		Where("sku IN ?", skus).
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
