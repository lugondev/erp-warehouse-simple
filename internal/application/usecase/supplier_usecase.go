package usecase

import (
	"context"

	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
	"github.com/lugondev/erp-warehouse-simple/internal/infrastructure/repository"
)

type SupplierUseCase struct {
	repo *repository.SupplierRepository
}

func NewSupplierUseCase(repo *repository.SupplierRepository) *SupplierUseCase {
	return &SupplierUseCase{repo: repo}
}

// CreateSupplier creates a new supplier
func (u *SupplierUseCase) CreateSupplier(ctx context.Context, supplier *entity.Supplier) error {
	return u.repo.CreateSupplier(ctx, supplier)
}

// UpdateSupplier updates an existing supplier
func (u *SupplierUseCase) UpdateSupplier(ctx context.Context, supplier *entity.Supplier) error {
	return u.repo.UpdateSupplier(ctx, supplier)
}

// GetSupplier gets a supplier by ID
func (u *SupplierUseCase) GetSupplier(ctx context.Context, id uint) (*entity.Supplier, error) {
	return u.repo.GetSupplierByID(ctx, id)
}

// DeleteSupplier deletes a supplier
func (u *SupplierUseCase) DeleteSupplier(ctx context.Context, id uint) error {
	return u.repo.DeleteSupplier(ctx, id)
}

// ListSuppliers lists suppliers with filters
func (u *SupplierUseCase) ListSuppliers(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]entity.Supplier, int64, error) {
	return u.repo.ListSuppliers(ctx, filter, page, pageSize)
}

// CreateProduct creates a new product
func (u *SupplierUseCase) CreateProduct(ctx context.Context, product *entity.Product) error {
	return u.repo.CreateProduct(ctx, product)
}

// AddProductToSupplier associates a product with a supplier
func (u *SupplierUseCase) AddProductToSupplier(ctx context.Context, supplierID, productID uint) error {
	// Verify supplier exists
	if _, err := u.repo.GetSupplierByID(ctx, supplierID); err != nil {
		return err
	}
	return u.repo.AddProductToSupplier(ctx, supplierID, productID)
}

// RemoveProductFromSupplier removes a product from a supplier
func (u *SupplierUseCase) RemoveProductFromSupplier(ctx context.Context, supplierID, productID uint) error {
	return u.repo.RemoveProductFromSupplier(ctx, supplierID, productID)
}

// CreateContract creates a new contract for a supplier
func (u *SupplierUseCase) CreateContract(ctx context.Context, contract *entity.Contract) error {
	// Verify supplier exists
	if _, err := u.repo.GetSupplierByID(ctx, contract.SupplierID); err != nil {
		return err
	}
	return u.repo.CreateContract(ctx, contract)
}

// UpdateContract updates an existing contract
func (u *SupplierUseCase) UpdateContract(ctx context.Context, contract *entity.Contract) error {
	return u.repo.UpdateContract(ctx, contract)
}

// GetContract gets a contract by ID
func (u *SupplierUseCase) GetContract(ctx context.Context, id uint) (*entity.Contract, error) {
	return u.repo.GetContractByID(ctx, id)
}

// AddSupplierRating adds a rating for a supplier
func (u *SupplierUseCase) AddSupplierRating(ctx context.Context, rating *entity.SupplierRating) error {
	// Verify supplier exists
	if _, err := u.repo.GetSupplierByID(ctx, rating.SupplierID); err != nil {
		return err
	}

	if rating.Score < 0 || rating.Score > 5 {
		return entity.ErrInvalidRating
	}

	return u.repo.AddSupplierRating(ctx, rating)
}

// GetSupplierRatings gets all ratings for a supplier
func (u *SupplierUseCase) GetSupplierRatings(ctx context.Context, supplierID uint) ([]entity.SupplierRating, error) {
	// Verify supplier exists
	if _, err := u.repo.GetSupplierByID(ctx, supplierID); err != nil {
		return nil, err
	}

	return u.repo.GetSupplierRatings(ctx, supplierID)
}
