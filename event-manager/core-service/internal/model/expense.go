package model

import "time"

// Структуры для работы с БД
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

// ExpenseShare представляет долю расхода для конкретного пользователя
type ExpenseShare struct {
	ShareID   int        `db:"share_id"`
	ExpenseID int        `db:"expense_id"`
	UserID    int        `db:"user_id"`
	Amount    float64    `db:"amount"`
	IsPaid    bool       `db:"is_paid"`
	PaidAt    *time.Time `db:"paid_at"`
}

// UserBalance представляет баланс пользователя в событии
type UserBalance struct {
	UserID       int     `db:"user_id"`
	Username     string  `db:"username"`
	Balance      float64 `db:"balance"`
	PaidAmount   float64 `db:"paid_amount"`
	UnpaidAmount float64 `db:"unpaid_amount"`
	TotalDue     float64 `db:"total_due"`
}

// Константы для методов разделения расходов
const (
	SplitMethodEqual   = "equal"   // Поровну между всеми
	SplitMethodPercent = "percent" // Процентное соотношение
	SplitMethodExact   = "exact"   // Точные суммы
)
