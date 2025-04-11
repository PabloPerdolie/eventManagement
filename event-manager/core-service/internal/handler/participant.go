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

// Create godoc
// @Summary Добавить участника в событие
// @Description Добавляет нового участника в событие
// @Tags participants
// @Accept json
// @Produce json
// @Param event_id path int true "ID события"
// @Param request body domain.EventParticipantCreateRequest true "Данные для добавления участника (необходимо указать либо user_id, либо username)"
// @Success 201 {object} domain.EventParticipantResponse "Информация о созданном участии"
// @Failure 400 {object} map[string]interface{} "Ошибка валидации или некорректный ID события"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /events/{event_id}/participants [post]
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

	// Проверка наличия либо user_id, либо username
	if req.UserID == 0 && req.Username == "" {
		h.logger.Errorw("Missing user identification", "error", "either user_id or username must be provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "either user_id or username must be provided"})
		return
	}

	participant, err := h.service.Create(c.Request.Context(), eventId, req)
	if err != nil {
		h.logger.Errorw("Failed to create event participant", "error", err, "eventId", eventId)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, participant)
}

// Delete godoc
// @Summary Удалить участника из события
// @Description Удаляет участника из события
// @Tags participants
// @Produce json
// @Param event_part_id path int true "ID участия в событии"
// @Success 204 "Участник успешно удален"
// @Failure 400 {object} map[string]interface{} "Некорректный ID участия"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /events/participants/{event_part_id} [delete]
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

// ListByUser godoc
// @Summary Получить список событий пользователя
// @Description Возвращает список событий, в которых пользователь является участником
// @Tags participants
// @Produce json
// @Param X-User-Id header string true "ID пользователя"
// @Param page query int false "Номер страницы (по умолчанию: 1)"
// @Param size query int false "Размер страницы (по умолчанию: 10)"
// @Success 200 {object} domain.EventParticipantsResponse "Список участий в событиях"
// @Failure 400 {object} map[string]interface{} "Некорректный ID пользователя"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /participants/user [get]
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

// ConfirmParticipation godoc
// @Summary Подтвердить участие в событии
// @Description Подтверждает участие пользователя в событии
// @Tags participants
// @Produce json
// @Param event_part_id path int true "ID участия в событии"
// @Success 204 "Участие успешно подтверждено"
// @Failure 400 {object} map[string]interface{} "Некорректный ID участия"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /events/participants/{event_part_id}/confirm [put]
func (h *ParticipantController) ConfirmParticipation(c *gin.Context) {
	eventPartIdStr := c.Param("event_part_id")
	eventPartId, err := strconv.Atoi(eventPartIdStr)
	if err != nil {
		h.logger.Errorw("Invalid event participant Id", "error", err, "id", eventPartIdStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event participant Id"})
		return
	}

	if err := h.service.ConfirmParticipation(c.Request.Context(), eventPartId); err != nil {
		h.logger.Errorw("Failed to confirm participation", "error", err, "eventPartId", eventPartId)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// DeclineParticipation godoc
// @Summary Отклонить участие в событии
// @Description Отклоняет участие пользователя в событии
// @Tags participants
// @Produce json
// @Param event_part_id path int true "ID участия в событии"
// @Success 204 "Участие успешно отклонено"
// @Failure 400 {object} map[string]interface{} "Некорректный ID участия"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /events/participants/{event_part_id}/decline [put]
func (h *ParticipantController) DeclineParticipation(c *gin.Context) {
	eventPartIdStr := c.Param("event_part_id")
	eventPartId, err := strconv.Atoi(eventPartIdStr)
	if err != nil {
		h.logger.Errorw("Invalid event participant Id", "error", err, "id", eventPartIdStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event participant Id"})
		return
	}

	if err := h.service.DeclineParticipation(c.Request.Context(), eventPartId); err != nil {
		h.logger.Errorw("Failed to decline participation", "error", err, "eventPartId", eventPartId)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
