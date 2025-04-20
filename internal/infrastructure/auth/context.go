package auth

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUserIDFromContext extracts the user ID from the Gin context
func GetUserIDFromContext(c *gin.Context) string {
	userID, exists := c.Get("user_id")
	if !exists {
		return ""
	}

	// Convert to string regardless of the type
	switch v := userID.(type) {
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case int:
		return strconv.Itoa(v)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

// GetUsernameFromContext extracts the username from the Gin context
func GetUsernameFromContext(c *gin.Context) string {
	username, exists := c.Get("username")
	if !exists {
		return ""
	}
	return username.(string)
}

// GetRoleFromContext extracts the role from the Gin context
func GetRoleFromContext(c *gin.Context) string {
	role, exists := c.Get("role")
	if !exists {
		return ""
	}
	return role.(string)
}
