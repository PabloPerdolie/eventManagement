package auth

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"time"

	"github.com/event-management/api-gateway/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Repository interface {
	CreateUser(ctx context.Context, user domain.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) error
	CheckEmailExists(ctx context.Context, email string) (bool, error)
}

type TokenStore interface {
	BlacklistToken(ctx context.Context, token string, expiry time.Duration) error
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
	StoreResetToken(ctx context.Context, userID uuid.UUID, token string, expiry time.Duration) error
	GetUserIDByResetToken(ctx context.Context, token string) (uuid.UUID, error)
	DeleteResetToken(ctx context.Context, token string) error
	DeleteRefreshToken(ctx context.Context, token string) error
}

type Service struct {
	repo                Repository
	tokenStore          TokenStore
	jwtSecret           string
	accessTokenExpiry   time.Duration
	refreshTokenExpiry  time.Duration
	passwordResetExpiry time.Duration
}

func New(repo Repository, tokenStore TokenStore, jwtSecret string, accessExp, refreshExp, resetExp time.Duration) Service {
	return Service{
		repo:                repo,
		tokenStore:          tokenStore,
		jwtSecret:           jwtSecret,
		accessTokenExpiry:   accessExp,
		refreshTokenExpiry:  refreshExp,
		passwordResetExpiry: resetExp,
	}
}

func (s Service) Register(ctx context.Context, req domain.UserRegisterRequest) (*domain.AuthResponse, error) {
	exists, err := s.repo.CheckEmailExists(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrEmailAlreadyExists
	}

	user := domain.User{
		ID:              uuid.New(),
		Email:           req.Email,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Role:            domain.RoleUser,
		IsActive:        true,
		IsEmailVerified: false,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = passwordHash

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	tokenPair, err := s.generateTokenPair(user.ID, string(user.Role))
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		User: domain.UserResponse{
			ID:              user.ID,
			Email:           user.Email,
			FirstName:       user.FirstName,
			LastName:        user.LastName,
			Role:            user.Role,
			IsEmailVerified: user.IsEmailVerified,
			CreatedAt:       user.CreatedAt,
		},
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}

// Login выполняет вход пользователя
func (s Service) Login(ctx context.Context, req domain.UserLoginRequest) (*domain.AuthResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	if !checkPasswordHash(req.Password, user.PasswordHash) {
		return nil, domain.ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, domain.ErrUserNotActive
	}

	tokenPair, err := s.generateTokenPair(user.ID, string(user.Role))
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		User: domain.UserResponse{
			ID:              user.ID,
			Email:           user.Email,
			FirstName:       user.FirstName,
			LastName:        user.LastName,
			Role:            user.Role,
			IsEmailVerified: user.IsEmailVerified,
			CreatedAt:       user.CreatedAt,
		},
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}

// RefreshToken обновляет токен доступа
func (s Service) RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenPair, error) {
	isBlacklisted, err := s.tokenStore.IsTokenBlacklisted(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	if isBlacklisted {
		return nil, domain.ErrInvalidToken
	}

	token, err := jwt.ParseWithClaims(refreshToken, &domain.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	claims, ok := token.Claims.(*domain.JWTClaims)
	if !ok || !token.Valid {
		return nil, domain.ErrInvalidToken
	}

	if err := s.tokenStore.BlacklistToken(ctx, refreshToken, s.refreshTokenExpiry); err != nil {
		return nil, err
	}

	user, err := s.repo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, domain.ErrUserNotActive
	}

	return s.generateTokenPair(user.ID, string(user.Role))
}

// Logout выполняет выход пользователя
func (s Service) Logout(ctx context.Context, token string) error {
	return s.tokenStore.BlacklistToken(ctx, token, s.refreshTokenExpiry)
}

// ValidateToken проверяет валидность токена доступа
func (s Service) ValidateToken(ctx context.Context, tokenString string) (*domain.JWTClaims, error) {
	isBlacklisted, err := s.tokenStore.IsTokenBlacklisted(ctx, tokenString)
	if err != nil {
		return nil, err
	}
	if isBlacklisted {
		return nil, domain.ErrInvalidToken
	}

	token, err := jwt.ParseWithClaims(tokenString, &domain.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	claims, ok := token.Claims.(*domain.JWTClaims)
	if !ok || !token.Valid {
		return nil, domain.ErrInvalidToken
	}

	return claims, nil
}

// GetUserInfo возвращает информацию о пользователе по его ID
func (s Service) GetUserInfo(ctx context.Context, userID uuid.UUID) (*domain.UserResponse, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &domain.UserResponse{
		ID:              user.ID,
		Email:           user.Email,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Role:            user.Role,
		IsEmailVerified: user.IsEmailVerified,
		CreatedAt:       user.CreatedAt,
	}, nil
}

// CreatePasswordResetToken создает токен для сброса пароля
func (s Service) CreatePasswordResetToken(ctx context.Context, email string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", domain.ErrUserNotFound
	}

	if !user.IsActive {
		return "", domain.ErrUserNotActive
	}

	resetToken := uuid.New().String()

	err = s.tokenStore.StoreResetToken(ctx, user.ID, resetToken, s.passwordResetExpiry)
	if err != nil {
		return "", err
	}

	return resetToken, nil
}

// generateTokenPair генерирует пару токенов (access и refresh)
func (s Service) generateTokenPair(userID uuid.UUID, role string) (*domain.TokenPair, error) {
	accessClaims := domain.JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "api-gateway",
			Subject:   userID.String(),
			ID:        uuid.New().String(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	refreshClaims := domain.JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "api-gateway",
			Subject:   userID.String(),
			ID:        uuid.New().String(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	return &domain.TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int64(s.accessTokenExpiry.Seconds()),
	}, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
