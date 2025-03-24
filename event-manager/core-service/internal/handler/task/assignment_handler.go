package task

import (
	"net/http"
	"strconv"

	"github.com/PabloPerdolie/event-manager/core-service/internal/model"
	"github.com/PabloPerdolie/event-manager/core-service/internal/service/task"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// AssignmentHandler handles task assignment-related HTTP requests
type AssignmentHandler interface {
	Create(c *gin.Context)
	GetByID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	ListByTask(c *gin.Context)
}

type assignmentHandler struct {
	service task.AssignmentService
	logger  *zap.SugaredLogger
}

// NewAssignmentHandler creates a new task assignment handler
func NewAssignmentHandler(service task.AssignmentService, logger *zap.SugaredLogger) AssignmentHandler {
	return &assignmentHandler{
		service: service,
		logger:  logger,
	}
}

// Create handles creating a new task assignment
func (h *assignmentHandler) Create(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		h.logger.Errorw("Invalid task ID", "error", err, "id", taskIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	var req model.TaskAssignmentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.TaskID = taskID

	id, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		h.logger.Errorw("Failed to create task assignment", "error", err, "taskId", taskID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// GetByID handles getting a task assignment by ID
func (h *assignmentHandler) GetByID(c *gin.Context) {
	taskIDStr := c.Param("id")
	assignmentIDStr := c.Param("assignment_id")

	_, err := uuid.Parse(taskIDStr)
	if err != nil {
		h.logger.Errorw("Invalid task ID", "error", err, "id", taskIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	assignmentID, err := uuid.Parse(assignmentIDStr)
	if err != nil {
		h.logger.Errorw("Invalid assignment ID", "error", err, "id", assignmentIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assignment ID"})
		return
	}

	assignment, err := h.service.GetByID(c.Request.Context(), assignmentID)
	if err != nil {
		h.logger.Errorw("Failed to get task assignment", "error", err, "id", assignmentID)
		c.JSON(http.StatusNotFound, gin.H{"error": "task assignment not found"})
		return
	}

	c.JSON(http.StatusOK, assignment)
}

// Update handles updating a task assignment
func (h *assignmentHandler) Update(c *gin.Context) {
	taskIDStr := c.Param("id")
	assignmentIDStr := c.Param("assignment_id")

	_, err := uuid.Parse(taskIDStr)
	if err != nil {
		h.logger.Errorw("Invalid task ID", "error", err, "id", taskIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	assignmentID, err := uuid.Parse(assignmentIDStr)
	if err != nil {
		h.logger.Errorw("Invalid assignment ID", "error", err, "id", assignmentIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assignment ID"})
		return
	}

	var req model.TaskAssignmentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(c.Request.Context(), assignmentID, req); err != nil {
		h.logger.Errorw("Failed to update task assignment", "error", err, "id", assignmentID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Delete handles deleting a task assignment
func (h *assignmentHandler) Delete(c *gin.Context) {
	taskIDStr := c.Param("id")
	assignmentIDStr := c.Param("assignment_id")

	_, err := uuid.Parse(taskIDStr)
	if err != nil {
		h.logger.Errorw("Invalid task ID", "error", err, "id", taskIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	assignmentID, err := uuid.Parse(assignmentIDStr)
	if err != nil {
		h.logger.Errorw("Invalid assignment ID", "error", err, "id", assignmentIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assignment ID"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), assignmentID); err != nil {
		h.logger.Errorw("Failed to delete task assignment", "error", err, "id", assignmentID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListByTask handles listing task assignments for a specific task
func (h *assignmentHandler) ListByTask(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		h.logger.Errorw("Invalid task ID", "error", err, "id", taskIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
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

	assignments, err := h.service.ListByTask(c.Request.Context(), taskID, page, size)
	if err != nil {
		h.logger.Errorw("Failed to list task assignments", "error", err, "taskId", taskID, "page", page, "size", size)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, assignments)
}
