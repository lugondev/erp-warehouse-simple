package usecase

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
	"github.com/lugondev/erp-warehouse-simple/internal/infrastructure/repository"
)

var (
	ErrInvalidSKU        = errors.New("invalid SKU format")
	ErrDuplicateSKU      = errors.New("SKU already exists")
	ErrItemNotFound      = errors.New("item not found")
	ErrCategoryNotFound  = errors.New("category not found")
	ErrInvalidPriceRange = errors.New("invalid price range")
)

type ItemUseCase struct {
	repo *repository.ItemRepository
}

func NewItemUseCase(repo *repository.ItemRepository) *ItemUseCase {
	return &ItemUseCase{repo: repo}
}

// validateItem validates item data
func (u *ItemUseCase) validateItem(ctx context.Context, item *entity.Item) error {
	// Validate SKU format (alphanumeric with optional hyphens)
	skuRegex := regexp.MustCompile(`^[A-Za-z0-9\-]+$`)
	if !skuRegex.MatchString(item.SKU) {
		return ErrInvalidSKU
	}

	// Check for duplicate SKU on create (when ID is empty)
	if item.ID == "" {
		existingItem, err := u.repo.GetItemBySKU(ctx, item.SKU)
		if err == nil && existingItem != nil {
			return ErrDuplicateSKU
		}
	}

	// Validate price is non-negative
	if item.Price < 0 {
		return ErrInvalidPriceRange
	}

	return nil
}

// CreateItem creates a new item
func (u *ItemUseCase) CreateItem(ctx context.Context, item *entity.Item) error {
	if err := u.validateItem(ctx, item); err != nil {
		return err
	}
	return u.repo.CreateItem(ctx, item)
}

// UpdateItem updates an existing item
func (u *ItemUseCase) UpdateItem(ctx context.Context, item *entity.Item) error {
	// Check if item exists
	existingItem, err := u.repo.GetItemByID(ctx, item.ID)
	if err != nil {
		return ErrItemNotFound
	}

	// If SKU is being changed, validate the new SKU
	if existingItem.SKU != item.SKU {
		if err := u.validateItem(ctx, item); err != nil {
			return err
		}
	} else {
		// Still validate other fields
		if item.Price < 0 {
			return ErrInvalidPriceRange
		}
	}

	return u.repo.UpdateItem(ctx, item)
}

// GetItem gets an item by ID
func (u *ItemUseCase) GetItem(ctx context.Context, id string) (*entity.Item, error) {
	item, err := u.repo.GetItemByID(ctx, id)
	if err != nil {
		return nil, ErrItemNotFound
	}
	return item, nil
}

// GetItemBySKU gets an item by SKU
func (u *ItemUseCase) GetItemBySKU(ctx context.Context, sku string) (*entity.Item, error) {
	item, err := u.repo.GetItemBySKU(ctx, sku)
	if err != nil {
		return nil, ErrItemNotFound
	}
	return item, nil
}

// DeleteItem deletes an item by ID
func (u *ItemUseCase) DeleteItem(ctx context.Context, id string) error {
	// Check if item exists
	_, err := u.repo.GetItemByID(ctx, id)
	if err != nil {
		return ErrItemNotFound
	}
	return u.repo.DeleteItem(ctx, id)
}

// ListItems lists items with filters
func (u *ItemUseCase) ListItems(ctx context.Context, filter *entity.ItemFilter, page, pageSize int) ([]entity.Item, int64, error) {
	// Validate page and pageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Validate price range if provided
	if filter != nil && filter.MinPrice != nil && filter.MaxPrice != nil {
		if *filter.MinPrice > *filter.MaxPrice {
			return nil, 0, ErrInvalidPriceRange
		}
	}

	return u.repo.ListItems(ctx, filter, page, pageSize)
}

// SearchItems searches for items based on a search term
func (u *ItemUseCase) SearchItems(ctx context.Context, searchTerm string, page, pageSize int) ([]entity.Item, int64, error) {
	// Validate page and pageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return u.repo.SearchItems(ctx, searchTerm, page, pageSize)
}

// CreateItemCategory creates a new item category
func (u *ItemUseCase) CreateItemCategory(ctx context.Context, category *entity.ItemCategory) error {
	// If parent ID is provided, check if parent exists
	if category.ParentID != nil && *category.ParentID != "" {
		_, err := u.repo.GetItemCategoryByID(ctx, *category.ParentID)
		if err != nil {
			return ErrCategoryNotFound
		}
	}

	return u.repo.CreateItemCategory(ctx, category)
}

// UpdateItemCategory updates an existing item category
func (u *ItemUseCase) UpdateItemCategory(ctx context.Context, category *entity.ItemCategory) error {
	// Check if category exists
	_, err := u.repo.GetItemCategoryByID(ctx, category.ID)
	if err != nil {
		return ErrCategoryNotFound
	}

	// If parent ID is provided, check if parent exists
	if category.ParentID != nil && *category.ParentID != "" {
		// Prevent circular reference
		if *category.ParentID == category.ID {
			return errors.New("category cannot be its own parent")
		}

		_, err := u.repo.GetItemCategoryByID(ctx, *category.ParentID)
		if err != nil {
			return ErrCategoryNotFound
		}
	}

	return u.repo.UpdateItemCategory(ctx, category)
}

// GetItemCategory gets an item category by ID
func (u *ItemUseCase) GetItemCategory(ctx context.Context, id string) (*entity.ItemCategory, error) {
	category, err := u.repo.GetItemCategoryByID(ctx, id)
	if err != nil {
		return nil, ErrCategoryNotFound
	}
	return category, nil
}

// DeleteItemCategory deletes an item category by ID
func (u *ItemUseCase) DeleteItemCategory(ctx context.Context, id string) error {
	// Check if category exists
	_, err := u.repo.GetItemCategoryByID(ctx, id)
	if err != nil {
		return ErrCategoryNotFound
	}

	// Repository will check if category has children or is used by items
	return u.repo.DeleteItemCategory(ctx, id)
}

// ListItemCategories lists all item categories
func (u *ItemUseCase) ListItemCategories(ctx context.Context) ([]entity.ItemCategory, error) {
	return u.repo.ListItemCategories(ctx)
}

// GetItemCategoriesTree gets item categories in a hierarchical structure
func (u *ItemUseCase) GetItemCategoriesTree(ctx context.Context) ([]entity.ItemCategory, error) {
	return u.repo.GetItemCategoriesTree(ctx)
}

// GetItemsByCategory gets items by category
func (u *ItemUseCase) GetItemsByCategory(ctx context.Context, categoryID string, page, pageSize int) ([]entity.Item, int64, error) {
	// Check if category exists
	_, err := u.repo.GetItemCategoryByID(ctx, categoryID)
	if err != nil {
		return nil, 0, ErrCategoryNotFound
	}

	// Validate page and pageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return u.repo.GetItemsByCategory(ctx, categoryID, page, pageSize)
}

// BulkCreateItems creates multiple items in a single transaction
func (u *ItemUseCase) BulkCreateItems(ctx context.Context, items []*entity.Item) error {
	// Validate all items first
	for i, item := range items {
		if err := u.validateItem(ctx, item); err != nil {
			return fmt.Errorf("validation failed for item at index %d: %w", i, err)
		}
	}

	return u.repo.BulkCreateItems(ctx, items)
}

// BulkUpdateItems updates multiple items in a single transaction
func (u *ItemUseCase) BulkUpdateItems(ctx context.Context, items []*entity.Item) error {
	// Validate all items first
	for i, item := range items {
		// Check if item exists
		_, err := u.repo.GetItemByID(ctx, item.ID)
		if err != nil {
			return fmt.Errorf("item at index %d not found: %s", i, item.ID)
		}

		// Validate item data
		if item.Price < 0 {
			return fmt.Errorf("invalid price for item at index %d: %s", i, item.ID)
		}
	}

	return u.repo.BulkUpdateItems(ctx, items)
}

// GetItemsByIDs gets items by their IDs
func (u *ItemUseCase) GetItemsByIDs(ctx context.Context, ids []string) ([]entity.Item, error) {
	return u.repo.GetItemsByIDs(ctx, ids)
}

// GetItemsBySKUs gets items by their SKUs
func (u *ItemUseCase) GetItemsBySKUs(ctx context.Context, skus []string) ([]entity.Item, error) {
	return u.repo.GetItemsBySKUs(ctx, skus)
}
