package user

import (
	"net/http"
	"strconv"

	"github.com/event-management/core-service/internal/model"
	"github.com/event-management/core-service/internal/service/user"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Handler handles user-related HTTP requests
type Handler interface {
	Create(c *gin.Context)
	GetByID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	List(c *gin.Context)
}

type handler struct {
	service user.Service
	logger  *zap.SugaredLogger
}

// NewHandler creates a new user handler
func NewHandler(service user.Service, logger *zap.SugaredLogger) Handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

// Create handles creating a new user
func (h *handler) Create(c *gin.Context) {
	var req model.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if req.Username == "" || req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and email are required"})
		return
	}

	// If role is not provided, set a default role
	if req.Role == "" {
		req.Role = model.UserRoleUser
	}

	id, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		h.logger.Errorw("Failed to create user", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// GetByID handles getting a user by ID
func (h *handler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Errorw("Invalid user ID", "error", err, "id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Errorw("Failed to get user", "error", err, "id", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Update handles updating a user
func (h *handler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Errorw("Invalid user ID", "error", err, "id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req model.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(c.Request.Context(), id, req); err != nil {
		h.logger.Errorw("Failed to update user", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Delete handles deleting a user
func (h *handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Errorw("Invalid user ID", "error", err, "id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		h.logger.Errorw("Failed to delete user", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// List handles listing users with pagination
func (h *handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	users, err := h.service.List(c.Request.Context(), page, size)
	if err != nil {
		h.logger.Errorw("Failed to list users", "error", err, "page", page, "size", size)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}
