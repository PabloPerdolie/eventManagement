package assembly

import (
	"github.com/PabloPerdolie/event-manager/communication-service/internal/config"
	"github.com/PabloPerdolie/event-manager/communication-service/internal/consumer"
	"github.com/PabloPerdolie/event-manager/communication-service/internal/handler"
	"github.com/PabloPerdolie/event-manager/communication-service/internal/repository"
	"github.com/PabloPerdolie/event-manager/communication-service/internal/routes"
	"github.com/PabloPerdolie/event-manager/communication-service/internal/service"
	"github.com/PabloPerdolie/event-manager/communication-service/pkg/postgres"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type ServiceLocator struct {
	Controllers routes.Controllers
	Consumer    *consumer.RabbitMQConsumer
	logger      *zap.SugaredLogger
}

func NewLocator(cfg *config.Config, logger *zap.SugaredLogger) (*ServiceLocator, error) {
	db, err := postgres.InitDB(logger, cfg.DatabaseURL)
	if err != nil {
		return nil, errors.WithMessage(err, "init db")
	}

	commentRepo := repository.New(db)

	healthService := service.New(*cfg)
	commentService := service.NewComment(commentRepo, logger)

	healthController := handler.New(healthService)
	commentController := handler.NewComment(commentService, logger)

	commentHandler := handler.NewCommentHandler(commentService, logger)
	commentConsumer := consumer.New(commentHandler, cfg, logger)

	go func() {
		if err := commentConsumer.Start(); err != nil {
			logger.Fatalw("Failed to start RabbitMQ consumer", "error", err)
		}
	}()

	controllers := routes.Controllers{
		HealthCtrl:  *healthController,
		CommentCtrl: commentController,
	}

	return &ServiceLocator{
		Controllers: controllers,
		Consumer:    commentConsumer,
		logger:      logger,
	}, nil
}

func (l *ServiceLocator) Close() {
	l.logger.Info("Cleaning up resources...")
	if l.Consumer != nil {
		l.Consumer.Stop()
	}
}
