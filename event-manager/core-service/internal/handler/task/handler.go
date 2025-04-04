package task

import (
	"net/http"
	"strconv"

	"github.com/PabloPerdolie/event-manager/core-service/internal/model"
	"github.com/PabloPerdolie/event-manager/core-service/internal/service/task"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Handler handles task-related HTTP requests
type Handler interface {
	Create(c *gin.Context)
	GetById(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	List(c *gin.Context)
}

type handler struct {
	service task.Service
	logger  *zap.SugaredLogger
}

// NewHandler creates a new task handler
func NewHandler(service task.Service, logger *zap.SugaredLogger) Handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

// Create handles creating a new task
func (h *handler) Create(c *gin.Context) {
	var req model.TaskCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if req.Title == "" || req.EventId == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title and event Id are required"})
		return
	}

	id, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		h.logger.Errorw("Failed to create task", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// GetById handles getting a task by Id
func (h *handler) GetById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Errorw("Invalid task Id", "error", err, "id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task Id"})
		return
	}

	task, err := h.service.GetById(c.Request.Context(), id)
	if err != nil {
		h.logger.Errorw("Failed to get task", "error", err, "id", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// Update handles updating a task
func (h *handler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Errorw("Invalid task Id", "error", err, "id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task Id"})
		return
	}

	var req model.TaskUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(c.Request.Context(), id, req); err != nil {
		h.logger.Errorw("Failed to update task", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Delete handles deleting a task
func (h *handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Errorw("Invalid task Id", "error", err, "id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task Id"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		h.logger.Errorw("Failed to delete task", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// List handles listing tasks with pagination and filtering
func (h *handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	// Get optional filters
	eventIdStr := c.Query("event_id")
	userIdStr := c.Query("user_id")

	var eventId *int
	if eventIdStr != "" {
		parsed, err := uuid.Parse(eventIdStr)
		if err == nil {
			eventId = &parsed
		} else {
			h.logger.Warnw("Invalid event Id filter", "error", err, "event_id", eventIdStr)
		}
	}

	var userId *int
	if userIdStr != "" {
		parsed, err := uuid.Parse(userIdStr)
		if err == nil {
			userId = &parsed
		} else {
			h.logger.Warnw("Invalid user Id filter", "error", err, "user_id", userIdStr)
		}
	}

	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	var tasks model.TasksResponse
	var err error

	// Apply filters if provided
	if eventId != nil {
		tasks, err = h.service.ListByEvent(c.Request.Context(), *eventId, page, size)
	} else if userId != nil {
		tasks, err = h.service.ListByUser(c.Request.Context(), *userId, page, size)
	} else {
		// Default to all tasks (for admin purposes)
		tasks, err = h.service.List(c.Request.Context(), page, size)
	}

	if err != nil {
		h.logger.Errorw("Failed to list tasks", "error", err, "page", page, "size", size, "event_id", eventId, "user_id", userId)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}
