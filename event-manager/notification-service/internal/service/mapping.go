package service

import (
	"fmt"
	"github.com/PabloPerdolie/event-manager/notification-service/internal/domain"
	"github.com/PabloPerdolie/event-manager/notification-service/internal/model"
)

var SupportedEvents = map[string]bool{
	"event_created": true,
	"task_assigned": true,
	"expense_added": true,
}

func GenerateEmailContent(msg domain.NotificationMessage) (*model.EmailContent, error) {
	if !SupportedEvents[msg.Event] {
		return nil, model.ErrUnsupportedEventType
	}

	userEmail, ok := msg.Data["user_email"].(string)
	if !ok || userEmail == "" {
		return nil, model.ErrInvalidNotificationData
	}

	var subject, body string

	switch msg.Event {
	case "event_created":
		title, ok := msg.Data["title"].(string)
		if !ok {
			return nil, model.ErrInvalidNotificationData
		}

		subject = "New Event Created"
		body = fmt.Sprintf("New event '%s' has been created.", title)

		if eventID, ok := msg.Data["event_id"].(float64); ok {
			body += fmt.Sprintf(" Event ID: %.0f", eventID)
		}

	case "task_assigned":
		taskName, ok := msg.Data["task_name"].(string)
		if !ok {
			return nil, model.ErrInvalidNotificationData
		}

		subject = "Task Assigned"
		body = fmt.Sprintf("You have been assigned a new task: '%s'.", taskName)

		if eventName, ok := msg.Data["event_name"].(string); ok {
			body += fmt.Sprintf(" Event: %s", eventName)
		}

	case "expense_added":
		amount, ok := msg.Data["amount"].(float64)
		if !ok {
			return nil, model.ErrInvalidNotificationData
		}

		description, _ := msg.Data["description"].(string)
		subject = "New Expense Added"
		body = fmt.Sprintf("A new expense of $%.2f has been added", amount)

		if description != "" {
			body += fmt.Sprintf(" for '%s'", description)
		}
		body += "."

		if eventName, ok := msg.Data["event_name"].(string); ok {
			body += fmt.Sprintf(" Event: %s", eventName)
		}

	default:
		return nil, model.ErrUnsupportedEventType
	}

	return &model.EmailContent{
		Subject: subject,
		Body:    body,
		To:      userEmail,
	}, nil
}
