package assembly

import (
	"github.com/PabloPerdolie/event-manager/core-service/internal/config"
	"github.com/PabloPerdolie/event-manager/core-service/internal/handler"
	"github.com/PabloPerdolie/event-manager/core-service/internal/repository"
	"github.com/PabloPerdolie/event-manager/core-service/internal/routes"
	"github.com/PabloPerdolie/event-manager/core-service/internal/service"
	"github.com/PabloPerdolie/event-manager/core-service/internal/service/event"
	"github.com/PabloPerdolie/event-manager/core-service/internal/service/task"
	"github.com/PabloPerdolie/event-manager/core-service/pkg/postgres"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type ServiceLocator struct {
	Controllers routes.Controllers
	DB          *sqlx.DB
	logger      *zap.SugaredLogger
}

func NewLocator(cfg *config.Config, logger *zap.SugaredLogger) (*ServiceLocator, error) {
	db, err := postgres.InitDB(logger, cfg.DatabaseURL)
	if err != nil {
		return nil, errors.WithMessage(err, "init db")
	}

	eventRepo := repository.NewEvent(db)
	taskRepo := repository.NewTask(db)
	assignmentRepo := repository.NewTaskAssignment(db)
	userRepo := repository.NewUser(db)
	participantRepo := repository.NewParticipant(db)

	healthService := service.NewHealthService(db, logger)

	eventService := event.NewService(eventRepo, participantRepo, logger)
	eventParticipantService := event.NewParticipantService(participantRepo, userRepo, logger)
	taskService := task.NewService(taskRepo, assignmentRepo, logger)

	commonService := service.NewService(taskService, eventParticipantService, eventService)

	healthCtrl := handler.NewHealthController(healthService, logger)
	eventCtrl := handler.NewEvent(commonService, eventService, logger)
	eventParticipantCtrl := handler.NewParticipantHandler(eventParticipantService, logger)
	taskCtrl := handler.NewTask(taskService, logger)

	controllers := routes.Controllers{
		HealthCtrl:           healthCtrl,
		EventCtrl:            eventCtrl,
		EventParticipantCtrl: eventParticipantCtrl,
		TaskCtrl:             taskCtrl,
		//ExpenseCtrl:          expenseCtrl,
		//ExpenseShareCtrl:     expenseShareCtrl,
	}

	return &ServiceLocator{
		Controllers: controllers,
		DB:          db,
		logger:      logger,
	}, nil
}

func (l *ServiceLocator) Close() {
	l.logger.Info("Cleaning up resources...")
	if l.DB != nil {
		l.DB.Close()
	}
}
