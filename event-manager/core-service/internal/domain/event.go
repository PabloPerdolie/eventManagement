package domain

import "time"

type EventCreateRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required"`
	Location    string    `json:"location"`
}

type EventUpdateRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Location    *string    `json:"location"`
}

type EventResponse struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Location    string    `json:"location"`
	Status      string    `json:"status"`
	CreatedBy   int       `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

type EventsResponse struct {
	Events []EventResponse `json:"events"`
	Total  int             `json:"total"`
}
