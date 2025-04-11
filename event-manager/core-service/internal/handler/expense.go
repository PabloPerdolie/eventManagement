package handler

import (
	"errors"
	"github.com/PabloPerdolie/event-manager/core-service/internal/domain"
	"github.com/PabloPerdolie/event-manager/core-service/internal/service/expense"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ExpenseController struct {
	service expense.Service
}

func NewExpenseController(service expense.Service) ExpenseController {
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

	// Получаем ID пользователя из контекста запроса
	userId := ctx.GetInt("user_id")
	if userId == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Устанавливаем создателя расхода
	req.CreatedBy = userId

	expenseId, err := c.service.CreateExpense(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": expenseId})
}

// Update godoc
// @Summary Update expense
// @Description Update an existing expense
// @Tags expenses
// @Accept json
// @Produce json
// @Param id path int true "Expense ID"
// @Param expense body domain.ExpenseUpdateRequest true "Expense update data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /expenses/{id} [put]
func (c ExpenseController) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid expense id"})
		return
	}

	var req domain.ExpenseUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.UpdateExpense(ctx, id, req); err != nil {
		if errors.Is(err, errors.New("expense not found")) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "expense not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "expense updated successfully"})
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
