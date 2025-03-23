package domain

import "time"

// CommentResponse represents a comment returned to the client
type CommentResponse struct {
	CommentId int       `json:"comment_id"`
	EventId   int       `json:"event_id"`
	SenderId  int       `json:"sender_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	IsRead    bool      `json:"is_read"`
}

// CommentCreateRequest represents a request to create a new comment
type CommentCreateRequest struct {
	EventId int    `json:"event_id" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// CommentListResponse represents a list of comments
type CommentListResponse struct {
	Comments []CommentResponse `json:"comments"`
	Total    int               `json:"total"`
}
