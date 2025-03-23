package routes

import (
	"github.com/PabloPerdolie/event-manager/core-service/internal/handler"
)

// Controllers содержит все контроллеры приложения
type Controllers struct {
	HealthCtrl           handler.HealthController
	UserCtrl             handler.UserController
	EventCtrl            handler.EventController
	EventParticipantCtrl handler.EventParticipantController
	TaskCtrl             handler.TaskController
	TaskAssignmentCtrl   handler.TaskAssignmentController
	ExpenseCtrl          handler.ExpenseController
	ExpenseShareCtrl     handler.ExpenseShareController
}
