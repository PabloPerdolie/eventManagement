package service

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// HealthService определяет методы для проверки здоровья системы
type HealthService interface {
	Check(ctx context.Context) (map[string]interface{}, error)
}

type healthService struct {
	db     *sqlx.DB
	logger *zap.SugaredLogger
}

// NewHealthService создает новый сервис проверки здоровья системы
func NewHealthService(db *sqlx.DB, logger *zap.SugaredLogger) HealthService {
	return &healthService{
		db:     db,
		logger: logger,
	}
}

// Check проверяет состояние сервиса и его зависимостей
func (s *healthService) Check(ctx context.Context) (map[string]interface{}, error) {
	status := map[string]interface{}{
		"service": "core-service",
		"time":    time.Now().Format(time.RFC3339),
	}

	// Проверка соединения с базой данных
	if err := s.db.PingContext(ctx); err != nil {
		s.logger.Errorw("Database connection failed", "error", err)
		status["database"] = "error"
		status["database_error"] = err.Error()
		return status, err
	}

	status["database"] = "ok"
	return status, nil
}
