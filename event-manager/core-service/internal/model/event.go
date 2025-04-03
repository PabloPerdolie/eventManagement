package model

import (
	"time"
)

type Event struct {
	EventID     int       `db:"event_id"`
	OrganizerID int       `db:"organizer_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	StartDate   time.Time `db:"start_date"`
	EndDate     time.Time `db:"end_date"`
	Location    *string   `db:"location"`
	Status      string    `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
}

type EventParticipant struct {
	EventParticipantID int        `db:"event_participant_id"`
	EventID            int        `db:"event_id"`
	UserID             int        `db:"user_id"`
	Role               string     `db:"role"`
	JoinedAt           *time.Time `db:"joined_at"`
	IsConfirmed        *bool      `db:"is_confirmed"`
}
