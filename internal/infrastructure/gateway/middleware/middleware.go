package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lugondev/erp-warehouse-simple/internal/infrastructure/auth"
	"golang.org/x/time/rate"
)

// Logger returns a middleware that logs request details
func Logger() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("[API-GATEWAY] %s | %d | %s | %s | %s %s | %s\n",
				param.TimeStamp.Format(time.RFC3339),
				param.StatusCode,
				param.Latency.String(),
				param.ClientIP,
				param.Method,
				param.Path,
				param.ErrorMessage,
			)
		},
		Output:    gin.DefaultWriter,
		SkipPaths: []string{"/health"},
	})
}

// CORS returns a middleware that handles CORS
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// Auth returns a middleware that handles authentication
func Auth(jwtService *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		tokenString := auth.ExtractTokenFromHeader(authHeader)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		claims, err := jwtService.ValidateAccessToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set user details to context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("permissions", claims.Permissions)

		c.Next()
	}
}

// RateLimit returns a middleware that limits request rate
func RateLimit(requestsPerSecond, burst int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(requestsPerSecond), burst)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// CircuitBreaker returns a middleware that implements circuit breaking
func CircuitBreaker() gin.HandlerFunc {
	// This is a simplified implementation
	// In a production environment, you would use a more robust circuit breaker
	// like github.com/sony/gobreaker or github.com/afex/hystrix-go
	return func(c *gin.Context) {
		// Circuit breaker logic would go here
		c.Next()
	}
}

// Tracing returns a middleware that adds tracing capabilities
func Tracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate a unique request ID
		requestID := uuid.New().String()
		c.Set("request_id", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)

		// Add the request ID to all outgoing requests
		c.Next()
	}
}
