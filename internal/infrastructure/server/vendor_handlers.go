package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lugondev/erp-warehouse-simple/internal/application/usecase"
	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
)

type VendorHandler struct {
	vendorUC *usecase.VendorUseCase
}

func NewVendorHandler(vendorUC *usecase.VendorUseCase) *VendorHandler {
	return &VendorHandler{
		vendorUC: vendorUC,
	}
}

// @Summary Create a new vendor
// @Description Create a new vendor with the provided details
// @Tags vendors
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param vendor body entity.Vendor true "Vendor details"
// @Success 201 {object} entity.Vendor
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /vendors [post]
func (h *VendorHandler) CreateVendor(c *gin.Context) {
	var vendor entity.Vendor
	if err := c.ShouldBindJSON(&vendor); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.vendorUC.CreateVendor(c.Request.Context(), &vendor); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, vendor)
}

// @Summary Get a vendor by ID
// @Description Get a vendor's details by its ID
// @Tags vendors
// @Security BearerAuth
// @Produce json
// @Param id path int true "Vendor ID"
// @Success 200 {object} entity.Vendor
// @Failure 400 {object} ErrorResponse "Invalid vendor ID"
// @Failure 404 {object} ErrorResponse "Vendor not found"
// @Router /vendors/{id} [get]
func (h *VendorHandler) GetVendor(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid vendor ID"})
		return
	}

	vendor, err := h.vendorUC.GetVendor(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, vendor)
}

// @Summary Update a vendor
// @Description Update a vendor's details
// @Tags vendors
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Vendor ID"
// @Param vendor body entity.Vendor true "Vendor details"
// @Success 200 {object} entity.Vendor
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /vendors/{id} [put]
func (h *VendorHandler) UpdateVendor(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid vendor ID"})
		return
	}

	var vendor entity.Vendor
	if err := c.ShouldBindJSON(&vendor); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	vendor.ID = uint(id)

	if err := h.vendorUC.UpdateVendor(c.Request.Context(), &vendor); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, vendor)
}

// @Summary Delete a vendor
// @Description Delete a vendor by ID
// @Tags vendors
// @Security BearerAuth
// @Produce json
// @Param id path int true "Vendor ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse "Invalid vendor ID"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /vendors/{id} [delete]
func (h *VendorHandler) DeleteVendor(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid vendor ID"})
		return
	}

	if err := h.vendorUC.DeleteVendor(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary List vendors
// @Description List vendors with optional filtering
// @Tags vendors
// @Security BearerAuth
// @Produce json
// @Param code query string false "Vendor code"
// @Param name query string false "Vendor name"
// @Param type query string false "Vendor type"
// @Param country query string false "Vendor country"
// @Param min_rating query number false "Minimum rating"
// @Param product_ids[] query array false "Product IDs"
// @Success 200 {array} entity.Vendor
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /vendors [get]
func (h *VendorHandler) ListVendors(c *gin.Context) {
	filter := entity.VendorFilter{
		Code:    c.Query("code"),
		Name:    c.Query("name"),
		Type:    c.Query("type"),
		Country: c.Query("country"),
	}

	if minRating := c.Query("min_rating"); minRating != "" {
		if rating, err := strconv.ParseFloat(minRating, 64); err == nil {
			filter.MinRating = &rating
		}
	}

	// Parse product IDs if provided
	if productIDsStr := c.QueryArray("product_ids[]"); len(productIDsStr) > 0 {
		productIDs := make([]uint, 0, len(productIDsStr))
		for _, idStr := range productIDsStr {
			if id, err := strconv.ParseUint(idStr, 10, 32); err == nil {
				productIDs = append(productIDs, uint(id))
			}
		}
		filter.ProductIDs = productIDs
	}

	vendors, err := h.vendorUC.ListVendors(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, vendors)
}

// @Summary Create a new product
// @Description Create a new product with the provided details
// @Tags vendors
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param product body entity.Product true "Product details"
// @Success 201 {object} entity.Product
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /vendors/products [post]
func (h *VendorHandler) CreateProduct(c *gin.Context) {
	var product entity.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.vendorUC.CreateProduct(c.Request.Context(), &product); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// @Summary Add product to vendor
// @Description Add a product to a vendor
// @Tags vendors
// @Security BearerAuth
// @Produce json
// @Param id path int true "Vendor ID"
// @Param productId path int true "Product ID"
// @Success 200 "OK"
// @Failure 400 {object} ErrorResponse "Invalid ID"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /vendors/{id}/products/{productId} [post]
func (h *VendorHandler) AddProductToVendor(c *gin.Context) {
	vendorID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid vendor ID"})
		return
	}

	productID, err := strconv.ParseUint(c.Param("productId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid product ID"})
		return
	}

	if err := h.vendorUC.AddProductToVendor(c.Request.Context(), uint(vendorID), uint(productID)); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// @Summary Remove product from vendor
// @Description Remove a product from a vendor
// @Tags vendors
// @Security BearerAuth
// @Produce json
// @Param id path int true "Vendor ID"
// @Param productId path int true "Product ID"
// @Success 200 "OK"
// @Failure 400 {object} ErrorResponse "Invalid ID"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /vendors/{id}/products/{productId} [delete]
func (h *VendorHandler) RemoveProductFromVendor(c *gin.Context) {
	vendorID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid vendor ID"})
		return
	}

	productID, err := strconv.ParseUint(c.Param("productId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid product ID"})
		return
	}

	if err := h.vendorUC.RemoveProductFromVendor(c.Request.Context(), uint(vendorID), uint(productID)); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// @Summary Create a vendor contract
// @Description Create a new contract for a vendor
// @Tags vendors
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Vendor ID"
// @Param contract body entity.Contract true "Contract details"
// @Success 201 {object} entity.Contract
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /vendors/{id}/contracts [post]
func (h *VendorHandler) CreateContract(c *gin.Context) {
	vendorID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid vendor ID"})
		return
	}

	var contract entity.Contract
	if err := c.ShouldBindJSON(&contract); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	contract.VendorID = uint(vendorID)

	if err := h.vendorUC.CreateContract(c.Request.Context(), &contract); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, contract)
}

// @Summary Update a contract
// @Description Update an existing contract
// @Tags vendors
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param contractId path int true "Contract ID"
// @Param contract body entity.Contract true "Contract details"
// @Success 200 {object} entity.Contract
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /vendors/contracts/{contractId} [put]
func (h *VendorHandler) UpdateContract(c *gin.Context) {
	contractID, err := strconv.ParseUint(c.Param("contractId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid contract ID"})
		return
	}

	var contract entity.Contract
	if err := c.ShouldBindJSON(&contract); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	contract.ID = uint(contractID)

	if err := h.vendorUC.UpdateContract(c.Request.Context(), &contract); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, contract)
}

// @Summary Get a contract by ID
// @Description Get a contract's details by its ID
// @Tags vendors
// @Security BearerAuth
// @Produce json
// @Param contractId path int true "Contract ID"
// @Success 200 {object} entity.Contract
// @Failure 400 {object} ErrorResponse "Invalid contract ID"
// @Failure 404 {object} ErrorResponse "Contract not found"
// @Router /vendors/contracts/{contractId} [get]
func (h *VendorHandler) GetContract(c *gin.Context) {
	contractID, err := strconv.ParseUint(c.Param("contractId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid contract ID"})
		return
	}

	contract, err := h.vendorUC.GetContract(c.Request.Context(), uint(contractID))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, contract)
}

// @Summary Add a vendor rating
// @Description Add a rating for a vendor
// @Tags vendors
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Vendor ID"
// @Param rating body entity.VendorRating true "Rating details"
// @Success 201 {object} entity.VendorRating
// @Failure 400 {object} ErrorResponse "Invalid input or rating"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /vendors/{id}/ratings [post]
func (h *VendorHandler) AddRating(c *gin.Context) {
	vendorID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid vendor ID"})
		return
	}

	var rating entity.VendorRating
	if err := c.ShouldBindJSON(&rating); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	rating.VendorID = uint(vendorID)
	rating.RatedByID = c.GetUint("user_id") // Assuming user_id is set in auth middleware

	if err := h.vendorUC.AddVendorRating(c.Request.Context(), &rating); err != nil {
		statusCode := http.StatusInternalServerError
		if err == entity.ErrInvalidRating {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, rating)
}

// @Summary Get vendor ratings
// @Description Get all ratings for a vendor
// @Tags vendors
// @Security BearerAuth
// @Produce json
// @Param id path int true "Vendor ID"
// @Success 200 {array} entity.VendorRating
// @Failure 400 {object} ErrorResponse "Invalid vendor ID"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /vendors/{id}/ratings [get]
func (h *VendorHandler) GetRatings(c *gin.Context) {
	vendorID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid vendor ID"})
		return
	}

	ratings, err := h.vendorUC.GetVendorRatings(c.Request.Context(), uint(vendorID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, ratings)
}
