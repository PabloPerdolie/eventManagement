package assembly

import (
	"github.com/PabloPerdolie/event-manager/notification-service/internal/config"
	"github.com/PabloPerdolie/event-manager/notification-service/internal/handler"
	"github.com/PabloPerdolie/event-manager/notification-service/internal/repository"
	"github.com/PabloPerdolie/event-manager/notification-service/internal/routes"
	"github.com/PabloPerdolie/event-manager/notification-service/internal/service"
	"github.com/PabloPerdolie/event-manager/notification-service/pkg/rabbitmq"
	"go.uber.org/zap"
)

type ServiceLocator struct {
	Controllers routes.Controllers
	Consumer    *rabbitmq.RabbitMQConsumer
	logger      *zap.SugaredLogger
}

func NewLocator(cfg *config.Config, logger *zap.SugaredLogger) (*ServiceLocator, error) {
	notifyRepo := repository.New(cfg.SMTP)

	notifyService := service.New(cfg.SMTP, notifyRepo)

	healthController := handler.New(notifyService)
	notifyController := handler.NewNotifyHandler(logger, notifyService)

	notifyConsumer := rabbitmq.New(notifyController, cfg, logger)
	go func() {
		if err := notifyConsumer.Start(); err != nil {
			logger.Fatalf("Failed to start RabbitMQ consumer: %v", err)
		}
	}()

	controllers := routes.Controllers{
		HealthCtrl: *healthController,
	}

	return &ServiceLocator{
		Controllers: controllers,
		Consumer:    notifyConsumer,
		logger:      logger,
	}, nil
}

func (l *ServiceLocator) Close() {
	l.logger.Info("Cleaning up resources...")
	l.Consumer.Stop()
}
