package assembly

import (
	"github.com/PabloPerdolie/event-manager/notification-service/internal/config"
	"github.com/PabloPerdolie/event-manager/notification-service/internal/handler"
	"github.com/PabloPerdolie/event-manager/notification-service/internal/service"
	"go.uber.org/zap"
)

type ServiceLocator struct {
	Service *service.Service
	Handler *handler.Handler
	logger  *zap.SugaredLogger
}

func NewServiceLocator(cfg *config.Config, logger *zap.SugaredLogger) (*ServiceLocator, error) {
	svc := service.New(cfg, logger)

	h := handler.New(svc, logger)

	return &ServiceLocator{
		Service: svc,
		Handler: h,
		logger:  logger,
	}, nil
}

func (l *ServiceLocator) Close() {

	l.logger.Info("Cleaning up resources...")
}
