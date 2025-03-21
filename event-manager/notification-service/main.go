package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PabloPerdolie/event-manager/notification-service/internal/assembly"
	"github.com/PabloPerdolie/event-manager/notification-service/internal/config"
	"github.com/PabloPerdolie/event-manager/notification-service/internal/consumer"
	"github.com/PabloPerdolie/event-manager/notification-service/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// @title Notification Service for Event Management App
// @version 1.0
// @description Service for sending notifications via email based on RabbitMQ messages
// @host localhost:8082
// @BasePath /api/v1
func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	if err := godotenv.Load(); err != nil {
		sugar.Warnf("Error loading .env file: %v", err)
	}

	cfg, err := config.New()
	if err != nil {
		sugar.Fatalf("Failed to load config: %v", err)
	}

	locator, err := assembly.NewServiceLocator(cfg, sugar)
	if err != nil {
		sugar.Fatalf("Failed to initialize service locator: %v", err)
	}
	defer locator.Close()

	rabbitConsumer := consumer.New(locator.Service, cfg, sugar)
	go func() {
		sugar.Info("Starting RabbitMQ consumer...")
		if err := rabbitConsumer.Start(); err != nil {
			sugar.Fatalf("Failed to start RabbitMQ consumer: %v", err)
		}
	}()

	router := gin.Default()
	routes.SetupRoutes(router, locator.Handler)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		sugar.Infof("Starting HTTP server on port %s...", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sugar.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	sugar.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		sugar.Fatalf("Server forced to shutdown: %v", err)
	}

	sugar.Info("Server exiting")
}
