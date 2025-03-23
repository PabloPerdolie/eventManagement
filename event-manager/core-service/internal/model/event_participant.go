package model

import (
	"time"

	"github.com/google/uuid"
)

// EventParticipant represents a participant in an event
type EventParticipant struct {
	ID          uuid.UUID `json:"id" db:"id"`
	EventID     uuid.UUID `json:"event_id" db:"event_id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	JoinedAt    time.Time `json:"joined_at" db:"joined_at"`
	IsConfirmed bool      `json:"is_confirmed" db:"is_confirmed"`
}

// EventParticipantCreateRequest represents the input for adding a participant to an event
type EventParticipantCreateRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
}

// EventParticipantUpdateRequest represents the input for updating a participant status
type EventParticipantUpdateRequest struct {
	IsConfirmed *bool `json:"is_confirmed"`
}

// EventParticipantResponse represents the output for event participant data
type EventParticipantResponse struct {
	ID          uuid.UUID    `json:"id"`
	EventID     uuid.UUID    `json:"event_id"`
	User        UserResponse `json:"user"`
	JoinedAt    time.Time    `json:"joined_at"`
	IsConfirmed bool         `json:"is_confirmed"`
}

// EventParticipantsResponse represents a list of event participants
type EventParticipantsResponse struct {
	Participants []EventParticipantResponse `json:"participants"`
	Total        int                        `json:"total"`
}
