package expense

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/PabloPerdolie/event-manager/core-service/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// ShareRepository defines expense share repository interface
type ShareRepository interface {
	Create(ctx context.Context, share model.ExpenseShare) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.ExpenseShare, error)
	Update(ctx context.Context, share model.ExpenseShare) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByExpense(ctx context.Context, expenseID uuid.UUID) error
	ListByExpense(ctx context.Context, expenseID uuid.UUID) ([]model.ExpenseShare, error)
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.ExpenseShare, int, error)
	GetUserDebts(ctx context.Context, userID uuid.UUID) (model.UserDebtSummary, error)
	GetEventUserDebts(ctx context.Context, eventID, userID uuid.UUID) (model.EventDebtSummary, error)
}

type shareRepository struct {
	db *sqlx.DB
}

// NewShareRepository creates a new expense share repository
func NewShareRepository(db *sqlx.DB) ShareRepository {
	return &shareRepository{db: db}
}

// Create creates a new expense share in the database
func (r *shareRepository) Create(ctx context.Context, share model.ExpenseShare) (uuid.UUID, error) {
	share.ID = uuid.New()

	query := `
		INSERT INTO expense_shares (id, expense_id, user_id, share_amount, is_paid, paid_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		share.ID,
		share.ExpenseID,
		share.UserID,
		share.ShareAmount,
		share.IsPaid,
		share.PaidAt,
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create expense share: %w", err)
	}

	return share.ID, nil
}

// GetByID retrieves an expense share by ID
func (r *shareRepository) GetByID(ctx context.Context, id uuid.UUID) (model.ExpenseShare, error) {
	var share model.ExpenseShare

	query := `
		SELECT id, expense_id, user_id, share_amount, is_paid, paid_at
		FROM expense_shares
		WHERE id = $1
	`

	err := r.db.GetContext(ctx, &share, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.ExpenseShare{}, fmt.Errorf("expense share not found: %w", err)
		}
		return model.ExpenseShare{}, fmt.Errorf("failed to get expense share: %w", err)
	}

	return share, nil
}

// Update updates an expense share in the database
func (r *shareRepository) Update(ctx context.Context, share model.ExpenseShare) error {
	query := `
		UPDATE expense_shares
		SET share_amount = $1, is_paid = $2, paid_at = $3
		WHERE id = $4
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		share.ShareAmount,
		share.IsPaid,
		share.PaidAt,
		share.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update expense share: %w", err)
	}

	return nil
}

// Delete deletes an expense share from the database
func (r *shareRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM expense_shares WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete expense share: %w", err)
	}

	return nil
}

// DeleteByExpense deletes all expense shares for a specific expense
func (r *shareRepository) DeleteByExpense(ctx context.Context, expenseID uuid.UUID) error {
	query := `DELETE FROM expense_shares WHERE expense_id = $1`

	_, err := r.db.ExecContext(ctx, query, expenseID)
	if err != nil {
		return fmt.Errorf("failed to delete expense shares: %w", err)
	}

	return nil
}

// ListByExpense retrieves all expense shares for a specific expense
func (r *shareRepository) ListByExpense(ctx context.Context, expenseID uuid.UUID) ([]model.ExpenseShare, error) {
	var shares []model.ExpenseShare

	query := `
		SELECT id, expense_id, user_id, share_amount, is_paid, paid_at
		FROM expense_shares
		WHERE expense_id = $1
		ORDER BY share_amount DESC
	`

	err := r.db.SelectContext(ctx, &shares, query, expenseID)
	if err != nil {
		return nil, fmt.Errorf("failed to list expense shares: %w", err)
	}

	return shares, nil
}

// ListByUser retrieves a list of expense shares for a specific user with pagination
func (r *shareRepository) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.ExpenseShare, int, error) {
	var shares []model.ExpenseShare
	var total int

	// Count total shares for the user
	countQuery := `SELECT COUNT(*) FROM expense_shares WHERE user_id = $1`
	err := r.db.GetContext(ctx, &total, countQuery, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count user expense shares: %w", err)
	}

	// Retrieve shares with pagination
	query := `
		SELECT es.id, es.expense_id, es.user_id, es.share_amount, es.is_paid, es.paid_at
		FROM expense_shares es
		JOIN expenses e ON es.expense_id = e.id
		WHERE es.user_id = $1
		ORDER BY e.expense_date DESC, e.created_at DESC
		LIMIT $2 OFFSET $3
	`

	err = r.db.SelectContext(ctx, &shares, query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list user expense shares: %w", err)
	}

	return shares, total, nil
}

// GetUserDebts retrieves a summary of a user's debts
func (r *shareRepository) GetUserDebts(ctx context.Context, userID uuid.UUID) (model.UserDebtSummary, error) {
	var summary model.UserDebtSummary
	summary.UserID = userID

	// Calculate total owed, paid, and unpaid amounts
	query := `
		SELECT 
			COALESCE(SUM(share_amount), 0) as total_owed,
			COALESCE(SUM(CASE WHEN is_paid THEN share_amount ELSE 0 END), 0) as total_paid,
			COALESCE(SUM(CASE WHEN NOT is_paid THEN share_amount ELSE 0 END), 0) as total_unpaid
		FROM expense_shares
		WHERE user_id = $1
	`

	type totalSummary struct {
		TotalOwed   float64 `db:"total_owed"`
		TotalPaid   float64 `db:"total_paid"`
		TotalUnpaid float64 `db:"total_unpaid"`
	}

	var totals totalSummary
	err := r.db.GetContext(ctx, &totals, query, userID)
	if err != nil {
		return summary, fmt.Errorf("failed to calculate user debt summary: %w", err)
	}

	summary.TotalOwed = totals.TotalOwed
	summary.TotalPaid = totals.TotalPaid
	summary.TotalUnpaid = totals.TotalUnpaid

	return summary, nil
}

// GetEventUserDebts retrieves a summary of a user's debts for a specific event
func (r *shareRepository) GetEventUserDebts(ctx context.Context, eventID, userID uuid.UUID) (model.EventDebtSummary, error) {
	var summary model.EventDebtSummary
	summary.EventID = eventID

	// Calculate total owed, paid, and unpaid amounts for the event
	query := `
		SELECT 
			COALESCE(SUM(es.share_amount), 0) as total_owed,
			COALESCE(SUM(CASE WHEN es.is_paid THEN es.share_amount ELSE 0 END), 0) as total_paid,
			COALESCE(SUM(CASE WHEN NOT es.is_paid THEN es.share_amount ELSE 0 END), 0) as total_unpaid
		FROM expense_shares es
		JOIN expenses e ON es.expense_id = e.id
		WHERE e.event_id = $1 AND es.user_id = $2
	`

	type totalSummary struct {
		TotalOwed   float64 `db:"total_owed"`
		TotalPaid   float64 `db:"total_paid"`
		TotalUnpaid float64 `db:"total_unpaid"`
	}

	var totals totalSummary
	err := r.db.GetContext(ctx, &totals, query, eventID, userID)
	if err != nil {
		return summary, fmt.Errorf("failed to calculate event user debt summary: %w", err)
	}

	summary.TotalOwed = totals.TotalOwed
	summary.TotalPaid = totals.TotalPaid
	summary.TotalUnpaid = totals.TotalUnpaid

	return summary, nil
}
