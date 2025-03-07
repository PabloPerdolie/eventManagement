package assembly

import (
	"github.com/event-management/api-gateway/internal/config"
	"github.com/event-management/api-gateway/internal/handler"
	"github.com/event-management/api-gateway/internal/middleware"
	"github.com/event-management/api-gateway/internal/repository"
	"github.com/event-management/api-gateway/internal/service"
	redis1 "github.com/event-management/api-gateway/internal/storage/redis"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// ServiceLocator содержит все слои приложения
type ServiceLocator struct {
	Config      *config.Config
	RedisClient *redis.Client
	Handler     *handler.Handler
	Middleware  *middleware.AuthMiddleware
	Logger      *zap.SugaredLogger
}

func NewServiceLocator(cfg *config.Config, logger *zap.SugaredLogger) (*ServiceLocator, error) {
	repo, err := repository.New(cfg)
	if err != nil {
		return nil, errors.WithMessage(err, "new repository")
	}

	client, err := redis1.NewClient(cfg.RedisURL)
	if err != nil {
		return nil, err
	}

	authService := service.New(repo, client, cfg, logger)
	handler := handler.New(authService, cfg, logger)
	middleware := middleware.NewAuthMiddleware(cfg.JWTSecretKey, client)

	return &ServiceLocator{
		Config: cfg,
		RedisClient: client,
		Handler: handler,
		Middleware: middleware,
		Logger: logger,
	}, nil
}

func (l *ServiceLocator) Close() {
	//if l.Repository != nil {
	//	l.Repository.Close()
	//}
	if l.RedisClient != nil {
		l.RedisClient.Close()
	}
}

