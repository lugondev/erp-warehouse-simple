package service

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
)

type AuditService struct {
	repo entity.AuditLogRepository
}

func NewAuditService(repo entity.AuditLogRepository) *AuditService {
	return &AuditService{repo: repo}
}

// LogUserAction creates an audit log entry for a user action
func (s *AuditService) LogUserAction(ctx context.Context, userID uint, action entity.ActionType, resource, detail string) error {
	c, ok := ctx.(*gin.Context)
	if !ok {
		c = &gin.Context{}
	}

	log := &entity.AuditLog{
		UserID:    userID,
		Action:    action,
		Resource:  resource,
		Detail:    detail,
		IP:        c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		CreatedAt: time.Now(),
	}

	return s.repo.Create(log)
}

// GetUserAuditLogs retrieves audit logs for a specific user
func (s *AuditService) GetUserAuditLogs(userID uint, page, pageSize int) ([]entity.AuditLog, error) {
	offset := (page - 1) * pageSize
	return s.repo.FindByUserID(userID, pageSize, offset)
}

// GetAuditLogsByAction retrieves audit logs for a specific action type
func (s *AuditService) GetAuditLogsByAction(action entity.ActionType, page, pageSize int) ([]entity.AuditLog, error) {
	offset := (page - 1) * pageSize
	return s.repo.FindByAction(action, pageSize, offset)
}

// GetAuditLogsByDateRange retrieves audit logs within a date range
func (s *AuditService) GetAuditLogsByDateRange(start, end time.Time, page, pageSize int) ([]entity.AuditLog, error) {
	offset := (page - 1) * pageSize
	return s.repo.FindByDateRange(start, end, pageSize, offset)
}

// ListAuditLogs retrieves all audit logs with pagination
func (s *AuditService) ListAuditLogs(page, pageSize int) ([]entity.AuditLog, error) {
	offset := (page - 1) * pageSize
	return s.repo.List(pageSize, offset)
}

// GetAuditLogsCount returns the total count of audit logs based on filter
func (s *AuditService) GetAuditLogsCount(filter map[string]interface{}) (int64, error) {
	return s.repo.Count(filter)
}

// CreateAuditLogMiddleware creates a middleware that logs user actions
func CreateAuditLogMiddleware(auditService *AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.Next()
			return
		}

		// Get the start time
		start := time.Now()

		// Process request
		c.Next()

		// After request
		latency := time.Since(start)

		// Determine action type based on request method
		var action entity.ActionType
		switch c.Request.Method {
		case "GET":
			action = entity.ActionRead
		case "POST":
			action = entity.ActionCreate
		case "PUT", "PATCH":
			action = entity.ActionUpdate
		case "DELETE":
			action = entity.ActionDelete
		}

		// Create audit log
		detail := fmt.Sprintf("Method: %s, Path: %s, Status: %d, Latency: %v",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			latency,
		)

		_ = auditService.LogUserAction(c,
			userID.(uint),
			action,
			c.Request.URL.Path,
			detail,
		)
	}
}
