package model

import (
	"time"

	"github.com/google/uuid"
)

// SplitMethod represents the method used to split an expense
type SplitMethod string

const (
	SplitMethodEqual  SplitMethod = "equal"
	SplitMethodCustom SplitMethod = "custom"
)

// Currency represents the currency of an expense
type Currency string

const (
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
	CurrencyRUB Currency = "RUB"
)

// Expense represents an expense for an event
type Expense struct {
	ID          uuid.UUID   `json:"id" db:"id"`
	EventID     uuid.UUID   `json:"event_id" db:"event_id"`
	Description string      `json:"description" db:"description"`
	Amount      float64     `json:"amount" db:"amount"`
	Currency    Currency    `json:"currency" db:"currency"`
	ExpenseDate time.Time   `json:"expense_date" db:"expense_date"`
	CreatedBy   uuid.UUID   `json:"created_by" db:"created_by"`
	SplitMethod SplitMethod `json:"split_method" db:"split_method"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
}

// ExpenseCreateRequest represents the input for creating a new expense
type ExpenseCreateRequest struct {
	Description string      `json:"description" binding:"required"`
	Amount      float64     `json:"amount" binding:"required,gt=0"`
	Currency    Currency    `json:"currency" binding:"required"`
	ExpenseDate time.Time   `json:"expense_date" binding:"required"`
	SplitMethod SplitMethod `json:"split_method" binding:"required"`
	Shares      []ExpenseShareCreateRequest `json:"shares,omitempty"`
}

// ExpenseUpdateRequest represents the input for updating an expense
type ExpenseUpdateRequest struct {
	Description *string      `json:"description"`
	Amount      *float64     `json:"amount" binding:"omitempty,gt=0"`
	Currency    *Currency    `json:"currency"`
	ExpenseDate *time.Time   `json:"expense_date"`
	SplitMethod *SplitMethod `json:"split_method"`
}

// ExpenseResponse represents the output for expense data
type ExpenseResponse struct {
	ID          uuid.UUID        `json:"id"`
	EventID     uuid.UUID        `json:"event_id"`
	Description string           `json:"description"`
	Amount      float64          `json:"amount"`
	Currency    Currency         `json:"currency"`
	ExpenseDate time.Time        `json:"expense_date"`
	CreatedBy   uuid.UUID        `json:"created_by"`
	SplitMethod SplitMethod      `json:"split_method"`
	CreatedAt   time.Time        `json:"created_at"`
	Shares      []ExpenseShare   `json:"shares,omitempty"`
}

// ExpensesResponse represents a list of expenses
type ExpensesResponse struct {
	Expenses []ExpenseResponse `json:"expenses"`
	Total    int               `json:"total"`
}
