package storage

import (
	"fmt"
	"time"

	"github.com/event-management/core-service/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// InitDB инициализирует подключение к базе данных
func InitDB(cfg *config.Config, logger *zap.SugaredLogger) (*sqlx.DB, error) {
	// Формирование строки подключения
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	// Подключение к PostgreSQL
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Проверка соединения
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Настройка пула соединений
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Second)

	logger.Infow("Connected to database",
		"host", cfg.Database.Host,
		"port", cfg.Database.Port,
		"dbname", cfg.Database.DBName,
	)

	return db, nil
}
