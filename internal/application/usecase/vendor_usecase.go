package usecase

import (
	"context"

	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
	"github.com/lugondev/erp-warehouse-simple/internal/infrastructure/repository"
)

type VendorUseCase struct {
	repo *repository.VendorRepository
}

func NewVendorUseCase(repo *repository.VendorRepository) *VendorUseCase {
	return &VendorUseCase{repo: repo}
}

// CreateVendor creates a new vendor
func (u *VendorUseCase) CreateVendor(ctx context.Context, vendor *entity.Vendor) error {
	return u.repo.Create(ctx, vendor)
}

// UpdateVendor updates an existing vendor
func (u *VendorUseCase) UpdateVendor(ctx context.Context, vendor *entity.Vendor) error {
	return u.repo.Update(ctx, vendor)
}

// GetVendor gets a vendor by ID
func (u *VendorUseCase) GetVendor(ctx context.Context, id uint) (*entity.Vendor, error) {
	return u.repo.FindByID(ctx, id)
}

// DeleteVendor deletes a vendor
func (u *VendorUseCase) DeleteVendor(ctx context.Context, id uint) error {
	return u.repo.Delete(ctx, id)
}

// ListVendors lists vendors with filters
func (u *VendorUseCase) ListVendors(ctx context.Context, filter entity.VendorFilter) ([]entity.Vendor, error) {
	return u.repo.List(ctx, filter)
}

// CreateProduct creates a new product
func (u *VendorUseCase) CreateProduct(ctx context.Context, product *entity.Product) error {
	return u.repo.CreateProduct(ctx, product)
}

// AddProductToVendor associates a product with a vendor
func (u *VendorUseCase) AddProductToVendor(ctx context.Context, vendorID, productID uint) error {
	// Verify vendor exists
	if _, err := u.repo.FindByID(ctx, vendorID); err != nil {
		return err
	}
	return u.repo.AddProductToVendor(ctx, vendorID, productID)
}

// RemoveProductFromVendor removes a product from a vendor
func (u *VendorUseCase) RemoveProductFromVendor(ctx context.Context, vendorID, productID uint) error {
	return u.repo.RemoveProductFromVendor(ctx, vendorID, productID)
}

// CreateContract creates a new contract for a vendor
func (u *VendorUseCase) CreateContract(ctx context.Context, contract *entity.Contract) error {
	// Verify vendor exists
	if _, err := u.repo.FindByID(ctx, contract.VendorID); err != nil {
		return err
	}
	return u.repo.CreateContract(ctx, contract)
}

// UpdateContract updates an existing contract
func (u *VendorUseCase) UpdateContract(ctx context.Context, contract *entity.Contract) error {
	return u.repo.UpdateContract(ctx, contract)
}

// GetContract gets a contract by ID
func (u *VendorUseCase) GetContract(ctx context.Context, id uint) (*entity.Contract, error) {
	return u.repo.FindContractByID(ctx, id)
}

// AddVendorRating adds a rating for a vendor
func (u *VendorUseCase) AddVendorRating(ctx context.Context, rating *entity.VendorRating) error {
	// Verify vendor exists
	if _, err := u.repo.FindByID(ctx, rating.VendorID); err != nil {
		return err
	}

	if rating.Score < 0 || rating.Score > 5 {
		return entity.ErrInvalidRating
	}

	return u.repo.CreateVendorRating(ctx, rating)
}

// GetVendorRatings gets all ratings for a vendor
func (u *VendorUseCase) GetVendorRatings(ctx context.Context, vendorID uint) ([]entity.VendorRating, error) {
	// Verify vendor exists
	if _, err := u.repo.FindByID(ctx, vendorID); err != nil {
		return nil, err
	}

	return u.repo.GetVendorRatings(ctx, vendorID)
}
