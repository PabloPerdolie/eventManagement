package handler

import (
	"context"
	"github.com/PabloPerdolie/event-manager/core-service/internal/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type EventCommonService interface {
	GetEventSummary(ctx context.Context, eventId int) (*domain.EventData, error)
}

type EventService interface {
	Create(ctx context.Context, userId int, req domain.EventCreateRequest) (int, error)
	Update(ctx context.Context, id int, req domain.EventUpdateRequest) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, page, size int) (*domain.EventsResponse, error)
	ListByOrganizer(ctx context.Context, organizerId int, page, size int) (*domain.EventsResponse, error)
	ListByParticipant(ctx context.Context, participantId int, page, size int) (*domain.EventsResponse, error)
}

type EventController struct {
	commonService EventCommonService
	service       EventService
	logger        *zap.SugaredLogger
}

func NewEvent(commonService EventCommonService, service EventService, logger *zap.SugaredLogger) EventController {
	return EventController{
		commonService: commonService,
		service:       service,
		logger:        logger,
	}
}

// Create godoc
// @Summary Создать новое событие
// @Description Создает новое событие с указанным пользователем в качестве организатора
// @Tags events
// @Accept json
// @Produce json
// @Param X-User-Id header string true "ID пользователя-организатора"
// @Param request body domain.EventCreateRequest true "Данные для создания события"
// @Success 201 {object} map[string]interface{} "Возвращает ID созданного события"
// @Failure 400 {object} map[string]interface{} "Ошибка валидации"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /events [post]
func (h *EventController) Create(c *gin.Context) {
	idStr := c.GetHeader("X-User-Id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Errorw("Invalid user id", "error", err, "user_id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req domain.EventCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	eventId, err := h.service.Create(c.Request.Context(), id, req)
	if err != nil {
		h.logger.Errorw("Failed to create event", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": eventId})
}

// Delete godoc
// @Summary Удалить событие
// @Description Удаляет событие по ID
// @Tags events
// @Produce json
// @Param event_id path int true "ID события"
// @Success 200 "Операция успешна"
// @Failure 400 {object} map[string]interface{} "Некорректный ID события"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /events/{event_id} [delete]
func (h *EventController) Delete(c *gin.Context) {
	idStr := c.Param("event_id")
	id, err := strconv.Atoi(idStr)
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

	c.Status(http.StatusOK)
}

// List godoc
// @Summary Получить список событий
// @Description Возвращает список событий с пагинацией, может фильтровать по участию пользователя
// @Tags events
// @Produce json
// @Param page query int false "Номер страницы (по умолчанию: 1)"
// @Param size query int false "Размер страницы (по умолчанию: 10)"
// @Param X-User-Id header string false "ID пользователя для фильтрации по участию"
// @Success 200 {object} domain.EventsResponse "Список событий"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /events [get]
func (h *EventController) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	userIDStr := c.GetHeader("X-User-Id")
	var userID *int
	if userIDStr != "" {
		parsed, err := strconv.Atoi(userIDStr)
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

	var events *domain.EventsResponse
	var err error

	if userID != nil {
		events, err = h.service.ListByParticipant(c.Request.Context(), *userID, page, size)
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

// EventSummary godoc
// @Summary Получить детальную информацию о событии
// @Description Возвращает детальную информацию о событии, включая список участников и задач
// @Tags events
// @Produce json
// @Param event_id path int true "ID события"
// @Success 200 {object} domain.EventData "Детальная информация о событии"
// @Failure 400 {object} map[string]interface{} "Некорректный ID события"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /events/{event_id} [get]
func (h *EventController) EventSummary(c *gin.Context) {
	idStr := c.Param("event_id")
	eventId, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Errorw("Invalid event ID", "error", err, "id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	summary, err := h.commonService.GetEventSummary(c.Request.Context(), eventId)
	if err != nil {
		h.logger.Errorw("Failed to get event summary", "error", err, "event_id", eventId)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}
