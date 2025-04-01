package domain

import "time"

type CommentResponse struct {
	CommentId int       `json:"comment_id"`
	EventId   int       `json:"event_id"`
	SenderId  int       `json:"sender_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	IsRead    bool      `json:"is_read"`
}

type CommentCreateRequest struct {
	EventId int    `json:"event_id" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type CommentListResponse struct {
	Comments []CommentResponse `json:"comments"`
	Total    int               `json:"total"`
}
