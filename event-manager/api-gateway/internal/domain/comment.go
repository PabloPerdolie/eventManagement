package domain

type CommentCreateRequest struct {
	SenderId int    `json:"sender_id"`
	EventId  int    `json:"event_id" binding:"required"`
	TaskId   *int   `json:"task_id"`
	Content  string `json:"content" binding:"required"`
}
