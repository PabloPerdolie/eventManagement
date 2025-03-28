package assembly

import (
	"github.com/PabloPerdolie/event-manager/core-service/internal/config"
	"github.com/PabloPerdolie/event-manager/core-service/internal/handler"
	"github.com/PabloPerdolie/event-manager/core-service/internal/routes"
	"github.com/PabloPerdolie/event-manager/core-service/internal/service"
	"github.com/PabloPerdolie/event-manager/core-service/internal/storage"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// ServiceLocator содержит ссылки на все компоненты приложения
type ServiceLocator struct {
	Controllers routes.Controllers
	DB          *sqlx.DB
	logger      *zap.SugaredLogger
}

// NewLocator создает и инициализирует новый ServiceLocator
func NewLocator(cfg *config.Config, logger *zap.SugaredLogger) (*ServiceLocator, error) {
	// Инициализация базы данных
	db, err := storage.InitDB(cfg, logger)
	if err != nil {
		return nil, err
	}

	// Инициализация репозиториев
	//repositories := repository.New(db)

	// Инициализация сервисов
	healthService := service.NewHealthService(db, logger)

	// Инициализация контроллеров
	healthCtrl := handler.NewHealthController(healthService, logger)

	// Формирование структуры контроллеров для маршрутизации
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

// Close освобождает ресурсы при завершении работы приложения
func (l *ServiceLocator) Close() {
	l.logger.Info("Cleaning up resources...")
	if l.DB != nil {
		l.DB.Close()
	}
}
