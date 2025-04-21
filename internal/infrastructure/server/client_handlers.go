package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lugondev/erp-warehouse-simple/internal/application/usecase"
	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
	"github.com/lugondev/erp-warehouse-simple/internal/infrastructure/server/middleware"
)

type ClientHandler struct {
	clientUC usecase.ClientUseCase
}

func NewClientHandler(clientUC usecase.ClientUseCase) *ClientHandler {
	return &ClientHandler{
		clientUC: clientUC,
	}
}

// RegisterRoutes registers all routes for the client handler
func (h *ClientHandler) RegisterRoutes(rg *gin.RouterGroup) {
	clients := rg.Group("/clients")
	{
		clients.POST("", middleware.PermissionMiddleware(entity.ClientCreate), h.CreateClient)
		clients.GET("", middleware.PermissionMiddleware(entity.ClientRead), h.ListClients)
		clients.GET("/:id", middleware.PermissionMiddleware(entity.ClientRead), h.GetClient)
		clients.PUT("/:id", middleware.PermissionMiddleware(entity.ClientUpdate), h.UpdateClient)
		clients.DELETE("/:id", middleware.PermissionMiddleware(entity.ClientDelete), h.DeleteClient)

		// Address management
		addresses := clients.Group("/:id/addresses")
		{
			addresses.POST("", middleware.PermissionMiddleware(entity.ClientUpdate), h.CreateAddress)
			addresses.GET("", middleware.PermissionMiddleware(entity.ClientRead), h.GetAddresses)
			addresses.PUT("/:addressId", middleware.PermissionMiddleware(entity.ClientUpdate), h.UpdateAddress)
			addresses.DELETE("/:addressId", middleware.PermissionMiddleware(entity.ClientUpdate), h.DeleteAddress)
		}

		// Order history
		clients.GET("/:id/history", middleware.PermissionMiddleware(entity.ClientRead), h.GetOrderHistory)
	}
}

// @Summary Create a new client
// @Description Create a new client with the provided details
// @Tags clients
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param client body entity.Client true "Client details"
// @Success 201 {object} entity.Client
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /clients [post]
func (h *ClientHandler) CreateClient(c *gin.Context) {
	var client entity.Client
	if err := c.ShouldBindJSON(&client); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.clientUC.CreateClient(&client); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, client)
}

// @Summary Get a client by ID
// @Description Get a client's details by their ID
// @Tags clients
// @Security BearerAuth
// @Produce json
// @Param id path int true "Client ID"
// @Success 200 {object} entity.Client
// @Failure 400 {object} ErrorResponse "Invalid client ID"
// @Failure 404 {object} ErrorResponse "Client not found"
// @Router /clients/{id} [get]
func (h *ClientHandler) GetClient(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid client ID"})
		return
	}

	client, err := h.clientUC.GetClientByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, client)
}

// @Summary Update a client
// @Description Update a client's details
// @Tags clients
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param client body entity.Client true "Client details"
// @Success 200 {object} entity.Client
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /clients/{id} [put]
func (h *ClientHandler) UpdateClient(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid client ID"})
		return
	}

	var client entity.Client
	if err := c.ShouldBindJSON(&client); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	client.ID = uint(id)

	if err := h.clientUC.UpdateClient(&client); err != nil {
		statusCode := http.StatusInternalServerError
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, client)
}

// @Summary Delete a client
// @Description Delete a client by ID
// @Tags clients
// @Security BearerAuth
// @Produce json
// @Param id path int true "Client ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse "Invalid client ID"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /clients/{id} [delete]
func (h *ClientHandler) DeleteClient(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid client ID"})
		return
	}

	if err := h.clientUC.DeleteClient(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary List clients
// @Description List clients with optional filtering
// @Tags clients
// @Security BearerAuth
// @Produce json
// @Param code query string false "Client code"
// @Param name query string false "Client name"
// @Param email query string false "Client email"
// @Param phone_number query string false "Client phone number"
// @Param city query string false "Client city"
// @Param country query string false "Client country"
// @Param type query string false "Client type"
// @Param loyalty_tier query string false "Client loyalty tier"
// @Success 200 {array} entity.Client
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /clients [get]
func (h *ClientHandler) ListClients(c *gin.Context) {
	filter := entity.ClientFilter{
		Code:        c.Query("code"),
		Name:        c.Query("name"),
		Email:       c.Query("email"),
		PhoneNumber: c.Query("phone_number"),
		City:        c.Query("city"),
		Country:     c.Query("country"),
	}

	if clientType := c.Query("type"); clientType != "" {
		filter.Type = &clientType
	}

	if tier := c.Query("loyalty_tier"); tier != "" {
		loyaltyTier := entity.ClientLoyaltyTier(tier)
		filter.LoyaltyTier = &loyaltyTier
	}

	clients, err := h.clientUC.ListClients(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, clients)
}

// @Summary Create a client address
// @Description Create a new address for a client
// @Tags clients
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param address body entity.ClientAddress true "Address details"
// @Success 201 {object} entity.ClientAddress
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /clients/{id}/addresses [post]
func (h *ClientHandler) CreateAddress(c *gin.Context) {
	clientID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid client ID"})
		return
	}

	var address entity.ClientAddress
	if err := c.ShouldBindJSON(&address); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	address.ClientID = uint(clientID)

	if err := h.clientUC.CreateAddress(&address); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, address)
}

// @Summary Get client addresses
// @Description Get all addresses for a client
// @Tags clients
// @Security BearerAuth
// @Produce json
// @Param id path int true "Client ID"
// @Success 200 {array} entity.ClientAddress
// @Failure 400 {object} ErrorResponse "Invalid client ID"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /clients/{id}/addresses [get]
func (h *ClientHandler) GetAddresses(c *gin.Context) {
	clientID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid client ID"})
		return
	}

	addresses, err := h.clientUC.GetAddressesByClientID(uint(clientID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, addresses)
}

// @Summary Update a client address
// @Description Update an address for a client
// @Tags clients
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param addressId path int true "Address ID"
// @Param address body entity.ClientAddress true "Address details"
// @Success 200 {object} entity.ClientAddress
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /clients/{id}/addresses/{addressId} [put]
func (h *ClientHandler) UpdateAddress(c *gin.Context) {
	clientID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid client ID"})
		return
	}

	addressID, err := strconv.ParseUint(c.Param("addressId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid address ID"})
		return
	}

	var address entity.ClientAddress
	if err := c.ShouldBindJSON(&address); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	address.ID = uint(addressID)
	address.ClientID = uint(clientID)

	if err := h.clientUC.UpdateAddress(&address); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, address)
}

// @Summary Delete a client address
// @Description Delete an address for a client
// @Tags clients
// @Security BearerAuth
// @Produce json
// @Param id path int true "Client ID"
// @Param addressId path int true "Address ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse "Invalid address ID"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /clients/{id}/addresses/{addressId} [delete]
func (h *ClientHandler) DeleteAddress(c *gin.Context) {
	addressID, err := strconv.ParseUint(c.Param("addressId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid address ID"})
		return
	}

	if err := h.clientUC.DeleteAddress(uint(addressID)); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Get client order history
// @Description Get a client's order history
// @Tags clients
// @Security BearerAuth
// @Produce json
// @Param id path int true "Client ID"
// @Success 200 {array} entity.SalesOrder
// @Failure 400 {object} ErrorResponse "Invalid client ID"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /clients/{id}/history [get]
func (h *ClientHandler) GetOrderHistory(c *gin.Context) {
	clientID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid client ID"})
		return
	}

	history, err := h.clientUC.GetOrderHistory(uint(clientID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}
