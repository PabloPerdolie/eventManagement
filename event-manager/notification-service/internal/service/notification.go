package service

import (
	"context"
	"github.com/PabloPerdolie/event-manager/notification-service/internal/domain"
	"github.com/pkg/errors"
	"strings"

	"github.com/PabloPerdolie/event-manager/notification-service/internal/config"
	"github.com/PabloPerdolie/event-manager/notification-service/internal/model"
)

type NotificationRepo interface {
	SendEmail(_ context.Context, content *model.EmailContent) error
}

type NotificationService struct {
	config     config.SMTPConfig
	notifyRepo NotificationRepo
}

func New(config config.SMTPConfig, notifyRepo NotificationRepo) *NotificationService {
	return &NotificationService{
		config:     config,
		notifyRepo: notifyRepo,
	}
}

func (s *NotificationService) ProcessNotification(ctx context.Context, message domain.NotificationMessage) error {
	emailContent, err := GenerateEmailContent(message)
	if err != nil {
		return errors.WithMessage(err, "generate email content")
	}

	if err := s.notifyRepo.SendEmail(ctx, emailContent); err != nil {
		return errors.WithMessage(err, "send email")
	}

	return nil
}

func (s *NotificationService) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"supported_events": strings.Join(s.getSupportedEvents(), ", "),
		"smtp_host":        s.config.Host,
		"smtp_sender":      s.config.Sender,
	}
}

func (s *NotificationService) getSupportedEvents() []string {
	events := make([]string, 0, len(SupportedEvents))
	for event := range SupportedEvents {
		events = append(events, event)
	}
	return events
}
