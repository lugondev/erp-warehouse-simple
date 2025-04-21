package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lugondev/erp-warehouse-simple/internal/application/usecase"
)

// @Summary List users
// @Description Get a list of all users
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {array} entity.User
// @Failure 500 {object} ErrorResponse
// @Router /users [get]
func (s *Server) handleListUsers(c *gin.Context) {
	users, err := s.userUC.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	// Don't return passwords in response
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, users)
}

// @Summary Get user by ID
// @Description Get a user's details by their ID
// @Tags users
// @Security BearerAuth
// @Param id path int true "User ID"
// @Produce json
// @Success 200 {object} entity.User
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /users/{id} [get]
func (s *Server) handleGetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid user ID"})
		return
	}

	user, err := s.userUC.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	// Don't return sensitive data in response
	user.Password = ""
	user.RefreshToken = ""
	user.PasswordResetToken = ""
	c.JSON(http.StatusOK, user)
}

// @Summary Update user
// @Description Update a user's details
// @Tags users
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body UpdateUserRequest true "User details"
// @Produce json
// @Success 200 {object} entity.User
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [put]
func (s *Server) handleUpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid user ID"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	input := &usecase.UpdateUserInput{
		ID:       uint(id),
		Username: req.Username,
		Email:    req.Email,
		RoleID:   req.RoleID,
		Status:   req.Status,
	}

	user, err := s.userUC.UpdateUser(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	// Don't return sensitive data in response
	user.Password = ""
	user.RefreshToken = ""
	user.PasswordResetToken = ""
	c.JSON(http.StatusOK, user)
}

// @Summary Delete user
// @Description Delete a user
// @Tags users
// @Security BearerAuth
// @Param id path int true "User ID"
// @Produce json
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [delete]
func (s *Server) handleDeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid user ID"})
		return
	}

	if err := s.userUC.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
