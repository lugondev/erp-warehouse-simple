package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lugondev/erp-warehouse-simple/internal/application/usecase"
	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
)

type SupplierHandler struct {
	supplierUseCase *usecase.SupplierUseCase
}

func NewSupplierHandler(supplierUseCase *usecase.SupplierUseCase) *SupplierHandler {
	return &SupplierHandler{
		supplierUseCase: supplierUseCase,
	}
}

// RegisterRoutes registers supplier routes
func (h *SupplierHandler) RegisterRoutes(router *gin.Engine) {
	suppliers := router.Group("/api/suppliers")
	{
		suppliers.POST("", h.CreateSupplier)
		suppliers.GET("", h.ListSuppliers)
		suppliers.GET("/:id", h.GetSupplier)
		suppliers.PUT("/:id", h.UpdateSupplier)
		suppliers.DELETE("/:id", h.DeleteSupplier)

		// Product management
		suppliers.POST("/products", h.CreateProduct)
		suppliers.POST("/:id/products/:productId", h.AddProductToSupplier)
		suppliers.DELETE("/:id/products/:productId", h.RemoveProductFromSupplier)

		// Contract management
		suppliers.POST("/:id/contracts", h.CreateContract)
		suppliers.PUT("/contracts/:contractId", h.UpdateContract)
		suppliers.GET("/contracts/:contractId", h.GetContract)

		// Rating management
		suppliers.POST("/:id/ratings", h.AddRating)
		suppliers.GET("/:id/ratings", h.GetRatings)
	}
}

// @Summary Create a new supplier
// @Description Create a new supplier with the provided information
// @Tags suppliers
// @Accept json
// @Produce json
// @Param supplier body entity.Supplier true "Supplier object"
// @Success 201 {object} entity.Supplier
// @Failure 400 {object} ErrorResponse
// @Router /api/suppliers [post]
func (h *SupplierHandler) CreateSupplier(c *gin.Context) {
	var supplier entity.Supplier
	if err := c.ShouldBindJSON(&supplier); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.supplierUseCase.CreateSupplier(c.Request.Context(), &supplier); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, supplier)
}

// @Summary List suppliers
// @Description Get a list of suppliers with optional filters
// @Tags suppliers
// @Accept json
// @Produce json
// @Param type query string false "Supplier type"
// @Param country query string false "Country"
// @Param name query string false "Name search"
// @Param code query string false "Code search"
// @Param min_rating query number false "Minimum rating"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {array} entity.Supplier
// @Router /api/suppliers [get]
func (h *SupplierHandler) ListSuppliers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	filter := make(map[string]interface{})
	if supplierType := c.Query("type"); supplierType != "" {
		filter["type"] = supplierType
	}
	if country := c.Query("country"); country != "" {
		filter["country"] = country
	}
	if name := c.Query("name"); name != "" {
		filter["name"] = name
	}
	if code := c.Query("code"); code != "" {
		filter["code"] = code
	}
	if minRating := c.Query("min_rating"); minRating != "" {
		if rating, err := strconv.ParseFloat(minRating, 64); err == nil {
			filter["min_rating"] = rating
		}
	}

	suppliers, total, err := h.supplierUseCase.ListSuppliers(c.Request.Context(), filter, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"suppliers": suppliers,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// @Summary Get a supplier by ID
// @Description Get detailed information about a supplier
// @Tags suppliers
// @Accept json
// @Produce json
// @Param id path int true "Supplier ID"
// @Success 200 {object} entity.Supplier
// @Failure 404 {object} ErrorResponse
// @Router /api/suppliers/{id} [get]
func (h *SupplierHandler) GetSupplier(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid supplier id"})
		return
	}

	supplier, err := h.supplierUseCase.GetSupplier(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, supplier)
}

// @Summary Update a supplier
// @Description Update an existing supplier's information
// @Tags suppliers
// @Accept json
// @Produce json
// @Param id path int true "Supplier ID"
// @Param supplier body entity.Supplier true "Updated supplier object"
// @Success 200 {object} entity.Supplier
// @Failure 400 {object} ErrorResponse
// @Router /api/suppliers/{id} [put]
func (h *SupplierHandler) UpdateSupplier(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid supplier id"})
		return
	}

	var supplier entity.Supplier
	if err := c.ShouldBindJSON(&supplier); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	supplier.ID = uint(id)

	if err := h.supplierUseCase.UpdateSupplier(c.Request.Context(), &supplier); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, supplier)
}

// @Summary Delete a supplier
// @Description Delete a supplier by ID
// @Tags suppliers
// @Accept json
// @Produce json
// @Param id path int true "Supplier ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Router /api/suppliers/{id} [delete]
func (h *SupplierHandler) DeleteSupplier(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid supplier id"})
		return
	}

	if err := h.supplierUseCase.DeleteSupplier(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Product Management Handlers

func (h *SupplierHandler) CreateProduct(c *gin.Context) {
	var product entity.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.supplierUseCase.CreateProduct(c.Request.Context(), &product); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (h *SupplierHandler) AddProductToSupplier(c *gin.Context) {
	supplierID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid supplier id"})
		return
	}

	productID, err := strconv.ParseUint(c.Param("productId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid product id"})
		return
	}

	err = h.supplierUseCase.AddProductToSupplier(c.Request.Context(), uint(supplierID), uint(productID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *SupplierHandler) RemoveProductFromSupplier(c *gin.Context) {
	supplierID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid supplier id"})
		return
	}

	productID, err := strconv.ParseUint(c.Param("productId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid product id"})
		return
	}

	err = h.supplierUseCase.RemoveProductFromSupplier(c.Request.Context(), uint(supplierID), uint(productID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// Contract Management Handlers

func (h *SupplierHandler) CreateContract(c *gin.Context) {
	supplierID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid supplier id"})
		return
	}

	var contract entity.Contract
	if err := c.ShouldBindJSON(&contract); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	contract.SupplierID = uint(supplierID)

	if err := h.supplierUseCase.CreateContract(c.Request.Context(), &contract); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, contract)
}

func (h *SupplierHandler) UpdateContract(c *gin.Context) {
	contractID, err := strconv.ParseUint(c.Param("contractId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid contract id"})
		return
	}

	var contract entity.Contract
	if err := c.ShouldBindJSON(&contract); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	contract.ID = uint(contractID)

	if err := h.supplierUseCase.UpdateContract(c.Request.Context(), &contract); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, contract)
}

func (h *SupplierHandler) GetContract(c *gin.Context) {
	contractID, err := strconv.ParseUint(c.Param("contractId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid contract id"})
		return
	}

	contract, err := h.supplierUseCase.GetContract(c.Request.Context(), uint(contractID))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, contract)
}

// Rating Management Handlers

func (h *SupplierHandler) AddRating(c *gin.Context) {
	supplierID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid supplier id"})
		return
	}

	var rating entity.SupplierRating
	if err := c.ShouldBindJSON(&rating); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	rating.SupplierID = uint(supplierID)

	// Get user ID from context (assuming it's set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "user not authenticated"})
		return
	}
	rating.RatedBy = userID.(uint)

	if err := h.supplierUseCase.AddSupplierRating(c.Request.Context(), &rating); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, rating)
}

func (h *SupplierHandler) GetRatings(c *gin.Context) {
	supplierID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid supplier id"})
		return
	}

	ratings, err := h.supplierUseCase.GetSupplierRatings(c.Request.Context(), uint(supplierID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, ratings)
}
