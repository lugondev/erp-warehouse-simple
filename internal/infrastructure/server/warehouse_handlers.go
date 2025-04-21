package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lugondev/erp-warehouse-simple/internal/application/usecase"
	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
)

type WarehouseHandler struct {
	warehouseUseCase *usecase.WarehouseUseCase
	inventoryUseCase *usecase.InventoryUseCase
}

func NewWarehouseHandler(warehouseUseCase *usecase.WarehouseUseCase, inventoryUseCase *usecase.InventoryUseCase) *WarehouseHandler {
	return &WarehouseHandler{
		warehouseUseCase: warehouseUseCase,
		inventoryUseCase: inventoryUseCase,
	}
}

// @Summary Create warehouse
// @Description Create a new warehouse
// @Tags Warehouses
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param warehouse body entity.Warehouse true "Warehouse info"
// @Success 201 {object} entity.Warehouse
// @Failure 400 {object} ErrorResponse
// @Router /warehouses [post]
func (h *WarehouseHandler) CreateWarehouse(c *gin.Context) {
	var warehouse entity.Warehouse
	if err := c.ShouldBindJSON(&warehouse); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.warehouseUseCase.CreateWarehouse(c.Request.Context(), &warehouse); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, warehouse)
}

// @Summary Update warehouse
// @Description Update an existing warehouse
// @Tags Warehouses
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Warehouse ID"
// @Param warehouse body entity.Warehouse true "Warehouse info"
// @Success 200 {object} entity.Warehouse
// @Failure 400,404 {object} ErrorResponse
// @Router /warehouses/{id} [put]
func (h *WarehouseHandler) UpdateWarehouse(c *gin.Context) {
	id := c.Param("id")
	var warehouse entity.Warehouse
	if err := c.ShouldBindJSON(&warehouse); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	warehouse.ID = id
	if err := h.warehouseUseCase.UpdateWarehouse(c.Request.Context(), &warehouse); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, warehouse)
}

// @Summary Get warehouse
// @Description Get warehouse by ID
// @Tags Warehouses
// @Security BearerAuth
// @Produce json
// @Param id path string true "Warehouse ID"
// @Success 200 {object} entity.Warehouse
// @Failure 404 {object} ErrorResponse
// @Router /warehouses/{id} [get]
func (h *WarehouseHandler) GetWarehouse(c *gin.Context) {
	id := c.Param("id")
	warehouse, err := h.warehouseUseCase.GetWarehouse(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, warehouse)
}

// @Summary List warehouses
// @Description List all warehouses with optional filters
// @Tags Warehouses
// @Security BearerAuth
// @Produce json
// @Param type query string false "Warehouse type"
// @Param status query string false "Warehouse status"
// @Success 200 {array} entity.Warehouse
// @Router /warehouses [get]
func (h *WarehouseHandler) ListWarehouses(c *gin.Context) {
	var filter entity.WarehouseFilter

	// Parse query parameters
	if typeStr := c.Query("type"); typeStr != "" {
		warehouseType := entity.WarehouseType(typeStr)
		filter.Type = &warehouseType
	}
	if statusStr := c.Query("status"); statusStr != "" {
		status := entity.WarehouseStatus(statusStr)
		filter.Status = &status
	}

	warehouses, err := h.warehouseUseCase.ListWarehouses(c.Request.Context(), &filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, warehouses)
}

// @Summary Delete warehouse
// @Description Delete a warehouse by ID
// @Tags Warehouses
// @Security BearerAuth
// @Produce json
// @Param id path string true "Warehouse ID"
// @Success 204 "No Content"
// @Failure 404,500 {object} ErrorResponse
// @Router /warehouses/{id} [delete]
func (h *WarehouseHandler) DeleteWarehouse(c *gin.Context) {
	id := c.Param("id")
	if err := h.warehouseUseCase.DeleteWarehouse(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
