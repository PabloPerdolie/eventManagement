package expense

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/event-management/core-service/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Repository defines expense repository interface
type Repository interface {
	Create(ctx context.Context, expense model.Expense) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Expense, error)
	Update(ctx context.Context, expense model.Expense) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByEvent(ctx context.Context, eventID uuid.UUID, limit, offset int) ([]model.Expense, int, error)
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Expense, int, error)
	GetEventTotalExpenses(ctx context.Context, eventID uuid.UUID) (float64, error)
}

type repository struct {
	db *sqlx.DB
}

// NewRepository creates a new expense repository
func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

// Create creates a new expense in the database
func (r *repository) Create(ctx context.Context, expense model.Expense) (uuid.UUID, error) {
	expense.ID = uuid.New()
	expense.CreatedAt = time.Now()

	query := `
		INSERT INTO expenses (id, event_id, description, amount, currency, expense_date, created_by, split_method, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		expense.ID,
		expense.EventID,
		expense.Description,
		expense.Amount,
		expense.Currency,
		expense.ExpenseDate,
		expense.CreatedBy,
		expense.SplitMethod,
		expense.CreatedAt,
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create expense: %w", err)
	}

	return expense.ID, nil
}

// GetByID retrieves an expense by ID
func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (model.Expense, error) {
	var expense model.Expense

	query := `
		SELECT id, event_id, description, amount, currency, expense_date, created_by, split_method, created_at
		FROM expenses
		WHERE id = $1
	`

	err := r.db.GetContext(ctx, &expense, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Expense{}, fmt.Errorf("expense not found: %w", err)
		}
		return model.Expense{}, fmt.Errorf("failed to get expense: %w", err)
	}

	return expense, nil
}

// Update updates an expense in the database
func (r *repository) Update(ctx context.Context, expense model.Expense) error {
	query := `
		UPDATE expenses
		SET description = $1, amount = $2, currency = $3, expense_date = $4, split_method = $5
		WHERE id = $6
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		expense.Description,
		expense.Amount,
		expense.Currency,
		expense.ExpenseDate,
		expense.SplitMethod,
		expense.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update expense: %w", err)
	}

	return nil
}

// Delete deletes an expense from the database
func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM expenses WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete expense: %w", err)
	}

	return nil
}

// ListByEvent retrieves a list of expenses for a specific event with pagination
func (r *repository) ListByEvent(ctx context.Context, eventID uuid.UUID, limit, offset int) ([]model.Expense, int, error) {
	var expenses []model.Expense
	var total int

	// Count total expenses for the event
	countQuery := `SELECT COUNT(*) FROM expenses WHERE event_id = $1`
	err := r.db.GetContext(ctx, &total, countQuery, eventID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count event expenses: %w", err)
	}

	// Retrieve expenses with pagination
	query := `
		SELECT id, event_id, description, amount, currency, expense_date, created_by, split_method, created_at
		FROM expenses
		WHERE event_id = $1
		ORDER BY expense_date DESC, created_at DESC
		LIMIT $2 OFFSET $3
	`

	err = r.db.SelectContext(ctx, &expenses, query, eventID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list event expenses: %w", err)
	}

	return expenses, total, nil
}

// ListByUser retrieves a list of expenses created by a specific user with pagination
func (r *repository) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Expense, int, error) {
	var expenses []model.Expense
	var total int

	// Count total expenses for the user
	countQuery := `SELECT COUNT(*) FROM expenses WHERE created_by = $1`
	err := r.db.GetContext(ctx, &total, countQuery, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count user expenses: %w", err)
	}

	// Retrieve expenses with pagination
	query := `
		SELECT id, event_id, description, amount, currency, expense_date, created_by, split_method, created_at
		FROM expenses
		WHERE created_by = $1
		ORDER BY expense_date DESC, created_at DESC
		LIMIT $2 OFFSET $3
	`

	err = r.db.SelectContext(ctx, &expenses, query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list user expenses: %w", err)
	}

	return expenses, total, nil
}

// GetEventTotalExpenses calculates the total amount of expenses for an event
func (r *repository) GetEventTotalExpenses(ctx context.Context, eventID uuid.UUID) (float64, error) {
	var total float64

	query := `
		SELECT COALESCE(SUM(amount), 0) as total
		FROM expenses
		WHERE event_id = $1
	`

	err := r.db.GetContext(ctx, &total, query, eventID)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate event total expenses: %w", err)
	}

	return total, nil
}
