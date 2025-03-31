package middleware

import (
	"context"
	"fmt"
	"github.com/PabloPerdolie/event-manager/api-gateway/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

type TokenCache interface {
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
}

type AuthMiddleware struct {
	jwtSecret  string
	tokenCache TokenCache
}

func NewAuthMiddleware(jwtSecret string, tokenCache TokenCache) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret:  jwtSecret,
		tokenCache: tokenCache,
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

		token, err := jwt.ParseWithClaims(tokenString, &domain.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(m.jwtSecret), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		isTokenBlacklisted, err := m.tokenCache.IsTokenBlacklisted(c, tokenString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}
		if isTokenBlacklisted {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*domain.JWTClaims); ok && token.Valid {
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
