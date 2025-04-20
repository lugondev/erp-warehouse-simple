package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lugondev/erp-warehouse-simple/internal/application/usecase"
	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
)

type CreateRoleRequest struct {
	Name        string              `json:"name" binding:"required" example:"manager"`
	Permissions []entity.Permission `json:"permissions" binding:"required" example:"[\"user:read\",\"user:create\"]"`
}

type UpdateRoleRequest struct {
	Name        string              `json:"name" binding:"required" example:"manager"`
	Permissions []entity.Permission `json:"permissions" binding:"required" example:"[\"user:read\",\"user:create\"]"`
}

// @Summary Create new role
// @Description Create a new role with specified permissions
// @Tags roles
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateRoleRequest true "Role details"
// @Success 201 {object} entity.Role
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /api/v1/roles [post]
func (s *Server) handleCreateRole(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !s.roleUC.ValidatePermissions(req.Permissions) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid permissions"})
		return
	}

	role, err := s.roleUC.CreateRole(&usecase.CreateRoleInput{
		Name:        req.Name,
		Permissions: req.Permissions,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, role)
}

// @Summary List all roles
// @Description Get a list of all roles
// @Tags roles
// @Security BearerAuth
// @Produce json
// @Success 200 {array} entity.Role
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /api/v1/roles [get]
func (s *Server) handleListRoles(c *gin.Context) {
	roles, err := s.roleUC.ListRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, roles)
}

// @Summary Get role by ID
// @Description Get role details by role ID
// @Tags roles
// @Security BearerAuth
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} entity.Role
// @Failure 400 {object} ErrorResponse "Invalid role ID"
// @Failure 404 {object} ErrorResponse "Role not found"
// @Router /api/v1/roles/{id} [get]
func (s *Server) handleGetRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	role, err := s.roleUC.GetRoleByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	c.JSON(http.StatusOK, role)
}

// @Summary Update role
// @Description Update role details and permissions
// @Tags roles
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Param request body UpdateRoleRequest true "Role details to update"
// @Success 200 {object} entity.Role
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /api/v1/roles/{id} [put]
func (s *Server) handleUpdateRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !s.roleUC.ValidatePermissions(req.Permissions) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid permissions"})
		return
	}

	role, err := s.roleUC.UpdateRole(&usecase.UpdateRoleInput{
		ID:          uint(id),
		Name:        req.Name,
		Permissions: req.Permissions,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, role)
}

// @Summary Delete role
// @Description Delete role by ID
// @Tags roles
// @Security BearerAuth
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} MessageResponse "Role deleted successfully"
// @Failure 400 {object} ErrorResponse "Invalid role ID"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /api/v1/roles/{id} [delete]
func (s *Server) handleDeleteRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	if err := s.roleUC.DeleteRole(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}
