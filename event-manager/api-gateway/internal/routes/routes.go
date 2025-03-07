package routes

import (
	"github.com/event-management/api-gateway/internal/config"
	"github.com/gin-contrib/cors"
	"net/http"

	"github.com/event-management/api-gateway/internal/handler"
	"github.com/event-management/api-gateway/internal/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine, cfg *config.Config, h *handler.Handler, authMiddleware *middleware.AuthMiddleware) {
	configureCORS(router, cfg.AllowedOrigin)


	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "API Gateway is running",
			"version": "1.0.0",
		})
	})
	
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	api := router.Group("/api/v1")
	{
		setupPublicRoutes(api, h)
		
		setupProtectedRoutes(api, h, authMiddleware)
		
		setupAdminRoutes(api, h, authMiddleware)
	}
	
	router.POST("/api/v1/auth/logout", authMiddleware.Authenticate(), h.Logout)
}

func setupPublicRoutes(api *gin.RouterGroup, h *handler.Handler) {
	auth := api.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.RefreshToken)
		auth.POST("/forgot-password", h.ForgotPassword)
		auth.POST("/reset-password", h.ResetPassword)
	}
}

func setupProtectedRoutes(api *gin.RouterGroup, h *handler.Handler, authMiddleware *middleware.AuthMiddleware) {
	// User management
	users := api.Group("/users", authMiddleware.Authenticate())
	{
		users.GET("/me", h.GetCurrentUser)
		users.PUT("/me", h.UpdateProfile)
		users.PUT("/me/password", h.ChangePassword)
		users.DELETE("/me", h.DeleteAccount)
	}
	
	setupServiceProxies(api, h, authMiddleware)
}

func setupAdminRoutes(api *gin.RouterGroup, h *handler.Handler, authMiddleware *middleware.AuthMiddleware) {
	admin := api.Group("/admin", authMiddleware.AuthenticateAdmin())
	{
		admin.GET("/users", h.GetAllUsers)
		admin.GET("/users/:id", h.GetUserById)
		admin.PUT("/users/:id", h.UpdateUser)
		admin.DELETE("/users/:id", h.DeleteUser)
	}
}

func setupServiceProxies(api *gin.RouterGroup, h *handler.Handler, authMiddleware *middleware.AuthMiddleware) {
	eventsProxy := api.Group("/events", authMiddleware.Authenticate())
	{
		eventsProxy.Any("/*path", h.ProxyToEventService)
	}
	
	notificationsProxy := api.Group("/notifications", authMiddleware.Authenticate())
	{
		notificationsProxy.Any("/*path", h.ProxyToNotificationService)
	}
	
	commentsProxy := api.Group("/comments", authMiddleware.Authenticate())
	{
		commentsProxy.Any("/*path", h.ProxyToCommunicationService)
	}
}

func configureCORS(router *gin.Engine, allowedOrigin string) {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{allowedOrigin}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))
}
