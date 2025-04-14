package repository

import (
	"context"
	"database/sql"
	"github.com/PabloPerdolie/event-manager/core-service/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Expense struct {
	db *sqlx.DB
}

func NewExpense(db *sqlx.DB) Expense {
	return Expense{
		db: db,
	}
}

func (r Expense) CreateExpense(ctx context.Context, expense model.Expense) (int, error) {
	query := `
		INSERT INTO expense (event_id, created_by, description, amount, currency, split_method, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP)
		RETURNING expense_id
	`
	var id int
	err := r.db.QueryRowContext(ctx, query,
		expense.EventID,
		expense.CreatedBy,
		expense.Description,
		expense.Amount,
		expense.Currency,
		expense.SplitMethod,
	).Scan(&id)
	if err != nil {
		return 0, errors.WithMessage(err, "create expense")
	}

	return id, nil
}

func (r Expense) DeleteExpense(ctx context.Context, id int) error {
	query := `DELETE FROM expense WHERE expense_id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.WithMessage(err, "delete expense")
	}

	return nil
}

func (r Expense) GetExpenseById(ctx context.Context, id int) (*model.Expense, error) {
	query := `
		SELECT expense_id, event_id, created_by, description, amount, currency, split_method, created_at
		FROM expense
		WHERE expense_id = $1
	`
	var expense model.Expense
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&expense.ExpenseID,
		&expense.EventID,
		&expense.CreatedBy,
		&expense.Description,
		&expense.Amount,
		&expense.Currency,
		&expense.SplitMethod,
		&expense.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("expense not found")
	}
	if err != nil {
		return nil, errors.WithMessage(err, "get expense by id")
	}

	return &expense, nil
}

func (r Expense) ListExpensesByEventId(ctx context.Context, eventId int, limit, offset int) ([]model.Expense, int, error) {
	countQuery := `SELECT COUNT(*) FROM expense WHERE event_id = $1`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, eventId).Scan(&total)
	if err != nil {
		return nil, 0, errors.WithMessage(err, "count expenses")
	}

	query := `
		SELECT expense_id, event_id, created_by, description, amount, currency, split_method, created_at
		FROM expense
		WHERE event_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, eventId, limit, offset)
	if err != nil {
		return nil, 0, errors.WithMessage(err, "list expenses by event")
	}
	defer rows.Close()

	expenses := make([]model.Expense, 0)
	for rows.Next() {
		var expense model.Expense
		err := rows.Scan(
			&expense.ExpenseID,
			&expense.EventID,
			&expense.CreatedBy,
			&expense.Description,
			&expense.Amount,
			&expense.Currency,
			&expense.SplitMethod,
			&expense.CreatedAt,
		)
		if err != nil {
			return nil, 0, errors.WithMessage(err, "scan expense")
		}
		expenses = append(expenses, expense)
	}

	if rows.Err() != nil {
		return nil, 0, errors.WithMessage(rows.Err(), "rows err")
	}

	return expenses, total, nil
}

func (r Expense) GetEventTotalExpenses(ctx context.Context, eventId int) (float64, error) {
	query := `SELECT COALESCE(SUM(amount), 0) FROM expense WHERE event_id = $1`
	var total float64
	err := r.db.QueryRowContext(ctx, query, eventId).Scan(&total)
	if err != nil {
		return 0, errors.WithMessage(err, "get event total expenses")
	}

	return total, nil
}
