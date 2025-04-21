package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lugondev/erp-warehouse-simple/internal/application/usecase"
	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
)

// Request/Response models
type RegisterRequest struct {
	Username string `json:"username" binding:"required" example:"johnsmith"`
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"secret123"`
	RoleID   uint   `json:"role_id" binding:"required" example:"2"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required" example:"secret123"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJ..."`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email" example:"john@example.com"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required" example:"reset_token_123"`
	NewPassword string `json:"new_password" binding:"required,min=6" example:"newpass123"`
}

type UpdateUserRequest struct {
	Username string            `json:"username" binding:"required" example:"johnsmith"`
	Email    string            `json:"email" binding:"required,email" example:"john@example.com"`
	RoleID   uint              `json:"role_id" binding:"required" example:"2"`
	Status   entity.UserStatus `json:"status" binding:"required" example:"active"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"Error message"`
}

type MessageResponse struct {
	Message string `json:"message" example:"Operation successful"`
}

type LoginResponse struct {
	AccessToken  string      `json:"access_token" example:"eyJhbGciOiJ..."`
	RefreshToken string      `json:"refresh_token" example:"eyJhbGciOiJ..."`
	User         entity.User `json:"user"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJ..."`
}

type AuditLogResponse struct {
	Total int64             `json:"total" example:"100"`
	Logs  []entity.AuditLog `json:"logs"`
}

// Auth handlers
// @Summary Register a new user
// @Description Register a new user with the provided details
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "User registration details"
// @Success 201 {object} entity.User
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /auth/register [post]
func (s *Server) handleRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := s.userUC.CreateUser(&usecase.CreateUserInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		RoleID:   req.RoleID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// @Summary User login
// @Description Authenticate user and return tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "User login credentials"
// @Success 200 {object} LoginResponse "Login successful"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Router /auth/login [post]
func (s *Server) handleLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := s.userUC.ValidateCredentials(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := s.jwtService.GenerateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, expiry, err := s.jwtService.GenerateRefreshToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	if err := s.userUC.UpdateRefreshToken(user.ID, refreshToken, expiry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          user,
	})
}

// @Summary User logout
// @Description Invalidate user's refresh token
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} MessageResponse "Logged out successfully"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /users/logout [post]
func (s *Server) handleLogout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	if err := s.userUC.UpdateRefreshToken(userID.(uint), "", time.Time{}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to invalidate refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// @Summary Refresh access token
// @Description Get new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} TokenResponse "New access token"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 401 {object} ErrorResponse "Invalid refresh token"
// @Router /auth/refresh-token [post]
func (s *Server) handleRefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := s.userUC.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	accessToken, err := s.jwtService.GenerateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}

// @Summary Request password reset
// @Description Send password reset token to user's email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ForgotPasswordRequest true "User email"
// @Success 200 {object} MessageResponse "Reset email sent"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /auth/forgot-password [post]
func (s *Server) handleForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := s.userUC.RequestPasswordReset(req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset email sent"})
}

// @Summary Reset password
// @Description Reset user password using reset token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ResetPasswordRequest true "Reset token and new password"
// @Success 200 {object} MessageResponse "Password reset successful"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Router /auth/reset-password [post]
func (s *Server) handleResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.userUC.ResetPassword(req.Token, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
}
