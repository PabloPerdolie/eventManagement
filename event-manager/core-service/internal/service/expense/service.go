package expense

import (
	"context"
	"fmt"
	"github.com/PabloPerdolie/event-manager/core-service/internal/domain"
	"github.com/PabloPerdolie/event-manager/core-service/internal/model"
	"github.com/PabloPerdolie/event-manager/core-service/internal/repository"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"math"
)

// Service предоставляет операции для работы с расходами
type Service interface {
	CreateExpense(ctx context.Context, req domain.ExpenseCreateRequest) (int, error)
	GetExpenseById(ctx context.Context, id int) (domain.ExpenseResponse, error)
	UpdateExpense(ctx context.Context, id int, req domain.ExpenseUpdateRequest) error
	DeleteExpense(ctx context.Context, id int) error
	ListExpensesByEvent(ctx context.Context, eventId int, page, size int) (domain.ExpensesResponse, error)
	GetEventTotalExpenses(ctx context.Context, eventId int) (float64, error)
	CreateExpenseShares(ctx context.Context, expenseId int, userIds []int, splitMethod string) error
	UpdateExpenseSharePaidStatus(ctx context.Context, shareId int, isPaid bool) error
	GetEventBalanceReport(ctx context.Context, eventId int) (domain.BalanceReportResponse, error)
}

type service struct {
	expenseRepo      repository.Expense
	expenseShareRepo repository.ExpenseShare
	logger           *zap.SugaredLogger
}

// NewService создает новый сервис для работы с расходами
func NewService(expenseRepo repository.Expense, expenseShareRepo repository.ExpenseShare, logger *zap.SugaredLogger) Service {
	return &service{
		expenseRepo:      expenseRepo,
		expenseShareRepo: expenseShareRepo,
		logger:           logger,
	}
}

// CreateExpense создает новый расход
func (s *service) CreateExpense(ctx context.Context, req domain.ExpenseCreateRequest) (int, error) {
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
		return 0, fmt.Errorf("failed to create expense: %w", err)
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

// CreateExpenseShares создает доли расходов для участников в зависимости от метода разделения
func (s *service) CreateExpenseShares(ctx context.Context, expenseId int, userIds []int, splitMethod string) error {
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
	case model.SplitMethodEqual:
		// Делим поровну между всеми участниками
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

	case model.SplitMethodPercent:
		// Для простоты реализуем как равное разделение
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
		return fmt.Errorf("unsupported split method: %s", splitMethod)
	}

	return nil
}

// GetExpenseById получает расход по ID
func (s *service) GetExpenseById(ctx context.Context, id int) (domain.ExpenseResponse, error) {
	expense, err := s.expenseRepo.GetExpenseById(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to get expense by ID", "error", err, "id", id)
		return domain.ExpenseResponse{}, errors.WithMessage(err, "failed to get expense")
	}

	// Получаем доли расходов
	shares, err := s.expenseShareRepo.ListExpenseSharesByExpenseId(ctx, id)
	if err != nil {
		s.logger.Warnw("Failed to get expense shares", "error", err, "expenseId", id)
		// Продолжаем, даже если не можем получить доли
		shares = []model.ExpenseShare{}
	}

	return domain.ExpenseResponse{
		ExpenseID:   expense.ExpenseID,
		EventID:     expense.EventID,
		Description: expense.Description,
		Amount:      expense.Amount,
		Currency:    expense.Currency,
		CreatedBy:   expense.CreatedBy,
		SplitMethod: expense.SplitMethod,
		CreatedAt:   expense.CreatedAt,
		Shares:      shares,
	}, nil
}

// UpdateExpense обновляет расход
func (s *service) UpdateExpense(ctx context.Context, id int, req domain.ExpenseUpdateRequest) error {
	expense, err := s.expenseRepo.GetExpenseById(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to get expense for update", "error", err, "id", id)
		return errors.WithMessage(err, "failed to get expense")
	}

	// Обновляем поля, если они предоставлены
	if req.Description != nil {
		expense.Description = *req.Description
	}

	if req.Amount != nil {
		expense.Amount = *req.Amount
	}

	if req.Currency != nil {
		expense.Currency = *req.Currency
	}

	if req.SplitMethod != nil {
		expense.SplitMethod = *req.SplitMethod
	}

	if err := s.expenseRepo.UpdateExpense(ctx, *expense); err != nil {
		s.logger.Errorw("Failed to update expense", "error", err, "id", id)
		return errors.WithMessage(err, "failed to update expense")
	}

	// Обновляем доли, если предоставлены участники и/или метод разделения
	if req.UserIDs != nil && len(*req.UserIDs) > 0 {
		splitMethod := expense.SplitMethod
		if req.SplitMethod != nil {
			splitMethod = *req.SplitMethod
		}

		err = s.CreateExpenseShares(ctx, id, *req.UserIDs, splitMethod)
		if err != nil {
			s.logger.Errorw("Failed to update expense shares", "error", err, "expenseId", id)
			return errors.WithMessage(err, "failed to update expense shares")
		}
	}

	return nil
}

// DeleteExpense удаляет расход
func (s *service) DeleteExpense(ctx context.Context, id int) error {
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

// ListExpensesByEvent получает список расходов для события с пагинацией
func (s *service) ListExpensesByEvent(ctx context.Context, eventId int, page, size int) (domain.ExpensesResponse, error) {
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

// GetEventTotalExpenses получает общую сумму расходов для события
func (s *service) GetEventTotalExpenses(ctx context.Context, eventId int) (float64, error) {
	total, err := s.expenseRepo.GetEventTotalExpenses(ctx, eventId)
	if err != nil {
		s.logger.Errorw("Failed to get event total expenses", "error", err, "eventId", eventId)
		return 0, errors.WithMessage(err, "failed to get event total expenses")
	}

	return total, nil
}

// UpdateExpenseSharePaidStatus обновляет статус оплаты доли расхода
func (s *service) UpdateExpenseSharePaidStatus(ctx context.Context, shareId int, isPaid bool) error {
	err := s.expenseShareRepo.UpdateExpenseSharePaidStatus(ctx, shareId, isPaid)
	if err != nil {
		s.logger.Errorw("Failed to update expense share paid status", "error", err, "shareId", shareId)
		return errors.WithMessage(err, "failed to update expense share paid status")
	}

	return nil
}

// GetEventBalanceReport получает отчет о балансе участников события
func (s *service) GetEventBalanceReport(ctx context.Context, eventId int) (domain.BalanceReportResponse, error) {
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
