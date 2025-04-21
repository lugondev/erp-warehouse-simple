package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lugondev/erp-warehouse-simple/internal/application/usecase"
	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
)

type StocksHandler struct {
	stocksUC *usecase.StocksUseCase
}

func NewStocksHandler(stocksUC *usecase.StocksUseCase) *StocksHandler {
	return &StocksHandler{
		stocksUC: stocksUC,
	}
}

// @Summary List stocks
// @Description List stocks with optional filtering
// @Tags stocks
// @Security BearerAuth
// @Produce json
// @Param sku_id query string false "SKU ID"
// @Param store_id query string false "Store ID"
// @Param batch_number query string false "Batch number"
// @Param lot_number query string false "Lot number"
// @Param zone_code query string false "Zone code"
// @Param bin_location query string false "Bin location"
// @Param shelf_number query string false "Shelf number"
// @Param min_quantity query number false "Minimum quantity"
// @Param max_quantity query number false "Maximum quantity"
// @Param expiry_date_from query string false "Expiry date from (RFC3339 format)"
// @Param expiry_date_to query string false "Expiry date to (RFC3339 format)"
// @Success 200 {array} entity.Stock
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /stocks [get]
func (h *StocksHandler) ListStocks(c *gin.Context) {
	filter := &entity.StockFilter{
		SKUID:       c.Query("sku_id"),
		StoreID:     c.Query("store_id"),
		BatchNumber: c.Query("batch_number"),
		LotNumber:   c.Query("lot_number"),
		ZoneCode:    c.Query("zone_code"),
		BinLocation: c.Query("bin_location"),
		ShelfNumber: c.Query("shelf_number"),
	}

	if minQty := c.Query("min_quantity"); minQty != "" {
		if qty, err := strconv.ParseFloat(minQty, 64); err == nil {
			filter.MinQuantity = qty
		}
	}

	if maxQty := c.Query("max_quantity"); maxQty != "" {
		if qty, err := strconv.ParseFloat(maxQty, 64); err == nil {
			filter.MaxQuantity = qty
		}
	}

	if expiryFrom := c.Query("expiry_date_from"); expiryFrom != "" {
		if t, err := time.Parse(time.RFC3339, expiryFrom); err == nil {
			filter.ExpiryDateFrom = t
		}
	}

	if expiryTo := c.Query("expiry_date_to"); expiryTo != "" {
		if t, err := time.Parse(time.RFC3339, expiryTo); err == nil {
			filter.ExpiryDateTo = t
		}
	}

	stocks, err := h.stocksUC.ListStocks(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, stocks)
}

// @Summary Check stock availability
// @Description Check stock availability for a specific SKU in a store
// @Tags stocks
// @Security BearerAuth
// @Produce json
// @Param sku_id query string true "SKU ID"
// @Param store_id query string true "Store ID"
// @Success 200 {object} entity.Stock
// @Failure 400 {object} ErrorResponse "Missing required parameters"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /stocks/check-stock [get]
func (h *StocksHandler) CheckStock(c *gin.Context) {
	skuID := c.Query("sku_id")
	storeID := c.Query("store_id")

	if skuID == "" || storeID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "sku_id and store_id are required"})
		return
	}

	stock, err := h.stocksUC.CheckStock(c.Request.Context(), skuID, storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, stock)
}

// @Summary Process stock entry
// @Description Process a stock entry (add, remove, transfer, adjust)
// @Tags stocks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param entry body entity.StockEntry true "Stock entry details"
// @Success 200 {object} entity.StockEntry
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 401 {object} ErrorResponse "User not authenticated"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /stocks/stock-entries [post]
func (h *StocksHandler) ProcessStockEntry(c *gin.Context) {
	var entry entity.StockEntry
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	userID := c.GetString("user_id") // Assuming user_id is set in auth middleware
	if userID == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "user not authenticated"})
		return
	}

	if err := h.stocksUC.ValidateStockEntry(&entry); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.stocksUC.ProcessStockEntry(c.Request.Context(), &entry, userID); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, entry)
}

// @Summary Process multiple stock entries
// @Description Process multiple stock entries in a single request
// @Tags stocks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param entries body []entity.StockEntry true "Array of stock entry details"
// @Success 200 "OK"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 401 {object} ErrorResponse "User not authenticated"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /stocks/batch-stock-entries [post]
func (h *StocksHandler) BatchStockEntry(c *gin.Context) {
	var entries []entity.StockEntry
	if err := c.ShouldBindJSON(&entries); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "user not authenticated"})
		return
	}

	// Validate all entries first
	for _, entry := range entries {
		if err := h.stocksUC.ValidateStockEntry(&entry); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}
	}

	if err := h.stocksUC.BatchStockEntry(c.Request.Context(), entries, userID); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// @Summary Update stock location
// @Description Update a stock's location information
// @Tags stocks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Stock ID"
// @Param location body object true "Location details"
// @Success 200 "OK"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /stocks/{id}/location [put]
func (h *StocksHandler) UpdateStockLocation(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		BinLocation string `json:"bin_location"`
		ShelfNumber string `json:"shelf_number"`
		ZoneCode    string `json:"zone_code"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.stocksUC.UpdateStockLocation(c.Request.Context(), id, req.BinLocation, req.ShelfNumber, req.ZoneCode); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// @Summary Get stock history
// @Description Get the history of a stock
// @Tags stocks
// @Security BearerAuth
// @Produce json
// @Param id path string true "Stock ID"
// @Success 200 {array} entity.StockEntry
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /stocks/{id}/history [get]
func (h *StocksHandler) GetStockHistory(c *gin.Context) {
	stockID := c.Param("id")

	history, err := h.stocksUC.GetStockHistory(c.Request.Context(), stockID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}
