package repository

import (
	"context"

	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"

	"gorm.io/gorm"
)

type VendorRepository struct {
	db *gorm.DB
}

func NewVendorRepository(db *gorm.DB) *VendorRepository {
	return &VendorRepository{db: db}
}

// Create creates a new vendor
func (r *VendorRepository) Create(ctx context.Context, vendor *entity.Vendor) error {
	return r.db.WithContext(ctx).Create(vendor).Error
}

// Update updates an existing vendor
func (r *VendorRepository) Update(ctx context.Context, vendor *entity.Vendor) error {
	return r.db.WithContext(ctx).Save(vendor).Error
}

// FindByID retrieves a vendor by ID
func (r *VendorRepository) FindByID(ctx context.Context, id uint) (*entity.Vendor, error) {
	var vendor entity.Vendor
	if err := r.db.WithContext(ctx).Preload("Products").Preload("Contracts").First(&vendor, id).Error; err != nil {
		return nil, err
	}
	return &vendor, nil
}

// FindByCode retrieves a vendor by code
func (r *VendorRepository) FindByCode(ctx context.Context, code string) (*entity.Vendor, error) {
	var vendor entity.Vendor
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&vendor).Error; err != nil {
		return nil, err
	}
	return &vendor, nil
}

// Delete deletes a vendor
func (r *VendorRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Vendor{}, id).Error
}

// List retrieves vendors with filters
func (r *VendorRepository) List(ctx context.Context, filter entity.VendorFilter) ([]entity.Vendor, error) {
	var vendors []entity.Vendor
	query := r.db.WithContext(ctx).Model(&entity.Vendor{})

	if filter.Code != "" {
		query = query.Where("code ILIKE ?", "%"+filter.Code+"%")
	}
	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.Country != "" {
		query = query.Where("country = ?", filter.Country)
	}
	if filter.MinRating != nil {
		query = query.Where("rating >= ?", *filter.MinRating)
	}
	if len(filter.ProductIDs) > 0 {
		query = query.Joins("JOIN vendor_products ON vendor_products.vendor_id = vendors.id").
			Where("vendor_products.product_id IN ?", filter.ProductIDs).
			Group("vendors.id")
	}

	if err := query.Preload("Products").Find(&vendors).Error; err != nil {
		return nil, err
	}

	return vendors, nil
}

// CreateProduct creates a new product
func (r *VendorRepository) CreateProduct(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

// UpdateProduct updates an existing product
func (r *VendorRepository) UpdateProduct(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

// DeleteProduct deletes a product
func (r *VendorRepository) DeleteProduct(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Product{}, id).Error
}

// FindProductByID retrieves a product by ID
func (r *VendorRepository) FindProductByID(ctx context.Context, id uint) (*entity.Product, error) {
	var product entity.Product
	if err := r.db.WithContext(ctx).First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

// FindProductByCode retrieves a product by code
func (r *VendorRepository) FindProductByCode(ctx context.Context, code string) (*entity.Product, error) {
	var product entity.Product
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

// ListProducts retrieves all products
func (r *VendorRepository) ListProducts(ctx context.Context) ([]entity.Product, error) {
	var products []entity.Product
	if err := r.db.WithContext(ctx).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// AddProductToVendor associates a product with a vendor
func (r *VendorRepository) AddProductToVendor(ctx context.Context, vendorID uint, productID uint) error {
	return r.db.WithContext(ctx).Exec(
		"INSERT INTO vendor_products (vendor_id, product_id) VALUES (?, ?)",
		vendorID, productID,
	).Error
}

// RemoveProductFromVendor removes product association from vendor
func (r *VendorRepository) RemoveProductFromVendor(ctx context.Context, vendorID uint, productID uint) error {
	return r.db.WithContext(ctx).Exec(
		"DELETE FROM vendor_products WHERE vendor_id = ? AND product_id = ?",
		vendorID, productID,
	).Error
}

// GetVendorProducts gets all products for a vendor
func (r *VendorRepository) GetVendorProducts(ctx context.Context, vendorID uint) ([]entity.Product, error) {
	var products []entity.Product
	err := r.db.WithContext(ctx).
		Joins("JOIN vendor_products ON vendor_products.product_id = products.id").
		Where("vendor_products.vendor_id = ?", vendorID).
		Find(&products).Error
	return products, err
}

// CreateContract creates a new contract
func (r *VendorRepository) CreateContract(ctx context.Context, contract *entity.Contract) error {
	return r.db.WithContext(ctx).Create(contract).Error
}

// UpdateContract updates an existing contract
func (r *VendorRepository) UpdateContract(ctx context.Context, contract *entity.Contract) error {
	return r.db.WithContext(ctx).Save(contract).Error
}

// DeleteContract deletes a contract
func (r *VendorRepository) DeleteContract(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Contract{}, id).Error
}

// FindContractByID retrieves a contract by ID
func (r *VendorRepository) FindContractByID(ctx context.Context, id uint) (*entity.Contract, error) {
	var contract entity.Contract
	if err := r.db.WithContext(ctx).First(&contract, id).Error; err != nil {
		return nil, err
	}
	return &contract, nil
}

// ListVendorContracts lists all contracts for a vendor
func (r *VendorRepository) ListVendorContracts(ctx context.Context, vendorID uint) ([]entity.Contract, error) {
	var contracts []entity.Contract
	err := r.db.WithContext(ctx).
		Where("vendor_id = ?", vendorID).
		Find(&contracts).Error
	return contracts, err
}

// CreateVendorRating adds a rating for a vendor
func (r *VendorRepository) CreateVendorRating(ctx context.Context, rating *entity.VendorRating) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(rating).Error; err != nil {
			return err
		}

		// Update vendor's average rating
		var avgRating float64
		err := tx.Model(&entity.VendorRating{}).
			Where("vendor_id = ?", rating.VendorID).
			Select("COALESCE(AVG(score), 0)").
			Scan(&avgRating).Error
		if err != nil {
			return err
		}

		return tx.Model(&entity.Vendor{}).
			Where("id = ?", rating.VendorID).
			Update("rating", avgRating).Error
	})
}

// GetVendorRatings retrieves ratings for a vendor
func (r *VendorRepository) GetVendorRatings(ctx context.Context, vendorID uint) ([]entity.VendorRating, error) {
	var ratings []entity.VendorRating
	err := r.db.WithContext(ctx).
		Where("vendor_id = ?", vendorID).
		Find(&ratings).Error
	return ratings, err
}

// GetVendorAverageRating gets the average rating for a vendor
func (r *VendorRepository) GetVendorAverageRating(ctx context.Context, vendorID uint) (float64, error) {
	var avgRating float64
	err := r.db.WithContext(ctx).Model(&entity.VendorRating{}).
		Where("vendor_id = ?", vendorID).
		Select("COALESCE(AVG(score), 0)").
		Scan(&avgRating).Error
	return avgRating, err
}
