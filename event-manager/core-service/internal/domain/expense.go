package domain

import (
	"github.com/PabloPerdolie/event-manager/core-service/internal/model"
	"time"
)

// UserBalance представляет баланс пользователя в событии
type UserBalance struct {
	UserID   int     `json:"user_id"`
	Username string  `json:"username"`
	Balance  float64 `json:"balance"`
}

// ExpenseCreateRequest запрос на создание расхода
type ExpenseCreateRequest struct {
	EventID     int     `json:"event_id" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Currency    string  `json:"currency" binding:"required"`
	SplitMethod string  `json:"split_method" binding:"required"`
	CreatedBy   int     `json:"created_by" binding:"required"`
	UserIDs     []int   `json:"user_ids" binding:"required"`
}

// ExpenseUpdateRequest запрос на обновление расхода
type ExpenseUpdateRequest struct {
	Description *string  `json:"description"`
	Amount      *float64 `json:"amount" binding:"omitempty,gt=0"`
	Currency    *string  `json:"currency"`
	SplitMethod *string  `json:"split_method"`
	UserIDs     *[]int   `json:"user_ids"`
}

// ExpenseResponse ответ с информацией о расходе
type ExpenseResponse struct {
	ExpenseID   int                  `json:"expense_id"`
	EventID     int                  `json:"event_id"`
	CreatedBy   int                  `json:"created_by"`
	Description string               `json:"description"`
	Amount      float64              `json:"amount"`
	Currency    string               `json:"currency"`
	SplitMethod string               `json:"split_method"`
	CreatedAt   time.Time            `json:"created_at"`
	Shares      []model.ExpenseShare `json:"shares,omitempty"`
}

// ExpensesResponse список расходов с пагинацией
type ExpensesResponse struct {
	Items      []ExpenseResponse `json:"items"`
	TotalCount int               `json:"total_count"`
}

// ExpenseShareUpdateRequest запрос на обновление статуса доли расхода
type ExpenseShareUpdateRequest struct {
	IsPaid bool `json:"is_paid" binding:"required"`
}

// BalanceReportResponse ответ с отчетом о балансах пользователей
type BalanceReportResponse struct {
	EventID      int           `json:"event_id"`
	TotalAmount  float64       `json:"total_amount"`
	UserBalances []UserBalance `json:"user_balances"`
}
