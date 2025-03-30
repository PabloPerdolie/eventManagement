package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/event-management/api-gateway/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	jwtSecret   string
	redisClient *redis.Client
}

func NewAuthMiddleware(jwtSecret string, redisClient *redis.Client) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret:   jwtSecret,
		redisClient: redisClient,
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

		if m.isTokenBlacklisted(c, tokenString) {
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

func (m *AuthMiddleware) isTokenBlacklisted(c *gin.Context, token string) bool {
	blacklistKey := fmt.Sprintf("blacklist:%s", token)

	result, err := m.redisClient.Exists(c, blacklistKey).Result()
	if err != nil {
		// If there's an error checking Redis, log it but allow the operation to continue
		return false
	}

	return result > 0
}

func (m *AuthMiddleware) BlacklistToken(ctx context.Context, token string, expiry time.Duration) error {
	blacklistKey := fmt.Sprintf("blacklist:%s", token)
	return m.redisClient.Set(ctx, blacklistKey, "revoked", expiry).Err()
}
