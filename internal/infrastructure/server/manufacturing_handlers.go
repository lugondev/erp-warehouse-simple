package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lugondev/erp-warehouse-simple/internal/application/usecase"
	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
)

type ManufacturingHandler struct {
	manufacturingUseCase *usecase.ManufacturingUseCase
}

func NewManufacturingHandler(manufacturingUseCase *usecase.ManufacturingUseCase) *ManufacturingHandler {
	return &ManufacturingHandler{
		manufacturingUseCase: manufacturingUseCase,
	}
}

// @Summary Create manufacturing facility
// @Description Create new manufacturing facility
// @Tags Manufacturing
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param facility body entity.ManufacturingFacility true "Facility details"
// @Success 201 {object} entity.ManufacturingFacility
// @Router /manufacturing/facilities [post]
func (h *ManufacturingHandler) CreateFacility(c *gin.Context) {
	var facility entity.ManufacturingFacility
	if err := c.ShouldBindJSON(&facility); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.manufacturingUseCase.CreateFacility(c.Request.Context(), &facility); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, facility)
}

// @Summary Get facility by ID
// @Description Get manufacturing facility details by ID
// @Tags Manufacturing
// @Security BearerAuth
// @Produce json
// @Param id path int true "Facility ID"
// @Success 200 {object} entity.ManufacturingFacility
// @Router /manufacturing/facilities/{id} [get]
func (h *ManufacturingHandler) GetFacility(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	facility, err := h.manufacturingUseCase.GetFacility(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "facility not found"})
		return
	}

	c.JSON(http.StatusOK, facility)
}

// @Summary List facilities
// @Description Get list of all manufacturing facilities
// @Tags Manufacturing
// @Security BearerAuth
// @Produce json
// @Success 200 {array} entity.ManufacturingFacility
// @Router /manufacturing/facilities [get]
func (h *ManufacturingHandler) ListFacilities(c *gin.Context) {
	facilities, err := h.manufacturingUseCase.ListFacilities(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, facilities)
}

// @Summary Create production order
// @Description Create new production order
// @Tags Manufacturing
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param order body entity.ProductionOrder true "Production order details"
// @Success 201 {object} entity.ProductionOrder
// @Router /manufacturing/orders [post]
func (h *ManufacturingHandler) CreateProductionOrder(c *gin.Context) {
	var order entity.ProductionOrder
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.manufacturingUseCase.CreateProductionOrder(c.Request.Context(), &order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// @Summary Start production
// @Description Start production for an order
// @Tags Manufacturing
// @Security BearerAuth
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {string} string "Production started"
// @Router /manufacturing/orders/{id}/start [post]
func (h *ManufacturingHandler) StartProduction(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	if err := h.manufacturingUseCase.StartProduction(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "production started"})
}

// @Summary Update production progress
// @Description Update production progress for an order
// @Tags Manufacturing
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param progress body map[string]int true "Progress details"
// @Success 200 {string} string "Progress updated"
// @Router /manufacturing/orders/{id}/progress [put]
func (h *ManufacturingHandler) UpdateProductionProgress(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	var progress struct {
		CompletedQty int `json:"completed_qty"`
		DefectQty    int `json:"defect_qty"`
	}

	if err := c.ShouldBindJSON(&progress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.manufacturingUseCase.UpdateProductionProgress(
		c.Request.Context(),
		uint(id),
		progress.CompletedQty,
		progress.DefectQty,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "progress updated"})
}

// @Summary Create BOM
// @Description Create bill of materials for a product
// @Tags Manufacturing
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param bom body map[string]interface{} true "BOM details with items"
// @Success 201 {object} entity.BillOfMaterial
// @Router /manufacturing/bom [post]
func (h *ManufacturingHandler) CreateBOM(c *gin.Context) {
	var request struct {
		BOM   entity.BillOfMaterial `json:"bom"`
		Items []entity.BOMItem      `json:"items"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.manufacturingUseCase.CreateBOM(c.Request.Context(), &request.BOM, request.Items); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, request.BOM)
}
