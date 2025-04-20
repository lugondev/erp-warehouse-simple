package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
	"github.com/lugondev/erp-warehouse-simple/internal/infrastructure/repository"
)

type ManufacturingUseCase struct {
	repo          *repository.ManufacturingRepository
	inventoryRepo *repository.InventoryRepository
}

func NewManufacturingUseCase(repo *repository.ManufacturingRepository, inventoryRepo *repository.InventoryRepository) *ManufacturingUseCase {
	return &ManufacturingUseCase{
		repo:          repo,
		inventoryRepo: inventoryRepo,
	}
}

// Facility management
func (uc *ManufacturingUseCase) CreateFacility(ctx context.Context, facility *entity.ManufacturingFacility) error {
	return uc.repo.CreateFacility(ctx, facility)
}

func (uc *ManufacturingUseCase) GetFacility(ctx context.Context, id uint) (*entity.ManufacturingFacility, error) {
	return uc.repo.GetFacility(ctx, id)
}

func (uc *ManufacturingUseCase) ListFacilities(ctx context.Context) ([]entity.ManufacturingFacility, error) {
	return uc.repo.ListFacilities(ctx)
}

func (uc *ManufacturingUseCase) UpdateFacility(ctx context.Context, facility *entity.ManufacturingFacility) error {
	return uc.repo.UpdateFacility(ctx, facility)
}

// Production order management
func (uc *ManufacturingUseCase) CreateProductionOrder(ctx context.Context, order *entity.ProductionOrder) error {
	// Validate product exists in inventory
	if _, err := uc.inventoryRepo.GetByID(ctx, fmt.Sprintf("%d", order.ProductID)); err != nil {
		return errors.New("invalid product ID")
	}

	// Validate facility exists
	if _, err := uc.repo.GetFacility(ctx, order.FacilityID); err != nil {
		return errors.New("invalid facility ID")
	}

	order.Status = entity.OrderStatusPending
	order.StartDate = time.Now()
	return uc.repo.CreateProductionOrder(ctx, order)
}

func (uc *ManufacturingUseCase) UpdateProductionProgress(ctx context.Context, orderID uint, completedQty, defectQty int) error {
	order, err := uc.repo.GetProductionOrder(ctx, orderID)
	if err != nil {
		return err
	}

	if order.Status != entity.OrderStatusInProcess {
		return errors.New("production order is not in process")
	}

	order.CompletedQty = completedQty
	order.DefectQty = defectQty

	// Check if production is complete
	if completedQty >= order.Quantity {
		order.Status = entity.OrderStatusCompleted

		// Update inventory
		if err := uc.updateInventoryOnCompletion(ctx, order); err != nil {
			return err
		}
	}

	return uc.repo.UpdateProductionOrder(ctx, order)
}

func (uc *ManufacturingUseCase) StartProduction(ctx context.Context, orderID uint) error {
	order, err := uc.repo.GetProductionOrder(ctx, orderID)
	if err != nil {
		return err
	}

	if order.Status != entity.OrderStatusPending {
		return errors.New("order is not in pending status")
	}

	// Calculate material requirements
	if err := uc.calculateMRP(ctx, order); err != nil {
		return err
	}

	order.Status = entity.OrderStatusInProcess
	order.StartDate = time.Now()
	return uc.repo.UpdateProductionOrder(ctx, order)
}

// BOM management
func (uc *ManufacturingUseCase) CreateBOM(ctx context.Context, bom *entity.BillOfMaterial, items []entity.BOMItem) error {
	if err := uc.repo.CreateBOM(ctx, bom); err != nil {
		return err
	}

	for _, item := range items {
		item.BOMID = bom.ID
		if err := uc.repo.AddBOMItem(ctx, &item); err != nil {
			return err
		}
	}

	return nil
}

// MRP calculation
func (uc *ManufacturingUseCase) calculateMRP(ctx context.Context, order *entity.ProductionOrder) error {
	// Get BOM for the product
	bom, err := uc.repo.GetBOM(ctx, order.ProductID)
	if err != nil {
		return err
	}

	// Get BOM items
	items, err := uc.repo.GetBOMItems(ctx, bom.ID)
	if err != nil {
		return err
	}

	// Calculate required quantities and check inventory
	for _, item := range items {
		inventory, err := uc.inventoryRepo.GetByID(ctx, fmt.Sprintf("%d", item.MaterialID))
		if err != nil {
			return err
		}

		requiredQty := item.QuantityNeeded * float64(order.Quantity)
		shortageQty := requiredQty - float64(inventory.Quantity)

		mrp := &entity.MRPCalculation{
			ProductionID:  order.ID,
			MaterialID:    item.MaterialID,
			RequiredQty:   requiredQty,
			AvailableQty:  float64(inventory.Quantity),
			ShortageQty:   shortageQty,
			UnitOfMeasure: item.UnitOfMeasure,
			CalculatedAt:  time.Now(),
		}

		if err := uc.repo.CreateMRPCalculation(ctx, mrp); err != nil {
			return err
		}
	}

	return nil
}

func (uc *ManufacturingUseCase) updateInventoryOnCompletion(ctx context.Context, order *entity.ProductionOrder) error {
	// Update finished product inventory
	product, err := uc.inventoryRepo.GetByID(ctx, fmt.Sprintf("%d", order.ProductID))
	if err != nil {
		return err
	}

	product.Quantity += float64(order.CompletedQty)
	if err := uc.inventoryRepo.UpdateQuantity(ctx, product.ID, float64(product.Quantity)); err != nil {
		return err
	}

	// Update raw materials inventory based on BOM
	mrpCalculations, err := uc.repo.GetMRPCalculations(ctx, order.ID)
	if err != nil {
		return err
	}

	for _, mrp := range mrpCalculations {
		material, err := uc.inventoryRepo.GetByID(ctx, fmt.Sprintf("%d", mrp.MaterialID))
		if err != nil {
			return err
		}

		// Deduct used materials
		usedQty := int(mrp.RequiredQty)
		material.Quantity -= float64(usedQty)

		if err := uc.inventoryRepo.UpdateQuantity(ctx, material.ID, float64(material.Quantity)); err != nil {
			return err
		}
	}

	return nil
}
