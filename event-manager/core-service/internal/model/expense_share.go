package model

import (
	"time"

	"github.com/google/uuid"
)

// ExpenseShare represents a share of an expense assigned to a user
type ExpenseShare struct {
	ID          uuid.UUID `json:"id" db:"id"`
	ExpenseID   uuid.UUID `json:"expense_id" db:"expense_id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	ShareAmount float64   `json:"share_amount" db:"share_amount"`
	IsPaid      bool      `json:"is_paid" db:"is_paid"`
	PaidAt      *time.Time `json:"paid_at,omitempty" db:"paid_at"`
}

// ExpenseShareCreateRequest represents the input for creating a new expense share
type ExpenseShareCreateRequest struct {
	UserID      uuid.UUID `json:"user_id" binding:"required"`
	ShareAmount float64   `json:"share_amount" binding:"required,gt=0"`
}

// ExpenseShareUpdateRequest represents the input for updating an expense share
type ExpenseShareUpdateRequest struct {
	ShareAmount *float64    `json:"share_amount" binding:"omitempty,gt=0"`
	IsPaid      *bool       `json:"is_paid"`
	PaidAt      *time.Time  `json:"paid_at"`
}

// ExpenseShareResponse represents the output for expense share data
type ExpenseShareResponse struct {
	ID          uuid.UUID    `json:"id"`
	ExpenseID   uuid.UUID    `json:"expense_id"`
	User        UserResponse `json:"user"`
	ShareAmount float64      `json:"share_amount"`
	IsPaid      bool         `json:"is_paid"`
	PaidAt      *time.Time   `json:"paid_at,omitempty"`
}

// ExpenseSharesResponse represents a list of expense shares
type ExpenseSharesResponse struct {
	Shares []ExpenseShareResponse `json:"shares"`
	Total  int                    `json:"total"`
}

// UserDebtSummary represents a summary of a user's debts
type UserDebtSummary struct {
	UserID       uuid.UUID              `json:"user_id"`
	User         UserResponse           `json:"user"`
	TotalOwed    float64                `json:"total_owed"`
	TotalPaid    float64                `json:"total_paid"`
	TotalUnpaid  float64                `json:"total_unpaid"`
	DebtByEvent  []EventDebtSummary     `json:"debt_by_event,omitempty"`
}

// EventDebtSummary represents a summary of debts for a specific event
type EventDebtSummary struct {
	EventID     uuid.UUID          `json:"event_id"`
	Event       EventResponse      `json:"event"`
	TotalOwed   float64            `json:"total_owed"`
	TotalPaid   float64            `json:"total_paid"`
	TotalUnpaid float64            `json:"total_unpaid"`
	Expenses    []ExpenseResponse  `json:"expenses,omitempty"`
}
