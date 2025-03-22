package handler

import (
	"github.com/PabloPerdolie/event-manager/notification-service/internal/domain"
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Service interface {
	GetStats() map[string]interface{}
}

type Handler struct {
	service Service
	logger  *zap.SugaredLogger
}

func New(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) HealthCheck(c *gin.Context) {
	health := domain.Health{
		Status:    "UP",
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   "1.0.0",
		Service:   "notification-service",
	}

	c.JSON(http.StatusOK, health)
}

func (h *Handler) ServiceInfo(c *gin.Context) {
	stats := h.service.GetStats()
	stats["service"] = "notification-service"
	stats["version"] = "1.0.0"
	stats["timestamp"] = time.Now().Format(time.RFC3339)

	c.JSON(http.StatusOK, stats)
}
