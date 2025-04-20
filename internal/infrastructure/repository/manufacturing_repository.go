package repository

import (
	"context"

	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"

	"gorm.io/gorm"
)

type ManufacturingRepository struct {
	db *gorm.DB
}

func NewManufacturingRepository(db *gorm.DB) *ManufacturingRepository {
	return &ManufacturingRepository{db: db}
}

// Facility methods
func (r *ManufacturingRepository) CreateFacility(ctx context.Context, facility *entity.ManufacturingFacility) error {
	return r.db.WithContext(ctx).Create(facility).Error
}

func (r *ManufacturingRepository) GetFacility(ctx context.Context, id uint) (*entity.ManufacturingFacility, error) {
	var facility entity.ManufacturingFacility
	if err := r.db.WithContext(ctx).First(&facility, id).Error; err != nil {
		return nil, err
	}
	return &facility, nil
}

func (r *ManufacturingRepository) ListFacilities(ctx context.Context) ([]entity.ManufacturingFacility, error) {
	var facilities []entity.ManufacturingFacility
	if err := r.db.WithContext(ctx).Find(&facilities).Error; err != nil {
		return nil, err
	}
	return facilities, nil
}

func (r *ManufacturingRepository) UpdateFacility(ctx context.Context, facility *entity.ManufacturingFacility) error {
	return r.db.WithContext(ctx).Save(facility).Error
}

func (r *ManufacturingRepository) DeleteFacility(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.ManufacturingFacility{}, id).Error
}

// Production Order methods
func (r *ManufacturingRepository) CreateProductionOrder(ctx context.Context, order *entity.ProductionOrder) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *ManufacturingRepository) GetProductionOrder(ctx context.Context, id uint) (*entity.ProductionOrder, error) {
	var order entity.ProductionOrder
	if err := r.db.WithContext(ctx).First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *ManufacturingRepository) ListProductionOrders(ctx context.Context, facilityID *uint) ([]entity.ProductionOrder, error) {
	var orders []entity.ProductionOrder
	query := r.db.WithContext(ctx)
	if facilityID != nil {
		query = query.Where("facility_id = ?", *facilityID)
	}
	if err := query.Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *ManufacturingRepository) UpdateProductionOrder(ctx context.Context, order *entity.ProductionOrder) error {
	return r.db.WithContext(ctx).Save(order).Error
}

// BOM methods
func (r *ManufacturingRepository) CreateBOM(ctx context.Context, bom *entity.BillOfMaterial) error {
	return r.db.WithContext(ctx).Create(bom).Error
}

func (r *ManufacturingRepository) GetBOM(ctx context.Context, id uint) (*entity.BillOfMaterial, error) {
	var bom entity.BillOfMaterial
	if err := r.db.WithContext(ctx).First(&bom, id).Error; err != nil {
		return nil, err
	}
	return &bom, nil
}

func (r *ManufacturingRepository) AddBOMItem(ctx context.Context, item *entity.BOMItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *ManufacturingRepository) GetBOMItems(ctx context.Context, bomID uint) ([]entity.BOMItem, error) {
	var items []entity.BOMItem
	if err := r.db.WithContext(ctx).Where("bom_id = ?", bomID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// MRP methods
func (r *ManufacturingRepository) CreateMRPCalculation(ctx context.Context, mrp *entity.MRPCalculation) error {
	return r.db.WithContext(ctx).Create(mrp).Error
}

func (r *ManufacturingRepository) GetMRPCalculations(ctx context.Context, productionID uint) ([]entity.MRPCalculation, error) {
	var calculations []entity.MRPCalculation
	if err := r.db.WithContext(ctx).Where("production_id = ?", productionID).Find(&calculations).Error; err != nil {
		return nil, err
	}
	return calculations, nil
}
