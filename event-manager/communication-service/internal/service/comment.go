package service

import (
	"context"
	"github.com/PabloPerdolie/event-manager/communication-service/internal/model"
)

type CommentRepo interface {
	Insert(ctx context.Context, comment model.Comment) (int, error)
}

type Comment struct {
	commentRepo CommentRepo
}

func NewComment(commentRepo CommentRepo) Comment {
	return Comment{
		commentRepo: commentRepo,
	}
}
