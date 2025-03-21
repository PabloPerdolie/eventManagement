package service

import (
	"github.com/PabloPerdolie/event-manager/notification-service/internal/config"
	"go.uber.org/zap"
)

type Service struct {
	Notification *NotificationService
	logger       *zap.SugaredLogger
}

func New(cfg *config.Config, logger *zap.SugaredLogger) *Service {
	notificationService := NewNotificationService(cfg, logger)

	return &Service{
		Notification: notificationService,
		logger:       logger,
	}
}
