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
	"github.com/yourusername/event-management/communication-service/internal/config"
	"github.com/yourusername/event-management/communication-service/internal/handler"
	"github.com/yourusername/event-management/communication-service/internal/repository"
	"github.com/yourusername/event-management/communication-service/internal/service"
	"go.uber.org/zap"
)

// @title Communication Service API
// @version 1.0
// @description API for comments and discussions
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

	// Initialize repository
	repo, err := repository.New(cfg)
	if err != nil {
		sugar.Fatalf("Failed to initialize repository: %v", err)
	}
	defer repo.Close()

	// Initialize service layer
	svc := service.New(repo, sugar)

	// Initialize handler layer with WebSocket hub
	h := handler.New(svc, sugar)

	// Initialize router
	router := gin.Default()
	
	// Register API routes
	api := router.Group("/api/v1")
	{
		comments := api.Group("/comments")
		{
			comments.GET("", h.GetComments)
			comments.POST("", h.CreateComment)
			comments.GET("/:id", h.GetComment)
			comments.PUT("/:id", h.UpdateComment)
			comments.DELETE("/:id", h.DeleteComment)
			
			// Reactions
			comments.POST("/:id/reactions", h.AddReaction)
			comments.DELETE("/:id/reactions/:reactionId", h.RemoveReaction)
			
			// Replies
			comments.GET("/:id/replies", h.GetReplies)
			comments.POST("/:id/replies", h.CreateReply)
		}
		
		// WebSocket endpoint for real-time updates
		router.GET("/ws", h.HandleWebSocket)
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
