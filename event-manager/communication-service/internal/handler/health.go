package handler

import (
	"github.com/PabloPerdolie/event-manager/communication-service/internal/domain"
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

// HealthCheck godoc
// @Summary Проверка состояния сервиса
// @Description Проверяет доступность и работоспособность сервиса
// @Tags health
// @Produce json
// @Success 200 {object} domain.Health "Информация о состоянии сервиса"
// @Router /health [get]
func (h *Handler) HealthCheck(c *gin.Context) {
	health := domain.Health{
		Status:    "UP",
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   "1.0.0",
		Service:   "notification-service",
	}

	c.JSON(http.StatusOK, health)
}

// ServiceInfo godoc
// @Summary Информация о сервисе
// @Description Возвращает детальную информацию о сервисе и его статистике
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{} "Статистика и информация о сервисе"
// @Router /info [get]
func (h *Handler) ServiceInfo(c *gin.Context) {
	stats := h.service.GetStats()
	stats["service"] = "notification-service"
	stats["version"] = "1.0.0"
	stats["timestamp"] = time.Now().Format(time.RFC3339)

	c.JSON(http.StatusOK, stats)
}
