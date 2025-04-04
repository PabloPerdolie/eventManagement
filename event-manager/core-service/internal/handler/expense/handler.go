package expense

import (
	"net/http"
	"strconv"

	"github.com/PabloPerdolie/event-manager/core-service/internal/model"
	"github.com/PabloPerdolie/event-manager/core-service/internal/service/expense"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Handler handles expense-related HTTP requests
type Handler interface {
	Create(c *gin.Context)
	GetById(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	List(c *gin.Context)
}

type handler struct {
	service expense.Service
	logger  *zap.SugaredLogger
}

// NewHandler creates a new expense handler
func NewHandler(service expense.Service, logger *zap.SugaredLogger) Handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

// Create handles creating a new expense
func (h *handler) Create(c *gin.Context) {
	var req model.ExpenseCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if req.Description == "" || req.Amount <= 0 || req.EventId == uuid.Nil || req.CreatedBy == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "description, amount, event Id, and created by are required"})
		return
	}

	id, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		h.logger.Errorw("Failed to create expense", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// GetById handles getting an expense by Id
func (h *handler) GetById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Errorw("Invalid expense Id", "error", err, "id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expense Id"})
		return
	}

	expense, err := h.service.GetById(c.Request.Context(), id)
	if err != nil {
		h.logger.Errorw("Failed to get expense", "error", err, "id", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "expense not found"})
		return
	}

	c.JSON(http.StatusOK, expense)
}

// Update handles updating an expense
func (h *handler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Errorw("Invalid expense Id", "error", err, "id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expense Id"})
		return
	}

	var req model.ExpenseUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(c.Request.Context(), id, req); err != nil {
		h.logger.Errorw("Failed to update expense", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Delete handles deleting an expense
func (h *handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Errorw("Invalid expense Id", "error", err, "id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expense Id"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		h.logger.Errorw("Failed to delete expense", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// List handles listing expenses with pagination and filtering
func (h *handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	// Get optional filters
	eventIdStr := c.Query("event_id")
	userIdStr := c.Query("user_id")

	var eventId *int
	if eventIdStr != "" {
		parsed, err := uuid.Parse(eventIdStr)
		if err == nil {
			eventId = &parsed
		} else {
			h.logger.Warnw("Invalid event Id filter", "error", err, "event_id", eventIdStr)
		}
	}

	var userId *int
	if userIdStr != "" {
		parsed, err := uuid.Parse(userIdStr)
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

	var expenses model.ExpensesResponse
	var err error

	// Apply filters if provided
	if eventId != nil {
		expenses, err = h.service.ListByEvent(c.Request.Context(), *eventId, page, size)
	} else if userId != nil {
		expenses, err = h.service.ListByUser(c.Request.Context(), *userId, page, size)
	} else {
		// Default to all expenses (only for admin purposes, should be restricted in the gateway)
		c.JSON(http.StatusBadRequest, gin.H{"error": "event_id or user_id parameter is required"})
		return
	}

	if err != nil {
		h.logger.Errorw("Failed to list expenses", "error", err, "page", page, "size", size, "event_id", eventId, "user_id", userId)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, expenses)
}
