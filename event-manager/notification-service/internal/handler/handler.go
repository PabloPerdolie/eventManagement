package handler

import (
	"github.com/PabloPerdolie/event-manager/notification-service/internal/service"
	"go.uber.org/zap"
)

// Handler handles HTTP requests
type Handler struct {
	service *service.Service
	logger  *zap.SugaredLogger
}

// New creates a new handler
func New(svc *service.Service, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		service: svc,
		logger:  logger,
	}
}
