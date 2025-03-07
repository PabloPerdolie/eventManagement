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
	"github.com/yourusername/event-management/core-service/internal/config"
	"github.com/yourusername/event-management/core-service/internal/handler"
	"github.com/yourusername/event-management/core-service/internal/repository"
	"github.com/yourusername/event-management/core-service/internal/service"
	"go.uber.org/zap"
)

// @title Event Management Core Service API
// @version 1.0
// @description API for managing events, budget, and tasks
// @host localhost:8080
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

	// Initialize service layer
	svc := service.New(repo, sugar)

	// Initialize handler layer
	h := handler.New(svc, sugar)

	// Initialize router
	router := gin.Default()
	
	// Register API routes
	api := router.Group("/api/v1")
	{
		events := api.Group("/events")
		{
			events.GET("", h.GetEvents)
			events.POST("", h.CreateEvent)
			events.GET("/:id", h.GetEvent)
			events.PUT("/:id", h.UpdateEvent)
			events.DELETE("/:id", h.DeleteEvent)
			
			// Budget routes
			events.GET("/:id/budget", h.GetEventBudget)
			events.POST("/:id/budget", h.AddBudgetItem)
			events.PUT("/:id/budget/:itemId", h.UpdateBudgetItem)
			events.DELETE("/:id/budget/:itemId", h.DeleteBudgetItem)
			
			// Task routes
			events.GET("/:id/tasks", h.GetEventTasks)
			events.POST("/:id/tasks", h.CreateTask)
			events.PUT("/:id/tasks/:taskId", h.UpdateTask)
			events.DELETE("/:id/tasks/:taskId", h.DeleteTask)
			
			// Participant routes
			events.GET("/:id/participants", h.GetEventParticipants)
			events.POST("/:id/participants", h.AddParticipant)
			events.DELETE("/:id/participants/:userId", h.RemoveParticipant)
		}
	}

	// Setup Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
