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

// ParticipantHandler handles event participant-related HTTP requests
type ParticipantHandler interface {
	Create(c *gin.Context)
	GetByID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	ListByEvent(c *gin.Context)
}

type participantHandler struct {
	service event.ParticipantService
	logger  *zap.SugaredLogger
}

// NewParticipantHandler creates a new event participant handler
func NewParticipantHandler(service event.ParticipantService, logger *zap.SugaredLogger) ParticipantHandler {
	return &participantHandler{
		service: service,
		logger:  logger,
	}
}

// Create handles creating a new event participant
func (h *participantHandler) Create(c *gin.Context) {
	eventIDStr := c.Param("id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		h.logger.Errorw("Invalid event ID", "error", err, "id", eventIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	var req model.EventParticipantCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.EventID = eventID

	id, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		h.logger.Errorw("Failed to create event participant", "error", err, "eventId", eventID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// GetByID handles getting an event participant by event ID and user ID
func (h *participantHandler) GetByID(c *gin.Context) {
	eventIDStr := c.Param("id")
	userIDStr := c.Param("user_id")
	
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		h.logger.Errorw("Invalid event ID", "error", err, "id", eventIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}
	
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.logger.Errorw("Invalid user ID", "error", err, "id", userIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	participant, err := h.service.GetByEventAndUserID(c.Request.Context(), eventID, userID)
	if err != nil {
		h.logger.Errorw("Failed to get event participant", "error", err, "eventId", eventID, "userId", userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "event participant not found"})
		return
	}

	c.JSON(http.StatusOK, participant)
}

// Update handles updating an event participant
func (h *participantHandler) Update(c *gin.Context) {
	eventIDStr := c.Param("id")
	userIDStr := c.Param("user_id")
	
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		h.logger.Errorw("Invalid event ID", "error", err, "id", eventIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}
	
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.logger.Errorw("Invalid user ID", "error", err, "id", userIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req model.EventParticipantUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set event and user IDs
	req.EventID = eventID
	req.UserID = userID

	if err := h.service.Update(c.Request.Context(), req); err != nil {
		h.logger.Errorw("Failed to update event participant", "error", err, "eventId", eventID, "userId", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Delete handles deleting an event participant
func (h *participantHandler) Delete(c *gin.Context) {
	eventIDStr := c.Param("id")
	userIDStr := c.Param("user_id")
	
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		h.logger.Errorw("Invalid event ID", "error", err, "id", eventIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}
	
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.logger.Errorw("Invalid user ID", "error", err, "id", userIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), eventID, userID); err != nil {
		h.logger.Errorw("Failed to delete event participant", "error", err, "eventId", eventID, "userId", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListByEvent handles listing event participants for a specific event
func (h *participantHandler) ListByEvent(c *gin.Context) {
	eventIDStr := c.Param("id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		h.logger.Errorw("Invalid event ID", "error", err, "id", eventIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	participants, err := h.service.ListByEvent(c.Request.Context(), eventID, page, size)
	if err != nil {
		h.logger.Errorw("Failed to list event participants", "error", err, "eventId", eventID, "page", page, "size", size)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, participants)
}
