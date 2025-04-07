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
