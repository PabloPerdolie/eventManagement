package auth

import (
	"context"
	"github.com/PabloPerdolie/event-manager/api-gateway/internal/model"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"time"

	"github.com/PabloPerdolie/event-manager/api-gateway/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Repository interface {
	CreateUser(ctx context.Context, user model.User) (int, error)
	GetUserById(ctx context.Context, id int) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	UpdateUser(ctx context.Context, user model.User) error
	DeleteUser(ctx context.Context, id int) error
	ListUsers(ctx context.Context, limit, offset int) ([]model.User, int, error)
}

type TokenRepo interface {
	BlacklistToken(ctx context.Context, token string, expiry time.Duration) error
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
	StoreResetToken(ctx context.Context, userId int, token string, expiry time.Duration) error
	GetUserIdByResetToken(ctx context.Context, token string) (int, error)
	DeleteResetToken(ctx context.Context, token string) error
	DeleteRefreshToken(ctx context.Context, token string) error
}

type Service struct {
	repo                Repository
	tokenStore          TokenRepo
	jwtSecret           string
	accessTokenExpiry   time.Duration
	refreshTokenExpiry  time.Duration
	passwordResetExpiry time.Duration
}

func New(repo Repository, tokenStore TokenRepo, jwtSecret string, accessExp, refreshExp, resetExp time.Duration) Service {
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
	exists, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return nil, errors.WithMessage(err, "get user by email")
	}
	if exists != nil {
		return nil, domain.ErrEmailAlreadyExists
	}

	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		return nil, errors.WithMessage(err, "hash password")
	}

	user := model.User{
		Email:        req.Email,
		PasswordHash: passwordHash,
		Username:     req.Username,
		Role:         domain.RoleUser,
		CreatedAt:    time.Now(),
	}

	id, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, errors.WithMessage(err, "create user")
	}

	tokenPair, err := s.generateTokenPair(id, string(user.Role))
	if err != nil {
		return nil, errors.WithMessage(err, "generate token pair")
	}

	return &domain.AuthResponse{
		User: domain.UserResponse{
			UserId:    id,
			Email:     user.Email,
			Username:  user.Username,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		},
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}

func (s Service) Login(ctx context.Context, req domain.UserLoginRequest) (*domain.AuthResponse, error) {
	user, err := s.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	if !checkPasswordHash(req.Password, user.PasswordHash) {
		return nil, domain.ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, domain.ErrUserNotActive
	}

	tokenPair, err := s.generateTokenPair(user.UserId, string(user.Role))
	if err != nil {
		return nil, errors.WithMessage(err, "generate token pair")
	}

	return &domain.AuthResponse{
		User: domain.UserResponse{
			UserId:    user.UserId,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
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
		return nil, errors.WithMessage(err, "is token blacklisted")
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
		return nil, errors.WithMessage(err, "blacklist token")
	}

	user, err := s.repo.GetUserById(ctx, claims.UserId)
	if err != nil {
		return nil, errors.WithMessage(err, "get user by id")
	}

	if !user.IsActive {
		return nil, domain.ErrUserNotActive
	}

	return s.generateTokenPair(user.UserId, string(user.Role))
}

// Logout выполняет выход пользователя
func (s Service) Logout(ctx context.Context, token string) error {
	return s.tokenStore.BlacklistToken(ctx, token, s.refreshTokenExpiry)
}

// ValidateToken проверяет валидность токена доступа
func (s Service) ValidateToken(ctx context.Context, tokenString string) (*domain.JWTClaims, error) {
	isBlacklisted, err := s.tokenStore.IsTokenBlacklisted(ctx, tokenString)
	if err != nil {
		return nil, errors.WithMessage(err, "is token black listed")
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

// GetUserInfo возвращает информацию о пользователе по его Id
func (s Service) GetUserInfo(ctx context.Context, userId int) (*domain.UserResponse, error) {
	user, err := s.repo.GetUserById(ctx, userId)
	if err != nil {
		return nil, errors.WithMessage(err, "get user by id")
	}

	return &domain.UserResponse{
		UserId:    user.UserId,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
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

	err = s.tokenStore.StoreResetToken(ctx, user.UserId, resetToken, s.passwordResetExpiry)
	if err != nil {
		return "", err
	}

	return resetToken, nil
}

// generateTokenPair генерирует пару токенов (access и refresh)
func (s Service) generateTokenPair(userId int, role string) (*domain.TokenPair, error) {
	accessClaims := domain.JWTClaims{
		UserId: userId,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "api-gateway",
			Subject:   string(rune(userId)),
			ID:        uuid.New().String(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, errors.WithMessage(err, "signed string")
	}

	refreshClaims := domain.JWTClaims{
		UserId: userId,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "api-gateway",
			Subject:   string(rune(userId)),
			ID:        uuid.New().String(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, errors.WithMessage(err, "signed string")
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
