package event

import (
	"net/http"
	"strconv"

	"github.com/PabloPerdolie/event-manager/core-service/internal/model"
	"github.com/PabloPerdolie/event-manager/core-service/internal/service/event"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ParticipantHandler handles event participant-related HTTP requests
type ParticipantHandler interface {
	Create(c *gin.Context)
	GetById(c *gin.Context)
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
	eventIdStr := c.Param("id")
	eventId, err := uuid.Parse(eventIdStr)
	if err != nil {
		h.logger.Errorw("Invalid event Id", "error", err, "id", eventIdStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event Id"})
		return
	}

	var req model.EventParticipantCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.EventId = eventId

	id, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		h.logger.Errorw("Failed to create event participant", "error", err, "eventId", eventId)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// GetById handles getting an event participant by event Id and user Id
func (h *participantHandler) GetById(c *gin.Context) {
	eventIdStr := c.Param("id")
	userIdStr := c.Param("user_id")

	eventId, err := uuid.Parse(eventIdStr)
	if err != nil {
		h.logger.Errorw("Invalid event Id", "error", err, "id", eventIdStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event Id"})
		return
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		h.logger.Errorw("Invalid user Id", "error", err, "id", userIdStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user Id"})
		return
	}

	participant, err := h.service.GetByEventAndUserId(c.Request.Context(), eventId, userId)
	if err != nil {
		h.logger.Errorw("Failed to get event participant", "error", err, "eventId", eventId, "userId", userId)
		c.JSON(http.StatusNotFound, gin.H{"error": "event participant not found"})
		return
	}

	c.JSON(http.StatusOK, participant)
}

// Update handles updating an event participant
func (h *participantHandler) Update(c *gin.Context) {
	eventIdStr := c.Param("id")
	userIdStr := c.Param("user_id")

	eventId, err := uuid.Parse(eventIdStr)
	if err != nil {
		h.logger.Errorw("Invalid event Id", "error", err, "id", eventIdStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event Id"})
		return
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		h.logger.Errorw("Invalid user Id", "error", err, "id", userIdStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user Id"})
		return
	}

	var req model.EventParticipantUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set event and user Ids
	req.EventId = eventId
	req.UserId = userId

	if err := h.service.Update(c.Request.Context(), req); err != nil {
		h.logger.Errorw("Failed to update event participant", "error", err, "eventId", eventId, "userId", userId)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Delete handles deleting an event participant
func (h *participantHandler) Delete(c *gin.Context) {
	eventIdStr := c.Param("id")
	userIdStr := c.Param("user_id")

	eventId, err := uuid.Parse(eventIdStr)
	if err != nil {
		h.logger.Errorw("Invalid event Id", "error", err, "id", eventIdStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event Id"})
		return
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		h.logger.Errorw("Invalid user Id", "error", err, "id", userIdStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user Id"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), eventId, userId); err != nil {
		h.logger.Errorw("Failed to delete event participant", "error", err, "eventId", eventId, "userId", userId)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListByEvent handles listing event participants for a specific event
func (h *participantHandler) ListByEvent(c *gin.Context) {
	eventIdStr := c.Param("id")
	eventId, err := uuid.Parse(eventIdStr)
	if err != nil {
		h.logger.Errorw("Invalid event Id", "error", err, "id", eventIdStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event Id"})
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

	participants, err := h.service.ListByEvent(c.Request.Context(), eventId, page, size)
	if err != nil {
		h.logger.Errorw("Failed to list event participants", "error", err, "eventId", eventId, "page", page, "size", size)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, participants)
}
