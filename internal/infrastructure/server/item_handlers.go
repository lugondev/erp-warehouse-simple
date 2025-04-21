package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lugondev/erp-warehouse-simple/internal/application/usecase"
	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
)

type ItemHandler struct {
	itemUseCase *usecase.ItemUseCase
}

func NewItemHandler(itemUseCase *usecase.ItemUseCase) *ItemHandler {
	return &ItemHandler{
		itemUseCase: itemUseCase,
	}
}

// @Summary Create a new item
// @Description Create a new item with the provided information
// @Tags items
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param item body entity.Item true "Item object"
// @Success 201 {object} entity.Item
// @Failure 400 {object} ErrorResponse
// @Router /items [post]
func (h *ItemHandler) CreateItem(c *gin.Context) {
	var item entity.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.itemUseCase.CreateItem(c.Request.Context(), &item); err != nil {
		statusCode := http.StatusInternalServerError
		switch err {
		case usecase.ErrInvalidSKU, usecase.ErrInvalidPriceRange:
			statusCode = http.StatusBadRequest
		case usecase.ErrDuplicateSKU:
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// @Summary Update an item
// @Description Update an existing item's information
// @Tags items
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Param item body entity.Item true "Updated item object"
// @Success 200 {object} entity.Item
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /items/{id} [put]
func (h *ItemHandler) UpdateItem(c *gin.Context) {
	id := c.Param("id")
	var item entity.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	item.ID = id

	if err := h.itemUseCase.UpdateItem(c.Request.Context(), &item); err != nil {
		statusCode := http.StatusInternalServerError
		switch err {
		case usecase.ErrInvalidSKU, usecase.ErrInvalidPriceRange:
			statusCode = http.StatusBadRequest
		case usecase.ErrDuplicateSKU:
			statusCode = http.StatusConflict
		case usecase.ErrItemNotFound:
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// @Summary Get an item by ID
// @Description Get detailed information about an item
// @Tags items
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Success 200 {object} entity.Item
// @Failure 404 {object} ErrorResponse
// @Router /items/{id} [get]
func (h *ItemHandler) GetItem(c *gin.Context) {
	id := c.Param("id")

	item, err := h.itemUseCase.GetItem(c.Request.Context(), id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == usecase.ErrItemNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// @Summary Get an item by SKU
// @Description Get detailed information about an item using its SKU
// @Tags items
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param sku path string true "Item SKU"
// @Success 200 {object} entity.Item
// @Failure 404 {object} ErrorResponse
// @Router /items/sku/{sku} [get]
func (h *ItemHandler) GetItemBySKU(c *gin.Context) {
	sku := c.Param("sku")

	item, err := h.itemUseCase.GetItemBySKU(c.Request.Context(), sku)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == usecase.ErrItemNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// @Summary Delete an item
// @Description Delete an item by ID
// @Tags items
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Success 204 "No Content"
// @Failure 404 {object} ErrorResponse
// @Router /items/{id} [delete]
func (h *ItemHandler) DeleteItem(c *gin.Context) {
	id := c.Param("id")

	if err := h.itemUseCase.DeleteItem(c.Request.Context(), id); err != nil {
		statusCode := http.StatusInternalServerError
		if err == usecase.ErrItemNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary List items
// @Description Get a list of items with optional filters
// @Tags items
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param sku query string false "SKU filter"
// @Param name query string false "Name filter"
// @Param category query string false "Category filter"
// @Param manufacturer_id query integer false "Manufacturer ID filter"
// @Param supplier_id query integer false "Supplier ID filter"
// @Param status query string false "Status filter (ACTIVE, INACTIVE, ARCHIVED)"
// @Param min_price query number false "Minimum price filter"
// @Param max_price query number false "Maximum price filter"
// @Param page query integer false "Page number"
// @Param page_size query integer false "Page size"
// @Success 200 {object} PaginatedResponse
// @Router /items [get]
func (h *ItemHandler) ListItems(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	filter := &entity.ItemFilter{
		SKU:      c.Query("sku"),
		Name:     c.Query("name"),
		Category: c.Query("category"),
	}

	if manufacturerID := c.Query("manufacturer_id"); manufacturerID != "" {
		if id, err := strconv.ParseUint(manufacturerID, 10, 32); err == nil {
			mID := uint(id)
			filter.ManufacturerID = &mID
		}
	}

	if supplierID := c.Query("supplier_id"); supplierID != "" {
		if id, err := strconv.ParseUint(supplierID, 10, 32); err == nil {
			sID := uint(id)
			filter.SupplierID = &sID
		}
	}

	if status := c.Query("status"); status != "" {
		itemStatus := entity.ItemStatus(status)
		filter.Status = &itemStatus
	}

	if minPrice := c.Query("min_price"); minPrice != "" {
		if price, err := strconv.ParseFloat(minPrice, 64); err == nil {
			filter.MinPrice = &price
		}
	}

	if maxPrice := c.Query("max_price"); maxPrice != "" {
		if price, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			filter.MaxPrice = &price
		}
	}

	items, total, err := h.itemUseCase.ListItems(c.Request.Context(), filter, page, pageSize)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == usecase.ErrInvalidPriceRange {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:      items,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// @Summary Search items
// @Description Search for items based on a search term
// @Tags items
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param q query string true "Search term"
// @Param page query integer false "Page number"
// @Param page_size query integer false "Page size"
// @Success 200 {object} PaginatedResponse
// @Router /items/search [get]
func (h *ItemHandler) SearchItems(c *gin.Context) {
	searchTerm := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	items, total, err := h.itemUseCase.SearchItems(c.Request.Context(), searchTerm, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:      items,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// @Summary Create a new item category
// @Description Create a new item category with the provided information
// @Tags item-categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param category body entity.ItemCategory true "Item category object"
// @Success 201 {object} entity.ItemCategory
// @Failure 400 {object} ErrorResponse
// @Router /item-categories [post]
func (h *ItemHandler) CreateItemCategory(c *gin.Context) {
	var category entity.ItemCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.itemUseCase.CreateItemCategory(c.Request.Context(), &category); err != nil {
		statusCode := http.StatusInternalServerError
		if err == usecase.ErrCategoryNotFound {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// @Summary Update an item category
// @Description Update an existing item category's information
// @Tags item-categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param category body entity.ItemCategory true "Updated category object"
// @Success 200 {object} entity.ItemCategory
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /item-categories/{id} [put]
func (h *ItemHandler) UpdateItemCategory(c *gin.Context) {
	id := c.Param("id")
	var category entity.ItemCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	category.ID = id

	if err := h.itemUseCase.UpdateItemCategory(c.Request.Context(), &category); err != nil {
		statusCode := http.StatusInternalServerError
		switch err {
		case usecase.ErrCategoryNotFound:
			statusCode = http.StatusNotFound
		default:
			if err.Error() == "category cannot be its own parent" {
				statusCode = http.StatusBadRequest
			}
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// @Summary Get an item category by ID
// @Description Get detailed information about an item category
// @Tags item-categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} entity.ItemCategory
// @Failure 404 {object} ErrorResponse
// @Router /item-categories/{id} [get]
func (h *ItemHandler) GetItemCategory(c *gin.Context) {
	id := c.Param("id")

	category, err := h.itemUseCase.GetItemCategory(c.Request.Context(), id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == usecase.ErrCategoryNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// @Summary Delete an item category
// @Description Delete an item category by ID
// @Tags item-categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /item-categories/{id} [delete]
func (h *ItemHandler) DeleteItemCategory(c *gin.Context) {
	id := c.Param("id")

	if err := h.itemUseCase.DeleteItemCategory(c.Request.Context(), id); err != nil {
		statusCode := http.StatusInternalServerError
		switch err {
		case usecase.ErrCategoryNotFound:
			statusCode = http.StatusNotFound
		default:
			if err.Error() == "cannot delete category with children" || err.Error() == "cannot delete category used by items" {
				statusCode = http.StatusBadRequest
			}
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary List item categories
// @Description Get a list of all item categories
// @Tags item-categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {array} entity.ItemCategory
// @Router /item-categories [get]
func (h *ItemHandler) ListItemCategories(c *gin.Context) {
	categories, err := h.itemUseCase.ListItemCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// @Summary Get item categories tree
// @Description Get item categories in a hierarchical structure
// @Tags item-categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {array} entity.ItemCategory
// @Router /item-categories/tree [get]
func (h *ItemHandler) GetItemCategoriesTree(c *gin.Context) {
	categories, err := h.itemUseCase.GetItemCategoriesTree(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// @Summary Get items by category
// @Description Get items belonging to a specific category
// @Tags items
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param page query integer false "Page number"
// @Param page_size query integer false "Page size"
// @Success 200 {object} PaginatedResponse
// @Failure 404 {object} ErrorResponse
// @Router /item-categories/{id}/items [get]
func (h *ItemHandler) GetItemsByCategory(c *gin.Context) {
	categoryID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	items, total, err := h.itemUseCase.GetItemsByCategory(c.Request.Context(), categoryID, page, pageSize)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == usecase.ErrCategoryNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:      items,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// @Summary Bulk create items
// @Description Create multiple items in a single request
// @Tags items
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param items body []entity.Item true "Array of item objects"
// @Success 201 {string} string "Items created successfully"
// @Failure 400 {object} ErrorResponse
// @Router /items/bulk [post]
func (h *ItemHandler) BulkCreateItems(c *gin.Context) {
	var items []*entity.Item
	if err := c.ShouldBindJSON(&items); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.itemUseCase.BulkCreateItems(c.Request.Context(), items); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Items created successfully"})
}

// @Summary Bulk update items
// @Description Update multiple items in a single request
// @Tags items
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param items body []entity.Item true "Array of item objects with IDs"
// @Success 200 {string} string "Items updated successfully"
// @Failure 400 {object} ErrorResponse
// @Router /items/bulk [put]
func (h *ItemHandler) BulkUpdateItems(c *gin.Context) {
	var items []*entity.Item
	if err := c.ShouldBindJSON(&items); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.itemUseCase.BulkUpdateItems(c.Request.Context(), items); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Items updated successfully"})
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data      interface{} `json:"data"`
	Total     int64       `json:"total"`
	Page      int         `json:"page"`
	PageSize  int         `json:"page_size"`
	TotalPage int64       `json:"total_page"`
}
