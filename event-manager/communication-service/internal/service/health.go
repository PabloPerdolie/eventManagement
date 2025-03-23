package service

import (
	"github.com/PabloPerdolie/event-manager/communication-service/internal/config"
)

type NotificationService struct {
	config config.Config
}

func New(config config.Config) *NotificationService {
	return &NotificationService{
		config: config,
	}
}

func (s *NotificationService) GetStats() map[string]interface{} {
	return map[string]interface{}{}
}
