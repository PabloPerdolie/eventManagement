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
		SELECT
    		(
                -- Сумма расходов, созданных пользователем
                COALESCE(SUM(CASE WHEN e.created_by = $2 THEN e.amount ELSE 0 END), 0) 
                -- Минус сумма неоплаченных долей пользователя
                - COALESCE(SUM(CASE WHEN es.user_id = $2 AND es.is_paid = false THEN es.amount ELSE 0 END), 0)
            ) AS balance
		FROM expense e
		LEFT JOIN expense_share es ON es.expense_id = e.expense_id
		WHERE e.event_id = $1
	`
	var balance float64
	err := r.db.QueryRowContext(ctx, query, eventId, userId).Scan(&balance)
	if err != nil {
		return 0, errors.WithMessage(err, "get user balance in event")
	}

	return balance, nil
}

func (r ExpenseShare) GetEventBalanceReport(ctx context.Context, eventId int) ([]domain.UserBalance, error) {
	query := `
		SELECT 
    		ep.user_id,
    		u.username,
    		(
                -- Сумма расходов, которые создал пользователь
                COALESCE(SUM(CASE WHEN e.created_by = ep.user_id THEN e.amount ELSE 0 END), 0) 
                -- Минус сумма долей, которые пользователь должен оплатить
                - COALESCE(SUM(CASE WHEN es.user_id = ep.user_id AND es.is_paid = false THEN es.amount ELSE 0 END), 0)
                -- Учитываем только неоплаченные доли, оплаченные уже не влияют на баланс
            ) AS balance
		FROM event_participant ep
		JOIN users u ON u.user_id = ep.user_id
		LEFT JOIN expense e ON e.event_id = ep.event_id
		LEFT JOIN expense_share es ON es.expense_id = e.expense_id AND es.user_id = ep.user_id
		WHERE ep.event_id = $1
		GROUP BY ep.user_id, u.username
		ORDER BY balance DESC
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
