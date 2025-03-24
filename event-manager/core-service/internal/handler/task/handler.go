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
	GetByID(c *gin.Context)
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
	if req.Title == "" || req.EventID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title and event ID are required"})
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

// GetByID handles getting a task by ID
func (h *handler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Errorw("Invalid task ID", "error", err, "id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	task, err := h.service.GetByID(c.Request.Context(), id)
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
		h.logger.Errorw("Invalid task ID", "error", err, "id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
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
		h.logger.Errorw("Invalid task ID", "error", err, "id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
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
	eventIDStr := c.Query("event_id")
	userIDStr := c.Query("user_id")

	var eventID *uuid.UUID
	if eventIDStr != "" {
		parsed, err := uuid.Parse(eventIDStr)
		if err == nil {
			eventID = &parsed
		} else {
			h.logger.Warnw("Invalid event ID filter", "error", err, "event_id", eventIDStr)
		}
	}

	var userID *uuid.UUID
	if userIDStr != "" {
		parsed, err := uuid.Parse(userIDStr)
		if err == nil {
			userID = &parsed
		} else {
			h.logger.Warnw("Invalid user ID filter", "error", err, "user_id", userIDStr)
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
	if eventID != nil {
		tasks, err = h.service.ListByEvent(c.Request.Context(), *eventID, page, size)
	} else if userID != nil {
		tasks, err = h.service.ListByUser(c.Request.Context(), *userID, page, size)
	} else {
		// Default to all tasks (for admin purposes)
		tasks, err = h.service.List(c.Request.Context(), page, size)
	}

	if err != nil {
		h.logger.Errorw("Failed to list tasks", "error", err, "page", page, "size", size, "event_id", eventID, "user_id", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}
