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
	ErrInvalidSKUCode    = errors.New("invalid SKU code format")
	ErrDuplicateSKUCode  = errors.New("SKU code already exists")
	ErrSKUNotFound       = errors.New("SKU not found")
	ErrCategoryNotFound  = errors.New("category not found")
	ErrInvalidPriceRange = errors.New("invalid price range")
)

type SKUUseCase struct {
	repo *repository.SKURepository
}

func NewSKUUseCase(repo *repository.SKURepository) *SKUUseCase {
	return &SKUUseCase{repo: repo}
}

// validateSKU validates SKU data
func (u *SKUUseCase) validateSKU(ctx context.Context, sku *entity.SKU) error {
	// Validate SKU code format (alphanumeric with optional hyphens)
	skuRegex := regexp.MustCompile(`^[A-Za-z0-9\-]+$`)
	if !skuRegex.MatchString(sku.SKUCode) {
		return ErrInvalidSKUCode
	}

	// Check for duplicate SKU code on create (when ID is empty)
	if sku.ID == "" {
		existingSKU, err := u.repo.GetSKUBySKUCode(ctx, sku.SKUCode)
		if err == nil && existingSKU != nil {
			return ErrDuplicateSKUCode
		}
	}

	// Validate price is non-negative
	if sku.Price < 0 {
		return ErrInvalidPriceRange
	}

	return nil
}

// CreateSKU creates a new SKU
func (u *SKUUseCase) CreateSKU(ctx context.Context, sku *entity.SKU) error {
	if err := u.validateSKU(ctx, sku); err != nil {
		return err
	}
	return u.repo.CreateSKU(ctx, sku)
}

// UpdateSKU updates an existing SKU
func (u *SKUUseCase) UpdateSKU(ctx context.Context, sku *entity.SKU) error {
	// Check if SKU exists
	existingSKU, err := u.repo.GetSKUByID(ctx, sku.ID)
	if err != nil {
		return ErrSKUNotFound
	}

	// If SKU code is being changed, validate the new SKU code
	if existingSKU.SKUCode != sku.SKUCode {
		if err := u.validateSKU(ctx, sku); err != nil {
			return err
		}
	} else {
		// Still validate other fields
		if sku.Price < 0 {
			return ErrInvalidPriceRange
		}
	}

	return u.repo.UpdateSKU(ctx, sku)
}

// GetSKU gets a SKU by ID
func (u *SKUUseCase) GetSKU(ctx context.Context, id string) (*entity.SKU, error) {
	sku, err := u.repo.GetSKUByID(ctx, id)
	if err != nil {
		return nil, ErrSKUNotFound
	}
	return sku, nil
}

// GetSKUBySKUCode gets a SKU by SKU code
func (u *SKUUseCase) GetSKUBySKUCode(ctx context.Context, skuCode string) (*entity.SKU, error) {
	sku, err := u.repo.GetSKUBySKUCode(ctx, skuCode)
	if err != nil {
		return nil, ErrSKUNotFound
	}
	return sku, nil
}

// DeleteSKU deletes a SKU by ID
func (u *SKUUseCase) DeleteSKU(ctx context.Context, id string) error {
	// Check if SKU exists
	_, err := u.repo.GetSKUByID(ctx, id)
	if err != nil {
		return ErrSKUNotFound
	}
	return u.repo.DeleteSKU(ctx, id)
}

// ListSKUs lists SKUs with filters
func (u *SKUUseCase) ListSKUs(ctx context.Context, filter *entity.SKUFilter, page, pageSize int) ([]entity.SKU, int64, error) {
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

	return u.repo.ListSKUs(ctx, filter, page, pageSize)
}

// SearchSKUs searches for SKUs based on a search term
func (u *SKUUseCase) SearchSKUs(ctx context.Context, searchTerm string, page, pageSize int) ([]entity.SKU, int64, error) {
	// Validate page and pageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return u.repo.SearchSKUs(ctx, searchTerm, page, pageSize)
}

// CreateSKUCategory creates a new SKU category
func (u *SKUUseCase) CreateSKUCategory(ctx context.Context, category *entity.SKUCategory) error {
	// If parent ID is provided, check if parent exists
	if category.ParentID != nil && *category.ParentID != "" {
		_, err := u.repo.GetSKUCategoryByID(ctx, *category.ParentID)
		if err != nil {
			return ErrCategoryNotFound
		}
	}

	return u.repo.CreateSKUCategory(ctx, category)
}

// UpdateSKUCategory updates an existing SKU category
func (u *SKUUseCase) UpdateSKUCategory(ctx context.Context, category *entity.SKUCategory) error {
	// Check if category exists
	_, err := u.repo.GetSKUCategoryByID(ctx, category.ID)
	if err != nil {
		return ErrCategoryNotFound
	}

	// If parent ID is provided, check if parent exists
	if category.ParentID != nil && *category.ParentID != "" {
		// Prevent circular reference
		if *category.ParentID == category.ID {
			return errors.New("category cannot be its own parent")
		}

		_, err := u.repo.GetSKUCategoryByID(ctx, *category.ParentID)
		if err != nil {
			return ErrCategoryNotFound
		}
	}

	return u.repo.UpdateSKUCategory(ctx, category)
}

// GetSKUCategory gets a SKU category by ID
func (u *SKUUseCase) GetSKUCategory(ctx context.Context, id string) (*entity.SKUCategory, error) {
	category, err := u.repo.GetSKUCategoryByID(ctx, id)
	if err != nil {
		return nil, ErrCategoryNotFound
	}
	return category, nil
}

// DeleteSKUCategory deletes a SKU category by ID
func (u *SKUUseCase) DeleteSKUCategory(ctx context.Context, id string) error {
	// Check if category exists
	_, err := u.repo.GetSKUCategoryByID(ctx, id)
	if err != nil {
		return ErrCategoryNotFound
	}

	// Repository will check if category has children or is used by SKUs
	return u.repo.DeleteSKUCategory(ctx, id)
}

// ListSKUCategories lists all SKU categories
func (u *SKUUseCase) ListSKUCategories(ctx context.Context) ([]entity.SKUCategory, error) {
	return u.repo.ListSKUCategories(ctx)
}

// GetSKUCategoriesTree gets SKU categories in a hierarchical structure
func (u *SKUUseCase) GetSKUCategoriesTree(ctx context.Context) ([]entity.SKUCategory, error) {
	return u.repo.GetSKUCategoriesTree(ctx)
}

// GetSKUsByCategory gets SKUs by category
func (u *SKUUseCase) GetSKUsByCategory(ctx context.Context, categoryID string, page, pageSize int) ([]entity.SKU, int64, error) {
	// Check if category exists
	_, err := u.repo.GetSKUCategoryByID(ctx, categoryID)
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

	return u.repo.GetSKUsByCategory(ctx, categoryID, page, pageSize)
}

// BulkCreateSKUs creates multiple SKUs in a single transaction
func (u *SKUUseCase) BulkCreateSKUs(ctx context.Context, skus []*entity.SKU) error {
	// Validate all SKUs first
	for i, sku := range skus {
		if err := u.validateSKU(ctx, sku); err != nil {
			return fmt.Errorf("validation failed for SKU at index %d: %w", i, err)
		}
	}

	return u.repo.BulkCreateSKUs(ctx, skus)
}

// BulkUpdateSKUs updates multiple SKUs in a single transaction
func (u *SKUUseCase) BulkUpdateSKUs(ctx context.Context, skus []*entity.SKU) error {
	// Validate all SKUs first
	for i, sku := range skus {
		// Check if SKU exists
		_, err := u.repo.GetSKUByID(ctx, sku.ID)
		if err != nil {
			return fmt.Errorf("SKU at index %d not found: %s", i, sku.ID)
		}

		// Validate SKU data
		if sku.Price < 0 {
			return fmt.Errorf("invalid price for SKU at index %d: %s", i, sku.ID)
		}
	}

	return u.repo.BulkUpdateSKUs(ctx, skus)
}

// GetSKUsByIDs gets SKUs by their IDs
func (u *SKUUseCase) GetSKUsByIDs(ctx context.Context, ids []string) ([]entity.SKU, error) {
	return u.repo.GetSKUsByIDs(ctx, ids)
}

// GetSKUsBySKUCodes gets SKUs by their SKU codes
func (u *SKUUseCase) GetSKUsBySKUCodes(ctx context.Context, skuCodes []string) ([]entity.SKU, error) {
	return u.repo.GetSKUsBySKUCodes(ctx, skuCodes)
}
