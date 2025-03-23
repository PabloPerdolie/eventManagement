package event

import (
	"net/http"
	"strconv"

	"github.com/event-management/core-service/internal/model"
	"github.com/event-management/core-service/internal/service/event"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Handler handles event-related HTTP requests
type Handler interface {
	Create(c *gin.Context)
	GetByID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	List(c *gin.Context)
}

type handler struct {
	service event.Service
	logger  *zap.SugaredLogger
}

// NewHandler creates a new event handler
func NewHandler(service event.Service, logger *zap.SugaredLogger) Handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

// Create handles creating a new event
func (h *handler) Create(c *gin.Context) {
	var req model.EventCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if req.Title == "" || req.StartDate.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title and start date are required"})
		return
	}

	id, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		h.logger.Errorw("Failed to create event", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// GetByID handles getting an event by ID
func (h *handler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Errorw("Invalid event ID", "error", err, "id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	event, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Errorw("Failed to get event", "error", err, "id", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	c.JSON(http.StatusOK, event)
}

// Update handles updating an event
func (h *handler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Errorw("Invalid event ID", "error", err, "id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	var req model.EventUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(c.Request.Context(), id, req); err != nil {
		h.logger.Errorw("Failed to update event", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Delete handles deleting an event
func (h *handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Errorw("Invalid event ID", "error", err, "id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		h.logger.Errorw("Failed to delete event", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// List handles listing events with pagination and filtering
func (h *handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	
	// Get optional user_id filter
	userIDStr := c.Query("user_id")
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

	var events model.EventsResponse
	var err error

	// If user_id is provided, filter events by user
	if userID != nil {
		events, err = h.service.ListByUser(c.Request.Context(), *userID, page, size)
	} else {
		events, err = h.service.List(c.Request.Context(), page, size)
	}

	if err != nil {
		h.logger.Errorw("Failed to list events", "error", err, "page", page, "size", size, "user_id", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}
