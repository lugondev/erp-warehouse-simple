package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
	"gorm.io/gorm"
)

type SKURepository struct {
	db *gorm.DB
}

func NewSKURepository(db *gorm.DB) *SKURepository {
	return &SKURepository{db: db}
}

// CreateSKU creates a new SKU
func (r *SKURepository) CreateSKU(ctx context.Context, sku *entity.SKU) error {
	if sku.ID == "" {
		sku.ID = uuid.New().String()
	}
	return r.db.WithContext(ctx).Create(sku).Error
}

// UpdateSKU updates an existing SKU
func (r *SKURepository) UpdateSKU(ctx context.Context, sku *entity.SKU) error {
	return r.db.WithContext(ctx).Save(sku).Error
}

// GetSKUByID retrieves a SKU by ID
func (r *SKURepository) GetSKUByID(ctx context.Context, id string) (*entity.SKU, error) {
	var sku entity.SKU
	if err := r.db.WithContext(ctx).
		Preload("Manufacturer").
		Preload("Vendor").
		First(&sku, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &sku, nil
}

// GetSKUBySKUCode retrieves a SKU by SKU code
func (r *SKURepository) GetSKUBySKUCode(ctx context.Context, skuCode string) (*entity.SKU, error) {
	var sku entity.SKU
	if err := r.db.WithContext(ctx).
		Preload("Manufacturer").
		Preload("Vendor").
		First(&sku, "sku_code = ?", skuCode).Error; err != nil {
		return nil, err
	}
	return &sku, nil
}

// DeleteSKU deletes a SKU by ID
func (r *SKURepository) DeleteSKU(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.SKU{}, "id = ?", id).Error
}

// ListSKUs retrieves SKUs with filters
func (r *SKURepository) ListSKUs(ctx context.Context, filter *entity.SKUFilter, page, pageSize int) ([]entity.SKU, int64, error) {
	var skus []entity.SKU
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.SKU{})

	// Apply filters
	if filter != nil {
		if filter.SKUCode != "" {
			query = query.Where("sku_code LIKE ?", "%"+filter.SKUCode+"%")
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
		if filter.VendorID != nil {
			query = query.Where("vendor_id = ?", filter.VendorID)
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

	// Count total SKUs
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).
		Preload("Manufacturer").
		Preload("Vendor").
		Find(&skus).Error; err != nil {
		return nil, 0, err
	}

	return skus, total, nil
}

// SearchSKUs searches for SKUs based on a search term
func (r *SKURepository) SearchSKUs(ctx context.Context, searchTerm string, page, pageSize int) ([]entity.SKU, int64, error) {
	var skus []entity.SKU
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.SKU{}).
		Where("sku_code ILIKE ? OR name ILIKE ? OR description ILIKE ?",
			"%"+searchTerm+"%", "%"+searchTerm+"%", "%"+searchTerm+"%")

	// Count total skus
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).
		Preload("Manufacturer").
		Preload("Vendor").
		Find(&skus).Error; err != nil {
		return nil, 0, err
	}

	return skus, total, nil
}

// CreateSKUCategory creates a new SKU category
func (r *SKURepository) CreateSKUCategory(ctx context.Context, category *entity.SKUCategory) error {
	if category.ID == "" {
		category.ID = uuid.New().String()
	}
	return r.db.WithContext(ctx).Create(category).Error
}

// UpdateSKUCategory updates an existing SKU category
func (r *SKURepository) UpdateSKUCategory(ctx context.Context, category *entity.SKUCategory) error {
	return r.db.WithContext(ctx).Save(category).Error
}

// GetSKUCategoryByID retrieves a SKU category by ID
func (r *SKURepository) GetSKUCategoryByID(ctx context.Context, id string) (*entity.SKUCategory, error) {
	var category entity.SKUCategory
	if err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		First(&category, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// DeleteSKUCategory deletes a SKU category by ID
func (r *SKURepository) DeleteSKUCategory(ctx context.Context, id string) error {
	// Check if category has children
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.SKUCategory{}).
		Where("parent_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("cannot delete category with children")
	}

	// Check if category is used by SKUs
	if err := r.db.WithContext(ctx).Model(&entity.SKU{}).
		Where("category = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("cannot delete category used by SKUs")
	}

	return r.db.WithContext(ctx).Delete(&entity.SKUCategory{}, "id = ?", id).Error
}

// ListSKUCategories retrieves all SKU categories
func (r *SKURepository) ListSKUCategories(ctx context.Context) ([]entity.SKUCategory, error) {
	var categories []entity.SKUCategory
	if err := r.db.WithContext(ctx).
		Preload("Parent").
		Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// GetSKUCategoriesTree retrieves SKU categories in a hierarchical structure
func (r *SKURepository) GetSKUCategoriesTree(ctx context.Context) ([]entity.SKUCategory, error) {
	var rootCategories []entity.SKUCategory
	if err := r.db.WithContext(ctx).
		Preload("Children.Children").
		Where("parent_id IS NULL").
		Find(&rootCategories).Error; err != nil {
		return nil, err
	}
	return rootCategories, nil
}

// GetSKUsByCategory retrieves SKUs by category
func (r *SKURepository) GetSKUsByCategory(ctx context.Context, categoryID string, page, pageSize int) ([]entity.SKU, int64, error) {
	var skus []entity.SKU
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.SKU{}).
		Where("category = ?", categoryID)

	// Count total SKUs
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).
		Preload("Manufacturer").
		Preload("Vendor").
		Find(&skus).Error; err != nil {
		return nil, 0, err
	}

	return skus, total, nil
}

// BulkCreateSKUs creates multiple SKUs in a single transaction
func (r *SKURepository) BulkCreateSKUs(ctx context.Context, skus []*entity.SKU) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, sku := range skus {
			if sku.ID == "" {
				sku.ID = uuid.New().String()
			}
			if err := tx.Create(sku).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// BulkUpdateSKUs updates multiple SKUs in a single transaction
func (r *SKURepository) BulkUpdateSKUs(ctx context.Context, skus []*entity.SKU) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, sku := range skus {
			if err := tx.Save(sku).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetSKUsByIDs retrieves SKUs by their IDs
func (r *SKURepository) GetSKUsByIDs(ctx context.Context, ids []string) ([]entity.SKU, error) {
	var skus []entity.SKU
	if err := r.db.WithContext(ctx).
		Preload("Manufacturer").
		Preload("Vendor").
		Where("id IN ?", ids).
		Find(&skus).Error; err != nil {
		return nil, err
	}
	return skus, nil
}

// GetSKUsBySKUCodes retrieves SKUs by their SKU codes
func (r *SKURepository) GetSKUsBySKUCodes(ctx context.Context, skuCodes []string) ([]entity.SKU, error) {
	var skus []entity.SKU
	if err := r.db.WithContext(ctx).
		Preload("Manufacturer").
		Preload("Vendor").
		Where("sku_code IN ?", skuCodes).
		Find(&skus).Error; err != nil {
		return nil, err
	}
	return skus, nil
}
