package domain

import "time"

type EventParticipantCreateRequest struct {
	EventTitle string `json:"event_title" binding:"required"`
	UserID     int    `json:"user_id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
}

type EventParticipantUpdateRequest struct {
	Role        *string `json:"role"`
	IsConfirmed *bool   `json:"is_confirmed"`
}

type EventParticipantResponse struct {
	Id          int          `json:"id"`
	EventID     int          `json:"event_id"`
	User        UserResponse `json:"user"`
	Role        string       `json:"role"`
	JoinedAt    *time.Time   `json:"joined_at,omitempty"`
	IsConfirmed *bool        `json:"is_confirmed,omitempty"`
}

type EventParticipantsResponse struct {
	Participants []EventParticipantResponse `json:"participants"`
	Total        int                        `json:"total"`
}
