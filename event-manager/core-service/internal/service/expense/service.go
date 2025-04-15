package expense

import (
	"context"
	"github.com/PabloPerdolie/event-manager/core-service/internal/domain"
	"github.com/PabloPerdolie/event-manager/core-service/internal/model"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"math"
)

type ExpenseRepository interface {
	CreateExpense(ctx context.Context, expense model.Expense) (int, error)
	GetExpenseById(ctx context.Context, id int) (*model.Expense, error)
	ListExpensesByEventId(ctx context.Context, eventId int, limit, offset int) ([]model.Expense, int, error)
	GetEventTotalExpenses(ctx context.Context, eventId int) (float64, error)
	DeleteExpense(ctx context.Context, id int) error
}

type ExpenseShareRepository interface {
	CreateExpenseShare(ctx context.Context, share model.ExpenseShare) (int, error)
	ListExpenseSharesByExpenseId(ctx context.Context, expenseId int) ([]model.ExpenseShare, error)
	UpdateExpenseSharePaidStatus(ctx context.Context, shareId int, isPaid bool) error
	DeleteExpenseSharesByExpenseId(ctx context.Context, expenseId int) error
	GetEventBalanceReport(ctx context.Context, eventId int) ([]domain.UserBalance, error)
}

type Service struct {
	expenseRepo      ExpenseRepository
	expenseShareRepo ExpenseShareRepository
	logger           *zap.SugaredLogger
}

func NewService(expenseRepo ExpenseRepository, expenseShareRepo ExpenseShareRepository, logger *zap.SugaredLogger) Service {
	return Service{
		expenseRepo:      expenseRepo,
		expenseShareRepo: expenseShareRepo,
		logger:           logger,
	}
}

func (s Service) CreateExpense(ctx context.Context, req domain.ExpenseCreateRequest) (int, error) {
	expense := model.Expense{
		EventID:     req.EventID,
		Description: req.Description,
		Amount:      req.Amount,
		Currency:    req.Currency,
		CreatedBy:   req.CreatedBy,
		SplitMethod: req.SplitMethod,
	}

	id, err := s.expenseRepo.CreateExpense(ctx, expense)
	if err != nil {
		s.logger.Errorw("Failed to create expense", "error", err, "eventId", req.EventID)
		return 0, errors.WithMessage(err, "create expense")
	}

	// Создаем доли расходов для участников
	if len(req.UserIDs) > 0 {
		err = s.CreateExpenseShares(ctx, id, req.UserIDs, req.SplitMethod)
		if err != nil {
			s.logger.Errorw("Failed to create expense shares", "error", err, "expenseId", id)
			// Мы не возвращаем ошибку, так как расход уже создан успешно
		}
	}

	return id, nil
}

func (s Service) CreateExpenseShares(ctx context.Context, expenseId int, userIds []int, splitMethod string) error {
	// Получаем расход для определения суммы для разделения
	expense, err := s.expenseRepo.GetExpenseById(ctx, expenseId)
	if err != nil {
		return errors.WithMessage(err, "failed to get expense for creating shares")
	}

	// Удаляем существующие доли для этого расхода
	err = s.expenseShareRepo.DeleteExpenseSharesByExpenseId(ctx, expenseId)
	if err != nil {
		return errors.WithMessage(err, "failed to delete existing expense shares")
	}

	// Создаем доли в зависимости от метода разделения
	switch splitMethod {
	case model.SplitMethodEqual, model.SplitMethodPercent:
		// Для equal и percent (пока) реализуем как равное разделение
		if len(userIds) == 0 {
			return errors.New("cannot split expense with no participants")
		}

		shareAmount := expense.Amount / float64(len(userIds))
		// Округляем до 2 знаков после запятой
		shareAmount = math.Round(shareAmount*100) / 100

		for _, userId := range userIds {
			share := model.ExpenseShare{
				ExpenseID: expenseId,
				UserID:    userId,
				Amount:    shareAmount,
				IsPaid:    false,
			}

			_, err := s.expenseShareRepo.CreateExpenseShare(ctx, share)
			if err != nil {
				s.logger.Warnw("Failed to create expense share", "error", err, "expenseId", expenseId, "userId", userId)
				// Продолжаем, даже если одна доля не создалась
			}
		}

	case model.SplitMethodExact:
		// Для точных сумм нам нужны дополнительные данные
		// Пока возвращаем ошибку
		return errors.New("exact split method requires additional data")

	default:
		return errors.Errorf("unsupported split method: %s", splitMethod)
	}

	return nil
}

func (s Service) ListExpensesByEvent(ctx context.Context, eventId int, page, size int) (domain.ExpensesResponse, error) {
	offset := (page - 1) * size
	expenses, total, err := s.expenseRepo.ListExpensesByEventId(ctx, eventId, size, offset)
	if err != nil {
		s.logger.Errorw("Failed to list expenses by event", "error", err, "eventId", eventId)
		return domain.ExpensesResponse{}, errors.WithMessage(err, "failed to list expenses")
	}

	// Преобразуем в ответ
	expenseResponses := make([]domain.ExpenseResponse, len(expenses))
	for i, expense := range expenses {
		// Получаем доли для каждого расхода
		shares, err := s.expenseShareRepo.ListExpenseSharesByExpenseId(ctx, expense.ExpenseID)
		if err != nil {
			s.logger.Warnw("Failed to get expense shares", "error", err, "expenseId", expense.ExpenseID)
			shares = []model.ExpenseShare{}
		}

		expenseResponses[i] = domain.ExpenseResponse{
			ExpenseID:   expense.ExpenseID,
			EventID:     expense.EventID,
			Description: expense.Description,
			Amount:      expense.Amount,
			Currency:    expense.Currency,
			CreatedBy:   expense.CreatedBy,
			SplitMethod: expense.SplitMethod,
			CreatedAt:   expense.CreatedAt,
			Shares:      shares,
		}
	}

	return domain.ExpensesResponse{
		Items:      expenseResponses,
		TotalCount: total,
	}, nil
}

func (s Service) DeleteExpense(ctx context.Context, id int) error {
	// Сначала проверяем, существует ли расход
	_, err := s.expenseRepo.GetExpenseById(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to get expense for deletion", "error", err, "id", id)
		return errors.WithMessage(err, "failed to get expense")
	}

	// Удаляем все доли расхода
	err = s.expenseShareRepo.DeleteExpenseSharesByExpenseId(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to delete expense shares", "error", err, "expenseId", id)
		// Продолжаем, так как хотим попробовать удалить сам расход
	}

	// Удаляем расход
	err = s.expenseRepo.DeleteExpense(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to delete expense", "error", err, "id", id)
		return errors.WithMessage(err, "failed to delete expense")
	}

	return nil
}

func (s Service) UpdateExpenseSharePaidStatus(ctx context.Context, shareId int, isPaid bool) error {
	err := s.expenseShareRepo.UpdateExpenseSharePaidStatus(ctx, shareId, isPaid)
	if err != nil {
		s.logger.Errorw("Failed to update expense share paid status", "error", err, "shareId", shareId)
		return errors.WithMessage(err, "failed to update expense share paid status")
	}

	return nil
}

func (s Service) GetEventBalanceReport(ctx context.Context, eventId int) (domain.BalanceReportResponse, error) {
	// Получаем общую сумму расходов для события
	totalAmount, err := s.expenseRepo.GetEventTotalExpenses(ctx, eventId)
	if err != nil {
		s.logger.Errorw("Failed to get event total expenses for report", "error", err, "eventId", eventId)
		return domain.BalanceReportResponse{}, errors.WithMessage(err, "failed to get event total expenses")
	}

	// Получаем балансы участников
	balances, err := s.expenseShareRepo.GetEventBalanceReport(ctx, eventId)
	if err != nil {
		s.logger.Errorw("Failed to get event balance report", "error", err, "eventId", eventId)
		return domain.BalanceReportResponse{}, errors.WithMessage(err, "failed to get event balance report")
	}

	// Преобразуем в доменную модель
	userBalances := make([]domain.UserBalance, len(balances))
	for i, balance := range balances {
		userBalances[i] = domain.UserBalance{
			UserID:   balance.UserID,
			Username: balance.Username,
			Balance:  balance.Balance,
		}
	}

	return domain.BalanceReportResponse{
		EventID:      eventId,
		TotalAmount:  totalAmount,
		UserBalances: userBalances,
	}, nil
}
