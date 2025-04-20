package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Get user audit logs
// @Description Get audit logs for a specific user
// @Tags audit
// @Security BearerAuth
// @Produce json
// @Param id path int true "User ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {array} entity.AuditLog
// @Failure 400 {object} ErrorResponse "Invalid user ID"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /api/v1/audit/logs/user/{id} [get]
func (s *Server) handleUserAuditLogs(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	logs, err := s.auditService.GetUserAuditLogs(uint(userID), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// @Summary List all audit logs
// @Description Get a paginated list of all audit logs
// @Tags audit
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} AuditLogResponse
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /api/v1/audit/logs [get]
func (s *Server) handleListAuditLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	logs, err := s.auditService.ListAuditLogs(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	count, err := s.auditService.GetAuditLogsCount(nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": count,
		"logs":  logs,
	})
}
