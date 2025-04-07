package routes

import (
	"github.com/PabloPerdolie/event-manager/communication-service/internal/handler"
	"github.com/gin-gonic/gin"
)

type Controllers struct {
	HealthCtrl  handler.Handler
	CommentCtrl *handler.CommentHandler
}

func SetupRoutes(router *gin.Engine, c *Controllers) {
	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", c.HealthCtrl.HealthCheck)
		v1.GET("/info", c.HealthCtrl.ServiceInfo)

		comments := v1.Group("/comments")
		{
			comments.GET("/event/:eventId", c.CommentCtrl.GetCommentsByEventId)
			comments.DELETE("/:id", c.CommentCtrl.DeleteComment)
			comments.PUT("/:id/read", c.CommentCtrl.MarkCommentAsRead)
		}
	}

	router.GET("/health", c.HealthCtrl.HealthCheck)
}
