package assembly

import (
	"github.com/PabloPerdolie/event-manager/core-service/internal/config"
	"github.com/PabloPerdolie/event-manager/core-service/internal/handler"
	"github.com/PabloPerdolie/event-manager/core-service/internal/routes"
	"github.com/PabloPerdolie/event-manager/core-service/internal/service"
	"github.com/PabloPerdolie/event-manager/core-service/internal/storage"
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
	db, err := postgres.InitDB(cfg.Database., logger)
	if err != nil {
		return nil, errors.WithMessage(err, "init db")
	}

	//repositories := repository.New(db)

	healthService := service.NewHealthService(db, logger)

	healthCtrl := handler.NewHealthController(healthService, logger)

	controllers := routes.Controllers{
		HealthCtrl: healthCtrl,
		//UserCtrl:             userCtrl,
		//EventCtrl:            eventCtrl,
		//EventParticipantCtrl: eventParticipantCtrl,
		//TaskCtrl:             taskCtrl,
		//TaskAssignmentCtrl:   taskAssignmentCtrl,
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
