package handler

import (
	"net/http"
	"strconv"

	"github.com/event-management/core-service/internal/model"
	"github.com/event-management/core-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// UserController обрабатывает запросы, связанные с пользователями
type UserController struct {
	service service.UserService
	logger  *zap.SugaredLogger
}

// NewUserController создает новый контроллер пользователей
func NewUserController(service service.UserService, logger *zap.SugaredLogger) UserController {
	return UserController{
		service: service,
		logger:  logger,
	}
}

// Create обрабатывает запрос на создание пользователя
func (c UserController) Create(ctx *gin.Context) {
	var req model.UserCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Errorw("Failed to bind user create request", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := c.service.Create(ctx.Request.Context(), req)
	if err != nil {
		c.logger.Errorw("Failed to create user", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": id})
}

// GetByID обрабатывает запрос на получение пользователя по ID
func (c UserController) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.logger.Errorw("Invalid user ID", "error", err, "id", idStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := c.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		c.logger.Errorw("Failed to get user", "error", err, "id", id)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// Update обрабатывает запрос на обновление пользователя
func (c UserController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.logger.Errorw("Invalid user ID", "error", err, "id", idStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req model.UserUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Errorw("Failed to bind user update request", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.Update(ctx.Request.Context(), id, req); err != nil {
		c.logger.Errorw("Failed to update user", "error", err, "id", id)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// Delete обрабатывает запрос на удаление пользователя
func (c UserController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.logger.Errorw("Invalid user ID", "error", err, "id", idStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		c.logger.Errorw("Failed to delete user", "error", err, "id", id)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// List обрабатывает запрос на получение списка пользователей
func (c UserController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))

	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	users, err := c.service.List(ctx.Request.Context(), page, size)
	if err != nil {
		c.logger.Errorw("Failed to list users", "error", err, "page", page, "size", size)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, users)
}
