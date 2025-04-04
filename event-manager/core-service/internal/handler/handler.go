package handler

import (
	"github.com/PabloPerdolie/event-manager/core-service/internal/handler/event"
	"github.com/PabloPerdolie/event-manager/core-service/internal/handler/expense"
	"github.com/PabloPerdolie/event-manager/core-service/internal/handler/task"
	"github.com/PabloPerdolie/event-manager/core-service/internal/handler/user"
	"github.com/PabloPerdolie/event-manager/core-service/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler contains all the HTTP handlers for the API
type Handler struct {
	User             user.Handler
	Event            event.Handler
	EventParticipant event.ParticipantHandler
	Task             task.Handler
	TaskAssignment   task.AssignmentHandler
	Expense          expense.Handler
	ExpenseShare     expense.ShareHandler
	logger           *zap.SugaredLogger
}

// New creates a new HTTP handler
func New(services *service.Service, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		User:             user.NewHandler(services.User, logger),
		Event:            event.NewHandler(services.Event, logger),
		EventParticipant: event.NewParticipantHandler(services.EventParticipant, logger),
		Task:             task.NewHandler(services.Task, logger),
		TaskAssignment:   task.NewAssignmentHandler(services.TaskAssignment, logger),
		Expense:          expense.NewHandler(services.Expense, logger),
		ExpenseShare:     expense.NewShareHandler(services.ExpenseShare, logger),
		logger:           logger,
	}
}

// InitRoutes initializes all routes for the API
func (h *Handler) InitRoutes(router *gin.Engine) {
	// API v1 group
	api := router.Group("/api/v1")
	{
		// User routes
		userRoutes := api.Group("/users")
		{
			userRoutes.GET("", h.User.List)
			userRoutes.POST("", h.User.Create)
			userRoutes.GET("/:id", h.User.GetById)
			userRoutes.PUT("/:id", h.User.Update)
			userRoutes.DELETE("/:id", h.User.Delete)
		}

		// Event routes
		eventRoutes := api.Group("/events")
		{
			eventRoutes.GET("", h.Event.List)
			eventRoutes.POST("", h.Event.Create)
			eventRoutes.GET("/:id", h.Event.GetById)
			eventRoutes.PUT("/:id", h.Event.Update)
			eventRoutes.DELETE("/:id", h.Event.Delete)

			// Event participants
			eventRoutes.GET("/:id/participants", h.EventParticipant.ListByEvent)
			eventRoutes.POST("/:id/participants", h.EventParticipant.Create)
			eventRoutes.GET("/:id/participants/:user_id", h.EventParticipant.GetById)
			eventRoutes.PUT("/:id/participants/:user_id", h.EventParticipant.Update)
			eventRoutes.DELETE("/:id/participants/:user_id", h.EventParticipant.Delete)
		}

		// Task routes
		taskRoutes := api.Group("/tasks")
		{
			taskRoutes.GET("", h.Task.List)
			taskRoutes.POST("", h.Task.Create)
			taskRoutes.GET("/:id", h.Task.GetById)
			taskRoutes.PUT("/:id", h.Task.Update)
			taskRoutes.DELETE("/:id", h.Task.Delete)

			// Task assignments
			taskRoutes.GET("/:id/assignments", h.TaskAssignment.ListByTask)
			taskRoutes.POST("/:id/assignments", h.TaskAssignment.Create)
			taskRoutes.GET("/:id/assignments/:assignment_id", h.TaskAssignment.GetById)
			taskRoutes.PUT("/:id/assignments/:assignment_id", h.TaskAssignment.Update)
			taskRoutes.DELETE("/:id/assignments/:assignment_id", h.TaskAssignment.Delete)
		}

		// Expense routes
		expenseRoutes := api.Group("/expenses")
		{
			expenseRoutes.GET("", h.Expense.List)
			expenseRoutes.POST("", h.Expense.Create)
			expenseRoutes.GET("/:id", h.Expense.GetById)
			expenseRoutes.PUT("/:id", h.Expense.Update)
			expenseRoutes.DELETE("/:id", h.Expense.Delete)

			// Expense shares
			expenseRoutes.GET("/:id/shares", h.ExpenseShare.ListByExpense)
			expenseRoutes.POST("/:id/shares", h.ExpenseShare.Create)
			expenseRoutes.GET("/:id/shares/:share_id", h.ExpenseShare.GetById)
			expenseRoutes.PUT("/:id/shares/:share_id", h.ExpenseShare.Update)
			expenseRoutes.DELETE("/:id/shares/:share_id", h.ExpenseShare.Delete)
		}
	}
}
