package api

import (
	"github.com/gin-gonic/gin"
)

func SetupAPIGatewayRoutes(router *gin.Engine) {
	router.GET("/notifications/status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "up",
			"service": "notification-service",
		})
	})
}
