package expense_worker

import (
	"context"
	"github.com/txix-open/isp-kit/worker"
	"go.uber.org/zap"
	"time"
)

const (
	defaultInterval = time.Duration(1) * time.Hour
)

func StartExpenseJob(
	ctx context.Context,
	expenseRepo ExpenseRepo,
	notifyPbl NotifyPublisher,
	logger *zap.SugaredLogger,
) *worker.Worker {
	expenseWorker := NewExpenseWorker(expenseRepo, notifyPbl, logger)
	w := worker.New(expenseWorker, worker.WithInterval(defaultInterval))
	w.Run(ctx)

	logger.Info("Started expense notification worker with 1 hour interval")
	return w
}
