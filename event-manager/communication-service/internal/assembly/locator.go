package assembly

import (
	"github.com/PabloPerdolie/event-manager/communication-service/internal/config"
	"github.com/PabloPerdolie/event-manager/communication-service/internal/consumer"
	"github.com/PabloPerdolie/event-manager/communication-service/internal/handler"
	"github.com/PabloPerdolie/event-manager/communication-service/internal/routes"
	"github.com/PabloPerdolie/event-manager/communication-service/internal/service"
	"go.uber.org/zap"
)

type ServiceLocator struct {
	Controllers routes.Controllers
	Consumer    *consumer.RabbitMQConsumer
	logger      *zap.SugaredLogger
}

func NewLocator(cfg *config.Config, logger *zap.SugaredLogger) (*ServiceLocator, error) {

	healthService := service.New(*cfg)

	healthController := handler.New(healthService)

	//notifyConsumer := consumer.New(notifyController, cfg, logger)
	//go func() {
	//	if err := notifyConsumer.Start(); err != nil {
	//		logger.Fatalf("Failed to start RabbitMQ consumer: %v", err)
	//	}
	//}()

	controllers := routes.Controllers{
		HealthCtrl: *healthController,
	}

	return &ServiceLocator{
		Controllers: controllers,
		//Consumer:    notifyConsumer,
		logger: logger,
	}, nil
}

func (l *ServiceLocator) Close() {
	l.logger.Info("Cleaning up resources...")
	l.Consumer.Stop()
}
