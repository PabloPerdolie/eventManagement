package handler

import (
	"context"
	"github.com/PabloPerdolie/event-manager/core-service/internal/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ParticipantService interface {
	ListByUser(ctx context.Context, userId int, page, size int) (*domain.EventParticipantsResponse, error)
	Create(ctx context.Context, eventID int, req domain.EventParticipantCreateRequest) (*domain.EventParticipantResponse, error)
	Delete(ctx context.Context, id int) error
	ConfirmParticipation(ctx context.Context, id int) error
	DeclineParticipation(ctx context.Context, id int) error
}

type ParticipantController struct {
	service ParticipantService
	logger  *zap.SugaredLogger
}

func NewParticipantHandler(service ParticipantService, logger *zap.SugaredLogger) ParticipantController {
	return ParticipantController{
		service: service,
		logger:  logger,
	}
}

func (h *ParticipantController) Create(c *gin.Context) {
	eventIdStr := c.Param("event_id")
	eventId, err := strconv.Atoi(eventIdStr)
	if err != nil {
		h.logger.Errorw("Invalid event Id", "error", err, "id", eventIdStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event Id"})
		return
	}

	var req domain.EventParticipantCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.service.Create(c.Request.Context(), eventId, req)
	if err != nil {
		h.logger.Errorw("Failed to create event participant", "error", err, "eventId", eventId)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *ParticipantController) Delete(c *gin.Context) {
	eventPartIdStr := c.Param("event_part_id")
	eventPartId, err := strconv.Atoi(eventPartIdStr)
	if err != nil {
		h.logger.Errorw("Invalid event Id", "error", err, "id", eventPartId)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event Id"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), eventPartId); err != nil {
		h.logger.Errorw("Failed to delete event participant", "error", err, "eventPartId", eventPartId)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ParticipantController) ListByUser(c *gin.Context) {
	userIdStr := c.GetHeader("X-User-Id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		h.logger.Errorw("Invalid user Id", "error", err, "id", userId)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user Id"})
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

	participants, err := h.service.ListByUser(c.Request.Context(), userId, page, size)
	if err != nil {
		h.logger.Errorw("Failed to list event participants", "error", err, "userId", userId, "page", page, "size", size)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, participants)
}

// todo confirm + decline participation
