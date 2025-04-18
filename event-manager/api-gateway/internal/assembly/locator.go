package assembly

import (
	"github.com/PabloPerdolie/event-manager/api-gateway/internal/config"
	"github.com/PabloPerdolie/event-manager/api-gateway/internal/handler"
	"github.com/PabloPerdolie/event-manager/api-gateway/internal/middleware"
	"github.com/PabloPerdolie/event-manager/api-gateway/internal/repository"
	"github.com/PabloPerdolie/event-manager/api-gateway/internal/routes"
	"github.com/PabloPerdolie/event-manager/api-gateway/internal/service/auth"
	"github.com/PabloPerdolie/event-manager/api-gateway/internal/service/comment"
	"github.com/PabloPerdolie/event-manager/api-gateway/internal/service/proxy"
	"github.com/PabloPerdolie/event-manager/api-gateway/pkg/postgres"
	rabbitmq "github.com/PabloPerdolie/event-manager/api-gateway/pkg/rabbitmq/publisher"
	redis1 "github.com/PabloPerdolie/event-manager/api-gateway/pkg/redis"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type ServiceLocator struct {
	Config      *config.Config
	RedisClient *redis.Client
	Controllers *routes.Controllers
	Middleware  *middleware.AuthMiddleware
	Logger      *zap.SugaredLogger
}

func NewServiceLocator(cfg *config.Config, logger *zap.SugaredLogger) (*ServiceLocator, error) {
	redisClient, err := redis1.NewClient(cfg.RedisURL)
	if err != nil {
		return nil, errors.WithMessage(err, "new redis client")
	}

	pbl, err := rabbitmq.New(cfg.CommentQueueName, cfg.RabbitMQURL)
	if err != nil {
		return nil, errors.WithMessage(err, "new publisher")
	}

	publisher := repository.NewPublisher(pbl)
	db, err := postgres.InitDB(logger, cfg.DatabaseURL)
	if err != nil {
		return nil, errors.WithMessage(err, "init db")
	}

	tokenCacheRepo := repository.NewClient(redisClient)
	userRepo := repository.NewPostgresRepository(db)

	commentService := comment.New(publisher)
	authService := auth.New(userRepo, &tokenCacheRepo, cfg.JWTSecretKey, cfg.JWTAccessExpiration, cfg.JWTRefreshExpiration, cfg.PasswordResetExpiration)
	proxyService := proxy.New(cfg.CoreServiceURL, logger)

	commentCtrl := handler.NewComment(commentService)
	authCtrl := handler.NewAuth(authService)
	proxyCtrl := handler.NewProxy(proxyService)

	controllers := routes.Controllers{
		AuthCtrl:    authCtrl,
		ProxyCtrl:   proxyCtrl,
		CommentCtrl: commentCtrl,
	}

	middleware := middleware.NewAuthMiddleware(authService)

	return &ServiceLocator{
		Config:      cfg,
		RedisClient: redisClient,
		Controllers: &controllers,
		Middleware:  middleware,
		Logger:      logger,
	}, nil
}

func (l *ServiceLocator) Close() {
	if l.RedisClient != nil {
		l.RedisClient.Close()
	}
}
