package model

import (
	"errors"
	"fmt"
)

// Custom errors
var (
	ErrInvalidNotificationData = errors.New("invalid notification data")
	ErrUnsupportedEventType    = errors.New("unsupported event type")
)

// NotificationMessage represents a message from the Core Service
type NotificationMessage struct {
	Event string                 `json:"event"`
	Data  map[string]interface{} `json:"data"`
}

// SupportedEvents maps event types to boolean values to check if an event type is supported
var SupportedEvents = map[string]bool{
	"event_created": true,
	"task_assigned": true,
	"expense_added": true,
}

type EmailContent struct {
	Subject string
	Body    string
	To      string
}

func (msg *NotificationMessage) GenerateEmailContent() (*EmailContent, error) {
	if !SupportedEvents[msg.Event] {
		return nil, ErrUnsupportedEventType
	}

	userEmail, ok := msg.Data["user_email"].(string)
	if !ok || userEmail == "" {
		return nil, ErrInvalidNotificationData
	}

	var subject, body string

	switch msg.Event {
	case "event_created":
		title, ok := msg.Data["title"].(string)
		if !ok {
			return nil, ErrInvalidNotificationData
		}

		subject = "New Event Created"
		body = fmt.Sprintf("New event '%s' has been created.", title)

		if eventID, ok := msg.Data["event_id"].(float64); ok {
			body += fmt.Sprintf(" Event ID: %.0f", eventID)
		}

	case "task_assigned":
		// Handle task assignment notification
		taskName, ok := msg.Data["task_name"].(string)
		if !ok {
			return nil, ErrInvalidNotificationData
		}

		subject = "Task Assigned"
		body = fmt.Sprintf("You have been assigned a new task: '%s'.", taskName)

		// Add event name if available
		if eventName, ok := msg.Data["event_name"].(string); ok {
			body += fmt.Sprintf(" Event: %s", eventName)
		}

	case "expense_added":
		// Handle expense addition notification
		amount, ok := msg.Data["amount"].(float64)
		if !ok {
			return nil, ErrInvalidNotificationData
		}

		description, _ := msg.Data["description"].(string)
		subject = "New Expense Added"
		body = fmt.Sprintf("A new expense of $%.2f has been added", amount)

		if description != "" {
			body += fmt.Sprintf(" for '%s'", description)
		}
		body += "."

		// Add event name if available
		if eventName, ok := msg.Data["event_name"].(string); ok {
			body += fmt.Sprintf(" Event: %s", eventName)
		}

	default:
		return nil, ErrUnsupportedEventType
	}

	return &EmailContent{
		Subject: subject,
		Body:    body,
		To:      userEmail,
	}, nil
}
