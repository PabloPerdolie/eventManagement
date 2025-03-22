package model

import (
	"errors"
)

var (
	ErrInvalidNotificationData = errors.New("invalid notification data")
	ErrUnsupportedEventType    = errors.New("unsupported event type")
)

type EmailContent struct {
	Subject string
	Body    string
	To      string
}
