package model

import "time"

//comment_id SERIAL PRIMARY KEY,
//event_id INT NOT NULL,
//sender_id INT NOT NULL,
//content TEXT NOT NULL,
//created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
//is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
//is_read BOOLEAN NOT NULL DEFAULT FALSE

type Comment struct {
	CommentId int
	EventId   int
	SenderId  int
	Content   string
	CreatedAt time.Time
	IsDeleted bool
	IsRead    bool
}
