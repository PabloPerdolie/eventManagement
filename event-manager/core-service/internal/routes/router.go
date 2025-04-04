package routes

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes настраивает все маршруты API
func SetupRoutes(router *gin.Engine, controllers *Controllers) {
	// API v1 группа
	api := router.Group("/api/v1")
	{
		// Маршрут проверки здоровья сервиса
		api.GET("/health", controllers.HealthCtrl.Check)

		// Маршруты событий
		events := api.Group("/events")
		{
			events.GET("", controllers.EventCtrl.List)
			events.POST("", controllers.EventCtrl.Create)
			events.GET("/:id", controllers.EventCtrl.GetById)
			events.PUT("/:id", controllers.EventCtrl.Update)
			events.DELETE("/:id", controllers.EventCtrl.Delete)

			// Маршруты участников событий
			participants := events.Group("/:id/participants")
			{
				participants.GET("", controllers.EventParticipantCtrl.ListByEvent)
				participants.POST("", controllers.EventParticipantCtrl.Create)
				participants.GET("/:user_id", controllers.EventParticipantCtrl.GetById)
				participants.PUT("/:user_id", controllers.EventParticipantCtrl.Update)
				participants.DELETE("/:user_id", controllers.EventParticipantCtrl.Delete)
			}
		}

		// Маршруты задач
		tasks := api.Group("/tasks")
		{
			tasks.GET("", controllers.TaskCtrl.List)
			tasks.POST("", controllers.TaskCtrl.Create)
			tasks.GET("/:id", controllers.TaskCtrl.GetById)
			tasks.PUT("/:id", controllers.TaskCtrl.Update)
			tasks.DELETE("/:id", controllers.TaskCtrl.Delete)

			// Маршруты назначений задач
			assignments := tasks.Group("/:id/assignments")
			{
				assignments.GET("", controllers.TaskAssignmentCtrl.ListByTask)
				assignments.POST("", controllers.TaskAssignmentCtrl.Create)
				assignments.GET("/:assignment_id", controllers.TaskAssignmentCtrl.GetById)
				assignments.PUT("/:assignment_id", controllers.TaskAssignmentCtrl.Update)
				assignments.DELETE("/:assignment_id", controllers.TaskAssignmentCtrl.Delete)
			}
		}

		// Маршруты расходов
		expenses := api.Group("/expenses")
		{
			expenses.GET("", controllers.ExpenseCtrl.List)
			expenses.POST("", controllers.ExpenseCtrl.Create)
			expenses.GET("/:id", controllers.ExpenseCtrl.GetById)
			expenses.PUT("/:id", controllers.ExpenseCtrl.Update)
			expenses.DELETE("/:id", controllers.ExpenseCtrl.Delete)

			// Маршруты долей расходов
			shares := expenses.Group("/:id/shares")
			{
				shares.GET("", controllers.ExpenseShareCtrl.ListByExpense)
				shares.POST("", controllers.ExpenseShareCtrl.Create)
				shares.GET("/:share_id", controllers.ExpenseShareCtrl.GetById)
				shares.PUT("/:share_id", controllers.ExpenseShareCtrl.Update)
				shares.DELETE("/:share_id", controllers.ExpenseShareCtrl.Delete)
			}
		}
	}
}
