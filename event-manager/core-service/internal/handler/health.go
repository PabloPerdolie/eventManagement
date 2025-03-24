package handler

import (
	"net/http"

	"github.com/PabloPerdolie/event-manager/core-service/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HealthController обрабатывает запросы проверки здоровья системы
type HealthController struct {
	service service.HealthService
	logger  *zap.SugaredLogger
}

// NewHealthController создает новый контроллер проверки здоровья
func NewHealthController(service service.HealthService, logger *zap.SugaredLogger) HealthController {
	return HealthController{
		service: service,
		logger:  logger,
	}
}

// Check обрабатывает запрос на проверку здоровья сервиса
func (c HealthController) Check(ctx *gin.Context) {
	status, err := c.service.Check(ctx.Request.Context())
	if err != nil {
		c.logger.Errorw("Health check failed", "error", err)
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   status,
	})
}
