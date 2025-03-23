package handler

import (
	"github.com/PabloPerdolie/event-manager/core-service/internal/model"
	"github.com/PabloPerdolie/event-manager/core-service/internal/service/event"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// EventController обрабатывает запросы, связанные с событиями
type EventController struct {
	service event.Service
	logger  *zap.SugaredLogger
}

// NewEventController создает новый контроллер событий
func NewEventController(service event.Service, logger *zap.SugaredLogger) EventController {
	return EventController{
		service: service,
		logger:  logger,
	}
}

// Create обрабатывает запрос на создание события
func (c EventController) Create(ctx *gin.Context) {
	var req model.EventCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Errorw("Failed to bind event create request", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := c.service.Create(ctx.Request.Context(), req)
	if err != nil {
		c.logger.Errorw("Failed to create event", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": id})
}

// GetByID обрабатывает запрос на получение события по ID
func (c EventController) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.logger.Errorw("Invalid event ID", "error", err, "id", idStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	event, err := c.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		c.logger.Errorw("Failed to get event", "error", err, "id", id)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	ctx.JSON(http.StatusOK, event)
}

// Update обрабатывает запрос на обновление события
func (c EventController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.logger.Errorw("Invalid event ID", "error", err, "id", idStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	var req model.EventUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Errorw("Failed to bind event update request", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.Update(ctx.Request.Context(), id, req); err != nil {
		c.logger.Errorw("Failed to update event", "error", err, "id", id)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// Delete обрабатывает запрос на удаление события
func (c EventController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.logger.Errorw("Invalid event ID", "error", err, "id", idStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		c.logger.Errorw("Failed to delete event", "error", err, "id", id)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// List обрабатывает запрос на получение списка событий
func (c EventController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	userIDStr := ctx.Query("user_id")

	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	var events model.EventsResponse
	var err error

	if userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.logger.Errorw("Invalid user ID", "error", err, "user_id", userIDStr)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
			return
		}
		events, err = c.service.ListByUser(ctx.Request.Context(), userID, page, size)
	} else {
		events, err = c.service.List(ctx.Request.Context(), page, size)
	}

	if err != nil {
		c.logger.Errorw("Failed to list events", "error", err, "page", page, "size", size, "user_id", userIDStr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, events)
}
