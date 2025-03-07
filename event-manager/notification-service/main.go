package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/yourusername/event-management/notification-service/internal/config"
	"github.com/yourusername/event-management/notification-service/internal/handler"
	"github.com/yourusername/event-management/notification-service/internal/repository"
	"github.com/yourusername/event-management/notification-service/internal/service"
	"go.uber.org/zap"
)

// @title Notification Service API
// @version 1.0
// @description API for sending notifications to users
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

	// Initialize repository
	repo, err := repository.New(cfg)
	if err != nil {
		sugar.Fatalf("Failed to initialize repository: %v", err)
	}
	defer repo.Close()

	// Initialize service layer with notification providers
	svc, err := service.New(repo, cfg, sugar)
	if err != nil {
		sugar.Fatalf("Failed to initialize service: %v", err)
	}

	// Initialize handler layer
	h := handler.New(svc, sugar)

	// Initialize router
	router := gin.Default()
	
	// Register API routes
	api := router.Group("/api/v1")
	{
		notifications := api.Group("/notifications")
		{
			notifications.POST("/send", h.SendNotification)
			notifications.GET("/history", h.GetNotificationHistory)
			notifications.GET("/history/:id", h.GetNotificationById)
			notifications.GET("/user/:userId", h.GetUserNotifications)
			
			// Channels management
			notifications.GET("/channels", h.GetNotificationChannels)
			notifications.POST("/channels", h.CreateNotificationChannel)
			notifications.PUT("/channels/:id", h.UpdateNotificationChannel)
			notifications.DELETE("/channels/:id", h.DeleteNotificationChannel)
			
			// Preferences management
			notifications.GET("/preferences/:userId", h.GetUserPreferences)
			notifications.PUT("/preferences/:userId", h.UpdateUserPreferences)
		}
		
		// Webhook endpoints for notification delivery status
		api.POST("/webhooks/delivery-status", h.HandleDeliveryStatus)
	}

	// Setup Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Register consumer for events from message queue
	go svc.ConsumeEvents()

	// Setup server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sugar.Infof("Starting server on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sugar.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-quit
	sugar.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		sugar.Fatalf("Server forced to shutdown: %v", err)
	}

	sugar.Info("Server exited properly")
}
