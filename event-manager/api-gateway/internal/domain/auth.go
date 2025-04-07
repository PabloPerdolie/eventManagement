package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

// UserRegisterRequest содержит данные для регистрации нового пользователя
// @Description Запрос на регистрацию нового пользователя
type UserRegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Username string `json:"username" binding:"required" example:"johndoe"`
	Password string `json:"password" binding:"required,min=8" example:"password123"`
}

// UserLoginRequest содержит учетные данные для входа пользователя
// @Description Запрос на вход пользователя в систему
type UserLoginRequest struct {
	Username string `json:"username" binding:"required" example:"johndoe"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// TokenPair содержит пару токенов доступа и обновления
// @Description Пара токенов для аутентификации пользователя
type TokenPair struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresIn    int64  `json:"expires_in" example:"3600"`
}

// UserResponse содержит информацию о пользователе
// @Description Информация о пользователе
type UserResponse struct {
	UserId    int       `json:"id" example:"1"`
	Username  string    `json:"username" example:"johndoe"`
	Email     string    `json:"email" example:"user@example.com"`
	Role      UserRole  `json:"role" example:"user"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
}

// AuthResponse содержит информацию о пользователе и токены
// @Description Ответ на успешную аутентификацию
type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string       `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresIn    int64        `json:"expires_in" example:"3600"`
}

// RefreshTokenRequest содержит токен обновления
// @Description Запрос на обновление токена доступа
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// JWTClaims содержит данные JWT токена
// @Description Утверждения в JWT токене
type JWTClaims struct {
	UserId int    `json:"user_id" example:"1"`
	Role   string `json:"role" example:"user"`
	jwt.RegisteredClaims
}
