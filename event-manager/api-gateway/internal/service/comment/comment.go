package comment

import (
	"context"
	"encoding/json"
	"github.com/PabloPerdolie/event-manager/api-gateway/internal/domain"
	"github.com/pkg/errors"
)

type Repository interface {
	Publish(ctx context.Context, data []byte) error
}

type Service struct {
	repo Repository
}

func New(repo Repository) Service {
	return Service{
		repo: repo,
	}
}

func (s Service) Create(ctx context.Context, req domain.CommentCreateRequest) error {
	data, err := json.Marshal(req)
	if err != nil {
		return errors.WithMessage(err, "marshal req")
	}

	err = s.repo.Publish(ctx, data)
	if err != nil {
		return errors.WithMessage(err, "publish comment req")
	}

	return nil
}
