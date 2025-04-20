package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lugondev/erp-warehouse-simple/internal/domain/entity"
)

type JWTService struct {
	accessTokenSecret  []byte
	refreshTokenSecret []byte
}

type Claims struct {
	jwt.RegisteredClaims
	UserID      uint                `json:"user_id"`
	Username    string              `json:"username"`
	Role        string              `json:"role"`
	Permissions []entity.Permission `json:"permissions"`
}

func NewJWTService(accessSecret, refreshSecret string) *JWTService {
	return &JWTService{
		accessTokenSecret:  []byte(accessSecret),
		refreshTokenSecret: []byte(refreshSecret),
	}
}

func (s *JWTService) GenerateAccessToken(user *entity.User) (string, error) {
	if !user.IsActive() {
		return "", errors.New("inactive user cannot generate token")
	}

	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID:      user.ID,
		Username:    user.Username,
		Role:        user.Role.Name,
		Permissions: user.Role.Permissions,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.accessTokenSecret)
}

func (s *JWTService) GenerateRefreshToken(user *entity.User) (string, time.Time, error) {
	expiry := time.Now().Add(7 * 24 * time.Hour) // 7 days

	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, err := token.SignedString(s.refreshTokenSecret)
	if err != nil {
		return "", time.Time{}, err
	}

	return refreshToken, expiry, nil
}

func (s *JWTService) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.accessTokenSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *JWTService) ValidateRefreshToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.refreshTokenSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid refresh token")
}

// ExtractTokenFromHeader extracts the token from the Authorization header
func ExtractTokenFromHeader(authHeader string) string {
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}
