package handler

import (
	"context"
	"github.com/PabloPerdolie/event-manager/communication-service/internal/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type CommentService interface {
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

// GetCommentsByEventId godoc
// @Summary Получить комментарии события
// @Description Возвращает все комментарии, связанные с указанным событием
// @Tags comments
// @Produce json
// @Param eventId path int true "ID события"
// @Success 200 {array} model.Comment "Список комментариев"
// @Failure 400 {object} map[string]interface{} "Некорректный ID события"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /comments/event/{eventId} [get]
func (h *CommentHandler) GetCommentsByEventId(c *gin.Context) {
	eventIdParam := c.Param("event_id")
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

// DeleteComment godoc
// @Summary Удалить комментарий
// @Description Удаляет комментарий по его идентификатору
// @Tags comments
// @Produce json
// @Param id path int true "ID комментария"
// @Success 204 "Комментарий успешно удален"
// @Failure 400 {object} map[string]interface{} "Некорректный ID комментария"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /comments/{id} [delete]
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

// MarkCommentAsRead godoc
// @Summary Отметить комментарий как прочитанный
// @Description Отмечает комментарий как прочитанный по его идентификатору
// @Tags comments
// @Produce json
// @Param id path int true "ID комментария"
// @Success 204 "Комментарий успешно отмечен как прочитанный"
// @Failure 400 {object} map[string]interface{} "Некорректный ID комментария"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /comments/{id}/read [put]
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
