package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PabloPerdolie/event-manager/core-service/internal/assembly"
	"github.com/PabloPerdolie/event-manager/core-service/internal/config"
	"github.com/PabloPerdolie/event-manager/core-service/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// @title Event Management Core Service API
// @version 1.0
// @description API for managing events, users, tasks and expenses
// @host localhost:8080
// @BasePath /api/v1
func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	// Загрузка переменных окружения из .env файла
	if err := godotenv.Load(); err != nil {
		sugar.Warnf("Error loading .env file: %v", err)
	}

	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		sugar.Fatalf("Failed to load config: %v", err)
	}

	// Инициализация service locator
	locator, err := assembly.NewLocator(cfg, sugar)
	if err != nil {
		sugar.Fatalf("Failed to initialize service locator: %v", err)
	}
	defer locator.Close()

	// Инициализация роутера
	router := gin.Default()

	// Настройка Swagger документации
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Настройка маршрутов API
	routes.SetupRoutes(router, &locator.Controllers)

	// Настройка HTTP-сервера
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Запуск сервера в отдельной горутине
	go func() {
		sugar.Infof("Starting server on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sugar.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Настройка обработки сигналов для корректного завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	sugar.Info("Shutting down server...")

	// Ожидание завершения текущих запросов
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		sugar.Fatalf("Server forced to shutdown: %v", err)
	}

	sugar.Info("Server exited properly")
}
