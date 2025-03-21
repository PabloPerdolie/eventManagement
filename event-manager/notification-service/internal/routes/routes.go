package routes

import (
	"github.com/PabloPerdolie/event-manager/notification-service/internal/handler"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, h *handler.Handler) {
	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", h.HealthCheck)
		v1.GET("/info", h.ServiceInfo)
	}

	router.GET("/health", h.HealthCheck)
}
