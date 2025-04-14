package handler

import (
	"context"
	"errors"
	"github.com/PabloPerdolie/event-manager/core-service/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ExpenseService interface {
	CreateExpense(ctx context.Context, req domain.ExpenseCreateRequest) (int, error)
	UpdateExpenseSharePaidStatus(ctx context.Context, shareId int, isPaid bool) error
	DeleteExpense(ctx context.Context, id int) error
}

type ExpenseController struct {
	service ExpenseService
}

func NewExpenseController(service ExpenseService) ExpenseController {
	return ExpenseController{
		service: service,
	}
}

// Create godoc
// @Summary Create a new expense
// @Description Create a new expense with participants
// @Tags expenses
// @Accept json
// @Produce json
// @Param expense body domain.ExpenseCreateRequest true "Expense creation data"
// @Success 201 {object} map[string]int
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /expenses [post]
func (c ExpenseController) Create(ctx *gin.Context) {
	var req domain.ExpenseCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idStr := ctx.GetHeader("X-User-Id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	// Устанавливаем создателя расхода
	req.CreatedBy = id

	expenseId, err := c.service.CreateExpense(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": expenseId})
}

// Delete godoc
// @Summary Delete expense
// @Description Delete an existing expense
// @Tags expenses
// @Produce json
// @Param id path int true "Expense ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /expenses/{id} [delete]
func (c ExpenseController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid expense id"})
		return
	}

	if err := c.service.DeleteExpense(ctx, id); err != nil {
		if errors.Is(err, errors.New("expense not found")) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "expense not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "expense deleted successfully"})
}

// UpdateExpenseSharePaidStatus godoc
// @Summary Update expense share paid status
// @Description Update the paid status of an expense share
// @Tags expenses
// @Accept json
// @Produce json
// @Param id path int true "Share ID"
// @Param request body domain.ExpenseShareUpdateRequest true "Paid status update data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /expenses/shares/{id}/paid-status [put]
func (c ExpenseController) UpdateExpenseSharePaidStatus(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный id доли расхода"})
		return
	}

	var req domain.ExpenseShareUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.UpdateExpenseSharePaidStatus(ctx, id, req.IsPaid); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "статус оплаты доли расхода успешно обновлен"})
}
