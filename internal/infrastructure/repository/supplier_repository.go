package repository

import (
	"context"

	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"

	"gorm.io/gorm"
)

type SupplierRepository struct {
	db *gorm.DB
}

func NewSupplierRepository(db *gorm.DB) *SupplierRepository {
	return &SupplierRepository{db: db}
}

// CreateSupplier creates a new supplier
func (r *SupplierRepository) CreateSupplier(ctx context.Context, supplier *entity.Supplier) error {
	return r.db.WithContext(ctx).Create(supplier).Error
}

// UpdateSupplier updates an existing supplier
func (r *SupplierRepository) UpdateSupplier(ctx context.Context, supplier *entity.Supplier) error {
	return r.db.WithContext(ctx).Save(supplier).Error
}

// GetSupplierByID retrieves a supplier by ID
func (r *SupplierRepository) GetSupplierByID(ctx context.Context, id uint) (*entity.Supplier, error) {
	var supplier entity.Supplier
	if err := r.db.WithContext(ctx).Preload("Products").Preload("Contracts").First(&supplier, id).Error; err != nil {
		return nil, err
	}
	return &supplier, nil
}

// DeleteSupplier deletes a supplier
func (r *SupplierRepository) DeleteSupplier(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Supplier{}, id).Error
}

// ListSuppliers retrieves suppliers with filters
func (r *SupplierRepository) ListSuppliers(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]entity.Supplier, int64, error) {
	var suppliers []entity.Supplier
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Supplier{})

	for key, value := range filter {
		switch key {
		case "type":
			query = query.Where("type = ?", value)
		case "country":
			query = query.Where("country = ?", value)
		case "name":
			query = query.Where("name ILIKE ?", "%"+value.(string)+"%")
		case "code":
			query = query.Where("code ILIKE ?", "%"+value.(string)+"%")
		case "min_rating":
			query = query.Where("rating >= ?", value)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Preload("Products").
		Find(&suppliers).Error; err != nil {
		return nil, 0, err
	}

	return suppliers, total, nil
}

// CreateProduct creates a new product
func (r *SupplierRepository) CreateProduct(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

// AddProductToSupplier associates a product with a supplier
func (r *SupplierRepository) AddProductToSupplier(ctx context.Context, supplierID uint, productID uint) error {
	return r.db.WithContext(ctx).Exec(
		"INSERT INTO supplier_products (supplier_id, product_id) VALUES (?, ?)",
		supplierID, productID,
	).Error
}

// RemoveProductFromSupplier removes product association from supplier
func (r *SupplierRepository) RemoveProductFromSupplier(ctx context.Context, supplierID uint, productID uint) error {
	return r.db.WithContext(ctx).Exec(
		"DELETE FROM supplier_products WHERE supplier_id = ? AND product_id = ?",
		supplierID, productID,
	).Error
}

// CreateContract creates a new contract
func (r *SupplierRepository) CreateContract(ctx context.Context, contract *entity.Contract) error {
	return r.db.WithContext(ctx).Create(contract).Error
}

// UpdateContract updates an existing contract
func (r *SupplierRepository) UpdateContract(ctx context.Context, contract *entity.Contract) error {
	return r.db.WithContext(ctx).Save(contract).Error
}

// GetContractByID retrieves a contract by ID
func (r *SupplierRepository) GetContractByID(ctx context.Context, id uint) (*entity.Contract, error) {
	var contract entity.Contract
	if err := r.db.WithContext(ctx).First(&contract, id).Error; err != nil {
		return nil, err
	}
	return &contract, nil
}

// AddSupplierRating adds a rating for a supplier
func (r *SupplierRepository) AddSupplierRating(ctx context.Context, rating *entity.SupplierRating) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(rating).Error; err != nil {
			return err
		}

		// Update supplier's average rating
		var avgRating float64
		err := tx.Model(&entity.SupplierRating{}).
			Where("supplier_id = ?", rating.SupplierID).
			Select("COALESCE(AVG(score), 0)").
			Scan(&avgRating).Error
		if err != nil {
			return err
		}

		return tx.Model(&entity.Supplier{}).
			Where("id = ?", rating.SupplierID).
			Update("rating", avgRating).Error
	})
}

// GetSupplierRatings retrieves ratings for a supplier
func (r *SupplierRepository) GetSupplierRatings(ctx context.Context, supplierID uint) ([]entity.SupplierRating, error) {
	var ratings []entity.SupplierRating
	err := r.db.WithContext(ctx).
		Where("supplier_id = ?", supplierID).
		Find(&ratings).Error
	return ratings, err
}
