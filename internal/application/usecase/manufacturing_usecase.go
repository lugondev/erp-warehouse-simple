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
	repo       *repository.ManufacturingRepository
	stocksRepo *repository.StocksRepository
}

func NewManufacturingUseCase(repo *repository.ManufacturingRepository, stocksRepo *repository.StocksRepository) *ManufacturingUseCase {
	return &ManufacturingUseCase{
		repo:       repo,
		stocksRepo: stocksRepo,
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
	// Validate product exists in stock by checking if it has any stock records
	stocks, err := uc.stocksRepo.List(ctx, &entity.StockFilter{SKUID: fmt.Sprintf("%d", order.ProductID)})
	if err != nil || len(stocks) == 0 {
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

		// Update stock levels
		if err := uc.updateStocksOnCompletion(ctx, order); err != nil {
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

	// Calculate required quantities and check stock levels
	for _, item := range items {
		stocks, err := uc.stocksRepo.List(ctx, &entity.StockFilter{SKUID: fmt.Sprintf("%d", item.MaterialID)})
		if err != nil || len(stocks) == 0 {
			return errors.New("material not found in stock")
		}

		// Sum up available quantity across all locations
		var availableQty float64
		for _, stock := range stocks {
			availableQty += stock.Quantity
		}

		requiredQty := item.QuantityNeeded * float64(order.Quantity)
		shortageQty := requiredQty - availableQty

		mrp := &entity.MRPCalculation{
			ProductionID:  order.ID,
			MaterialID:    item.MaterialID,
			RequiredQty:   requiredQty,
			AvailableQty:  availableQty,
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

func (uc *ManufacturingUseCase) updateStocksOnCompletion(ctx context.Context, order *entity.ProductionOrder) error {
	// Update finished product stock
	skuID := fmt.Sprintf("%d", order.ProductID)
	stocks, err := uc.stocksRepo.List(ctx, &entity.StockFilter{SKUID: skuID})
	if err != nil {
		return err
	}

	var stock *entity.Stock
	if len(stocks) == 0 {
		// Create new stock record for the finished product
		stock = &entity.Stock{
			SKUID:    skuID,
			Quantity: float64(order.CompletedQty),
		}
		if err := uc.stocksRepo.CreateOrUpdateStock(ctx, stock); err != nil {
			return err
		}
	} else {
		// Update existing stock record
		stock = &stocks[0]
		stock.Quantity += float64(order.CompletedQty)
		if err := uc.stocksRepo.CreateOrUpdateStock(ctx, stock); err != nil {
			return err
		}
	}

	// Update raw materials stock based on BOM
	mrpCalculations, err := uc.repo.GetMRPCalculations(ctx, order.ID)
	if err != nil {
		return err
	}

	for _, mrp := range mrpCalculations {
		materialID := fmt.Sprintf("%d", mrp.MaterialID)
		stocks, err := uc.stocksRepo.List(ctx, &entity.StockFilter{SKUID: materialID})
		if err != nil || len(stocks) == 0 {
			return errors.New("material not found in stock")
		}

		// Deduct from the first available stock location
		stock = &stocks[0]
		stock.Quantity -= mrp.RequiredQty

		if err := uc.stocksRepo.CreateOrUpdateStock(ctx, stock); err != nil {
			return err
		}
	}

	return nil
}
