package handler

import (
	"github.com/event-management/api-gateway/internal/config"
	"github.com/event-management/api-gateway/internal/service"
	"go.uber.org/zap"
)

type Handler struct {
	service *service.Service
	config  *config.Config
	logger  *zap.SugaredLogger
}

func New(service *service.Service, config *config.Config, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		service: service,
		config:  config,
		logger:  logger,
	}
}
