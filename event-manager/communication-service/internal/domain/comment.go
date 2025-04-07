package domain

type CreateCommentMessage struct {
	EventId  int    `json:"event_id"`
	SenderId int    `json:"sender_id"`
	Content  string `json:"content"`
	TaskId   *int   `json:"task_id"`
}
