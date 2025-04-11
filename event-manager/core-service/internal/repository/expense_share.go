package repository

import (
	"context"
	"database/sql"
	"github.com/PabloPerdolie/event-manager/core-service/internal/domain"
	"github.com/PabloPerdolie/event-manager/core-service/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
)

type ExpenseShare struct {
	db *sqlx.DB
}

func NewExpenseShare(db *sqlx.DB) ExpenseShare {
	return ExpenseShare{
		db: db,
	}
}

// CreateExpenseShare создает новую запись о доле расхода
func (r ExpenseShare) CreateExpenseShare(ctx context.Context, share model.ExpenseShare) (int, error) {
	query := `
		INSERT INTO expense_share (expense_id, user_id, amount, is_paid, paid_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING share_id
	`
	var id int
	err := r.db.QueryRowContext(ctx, query,
		share.ExpenseID,
		share.UserID,
		share.Amount,
		share.IsPaid,
		share.PaidAt,
	).Scan(&id)
	if err != nil {
		return 0, errors.WithMessage(err, "create expense share")
	}

	return id, nil
}

// GetExpenseShareById получает долю расхода по ID
func (r ExpenseShare) GetExpenseShareById(ctx context.Context, id int) (*model.ExpenseShare, error) {
	query := `
		SELECT share_id, expense_id, user_id, amount, is_paid, paid_at
		FROM expense_share
		WHERE share_id = $1
	`
	var share model.ExpenseShare
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&share.ShareID,
		&share.ExpenseID,
		&share.UserID,
		&share.Amount,
		&share.IsPaid,
		&share.PaidAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("expense share not found")
	}
	if err != nil {
		return nil, errors.WithMessage(err, "get expense share by id")
	}

	return &share, nil
}

// ListExpenseSharesByExpenseId получает все доли для определенного расхода
func (r ExpenseShare) ListExpenseSharesByExpenseId(ctx context.Context, expenseId int) ([]model.ExpenseShare, error) {
	query := `
		SELECT share_id, expense_id, user_id, amount, is_paid, paid_at
		FROM expense_share
		WHERE expense_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, expenseId)
	if err != nil {
		return nil, errors.WithMessage(err, "list expense shares by expense id")
	}
	defer rows.Close()

	shares := make([]model.ExpenseShare, 0)
	for rows.Next() {
		var share model.ExpenseShare
		err := rows.Scan(
			&share.ShareID,
			&share.ExpenseID,
			&share.UserID,
			&share.Amount,
			&share.IsPaid,
			&share.PaidAt,
		)
		if err != nil {
			return nil, errors.WithMessage(err, "scan expense share")
		}
		shares = append(shares, share)
	}

	if rows.Err() != nil {
		return nil, errors.WithMessage(rows.Err(), "rows err")
	}

	return shares, nil
}

// ListExpenseSharesByUserId получает все доли расходов конкретного пользователя
func (r ExpenseShare) ListExpenseSharesByUserId(ctx context.Context, userId int) ([]model.ExpenseShare, error) {
	query := `
		SELECT share_id, expense_id, user_id, amount, is_paid, paid_at
		FROM expense_share
		WHERE user_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, errors.WithMessage(err, "list expense shares by user id")
	}
	defer rows.Close()

	shares := make([]model.ExpenseShare, 0)
	for rows.Next() {
		var share model.ExpenseShare
		err := rows.Scan(
			&share.ShareID,
			&share.ExpenseID,
			&share.UserID,
			&share.Amount,
			&share.IsPaid,
			&share.PaidAt,
		)
		if err != nil {
			return nil, errors.WithMessage(err, "scan expense share")
		}
		shares = append(shares, share)
	}

	if rows.Err() != nil {
		return nil, errors.WithMessage(rows.Err(), "rows err")
	}

	return shares, nil
}

// UpdateExpenseSharePaidStatus обновляет статус оплаты доли расхода
func (r ExpenseShare) UpdateExpenseSharePaidStatus(ctx context.Context, shareId int, isPaid bool) error {
	var paidAt interface{} = nil
	if isPaid {
		paidAt = time.Now()
	}

	query := `
		UPDATE expense_share
		SET is_paid = $1, paid_at = $2
		WHERE share_id = $3
	`
	_, err := r.db.ExecContext(ctx, query, isPaid, paidAt, shareId)
	if err != nil {
		return errors.WithMessage(err, "update expense share paid status")
	}

	return nil
}

// DeleteExpenseSharesByExpenseId удаляет все доли для указанного расхода
func (r ExpenseShare) DeleteExpenseSharesByExpenseId(ctx context.Context, expenseId int) error {
	query := `DELETE FROM expense_share WHERE expense_id = $1`
	_, err := r.db.ExecContext(ctx, query, expenseId)
	if err != nil {
		return errors.WithMessage(err, "delete expense shares")
	}

	return nil
}

// GetUserBalanceInEvent рассчитывает баланс пользователя в событии
// Положительный баланс означает, что пользователь должен получить деньги
// Отрицательный - должен заплатить
func (r ExpenseShare) GetUserBalanceInEvent(ctx context.Context, userId int, eventId int) (float64, error) {
	query := `
		WITH user_expenses AS (
			SELECT e.expense_id, e.amount, e.created_by
			FROM expense e
			WHERE e.event_id = $1
		),
		user_shares AS (
			SELECT es.expense_id, es.amount
			FROM expense_share es
			JOIN user_expenses ue ON es.expense_id = ue.expense_id
			WHERE es.user_id = $2
		)
		SELECT 
			COALESCE((SELECT SUM(amount) FROM user_expenses WHERE created_by = $2), 0) -
			COALESCE((SELECT SUM(amount) FROM user_shares), 0) as balance
	`
	var balance float64
	err := r.db.QueryRowContext(ctx, query, eventId, userId).Scan(&balance)
	if err != nil {
		return 0, errors.WithMessage(err, "get user balance in event")
	}

	return balance, nil
}

// GetEventBalanceReport формирует отчет о балансе всех участников события
func (r ExpenseShare) GetEventBalanceReport(ctx context.Context, eventId int) ([]domain.UserBalance, error) {
	query := `
		WITH event_users AS (
			SELECT DISTINCT ep.user_id
			FROM event_participant ep
			WHERE ep.event_id = $1
		),
		user_expenses AS (
			SELECT e.expense_id, e.amount, e.created_by
			FROM expense e
			WHERE e.event_id = $1
		),
		user_shares AS (
			SELECT es.user_id, es.expense_id, es.amount
			FROM expense_share es
			JOIN user_expenses ue ON es.expense_id = ue.expense_id
		),
		user_balances AS (
			SELECT 
				eu.user_id,
				COALESCE((SELECT SUM(amount) FROM user_expenses WHERE created_by = eu.user_id), 0) -
				COALESCE((SELECT SUM(amount) FROM user_shares WHERE user_id = eu.user_id), 0) as balance
			FROM event_users eu
		)
		SELECT ub.user_id, u.username, ub.balance
		FROM user_balances ub
		JOIN users u ON ub.user_id = u.user_id
		ORDER BY ub.balance DESC
	`

	rows, err := r.db.QueryContext(ctx, query, eventId)
	if err != nil {
		return nil, errors.WithMessage(err, "get event balance report")
	}
	defer rows.Close()

	balances := make([]domain.UserBalance, 0)
	for rows.Next() {
		var balance domain.UserBalance
		err := rows.Scan(
			&balance.UserID,
			&balance.Username,
			&balance.Balance,
		)
		if err != nil {
			return nil, errors.WithMessage(err, "scan user balance")
		}
		balances = append(balances, balance)
	}

	if rows.Err() != nil {
		return nil, errors.WithMessage(rows.Err(), "rows err")
	}

	return balances, nil
}
