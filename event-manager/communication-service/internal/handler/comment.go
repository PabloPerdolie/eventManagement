package handler

import (
	"context"
	"github.com/PabloPerdolie/event-manager/communication-service/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type CommentService interface {
	GetCommentById(ctx context.Context, id int) (model.Comment, error)
	GetCommentsByEventId(ctx context.Context, eventId int) ([]model.Comment, error)
	DeleteComment(ctx context.Context, id int) error
	MarkCommentAsRead(ctx context.Context, id int) error
}

type CommentHandler struct {
	service CommentService
	logger  *zap.SugaredLogger
}

func NewComment(service CommentService, logger *zap.SugaredLogger) *CommentHandler {
	return &CommentHandler{
		service: service,
		logger:  logger,
	}
}

func (h *CommentHandler) GetCommentById(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id parameter"})
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
		return
	}

	comment, err := h.service.GetCommentById(c, id)
	if err != nil {
		if errors.Is(err, errors.New("comment not found")) {
			c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
			return
		}
		h.logger.Errorw("failed to get comment", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, comment)
}

func (h *CommentHandler) GetCommentsByEventId(c *gin.Context) {
	eventIdParam := c.Param("eventId")
	if eventIdParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing eventId parameter"})
		return
	}

	eventId, err := strconv.Atoi(eventIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid eventId parameter"})
		return
	}

	comments, err := h.service.GetCommentsByEventId(c, eventId)
	if err != nil {
		h.logger.Errorw("failed to get comments", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, comments)
}

func (h *CommentHandler) DeleteComment(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id parameter"})
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
		return
	}

	err = h.service.DeleteComment(c, id)
	if err != nil {
		h.logger.Errorw("failed to delete comment", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *CommentHandler) MarkCommentAsRead(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id parameter"})
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
		return
	}

	err = h.service.MarkCommentAsRead(c, id)
	if err != nil {
		h.logger.Errorw("failed to mark comment as read", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}
