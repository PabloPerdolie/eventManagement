package expense

//
//import (
//	"net/http"
//	"strconv"
//
//	"github.com/PabloPerdolie/event-manager/core-service/internal/model"
//	"github.com/PabloPerdolie/event-manager/core-service/internal/service/expense"
//	"github.com/gin-gonic/gin"
//	"github.com/google/uuid"
//	"go.uber.org/zap"
//)
//
//// ShareHandler handles expense share-related HTTP requests
//type ShareHandler interface {
//	Create(c *gin.Context)
//	GetByID(c *gin.Context)
//	Update(c *gin.Context)
//	Delete(c *gin.Context)
//	ListByExpense(c *gin.Context)
//}
//
//type shareHandler struct {
//	service expense.ShareService
//	logger  *zap.SugaredLogger
//}
//
//// NewShareHandler creates a new expense share handler
//func NewShareHandler(service expense.ShareService, logger *zap.SugaredLogger) ShareHandler {
//	return &shareHandler{
//		service: service,
//		logger:  logger,
//	}
//}
//
//// Create handles bulk creating expense shares
//func (h *shareHandler) Create(c *gin.Context) {
//	expenseIDStr := c.Param("id")
//	expenseID, err := uuid.Parse(expenseIDStr)
//	if err != nil {
//		h.logger.Errorw("Invalid expense ID", "error", err, "id", expenseIDStr)
//		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expense ID"})
//		return
//	}
//
//	// For creating shares, we'll use a different request structure
//	// that allows creating multiple shares at once
//	var req struct {
//		ParticipantIDs []uuid.UUID              `json:"participant_ids"`
//		SplitMethod    model.ExpenseSplitMethod `json:"split_method"`
//	}
//
//	if err := c.ShouldBindJSON(&req); err != nil {
//		h.logger.Errorw("Failed to bind request", "error", err)
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	// Validate request
//	if len(req.ParticipantIDs) == 0 {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "participant IDs are required"})
//		return
//	}
//
//	if err := h.service.CreateShares(c.Request.Context(), expenseID, req.ParticipantIDs, req.SplitMethod); err != nil {
//		h.logger.Errorw("Failed to create expense shares", "error", err, "expenseId", expenseID)
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.Status(http.StatusCreated)
//}
//
//// GetByID handles getting an expense share by ID
//func (h *shareHandler) GetByID(c *gin.Context) {
//	expenseIDStr := c.Param("id")
//	shareIDStr := c.Param("share_id")
//
//	_, err := uuid.Parse(expenseIDStr)
//	if err != nil {
//		h.logger.Errorw("Invalid expense ID", "error", err, "id", expenseIDStr)
//		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expense ID"})
//		return
//	}
//
//	shareID, err := uuid.Parse(shareIDStr)
//	if err != nil {
//		h.logger.Errorw("Invalid share ID", "error", err, "id", shareIDStr)
//		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid share ID"})
//		return
//	}
//
//	share, err := h.service.GetByID(c.Request.Context(), shareID)
//	if err != nil {
//		h.logger.Errorw("Failed to get expense share", "error", err, "id", shareID)
//		c.JSON(http.StatusNotFound, gin.H{"error": "expense share not found"})
//		return
//	}
//
//	c.JSON(http.StatusOK, share)
//}
//
//// Update handles updating an expense share
//func (h *shareHandler) Update(c *gin.Context) {
//	expenseIDStr := c.Param("id")
//	shareIDStr := c.Param("share_id")
//
//	_, err := uuid.Parse(expenseIDStr)
//	if err != nil {
//		h.logger.Errorw("Invalid expense ID", "error", err, "id", expenseIDStr)
//		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expense ID"})
//		return
//	}
//
//	shareID, err := uuid.Parse(shareIDStr)
//	if err != nil {
//		h.logger.Errorw("Invalid share ID", "error", err, "id", shareIDStr)
//		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid share ID"})
//		return
//	}
//
//	var req model.ExpenseShareUpdateRequest
//	if err := c.ShouldBindJSON(&req); err != nil {
//		h.logger.Errorw("Failed to bind request", "error", err)
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	if err := h.service.Update(c.Request.Context(), shareID, req); err != nil {
//		h.logger.Errorw("Failed to update expense share", "error", err, "id", shareID)
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.Status(http.StatusNoContent)
//}
//
//// Delete handles deleting an expense share
//func (h *shareHandler) Delete(c *gin.Context) {
//	// In practice, we may not want to allow direct deletion of individual shares
//	// as they're usually managed as a group. But we'll add this endpoint for completeness.
//	c.JSON(http.StatusNotImplemented, gin.H{"error": "deletion of individual expense shares is not supported"})
//}
//
//// ListByExpense handles listing expense shares for a specific expense
//func (h *shareHandler) ListByExpense(c *gin.Context) {
//	expenseIDStr := c.Param("id")
//	expenseID, err := uuid.Parse(expenseIDStr)
//	if err != nil {
//		h.logger.Errorw("Invalid expense ID", "error", err, "id", expenseIDStr)
//		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expense ID"})
//		return
//	}
//
//	shares, err := h.service.ListByExpense(c.Request.Context(), expenseID)
//	if err != nil {
//		h.logger.Errorw("Failed to list expense shares", "error", err, "expenseId", expenseID)
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"shares": shares})
//}
