package repository

import (
	"context"
	"fmt"
	"github.com/PabloPerdolie/event-manager/notification-service/internal/config"
	"github.com/PabloPerdolie/event-manager/notification-service/internal/model"
	"github.com/pkg/errors"
	"net/smtp"
)

type Notification struct {
	smtpCfg config.SMTPConfig
}

func New(smtpConfig config.SMTPConfig) Notification {
	return Notification{
		smtpCfg: smtpConfig,
	}
}

func (r Notification) SendEmail(_ context.Context, content *model.EmailContent) error {
	smtpConfig := r.smtpCfg

	smtpAddress := smtpConfig.GetSMTPAddress()

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
		return errors.WithMessage(err, "send mail")
	}

	return nil
}
