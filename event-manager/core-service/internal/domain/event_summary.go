package domain

import "github.com/PabloPerdolie/event-manager/core-service/internal/model"

type EventData struct {
	EventParticipants EventParticipantsResponse
	EventData         EventResponse
	Tasks             TasksResponse
	Comments          model.CommunicationServiceResponse
	// expenses
}
