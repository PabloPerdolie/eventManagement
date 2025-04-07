package handler

import (
	"context"
	"github.com/PabloPerdolie/event-manager/core-service/internal/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TaskService interface {
	Create(ctx context.Context, req domain.TaskCreateRequest) (*domain.TaskResponse, error)
	Update(ctx context.Context, id int, req domain.TaskUpdateRequest) error
	Delete(ctx context.Context, id int) error
	ListByEvent(ctx context.Context, eventId int, page, size int) (*domain.TasksResponse, error)
	ListByUser(ctx context.Context, userId int, page, size int) (*domain.TasksResponse, error)
	UpdateStatus(ctx context.Context, id int, status domain.TaskStatus) error
}
type TaskController struct {
	service TaskService
	logger  *zap.SugaredLogger
}

func NewTask(service TaskService, logger *zap.SugaredLogger) TaskController {
	return TaskController{
		service: service,
		logger:  logger,
	}
}

// Create godoc
// @Summary Создать новую задачу
// @Description Создает новую задачу в системе
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body domain.TaskCreateRequest true "Данные для создания задачи"
// @Success 201 {object} map[string]interface{} "Возвращает ID созданной задачи"
// @Failure 400 {object} map[string]interface{} "Ошибка валидации"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /tasks [post]
func (h *TaskController) Create(c *gin.Context) {
	var req domain.TaskCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

// Update godoc
// @Summary Обновить задачу
// @Description Обновляет существующую задачу по ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param task_id path int true "ID задачи"
// @Param request body domain.TaskUpdateRequest true "Данные для обновления задачи"
// @Success 204 "Задача успешно обновлена"
// @Failure 400 {object} map[string]interface{} "Ошибка валидации или некорректный ID задачи"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /tasks/{task_id} [put]
func (h *TaskController) Update(c *gin.Context) {
	idStr := c.Param("task_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Errorw("Invalid task Id", "error", err, "task_id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task Id"})
		return
	}

	var req domain.TaskUpdateRequest
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

// Delete godoc
// @Summary Удалить задачу
// @Description Удаляет задачу по ID
// @Tags tasks
// @Produce json
// @Param task_id path int true "ID задачи"
// @Success 204 "Задача успешно удалена"
// @Failure 400 {object} map[string]interface{} "Некорректный ID задачи"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /tasks/{task_id} [delete]
func (h *TaskController) Delete(c *gin.Context) {
	idStr := c.Param("task_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Errorw("Invalid task Id", "error", err, "task_id", idStr)
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

// List godoc
// @Summary Получить список задач
// @Description Возвращает список задач с пагинацией, может фильтровать по событию или пользователю
// @Tags tasks
// @Produce json
// @Param page query int false "Номер страницы (по умолчанию: 1)"
// @Param size query int false "Размер страницы (по умолчанию: 10)"
// @Param event_id query int false "ID события для фильтрации задач"
// @Param X-User-Id header string false "ID пользователя для фильтрации задач"
// @Success 200 {object} domain.TasksResponse "Список задач"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /tasks [get]
func (h *TaskController) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	eventIdStr := c.Query("event_id")
	userIdStr := c.GetHeader("X-User-Id")

	var eventId *int
	if eventIdStr != "" {
		parsed, err := strconv.Atoi(eventIdStr)
		if err == nil {
			eventId = &parsed
		} else {
			h.logger.Warnw("Invalid event Id filter", "error", err, "event_id", eventIdStr)
		}
	}

	var userId *int
	if userIdStr != "" {
		parsed, err := strconv.Atoi(userIdStr)
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

	var tasks *domain.TasksResponse
	var err error

	if eventId != nil {
		tasks, err = h.service.ListByEvent(c.Request.Context(), *eventId, page, size)
	} else if userId != nil {
		tasks, err = h.service.ListByUser(c.Request.Context(), *userId, page, size)
	} else {
		// Default to all tasks (for admin purposes)
		//tasks, err = h.service.List(c.Request.Context(), page, size)
	}

	if err != nil {
		h.logger.Errorw("Failed to list tasks", "error", err, "page", page, "size", size, "event_id", eventId, "user_id", userId)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}
