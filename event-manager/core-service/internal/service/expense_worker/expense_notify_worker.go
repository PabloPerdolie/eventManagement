package expense_worker

import (
	"context"
	"encoding/json"
	"github.com/PabloPerdolie/event-manager/core-service/internal/model"
	"go.uber.org/zap"
)

type ExpenseRepo interface {
	ListAllUnpaidExpenseShares(ctx context.Context) ([]model.ExpenseShareUnpaid, error)
}

type NotifyPublisher interface {
	Publish(ctx context.Context, data []byte) error
}

type ExpenseWorker struct {
	expenseRepo ExpenseRepo
	notifyPbl   NotifyPublisher
	logger      *zap.SugaredLogger
}

func NewExpenseWorker(expenseRepo ExpenseRepo, notifyPbl NotifyPublisher, logger *zap.SugaredLogger) ExpenseWorker {
	return ExpenseWorker{
		expenseRepo: expenseRepo,
		notifyPbl:   notifyPbl,
		logger:      logger,
	}
}

func (w ExpenseWorker) Do(ctx context.Context) {
	w.logger.Info("Starting unpaid expense shares notification job")

	shares, err := w.expenseRepo.ListAllUnpaidExpenseShares(ctx)
	if err != nil {
		w.logger.Errorw("Failed to get unpaid expense shares", "error", err)
		return
	}

	w.logger.Infow("Found unpaid expense shares", "count", len(shares))

	for _, share := range shares {
		w.do(ctx, share)
	}

	w.logger.Info("Completed unpaid expense shares notification job")
}

func (w ExpenseWorker) do(ctx context.Context, share model.ExpenseShareUnpaid) {
	notification := map[string]any{
		"event": "unpaid_expense_share",
		"data": map[string]any{
			"user_email":          share.Email,
			"total_unpaid_amount": share.TotalUnpaidAmount,
			"event_id":            share.EventID,
		},
	}

	data, err := json.Marshal(notification)
	if err != nil {
		w.logger.Errorw("Failed to marshal notification", "error", err, "notification", notification)
		return
	}

	err = w.notifyPbl.Publish(ctx, data)
	if err != nil {
		w.logger.Errorw("Failed to publish notification", "error", err, "notification", notification)
		return
	}

	w.logger.Infow("Sent notification for unpaid expense share",
		"email", share.Email,
		"amount", share.TotalUnpaidAmount)
}
