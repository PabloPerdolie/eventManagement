package model

import (
	"time"

	"github.com/google/uuid"
)

// BudgetItem represents a budget item for an event
type BudgetItem struct {
	ID          uuid.UUID `json:"id" db:"id"`
	EventID     uuid.UUID `json:"event_id" db:"event_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Amount      float64   `json:"amount" db:"amount"`
	Category    string    `json:"category" db:"category"`
	PaidBy      uuid.UUID `json:"paid_by" db:"paid_by"`
	SplitType   string    `json:"split_type" db:"split_type"` // Equal, Percentage, Custom
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// BudgetItemCreateRequest represents the input for creating a new budget item
type BudgetItemCreateRequest struct {
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount" binding:"required"`
	Category    string    `json:"category"`
	PaidBy      uuid.UUID `json:"paid_by" binding:"required"`
	SplitType   string    `json:"split_type" binding:"required"`
	SplitDetails []SplitDetail `json:"split_details,omitempty"`
}

// BudgetItemUpdateRequest represents the input for updating a budget item
type BudgetItemUpdateRequest struct {
	Name        *string   `json:"name"`
	Description *string   `json:"description"`
	Amount      *float64  `json:"amount"`
	Category    *string   `json:"category"`
	PaidBy      *uuid.UUID `json:"paid_by"`
	SplitType   *string   `json:"split_type"`
	SplitDetails []SplitDetail `json:"split_details,omitempty"`
}

// SplitDetail represents how a budget item is split between participants
type SplitDetail struct {
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	Percentage  float64   `json:"percentage" db:"percentage"`
	FixedAmount float64   `json:"fixed_amount" db:"fixed_amount"`
}

// BudgetResponse represents the budget summary for an event
type BudgetResponse struct {
	TotalAmount  float64      `json:"total_amount"`
	Items        []BudgetItem `json:"items"`
	UserBalances []Balance    `json:"user_balances"`
}

// Balance represents a user's balance in an event
type Balance struct {
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
	Paid   float64   `json:"paid"`
	Owes   float64   `json:"owes"`
	Net    float64   `json:"net"`
}
