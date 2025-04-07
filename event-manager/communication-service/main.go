package main

import (
	"context"
	"github.com/PabloPerdolie/event-manager/communication-service/internal/assembly"
	"github.com/PabloPerdolie/event-manager/communication-service/internal/config"
	"github.com/PabloPerdolie/event-manager/communication-service/internal/routes"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	_ "github.com/PabloPerdolie/event-manager/communication-service/docs"
)

// @title Communication Service for Event Management App
// @version 1.0
// @description Service for handling event communication, including comments and notifications
// @host localhost:8083
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

	locator, err := assembly.NewLocator(cfg, sugar)
	if err != nil {
		sugar.Fatalf("Failed to initialize service locator: %v", err)
	}
	defer locator.Close()

	router := gin.Default()

	// Добавляем endpoint для Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.SetupRoutes(router, &locator.Controllers)

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
