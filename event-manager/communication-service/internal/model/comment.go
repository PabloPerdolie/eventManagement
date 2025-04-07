package model

import "time"

type Comment struct {
	CommentId int
	EventId   int
	SenderId  int
	TaskId    *int
	Content   string
	CreatedAt time.Time
	IsDeleted bool
	IsRead    bool
}
