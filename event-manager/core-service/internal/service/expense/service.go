package expense

import (
	"context"
	"errors"
	"fmt"
	"github.com/PabloPerdolie/event-manager/core-service/internal/model"
	"github.com/PabloPerdolie/event-manager/core-service/internal/repository/expense"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"math"
)

// Service provides expense-related operations
type Service interface {
	Create(ctx context.Context, req model.ExpenseCreateRequest) (int, error)
	GetById(ctx context.Context, id int) (model.ExpenseResponse, error)
	Update(ctx context.Context, id int, req model.ExpenseUpdateRequest) error
	Delete(ctx context.Context, id int) error
	ListByEvent(ctx context.Context, eventId int, page, size int) (model.ExpensesResponse, error)
	ListByUser(ctx context.Context, userId int, page, size int) (model.ExpensesResponse, error)
	GetEventTotalExpenses(ctx context.Context, eventId int) (float64, error)
	CreateShares(ctx context.Context, expenseId int, participants []int, method model.ExpenseSplitMethod) error
}

type service struct {
	expenseRepo      expense.Repository
	expenseShareRepo expense.ShareRepository
	logger           *zap.SugaredLogger
}

// NewService creates a new expense service
func NewService(expenseRepo expense.Repository, expenseShareRepo expense.ShareRepository, logger *zap.SugaredLogger) Service {
	return &service{
		expenseRepo:      expenseRepo,
		expenseShareRepo: expenseShareRepo,
		logger:           logger,
	}
}

// Create creates a new expense
func (s *service) Create(ctx context.Context, req model.ExpenseCreateRequest) (int, error) {
	expense := model.Expense{
		EventId:     req.EventId,
		Description: req.Description,
		Amount:      req.Amount,
		Currency:    req.Currency,
		ExpenseDate: req.ExpenseDate,
		CreatedBy:   req.CreatedBy,
		SplitMethod: req.SplitMethod,
	}

	id, err := s.expenseRepo.Create(ctx, expense)
	if err != nil {
		s.logger.Errorw("Failed to create expense", "error", err, "eventId", req.EventId)
		return uuid.Nil, fmt.Errorf("failed to create expense: %w", err)
	}

	// Create expense shares if participants are provided
	if len(req.ParticipantIds) > 0 {
		err = s.CreateShares(ctx, id, req.ParticipantIds, req.SplitMethod)
		if err != nil {
			s.logger.Errorw("Failed to create expense shares", "error", err, "expenseId", id)
			// We don't return an error here, as the expense was already created successfully
		}
	}

	return id, nil
}

// CreateShares creates expense shares for participants based on the split method
func (s *service) CreateShares(ctx context.Context, expenseId int, participants []int, method model.ExpenseSplitMethod) error {
	// Get the expense to determine the amount to split
	expense, err := s.expenseRepo.GetById(ctx, expenseId)
	if err != nil {
		return errors.WithMessage(err, "")("failed to get expense for creating shares: %w", err)
	}

	// Delete any existing shares for this expense
	err = s.expenseShareRepo.DeleteByExpense(ctx, expenseId)
	if err != nil {
		return errors.WithMessage(err, "")("failed to delete existing expense shares: %w", err)
	}

	// Create shares based on the split method
	switch method {
	case model.ExpenseSplitMethodEqual:
		// Split equally among all participants
		if len(participants) == 0 {
			return errors.New("cannot split expense with no participants")
		}

		shareAmount := expense.Amount / float64(len(participants))
		// Round to 2 decimal places
		shareAmount = math.Round(shareAmount*100) / 100

		for _, participantId := range participants {
			share := model.ExpenseShare{
				ExpenseId:   expenseId,
				UserId:      participantId,
				ShareAmount: shareAmount,
				IsPaid:      false,
			}

			_, err := s.expenseShareRepo.Create(ctx, share)
			if err != nil {
				s.logger.Warnw("Failed to create expense share", "error", err, "expenseId", expenseId, "userId", participantId)
				// Continue even if one share fails
			}
		}

	case model.ExpenseSplitMethodPercentage:
		// This would typically be handled differently with custom percentages
		// For simplicity, we'll do equal percentages here
		if len(participants) == 0 {
			return errors.New("cannot split expense with no participants")
		}

		shareAmount := expense.Amount / float64(len(participants))
		// Round to 2 decimal places
		shareAmount = math.Round(shareAmount*100) / 100

		for _, participantId := range participants {
			share := model.ExpenseShare{
				ExpenseId:   expenseId,
				UserId:      participantId,
				ShareAmount: shareAmount,
				IsPaid:      false,
			}

			_, err := s.expenseShareRepo.Create(ctx, share)
			if err != nil {
				s.logger.Warnw("Failed to create expense share", "error", err, "expenseId", expenseId, "userId", participantId)
				// Continue even if one share fails
			}
		}

	case model.ExpenseSplitMethodCustom:
		// For custom splits, we would need additional data about how much each person pays
		// For now, we'll return an error
		return errors.New("custom split method requires additional data")

	default:
		return errors.WithMessage(err, "")("unsupported split method: %s", method)
	}

	return nil
}

// GetById retrieves an expense by Id
func (s *service) GetById(ctx context.Context, id int) (model.ExpenseResponse, error) {
	expense, err := s.expenseRepo.GetById(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to get expense by Id", "error", err, "id", id)
		return model.ExpenseResponse{}, fmt.Errorf("failed to get expense: %w", err)
	}

	// Get expense shares
	shares, err := s.expenseShareRepo.ListByExpense(ctx, id)
	if err != nil {
		s.logger.Warnw("Failed to get expense shares", "error", err, "expenseId", id)
		// Continue even if we can't get shares
	}

	// Extract participant Ids and create share responses
	participantIds := make([]int, len(shares))
	shareResponses := make([]model.ExpenseShareResponse, len(shares))
	for i, share := range shares {
		participantIds[i] = share.UserId
		shareResponses[i] = model.ExpenseShareResponse{
			Id:          share.Id,
			ExpenseId:   share.ExpenseId,
			UserId:      share.UserId,
			ShareAmount: share.ShareAmount,
			IsPaid:      share.IsPaid,
			PaidAt:      share.PaidAt,
		}
	}

	return model.ExpenseResponse{
		Id:             expense.Id,
		EventId:        expense.EventId,
		Description:    expense.Description,
		Amount:         expense.Amount,
		Currency:       expense.Currency,
		ExpenseDate:    expense.ExpenseDate,
		CreatedBy:      expense.CreatedBy,
		SplitMethod:    expense.SplitMethod,
		CreatedAt:      expense.CreatedAt,
		ParticipantIds: participantIds,
		Shares:         shareResponses,
	}, nil
}

// Update updates an expense
func (s *service) Update(ctx context.Context, id int, req model.ExpenseUpdateRequest) error {
	expense, err := s.expenseRepo.GetById(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to get expense for update", "error", err, "id", id)
		return errors.WithMessage(err, "")("failed to get expense: %w", err)
	}

	// Update fields if provided
	if req.Description != nil {
		expense.Description = *req.Description
	}

	if req.Amount != nil {
		expense.Amount = *req.Amount
	}

	if req.Currency != nil {
		expense.Currency = *req.Currency
	}

	if req.ExpenseDate != nil {
		expense.ExpenseDate = *req.ExpenseDate
	}

	if req.SplitMethod != nil {
		expense.SplitMethod = *req.SplitMethod
	}

	if err := s.expenseRepo.Update(ctx, expense); err != nil {
		s.logger.Errorw("Failed to update expense", "error", err, "id", id)
		return errors.WithMessage(err, "")("failed to update expense: %w", err)
	}

	// Update shares if participants and/or split method are provided
	if req.ParticipantIds != nil && len(*req.ParticipantIds) > 0 {
		splitMethod := expense.SplitMethod
		if req.SplitMethod != nil {
			splitMethod = *req.SplitMethod
		}

		err = s.CreateShares(ctx, id, *req.ParticipantIds, splitMethod)
		if err != nil {
			s.logger.Errorw("Failed to update expense shares", "error", err, "expenseId", id)
			return errors.WithMessage(err, "")("failed to update expense shares: %w", err)
		}
	}

	return nil
}

// Delete deletes an expense
func (s *service) Delete(ctx context.Context, id int) error {
	// First check if the expense exists
	_, err := s.expenseRepo.GetById(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to get expense for deletion", "error", err, "id", id)
		return errors.WithMessage(err, "")("failed to get expense: %w", err)
	}

	// Delete associated shares
	err = s.expenseShareRepo.DeleteByExpense(ctx, id)
	if err != nil {
		s.logger.Warnw("Failed to delete expense shares", "error", err, "expenseId", id)
		// Continue even if share deletion fails
	}

	// Delete the expense
	if err := s.expenseRepo.Delete(ctx, id); err != nil {
		s.logger.Errorw("Failed to delete expense", "error", err, "id", id)
		return errors.WithMessage(err, "")("failed to delete expense: %w", err)
	}

	return nil
}

// ListByEvent retrieves a list of expenses for a specific event with pagination
func (s *service) ListByEvent(ctx context.Context, eventId int, page, size int) (model.ExpensesResponse, error) {
	// Set default pagination values if not provided
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	offset := (page - 1) * size
	expenses, total, err := s.expenseRepo.ListByEvent(ctx, eventId, size, offset)
	if err != nil {
		s.logger.Errorw("Failed to list event expenses", "error", err, "eventId", eventId, "page", page, "size", size)
		return model.ExpensesResponse{}, fmt.Errorf("failed to list event expenses: %w", err)
	}

	// Convert to response objects
	expenseResponses := make([]model.ExpenseResponse, len(expenses))
	for i, expense := range expenses {
		// Get expense shares
		shares, err := s.expenseShareRepo.ListByExpense(ctx, expense.Id)
		if err != nil {
			s.logger.Warnw("Failed to get expense shares", "error", err, "expenseId", expense.Id)
			// Continue even if we can't get shares
			expenseResponses[i] = model.ExpenseResponse{
				Id:          expense.Id,
				EventId:     expense.EventId,
				Description: expense.Description,
				Amount:      expense.Amount,
				Currency:    expense.Currency,
				ExpenseDate: expense.ExpenseDate,
				CreatedBy:   expense.CreatedBy,
				SplitMethod: expense.SplitMethod,
				CreatedAt:   expense.CreatedAt,
			}
			continue
		}

		// Extract participant Ids and create share responses
		participantIds := make([]int, len(shares))
		shareResponses := make([]model.ExpenseShareResponse, len(shares))
		for j, share := range shares {
			participantIds[j] = share.UserId
			shareResponses[j] = model.ExpenseShareResponse{
				Id:          share.Id,
				ExpenseId:   share.ExpenseId,
				UserId:      share.UserId,
				ShareAmount: share.ShareAmount,
				IsPaid:      share.IsPaid,
				PaidAt:      share.PaidAt,
			}
		}

		expenseResponses[i] = model.ExpenseResponse{
			Id:             expense.Id,
			EventId:        expense.EventId,
			Description:    expense.Description,
			Amount:         expense.Amount,
			Currency:       expense.Currency,
			ExpenseDate:    expense.ExpenseDate,
			CreatedBy:      expense.CreatedBy,
			SplitMethod:    expense.SplitMethod,
			CreatedAt:      expense.CreatedAt,
			ParticipantIds: participantIds,
			Shares:         shareResponses,
		}
	}

	return model.ExpensesResponse{
		Expenses: expenseResponses,
		Total:    total,
	}, nil
}

// ListByUser retrieves a list of expenses created by a specific user with pagination
func (s *service) ListByUser(ctx context.Context, userId int, page, size int) (model.ExpensesResponse, error) {
	// Set default pagination values if not provided
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	offset := (page - 1) * size
	expenses, total, err := s.expenseRepo.ListByUser(ctx, userId, size, offset)
	if err != nil {
		s.logger.Errorw("Failed to list user expenses", "error", err, "userId", userId, "page", page, "size", size)
		return model.ExpensesResponse{}, fmt.Errorf("failed to list user expenses: %w", err)
	}

	// Convert to response objects
	expenseResponses := make([]model.ExpenseResponse, len(expenses))
	for i, expense := range expenses {
		// Get expense shares
		shares, err := s.expenseShareRepo.ListByExpense(ctx, expense.Id)
		if err != nil {
			s.logger.Warnw("Failed to get expense shares", "error", err, "expenseId", expense.Id)
			// Continue even if we can't get shares
			expenseResponses[i] = model.ExpenseResponse{
				Id:          expense.Id,
				EventId:     expense.EventId,
				Description: expense.Description,
				Amount:      expense.Amount,
				Currency:    expense.Currency,
				ExpenseDate: expense.ExpenseDate,
				CreatedBy:   expense.CreatedBy,
				SplitMethod: expense.SplitMethod,
				CreatedAt:   expense.CreatedAt,
			}
			continue
		}

		// Extract participant Ids and create share responses
		participantIds := make([]int, len(shares))
		shareResponses := make([]model.ExpenseShareResponse, len(shares))
		for j, share := range shares {
			participantIds[j] = share.UserId
			shareResponses[j] = model.ExpenseShareResponse{
				Id:          share.Id,
				ExpenseId:   share.ExpenseId,
				UserId:      share.UserId,
				ShareAmount: share.ShareAmount,
				IsPaid:      share.IsPaid,
				PaidAt:      share.PaidAt,
			}
		}

		expenseResponses[i] = model.ExpenseResponse{
			Id:             expense.Id,
			EventId:        expense.EventId,
			Description:    expense.Description,
			Amount:         expense.Amount,
			Currency:       expense.Currency,
			ExpenseDate:    expense.ExpenseDate,
			CreatedBy:      expense.CreatedBy,
			SplitMethod:    expense.SplitMethod,
			CreatedAt:      expense.CreatedAt,
			ParticipantIds: participantIds,
			Shares:         shareResponses,
		}
	}

	return model.ExpensesResponse{
		Expenses: expenseResponses,
		Total:    total,
	}, nil
}

// GetEventTotalExpenses calculates the total amount of expenses for an event
func (s *service) GetEventTotalExpenses(ctx context.Context, eventId int) (float64, error) {
	total, err := s.expenseRepo.GetEventTotalExpenses(ctx, eventId)
	if err != nil {
		s.logger.Errorw("Failed to calculate event total expenses", "error", err, "eventId", eventId)
		return 0, fmt.Errorf("failed to calculate event total expenses: %w", err)
	}

	return total, nil
}
