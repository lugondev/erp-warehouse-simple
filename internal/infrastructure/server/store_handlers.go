package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lugondev/erp-warehouse-simple/internal/application/usecase"
	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
)

type StoreHandler struct {
	storeUC  *usecase.StoreUseCase
	stocksUC *usecase.StocksUseCase
}

func NewStoreHandler(storeUC *usecase.StoreUseCase, stocksUC *usecase.StocksUseCase) *StoreHandler {
	return &StoreHandler{
		storeUC:  storeUC,
		stocksUC: stocksUC,
	}
}

// @Summary Create a new store
// @Description Create a new store with the provided details
// @Tags stores
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param store body entity.Store true "Store details"
// @Success 201 {object} entity.Store
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /stores [post]
func (h *StoreHandler) CreateStore(c *gin.Context) {
	var store entity.Store
	if err := c.ShouldBindJSON(&store); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.storeUC.CreateStore(c.Request.Context(), &store); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, store)
}

// @Summary Get a store by ID
// @Description Get a store's details by its ID
// @Tags stores
// @Security BearerAuth
// @Produce json
// @Param id path string true "Store ID"
// @Success 200 {object} entity.Store
// @Failure 404 {object} ErrorResponse "Store not found"
// @Router /stores/{id} [get]
func (h *StoreHandler) GetStore(c *gin.Context) {
	id := c.Param("id")

	store, err := h.storeUC.GetStore(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, store)
}

// @Summary Update a store
// @Description Update a store's details
// @Tags stores
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Store ID"
// @Param store body entity.Store true "Store details"
// @Success 200 {object} entity.Store
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /stores/{id} [put]
func (h *StoreHandler) UpdateStore(c *gin.Context) {
	id := c.Param("id")
	var store entity.Store
	if err := c.ShouldBindJSON(&store); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	store.ID = id

	if err := h.storeUC.UpdateStore(c.Request.Context(), &store); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, store)
}

// @Summary Delete a store
// @Description Delete a store by ID
// @Tags stores
// @Security BearerAuth
// @Produce json
// @Param id path string true "Store ID"
// @Success 204 "No Content"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /stores/{id} [delete]
func (h *StoreHandler) DeleteStore(c *gin.Context) {
	id := c.Param("id")

	if err := h.storeUC.DeleteStore(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary List stores
// @Description List stores with optional filtering
// @Tags stores
// @Security BearerAuth
// @Produce json
// @Param name query string false "Store name"
// @Param type query string false "Store type"
// @Param status query string false "Store status"
// @Param manager_id query int false "Manager ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} PaginatedResponse
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /stores [get]
func (h *StoreHandler) ListStores(c *gin.Context) {
	var filter entity.StoreFilter

	// Parse filters from query parameters
	filter.Name = c.Query("name")

	if typeStr := c.Query("type"); typeStr != "" {
		storeType := entity.StoreType(typeStr)
		filter.Type = &storeType
	}

	if statusStr := c.Query("status"); statusStr != "" {
		storeStatus := entity.StoreStatus(statusStr)
		filter.Status = &storeStatus
	}

	if managerIDStr := c.Query("manager_id"); managerIDStr != "" {
		if managerID, err := strconv.ParseUint(managerIDStr, 10, 32); err == nil {
			mID := uint(managerID)
			filter.ManagerID = &mID
		}
	}

	stores, err := h.storeUC.ListStores(c.Request.Context(), &filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// Manual pagination since ListStores doesn't support it at the use case level
	total := int64(len(stores))
	start := (page - 1) * pageSize
	end := start + pageSize
	if start >= len(stores) {
		start = len(stores)
	}
	if end > len(stores) {
		end = len(stores)
	}

	pagedStores := stores[start:end]

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:      pagedStores,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// @Summary Get store stocks
// @Description Get all stocks in a store
// @Tags stores
// @Security BearerAuth
// @Produce json
// @Param id path string true "Store ID"
// @Success 200 {array} entity.Stock
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /stores/{id}/stocks [get]
func (h *StoreHandler) GetStoreStocks(c *gin.Context) {
	storeID := c.Param("id")

	stocks, err := h.storeUC.GetStoreStocks(c.Request.Context(), storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, stocks)
}

// @Summary Get store stock value
// @Description Get the total value of stocks in a store
// @Tags stores
// @Security BearerAuth
// @Produce json
// @Param id path string true "Store ID"
// @Success 200 {object} map[string]float64
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /stores/{id}/stock-value [get]
func (h *StoreHandler) GetStoreStockValue(c *gin.Context) {
	storeID := c.Param("id")

	value, err := h.storeUC.GetStoreStockValue(c.Request.Context(), storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"value": value})
}

// @Summary Update store status
// @Description Update a store's status
// @Tags stores
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Store ID"
// @Param status body object true "Status details"
// @Success 200 "OK"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /stores/{id}/status [put]
func (h *StoreHandler) UpdateStoreStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status entity.StoreStatus `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.storeUC.UpdateStatus(c.Request.Context(), id, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// @Summary Assign manager to store
// @Description Assign a manager to a store
// @Tags stores
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Store ID"
// @Param manager body object true "Manager details"
// @Success 200 "OK"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /stores/{id}/manager [post]
func (h *StoreHandler) AssignManager(c *gin.Context) {
	storeID := c.Param("id")
	var req struct {
		ManagerID uint `json:"manager_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.storeUC.AssignManager(c.Request.Context(), storeID, req.ManagerID); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
