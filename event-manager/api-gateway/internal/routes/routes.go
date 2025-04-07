package routes

import (
	"github.com/PabloPerdolie/event-manager/api-gateway/internal/config"
	"github.com/gin-contrib/cors"
	"net/http"

	"github.com/PabloPerdolie/event-manager/api-gateway/internal/handler"
	"github.com/PabloPerdolie/event-manager/api-gateway/internal/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Controllers struct {
	AuthCtrl  handler.Auth
	ProxyCtrl handler.Proxy
}

func SetupRoutes(router *gin.Engine, cfg *config.Config, c *Controllers, authMiddleware *middleware.AuthMiddleware) {
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
		setupPublicRoutes(api, c)

		setupProtectedRoutes(api, c, authMiddleware)

		//setupAdminRoutes(api, c, authMiddleware)
	}

	router.POST("/api/v1/auth/logout", authMiddleware.Authenticate(), c.AuthCtrl.Logout)
}

func setupPublicRoutes(api *gin.RouterGroup, c *Controllers) {
	auth := api.Group("/auth")
	{
		auth.POST("/register", c.AuthCtrl.Register)
		auth.POST("/login", c.AuthCtrl.Login)
		auth.POST("/refresh", c.AuthCtrl.RefreshToken)
		//auth.POST("/forgot-password", c.AuthCtrl.ForgotPassword)
		//auth.POST("/reset-password", c.AuthCtrl.ResetPassword)
	}
}

func setupProtectedRoutes(api *gin.RouterGroup, c *Controllers, authMiddleware *middleware.AuthMiddleware) {
	// User management
	//users := api.Group("/users", authMiddleware.Authenticate())
	//{
	//	users.GET("/me", h.GetCurrentUser)
	//	users.PUT("/me", h.UpdateProfile)
	//	users.PUT("/me/password", h.ChangePassword)
	//	users.DELETE("/me", h.DeleteAccount)
	//
	//	// User's comments
	//	users.GET("/me/comments", h.GetUserComments)
	//}

	setupServiceProxies(api, c, authMiddleware)
}

//
//func setupAdminRoutes(api *gin.RouterGroup, h *handler.Handler, authMiddleware *middleware.AuthMiddleware) {
//	admin := api.Group("/admin", authMiddleware.AuthenticateAdmin())
//	{
//		admin.GET("/users", h.GetAllUsers)
//		admin.GET("/users/:id", h.GetUserById)
//		admin.PUT("/users/:id", h.UpdateUser)
//		admin.DELETE("/users/:id", h.DeleteUser)
//	}
//}

func setupServiceProxies(api *gin.RouterGroup, c *Controllers, authMiddleware *middleware.AuthMiddleware) {
	eventsProxy := api.Group("/events", authMiddleware.Authenticate())
	{
		eventsProxy.Any("/*path", c.ProxyCtrl.ProxyToEventService)
	}

	commentsProxy := api.Group("/comments", authMiddleware.Authenticate())
	{
		commentsProxy.GET("/*path", c.ProxyCtrl.ProxyToCommunicationService)
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
