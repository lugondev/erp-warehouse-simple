package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lugondev/erp-warehouse-simple/internal/application/usecase"
	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
)

type SKUHandler struct {
	skuUseCase *usecase.SKUUseCase
}

func NewSKUHandler(skuUseCase *usecase.SKUUseCase) *SKUHandler {
	return &SKUHandler{
		skuUseCase: skuUseCase,
	}
}

// @Summary Create a new SKU
// @Description Create a new SKU with the provided details
// @Tags skus
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param sku body entity.SKU true "SKU details"
// @Success 201 {object} entity.SKU
// @Failure 400 {object} ErrorResponse "Invalid input or SKU code"
// @Failure 409 {object} ErrorResponse "Duplicate SKU code"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /skus [post]
func (h *SKUHandler) CreateSKU(c *gin.Context) {
	var sku entity.SKU
	if err := c.ShouldBindJSON(&sku); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.skuUseCase.CreateSKU(c.Request.Context(), &sku); err != nil {
		statusCode := http.StatusInternalServerError
		switch err {
		case usecase.ErrInvalidSKUCode:
			statusCode = http.StatusBadRequest
		case usecase.ErrDuplicateSKUCode:
			statusCode = http.StatusConflict
		case usecase.ErrInvalidPriceRange:
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sku)
}

// @Summary Update an SKU
// @Description Update an existing SKU with the provided details
// @Tags skus
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "SKU ID"
// @Param sku body entity.SKU true "SKU details"
// @Success 200 {object} entity.SKU
// @Failure 400 {object} ErrorResponse "Invalid input or SKU code"
// @Failure 404 {object} ErrorResponse "SKU not found"
// @Failure 409 {object} ErrorResponse "Duplicate SKU code"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /skus/{id} [put]
func (h *SKUHandler) UpdateSKU(c *gin.Context) {
	id := c.Param("id")
	var sku entity.SKU
	if err := c.ShouldBindJSON(&sku); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	sku.ID = id

	if err := h.skuUseCase.UpdateSKU(c.Request.Context(), &sku); err != nil {
		statusCode := http.StatusInternalServerError
		switch err {
		case usecase.ErrInvalidSKUCode:
			statusCode = http.StatusBadRequest
		case usecase.ErrDuplicateSKUCode:
			statusCode = http.StatusConflict
		case usecase.ErrSKUNotFound:
			statusCode = http.StatusNotFound
		case usecase.ErrInvalidPriceRange:
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, sku)
}

// @Summary Get an SKU by ID
// @Description Get an SKU's details by its ID
// @Tags skus
// @Security BearerAuth
// @Produce json
// @Param id path string true "SKU ID"
// @Success 200 {object} entity.SKU
// @Failure 404 {object} ErrorResponse "SKU not found"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /skus/{id} [get]
func (h *SKUHandler) GetSKU(c *gin.Context) {
	id := c.Param("id")

	sku, err := h.skuUseCase.GetSKU(c.Request.Context(), id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == usecase.ErrSKUNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, sku)
}

// @Summary Get an SKU by code
// @Description Get an SKU's details by its SKU code
// @Tags skus
// @Security BearerAuth
// @Produce json
// @Param code path string true "SKU Code"
// @Success 200 {object} entity.SKU
// @Failure 404 {object} ErrorResponse "SKU not found"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /skus/code/{code} [get]
func (h *SKUHandler) GetSKUByCode(c *gin.Context) {
	code := c.Param("code")

	sku, err := h.skuUseCase.GetSKUBySKUCode(c.Request.Context(), code)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == usecase.ErrSKUNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, sku)
}

// @Summary Delete an SKU
// @Description Delete an SKU by its ID
// @Tags skus
// @Security BearerAuth
// @Produce json
// @Param id path string true "SKU ID"
// @Success 204 "No Content"
// @Failure 404 {object} ErrorResponse "SKU not found"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /skus/{id} [delete]
func (h *SKUHandler) DeleteSKU(c *gin.Context) {
	id := c.Param("id")

	if err := h.skuUseCase.DeleteSKU(c.Request.Context(), id); err != nil {
		statusCode := http.StatusInternalServerError
		if err == usecase.ErrSKUNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary List SKUs
// @Description List SKUs with optional filtering
// @Tags skus
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param sku_code query string false "SKU code"
// @Param name query string false "SKU name"
// @Param category query string false "Category"
// @Param status query string false "Status"
// @Param vendor_id query int false "Vendor ID"
// @Param manufacturer_id query int false "Manufacturer ID"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Success 200 {object} PaginatedResponse
// @Failure 400 {object} ErrorResponse "Invalid price range"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /skus [get]
func (h *SKUHandler) ListSKUs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	filter := &entity.SKUFilter{
		SKUCode:  c.Query("sku_code"),
		Name:     c.Query("name"),
		Category: c.Query("category"),
	}

	if status := c.Query("status"); status != "" {
		skuStatus := entity.SKUStatus(status)
		filter.Status = &skuStatus
	}

	if vendorID := c.Query("vendor_id"); vendorID != "" {
		if id, err := strconv.ParseUint(vendorID, 10, 32); err == nil {
			vID := uint(id)
			filter.VendorID = &vID
		}
	}

	if manufacturerID := c.Query("manufacturer_id"); manufacturerID != "" {
		if id, err := strconv.ParseUint(manufacturerID, 10, 32); err == nil {
			mID := uint(id)
			filter.ManufacturerID = &mID
		}
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

	skus, total, err := h.skuUseCase.ListSKUs(c.Request.Context(), filter, page, pageSize)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == usecase.ErrInvalidPriceRange {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:      skus,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// @Summary Search SKUs
// @Description Search for SKUs by a search term
// @Tags skus
// @Security BearerAuth
// @Produce json
// @Param q query string true "Search term"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} PaginatedResponse
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /skus/search [get]
func (h *SKUHandler) SearchSKUs(c *gin.Context) {
	searchTerm := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	skus, total, err := h.skuUseCase.SearchSKUs(c.Request.Context(), searchTerm, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:      skus,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// @Summary Create a new SKU category
// @Description Create a new SKU category with the provided details
// @Tags sku-categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param category body entity.SKUCategory true "Category details"
// @Success 201 {object} entity.SKUCategory
// @Failure 400 {object} ErrorResponse "Invalid input or parent category not found"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /sku-categories [post]
func (h *SKUHandler) CreateSKUCategory(c *gin.Context) {
	var category entity.SKUCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.skuUseCase.CreateSKUCategory(c.Request.Context(), &category); err != nil {
		statusCode := http.StatusInternalServerError
		if err == usecase.ErrCategoryNotFound {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// @Summary Update an SKU category
// @Description Update an existing SKU category with the provided details
// @Tags sku-categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param category body entity.SKUCategory true "Category details"
// @Success 200 {object} entity.SKUCategory
// @Failure 400 {object} ErrorResponse "Invalid input or category cannot be its own parent"
// @Failure 404 {object} ErrorResponse "Category not found"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /sku-categories/{id} [put]
func (h *SKUHandler) UpdateSKUCategory(c *gin.Context) {
	id := c.Param("id")
	var category entity.SKUCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	category.ID = id

	if err := h.skuUseCase.UpdateSKUCategory(c.Request.Context(), &category); err != nil {
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

// @Summary Get an SKU category by ID
// @Description Get an SKU category's details by its ID
// @Tags sku-categories
// @Security BearerAuth
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} entity.SKUCategory
// @Failure 404 {object} ErrorResponse "Category not found"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /sku-categories/{id} [get]
func (h *SKUHandler) GetSKUCategory(c *gin.Context) {
	id := c.Param("id")

	category, err := h.skuUseCase.GetSKUCategory(c.Request.Context(), id)
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

// @Summary Delete an SKU category
// @Description Delete an SKU category by its ID
// @Tags sku-categories
// @Security BearerAuth
// @Produce json
// @Param id path string true "Category ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse "Cannot delete category with children or used by SKUs"
// @Failure 404 {object} ErrorResponse "Category not found"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /sku-categories/{id} [delete]
func (h *SKUHandler) DeleteSKUCategory(c *gin.Context) {
	id := c.Param("id")

	if err := h.skuUseCase.DeleteSKUCategory(c.Request.Context(), id); err != nil {
		statusCode := http.StatusInternalServerError
		switch err {
		case usecase.ErrCategoryNotFound:
			statusCode = http.StatusNotFound
		default:
			if err.Error() == "cannot delete category with children" || err.Error() == "cannot delete category used by SKUs" {
				statusCode = http.StatusBadRequest
			}
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary List all SKU categories
// @Description Get a list of all SKU categories
// @Tags sku-categories
// @Security BearerAuth
// @Produce json
// @Success 200 {array} entity.SKUCategory
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /sku-categories [get]
func (h *SKUHandler) ListSKUCategories(c *gin.Context) {
	categories, err := h.skuUseCase.ListSKUCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// @Summary Get SKU categories tree
// @Description Get SKU categories in a hierarchical tree structure
// @Tags sku-categories
// @Security BearerAuth
// @Produce json
// @Success 200 {array} entity.SKUCategory
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /sku-categories/tree [get]
func (h *SKUHandler) GetSKUCategoriesTree(c *gin.Context) {
	categories, err := h.skuUseCase.GetSKUCategoriesTree(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// @Summary Get SKUs by category
// @Description Get SKUs belonging to a specific category
// @Tags sku-categories
// @Security BearerAuth
// @Produce json
// @Param id path string true "Category ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} PaginatedResponse
// @Failure 404 {object} ErrorResponse "Category not found"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /sku-categories/{id}/skus [get]
func (h *SKUHandler) GetSKUsByCategory(c *gin.Context) {
	categoryID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	skus, total, err := h.skuUseCase.GetSKUsByCategory(c.Request.Context(), categoryID, page, pageSize)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == usecase.ErrCategoryNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:      skus,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// @Summary Bulk create SKUs
// @Description Create multiple SKUs in a single request
// @Tags skus
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param skus body []entity.SKU true "Array of SKU details"
// @Success 201 {object} map[string]string
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Router /skus/bulk [post]
func (h *SKUHandler) BulkCreateSKUs(c *gin.Context) {
	var skus []*entity.SKU
	if err := c.ShouldBindJSON(&skus); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.skuUseCase.BulkCreateSKUs(c.Request.Context(), skus); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "SKUs created successfully"})
}

// @Summary Bulk update SKUs
// @Description Update multiple SKUs in a single request
// @Tags skus
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param skus body []entity.SKU true "Array of SKU details"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Router /skus/bulk [put]
func (h *SKUHandler) BulkUpdateSKUs(c *gin.Context) {
	var skus []*entity.SKU
	if err := c.ShouldBindJSON(&skus); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.skuUseCase.BulkUpdateSKUs(c.Request.Context(), skus); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SKUs updated successfully"})
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data      interface{} `json:"data"`
	Total     int64       `json:"total"`
	Page      int         `json:"page"`
	PageSize  int         `json:"page_size"`
	TotalPage int64       `json:"total_page"`
}
