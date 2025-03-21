package service

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/PabloPerdolie/event-manager/notification-service/internal/config"
	"github.com/PabloPerdolie/event-manager/notification-service/internal/model"
	"go.uber.org/zap"
)

type NotificationService struct {
	config *config.Config
	logger *zap.SugaredLogger
}

func NewNotificationService(cfg *config.Config, logger *zap.SugaredLogger) *NotificationService {
	return &NotificationService{
		config: cfg,
		logger: logger,
	}
}

func (s *NotificationService) ProcessNotification(message *model.NotificationMessage) error {
	s.logger.Infof("Processing notification of type: %s", message.Event)

	emailContent, err := message.GenerateEmailContent()
	if err != nil {
		s.logger.Errorf("Failed to generate email content: %v", err)
		return err
	}

	if err := s.SendEmail(emailContent); err != nil {
		s.logger.Errorf("Failed to send email: %v", err)
		return err
	}

	s.logger.Infof("Email notification sent to %s for event %s", emailContent.To, message.Event)
	return nil
}

func (s *NotificationService) SendEmail(content *model.EmailContent) error {
	smtpConfig := s.config.SMTP

	smtpAddress := smtpConfig.GetSMTPAddress()
	s.logger.Infof("Attempting to connect to SMTP server at %s", smtpAddress)

	to := []string{content.To}
	message := []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"\r\n"+
			"%s\r\n",
		smtpConfig.Sender,
		content.To,
		content.Subject,
		content.Body,
	))

	err := smtp.SendMail(
		smtpAddress,       // Должно быть "localhost:1025"
		nil,               // Без аутентификации для MailHog
		smtpConfig.Sender, // Отправитель
		to,                // Получатель
		message,           // Тело письма
	)

	if err != nil {
		s.logger.Infof("Failed to send email to %s: %v", content.To, err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	s.logger.Infof("Successfully sent email to %s", content.To)
	return nil
}

func (s *NotificationService) GetSupportedEvents() []string {
	events := make([]string, 0, len(model.SupportedEvents))
	for event := range model.SupportedEvents {
		events = append(events, event)
	}
	return events
}

func (s *NotificationService) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"supported_events": strings.Join(s.GetSupportedEvents(), ", "),
		"smtp_host":        s.config.SMTP.Host,
		"smtp_sender":      s.config.SMTP.Sender,
	}
}
