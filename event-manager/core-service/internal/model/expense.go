package model

import "time"

type Expense struct {
	ExpenseID   int       `db:"expense_id"`
	EventID     int       `db:"event_id"`
	CreatedBy   int       `db:"created_by"`
	Description string    `db:"description"`
	Amount      float64   `db:"amount"`
	Currency    string    `db:"currency"`
	SplitMethod string    `db:"split_method"`
	CreatedAt   time.Time `db:"created_at"`
}
