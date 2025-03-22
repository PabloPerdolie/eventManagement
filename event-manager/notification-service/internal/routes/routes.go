package routes

import (
	"github.com/PabloPerdolie/event-manager/notification-service/internal/handler"
	"github.com/gin-gonic/gin"
)

type Controllers struct {
	HealthCtrl handler.Handler
}

func SetupRoutes(router *gin.Engine, c *Controllers) {
	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", c.HealthCtrl.HealthCheck)
		v1.GET("/info", c.HealthCtrl.ServiceInfo)
	}

	router.GET("/health", c.HealthCtrl.HealthCheck)
}
