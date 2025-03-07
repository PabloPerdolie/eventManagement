package service

import (
	"github.com/event-management/api-gateway/internal/config"
	"github.com/event-management/api-gateway/internal/repository"
	"github.com/event-management/api-gateway/internal/service/auth"
	"github.com/event-management/api-gateway/internal/service/proxy"
	"github.com/event-management/api-gateway/internal/service/user"
	"github.com/event-management/api-gateway/internal/storage/redis"
	"go.uber.org/zap"
)

// Service contains all the services for the application
type Service struct {
	Auth     auth.Service
	User     *user.Service
	Proxy    *proxy.Service
	logger   *zap.SugaredLogger
}

// New creates a new service
func New(repo *repository.Repository, redisClient *redis.Client, cfg *config.Config, logger *zap.SugaredLogger) *Service {
	authService := auth.New(repo.User, redisClient, cfg.JWTSecretKey, cfg.JWTAccessExpiration, cfg.JWTRefreshExpiration, cfg.PasswordResetExpiration)
	userService := user.New(repo.User, logger)
	proxyService := proxy.New(cfg.CoreServiceURL, cfg.NotificationServiceURL, cfg.CommunicationServiceURL, logger)

	return &Service{
		Auth:     authService,
		User:     userService,
		Proxy:    proxyService,
		logger:   logger,
	}
}
