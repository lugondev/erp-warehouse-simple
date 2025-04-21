package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lugondev/erp-warehouse-simple/internal/application/usecase"
	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
)

type InventoryHandler struct {
	inventoryUseCase *usecase.InventoryUseCase
}

func NewInventoryHandler(inventoryUseCase *usecase.InventoryUseCase) *InventoryHandler {
	return &InventoryHandler{
		inventoryUseCase: inventoryUseCase,
	}
}

// @Summary Process stock entry
// @Description Process a stock entry (in/out)
// @Tags Inventory
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param entry body entity.StockEntry true "Stock entry details"
// @Success 200 {object} entity.StockEntry
// @Failure 400,404 {object} ErrorResponse
// @Router /inventory/stock-entries [post]
func (h *InventoryHandler) ProcessStockEntry(c *gin.Context) {
	var entry entity.StockEntry
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	userID := getUserIDFromContext(c)
	if err := h.inventoryUseCase.ProcessStockEntry(c.Request.Context(), &entry, userID); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, entry)
}

// @Summary Check stock
// @Description Check stock level for a product in a warehouse
// @Tags Inventory
// @Security BearerAuth
// @Produce json
// @Param warehouse_id query string true "Warehouse ID"
// @Param product_id query string true "Product ID"
// @Success 200 {object} entity.Inventory
// @Failure 400,404 {object} ErrorResponse
// @Router /inventory/check-stock [get]
func (h *InventoryHandler) CheckStock(c *gin.Context) {
	warehouseID := c.Query("warehouse_id")
	productID := c.Query("product_id")

	if warehouseID == "" || productID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "warehouse_id and product_id are required"})
		return
	}

	inventory, err := h.inventoryUseCase.CheckStock(c.Request.Context(), productID, warehouseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, inventory)
}

// @Summary List inventory
// @Description List inventory with optional filters
// @Tags Inventory
// @Security BearerAuth
// @Produce json
// @Param warehouse_id query string false "Filter by warehouse ID"
// @Param product_id query string false "Filter by product ID"
// @Param batch_number query string false "Filter by batch number"
// @Success 200 {array} entity.Inventory
// @Router /inventory [get]
func (h *InventoryHandler) ListInventory(c *gin.Context) {
	filter := &entity.InventoryFilter{
		WarehouseID: c.Query("warehouse_id"),
		ProductID:   c.Query("product_id"),
		BatchNumber: c.Query("batch_number"),
		LotNumber:   c.Query("lot_number"),
	}

	inventories, err := h.inventoryUseCase.ListInventory(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, inventories)
}

// @Summary Update stock location
// @Description Update stock location details
// @Tags Inventory
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Inventory ID"
// @Param location body entity.LocationUpdate true "Location details"
// @Success 200
// @Failure 400,404 {object} ErrorResponse
// @Router /inventory/{id}/location [put]
func (h *InventoryHandler) UpdateStockLocation(c *gin.Context) {
	id := c.Param("id")
	var location entity.LocationUpdate
	if err := c.ShouldBindJSON(&location); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err := h.inventoryUseCase.UpdateStockLocation(
		c.Request.Context(),
		id,
		location.BinLocation,
		location.ShelfNumber,
		location.ZoneCode,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// @Summary Process batch stock entry
// @Description Process multiple stock entries in batch
// @Tags Inventory
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param entries body []entity.StockEntry true "Array of stock entries"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Router /inventory/batch-stock-entries [post]
func (h *InventoryHandler) BatchStockEntry(c *gin.Context) {
	var entries []entity.StockEntry
	if err := c.ShouldBindJSON(&entries); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	userID := getUserIDFromContext(c)
	if err := h.inventoryUseCase.BatchStockEntry(c.Request.Context(), entries, userID); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// @Summary Get inventory history
// @Description Get history of stock movements for an inventory item
// @Tags Inventory
// @Security BearerAuth
// @Produce json
// @Param id path string true "Inventory ID"
// @Success 200 {array} entity.InventoryHistory
// @Failure 404 {object} ErrorResponse
// @Router /inventory/{id}/history [get]
func (h *InventoryHandler) GetInventoryHistory(c *gin.Context) {
	id := c.Param("id")
	history, err := h.inventoryUseCase.GetInventoryHistory(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}

func getUserIDFromContext(c *gin.Context) string {
	// In a real application, this would get the authenticated user's ID
	// For now, return a dummy value
	return "system"
}
