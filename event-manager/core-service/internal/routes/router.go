package routes

import (
	"github.com/PabloPerdolie/event-manager/core-service/internal/handler"
	"github.com/gin-gonic/gin"
)

type Controllers struct {
	HealthCtrl           handler.HealthController
	EventCtrl            handler.EventController
	EventParticipantCtrl handler.ParticipantController
	TaskCtrl             handler.TaskController
	ExpenseCtrl          handler.ExpenseController
}

func SetupRoutes(router *gin.Engine, controllers *Controllers) {
	// API v1 группа
	api := router.Group("/api/v1")
	{
		api.GET("/health", controllers.HealthCtrl.Check)

		// Маршруты событий
		events := api.Group("/events")
		{
			events.GET("", controllers.EventCtrl.List)
			events.POST("", controllers.EventCtrl.Create)
			//events.PUT("/:id", controllers.EventCtrl.Update)
			events.DELETE("/:event_id", controllers.EventCtrl.Delete)
			events.GET("/:event_id", controllers.EventCtrl.EventSummary)

			// Маршруты участников событий
			participants := events.Group("/:event_id/participants")
			{
				participants.POST("", controllers.EventParticipantCtrl.Create)
				participants.DELETE("/:user_id", controllers.EventParticipantCtrl.Delete)
			}
		}

		// Маршруты задач
		tasks := api.Group("/tasks")
		{
			tasks.GET("", controllers.TaskCtrl.List)
			tasks.POST("", controllers.TaskCtrl.Create)
			tasks.PUT("/:task_id", controllers.TaskCtrl.Update)
			tasks.DELETE("/:task_id", controllers.TaskCtrl.Delete)
		}

		// Маршруты расходов - создание, обновление и удаление
		expenses := api.Group("/expenses")
		{
			expenses.POST("", controllers.ExpenseCtrl.Create)
			expenses.DELETE("/:id", controllers.ExpenseCtrl.Delete)

			// Маршруты для долей расходов
			expenseShares := api.Group("/shares")
			{
				expenseShares.PUT("/:id/paid-status", controllers.ExpenseCtrl.UpdateExpenseSharePaidStatus)
			}
		}
	}
}
