package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/event-management/core-service/internal/config"
	"github.com/event-management/core-service/internal/handler"
	"github.com/event-management/core-service/internal/repository"
	"github.com/event-management/core-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	// Инициализация логгера
	logger, err := initLogger()
	if err != nil {
		fmt.Printf("Error initializing logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		sugar.Fatalw("Failed to load config", "error", err)
	}

	// Подключение к базе данных
	db, err := initDB(cfg, sugar)
	if err != nil {
		sugar.Fatalw("Failed to initialize database", "error", err)
	}
	defer db.Close()

	// Инициализация репозиториев
	repos := repository.New(db)

	// Инициализация сервисов
	services := service.New(repos, sugar)

	// Инициализация обработчиков запросов
	handlers := handler.New(services, sugar)

	// Инициализация HTTP-сервера с Gin
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Регистрация маршрутов
	handlers.InitRoutes(router)

	// Запуск HTTP-сервера
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Запуск сервера в горутине
	go func() {
		sugar.Infow("Starting server", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sugar.Fatalw("Failed to start server", "error", err)
		}
	}()

	// Обработка сигналов для корректного завершения работы
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	sugar.Info("Shutting down server...")

	// Установка тайм-аута для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		sugar.Fatalw("Server forced to shutdown", "error", err)
	}

	sugar.Info("Server exiting")
}

// initLogger инициализирует zap-логгер
func initLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	return config.Build()
}

// initDB инициализирует подключение к базе данных
func initDB(cfg *config.Config, logger *zap.SugaredLogger) (*sqlx.DB, error) {
	// Формирование строки подключения
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.DBName,
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
