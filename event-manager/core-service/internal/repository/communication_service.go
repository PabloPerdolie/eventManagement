package repository

import (
	"context"
	"encoding/json"
	"github.com/PabloPerdolie/event-manager/core-service/internal/model"
	"github.com/PabloPerdolie/event-manager/core-service/pkg/http/client"
	"github.com/pkg/errors"
	"strconv"
)

const (
	getEventById = "/api/v1/comments/event/"
)

type CommunicationServiceRepo struct {
	cli *client.Client
}

func NewCommunicationService(cli *client.Client) CommunicationServiceRepo {
	return CommunicationServiceRepo{cli: cli}
}

func (r CommunicationServiceRepo) GetByEventId(_ context.Context, eventId int) (*model.CommunicationServiceResponse, error) {
	id := strconv.Itoa(eventId)
	resp, err := r.cli.Invoke(getEventById + id)
	if err != nil {
		return nil, errors.WithMessagef(err, "invoke request by endpoint: %s", getEventById)
	}

	var comments model.CommunicationServiceResponse
	err = json.Unmarshal(resp, &comments)
	if err != nil {
		return nil, errors.WithMessage(err, "unmarshal response")
	}

	return &comments, nil
}
