package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lugondev/erp-warehouse-simple/internal/application/usecase"
	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
)

type CustomerHandler struct {
	customerUseCase usecase.CustomerUseCase
}

func NewCustomerHandler(customerUseCase usecase.CustomerUseCase) *CustomerHandler {
	return &CustomerHandler{
		customerUseCase: customerUseCase,
	}
}

// RegisterRoutes registers customer-related routes
func (h *CustomerHandler) RegisterRoutes(router *gin.RouterGroup) {
	customerRoutes := router.Group("/customers")
	{
		customerRoutes.POST("", h.CreateCustomer)
		customerRoutes.GET("", h.ListCustomers)
		customerRoutes.GET("/:id", h.GetCustomer)
		customerRoutes.PUT("/:id", h.UpdateCustomer)
		customerRoutes.DELETE("/:id", h.DeleteCustomer)

		// Address routes
		customerRoutes.POST("/:id/addresses", h.CreateAddress)
		customerRoutes.GET("/:id/addresses", h.GetAddresses)
		customerRoutes.PUT("/addresses/:addressId", h.UpdateAddress)
		customerRoutes.DELETE("/addresses/:addressId", h.DeleteAddress)

		// Order history routes
		customerRoutes.GET("/:id/orders", h.GetOrderHistory)

		// Debt management routes
		customerRoutes.GET("/:id/debt", h.GetCustomerDebt)
		customerRoutes.PUT("/:id/debt", h.UpdateCustomerDebt)

		// Loyalty management routes
		customerRoutes.PUT("/:id/loyalty/points", h.UpdateLoyaltyPoints)
		customerRoutes.PUT("/:id/loyalty/tier", h.UpdateLoyaltyTier)
		customerRoutes.GET("/:id/loyalty/calculate-tier", h.CalculateLoyaltyTier)
	}
}

// CreateCustomer godoc
// @Summary Create a new customer
// @Description Create a new customer with the provided details
// @Tags customers
// @Accept json
// @Produce json
// @Param customer body entity.Customer true "Customer details"
// @Success 201 {object} entity.Customer
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /customers [post]
func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var customer entity.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body: " + err.Error()})
		return
	}

	if err := h.customerUseCase.CreateCustomer(&customer); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create customer: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, customer)
}

// GetCustomer godoc
// @Summary Get a customer by ID
// @Description Get a customer's details by their ID
// @Tags customers
// @Produce json
// @Param id path int true "Customer ID"
// @Success 200 {object} entity.Customer
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /customers/{id} [get]
func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid customer ID"})
		return
	}

	customer, err := h.customerUseCase.GetCustomerByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Customer not found: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// UpdateCustomer godoc
// @Summary Update a customer
// @Description Update a customer's details
// @Tags customers
// @Accept json
// @Produce json
// @Param id path int true "Customer ID"
// @Param customer body entity.Customer true "Updated customer details"
// @Success 200 {object} entity.Customer
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /customers/{id} [put]
func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid customer ID"})
		return
	}

	var customer entity.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body: " + err.Error()})
		return
	}

	customer.ID = uint(id)
	if err := h.customerUseCase.UpdateCustomer(&customer); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update customer: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// DeleteCustomer godoc
// @Summary Delete a customer
// @Description Delete a customer by their ID
// @Tags customers
// @Produce json
// @Param id path int true "Customer ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /customers/{id} [delete]
func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid customer ID"})
		return
	}

	if err := h.customerUseCase.DeleteCustomer(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete customer: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListCustomers godoc
// @Summary List customers
// @Description List customers with optional filtering
// @Tags customers
// @Produce json
// @Param code query string false "Filter by code"
// @Param name query string false "Filter by name"
// @Param type query string false "Filter by type"
// @Param email query string false "Filter by email"
// @Param phone_number query string false "Filter by phone number"
// @Param loyalty_tier query string false "Filter by loyalty tier"
// @Param city query string false "Filter by city"
// @Param country query string false "Filter by country"
// @Success 200 {array} entity.Customer
// @Failure 500 {object} ErrorResponse
// @Router /customers [get]
func (h *CustomerHandler) ListCustomers(c *gin.Context) {
	filter := entity.CustomerFilter{
		Code:        c.Query("code"),
		Name:        c.Query("name"),
		Email:       c.Query("email"),
		PhoneNumber: c.Query("phone_number"),
		City:        c.Query("city"),
		Country:     c.Query("country"),
	}

	// Handle type filter
	if typeStr := c.Query("type"); typeStr != "" {
		customerType := entity.CustomerType(typeStr)
		filter.Type = &customerType
	}

	// Handle loyalty tier filter
	if tierStr := c.Query("loyalty_tier"); tierStr != "" {
		loyaltyTier := entity.CustomerLoyaltyTier(tierStr)
		filter.LoyaltyTier = &loyaltyTier
	}

	customers, err := h.customerUseCase.ListCustomers(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to list customers: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, customers)
}

// CreateAddress godoc
// @Summary Create a new address for a customer
// @Description Create a new address for a customer
// @Tags customers
// @Accept json
// @Produce json
// @Param id path int true "Customer ID"
// @Param address body entity.CustomerAddress true "Address details"
// @Success 201 {object} entity.CustomerAddress
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /customers/{id}/addresses [post]
func (h *CustomerHandler) CreateAddress(c *gin.Context) {
	customerID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid customer ID"})
		return
	}

	var address entity.CustomerAddress
	if err := c.ShouldBindJSON(&address); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body: " + err.Error()})
		return
	}

	address.CustomerID = uint(customerID)
	if err := h.customerUseCase.CreateAddress(&address); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create address: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, address)
}

// GetAddresses godoc
// @Summary Get all addresses for a customer
// @Description Get all addresses for a customer
// @Tags customers
// @Produce json
// @Param id path int true "Customer ID"
// @Success 200 {array} entity.CustomerAddress
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /customers/{id}/addresses [get]
func (h *CustomerHandler) GetAddresses(c *gin.Context) {
	customerID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid customer ID"})
		return
	}

	addresses, err := h.customerUseCase.GetAddressesByCustomerID(uint(customerID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get addresses: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, addresses)
}

// UpdateAddress godoc
// @Summary Update a customer address
// @Description Update a customer address
// @Tags customers
// @Accept json
// @Produce json
// @Param addressId path int true "Address ID"
// @Param address body entity.CustomerAddress true "Updated address details"
// @Success 200 {object} entity.CustomerAddress
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /customers/addresses/{addressId} [put]
func (h *CustomerHandler) UpdateAddress(c *gin.Context) {
	addressID, err := strconv.ParseUint(c.Param("addressId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid address ID"})
		return
	}

	var address entity.CustomerAddress
	if err := c.ShouldBindJSON(&address); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body: " + err.Error()})
		return
	}

	address.ID = uint(addressID)
	if err := h.customerUseCase.UpdateAddress(&address); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update address: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, address)
}

// DeleteAddress godoc
// @Summary Delete a customer address
// @Description Delete a customer address
// @Tags customers
// @Produce json
// @Param addressId path int true "Address ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /customers/addresses/{addressId} [delete]
func (h *CustomerHandler) DeleteAddress(c *gin.Context) {
	addressID, err := strconv.ParseUint(c.Param("addressId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid address ID"})
		return
	}

	if err := h.customerUseCase.DeleteAddress(uint(addressID)); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete address: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetOrderHistory godoc
// @Summary Get a customer's order history
// @Description Get a summary of a customer's order history
// @Tags customers
// @Produce json
// @Param id path int true "Customer ID"
// @Success 200 {object} entity.CustomerOrderHistory
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /customers/{id}/orders [get]
func (h *CustomerHandler) GetOrderHistory(c *gin.Context) {
	customerID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid customer ID"})
		return
	}

	history, err := h.customerUseCase.GetOrderHistory(uint(customerID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get order history: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}

// GetCustomerDebt godoc
// @Summary Get a customer's debt information
// @Description Get a customer's debt information
// @Tags customers
// @Produce json
// @Param id path int true "Customer ID"
// @Success 200 {object} entity.CustomerDebt
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /customers/{id}/debt [get]
func (h *CustomerHandler) GetCustomerDebt(c *gin.Context) {
	customerID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid customer ID"})
		return
	}

	debt, err := h.customerUseCase.GetCustomerDebt(uint(customerID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get customer debt: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, debt)
}

// UpdateCustomerDebt godoc
// @Summary Update a customer's debt amount
// @Description Update a customer's debt amount
// @Tags customers
// @Accept json
// @Produce json
// @Param id path int true "Customer ID"
// @Param debt body map[string]float64 true "Debt amount"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /customers/{id}/debt [put]
func (h *CustomerHandler) UpdateCustomerDebt(c *gin.Context) {
	customerID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid customer ID"})
		return
	}

	var requestBody struct {
		Amount float64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body: " + err.Error()})
		return
	}

	if err := h.customerUseCase.UpdateCustomerDebt(uint(customerID), requestBody.Amount); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update customer debt: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer debt updated successfully"})
}

// UpdateLoyaltyPoints godoc
// @Summary Update a customer's loyalty points
// @Description Update a customer's loyalty points
// @Tags customers
// @Accept json
// @Produce json
// @Param id path int true "Customer ID"
// @Param points body map[string]int true "Loyalty points"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /customers/{id}/loyalty/points [put]
func (h *CustomerHandler) UpdateLoyaltyPoints(c *gin.Context) {
	customerID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid customer ID"})
		return
	}

	var requestBody struct {
		Points int `json:"points" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body: " + err.Error()})
		return
	}

	if err := h.customerUseCase.UpdateLoyaltyPoints(uint(customerID), requestBody.Points); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update loyalty points: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Loyalty points updated successfully"})
}

// UpdateLoyaltyTier godoc
// @Summary Update a customer's loyalty tier
// @Description Update a customer's loyalty tier
// @Tags customers
// @Accept json
// @Produce json
// @Param id path int true "Customer ID"
// @Param tier body map[string]string true "Loyalty tier"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /customers/{id}/loyalty/tier [put]
func (h *CustomerHandler) UpdateLoyaltyTier(c *gin.Context) {
	customerID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid customer ID"})
		return
	}

	var requestBody struct {
		Tier string `json:"tier" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body: " + err.Error()})
		return
	}

	tier := entity.CustomerLoyaltyTier(requestBody.Tier)
	if err := h.customerUseCase.UpdateLoyaltyTier(uint(customerID), tier); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update loyalty tier: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Loyalty tier updated successfully"})
}

// CalculateLoyaltyTier godoc
// @Summary Calculate a customer's loyalty tier
// @Description Calculate a customer's loyalty tier based on points and purchase history
// @Tags customers
// @Produce json
// @Param id path int true "Customer ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /customers/{id}/loyalty/calculate-tier [get]
func (h *CustomerHandler) CalculateLoyaltyTier(c *gin.Context) {
	customerID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid customer ID"})
		return
	}

	tier, err := h.customerUseCase.CalculateLoyaltyTier(uint(customerID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to calculate loyalty tier: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tier": tier})
}
