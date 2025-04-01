package middleware

import (
	"context"
	"github.com/PabloPerdolie/event-manager/api-gateway/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type AuthService interface {
	ValidateToken(ctx context.Context, tokenString string) (*domain.JWTClaims, error)
}

type AuthMiddleware struct {
	authService AuthService
}

func NewAuthMiddleware(authService AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := m.authService.ValidateToken(c, tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		if claims != nil {
			c.Set("user_id", claims.UserId)
			c.Set("user_role", claims.Role)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
	}
}

func (m *AuthMiddleware) AuthenticateAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		m.Authenticate()(c)

		if c.IsAborted() {
			return
		}

		role, exists := c.Get("user_role")
		if !exists || role.(domain.UserRole) != domain.RoleAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin privileges required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
