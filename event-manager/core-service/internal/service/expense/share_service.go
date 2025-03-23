package expense

import (
	"context"
	"fmt"
	"time"

	"github.com/event-management/core-service/internal/model"
	"github.com/event-management/core-service/internal/repository/expense"
	"github.com/event-management/core-service/internal/repository/user"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ShareService provides expense share-related operations
type ShareService interface {
	GetByID(ctx context.Context, id uuid.UUID) (model.ExpenseShareResponse, error)
	Update(ctx context.Context, id uuid.UUID, req model.ExpenseShareUpdateRequest) error
	MarkPaid(ctx context.Context, id uuid.UUID) error
	ListByExpense(ctx context.Context, expenseID uuid.UUID) ([]model.ExpenseShareResponse, error)
	ListByUser(ctx context.Context, userID uuid.UUID, page, size int) (model.ExpenseSharesResponse, error)
	GetUserDebts(ctx context.Context, userID uuid.UUID) (model.UserDebtSummary, error)
	GetEventUserDebts(ctx context.Context, eventID, userID uuid.UUID) (model.EventDebtSummary, error)
}

type shareService struct {
	shareRepo   expense.ShareRepository
	expenseRepo expense.Repository
	userRepo    user.Repository
	logger      *zap.SugaredLogger
}

// NewShareService creates a new expense share service
func NewShareService(shareRepo expense.ShareRepository, expenseRepo expense.Repository, userRepo user.Repository, logger *zap.SugaredLogger) ShareService {
	return &shareService{
		shareRepo:   shareRepo,
		expenseRepo: expenseRepo,
		userRepo:    userRepo,
		logger:      logger,
	}
}

// GetByID retrieves an expense share by ID
func (s *shareService) GetByID(ctx context.Context, id uuid.UUID) (model.ExpenseShareResponse, error) {
	share, err := s.shareRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to get expense share by ID", "error", err, "id", id)
		return model.ExpenseShareResponse{}, fmt.Errorf("failed to get expense share: %w", err)
	}

	// Get expense details
	expense, err := s.expenseRepo.GetByID(ctx, share.ExpenseID)
	if err != nil {
		s.logger.Warnw("Failed to get expense details for share", "error", err, "expenseId", share.ExpenseID)
		// Continue even if we can't get expense details
	}

	// Get user details
	user, err := s.userRepo.GetByID(ctx, share.UserID)
	if err != nil {
		s.logger.Warnw("Failed to get user details for share", "error", err, "userId", share.UserID)
		// Continue even if we can't get user details
	}

	return model.ExpenseShareResponse{
		ID:          share.ID,
		ExpenseID:   share.ExpenseID,
		ExpenseDesc: expense.Description,
		UserID:      share.UserID,
		Username:    user.Username,
		FullName:    fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		ShareAmount: share.ShareAmount,
		Currency:    expense.Currency,
		IsPaid:      share.IsPaid,
		PaidAt:      share.PaidAt,
	}, nil
}

// Update updates an expense share
func (s *shareService) Update(ctx context.Context, id uuid.UUID, req model.ExpenseShareUpdateRequest) error {
	share, err := s.shareRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to get expense share for update", "error", err, "id", id)
		return fmt.Errorf("failed to get expense share: %w", err)
	}

	// Update fields if provided
	if req.ShareAmount != nil {
		share.ShareAmount = *req.ShareAmount
	}

	if req.IsPaid != nil {
		share.IsPaid = *req.IsPaid
		if *req.IsPaid && share.PaidAt == nil {
			now := time.Now()
			share.PaidAt = &now
		} else if !*req.IsPaid {
			share.PaidAt = nil
		}
	}

	if err := s.shareRepo.Update(ctx, share); err != nil {
		s.logger.Errorw("Failed to update expense share", "error", err, "id", id)
		return fmt.Errorf("failed to update expense share: %w", err)
	}

	return nil
}

// MarkPaid marks an expense share as paid
func (s *shareService) MarkPaid(ctx context.Context, id uuid.UUID) error {
	share, err := s.shareRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to get expense share for marking paid", "error", err, "id", id)
		return fmt.Errorf("failed to get expense share: %w", err)
	}

	// If already paid, do nothing
	if share.IsPaid {
		return nil
	}

	now := time.Now()
	share.IsPaid = true
	share.PaidAt = &now

	if err := s.shareRepo.Update(ctx, share); err != nil {
		s.logger.Errorw("Failed to mark expense share as paid", "error", err, "id", id)
		return fmt.Errorf("failed to update expense share: %w", err)
	}

	return nil
}

// ListByExpense retrieves all expense shares for a specific expense
func (s *shareService) ListByExpense(ctx context.Context, expenseID uuid.UUID) ([]model.ExpenseShareResponse, error) {
	shares, err := s.shareRepo.ListByExpense(ctx, expenseID)
	if err != nil {
		s.logger.Errorw("Failed to list expense shares", "error", err, "expenseId", expenseID)
		return nil, fmt.Errorf("failed to list expense shares: %w", err)
	}

	// Get expense details
	expense, err := s.expenseRepo.GetByID(ctx, expenseID)
	if err != nil {
		s.logger.Warnw("Failed to get expense details for shares", "error", err, "expenseId", expenseID)
		// Continue even if we can't get expense details
	}

	// Convert to response objects
	shareResponses := make([]model.ExpenseShareResponse, len(shares))
	for i, share := range shares {
		// Get user details
		user, err := s.userRepo.GetByID(ctx, share.UserID)
		if err != nil {
			s.logger.Warnw("Failed to get user details for share", "error", err, "userId", share.UserID)
			// Continue with minimal user info
			shareResponses[i] = model.ExpenseShareResponse{
				ID:          share.ID,
				ExpenseID:   share.ExpenseID,
				ExpenseDesc: expense.Description,
				UserID:      share.UserID,
				ShareAmount: share.ShareAmount,
				Currency:    expense.Currency,
				IsPaid:      share.IsPaid,
				PaidAt:      share.PaidAt,
			}
			continue
		}

		shareResponses[i] = model.ExpenseShareResponse{
			ID:          share.ID,
			ExpenseID:   share.ExpenseID,
			ExpenseDesc: expense.Description,
			UserID:      share.UserID,
			Username:    user.Username,
			FullName:    fmt.Sprintf("%s %s", user.FirstName, user.LastName),
			ShareAmount: share.ShareAmount,
			Currency:    expense.Currency,
			IsPaid:      share.IsPaid,
			PaidAt:      share.PaidAt,
		}
	}

	return shareResponses, nil
}

// ListByUser retrieves a list of expense shares for a specific user with pagination
func (s *shareService) ListByUser(ctx context.Context, userID uuid.UUID, page, size int) (model.ExpenseSharesResponse, error) {
	// Set default pagination values if not provided
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	offset := (page - 1) * size
	shares, total, err := s.shareRepo.ListByUser(ctx, userID, size, offset)
	if err != nil {
		s.logger.Errorw("Failed to list user expense shares", "error", err, "userId", userID, "page", page, "size", size)
		return model.ExpenseSharesResponse{}, fmt.Errorf("failed to list user expense shares: %w", err)
	}

	// Get user details
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.Warnw("Failed to get user details for shares", "error", err, "userId", userID)
		// Continue with minimal info
	}

	// Convert to response objects
	shareResponses := make([]model.ExpenseShareResponse, len(shares))
	for i, share := range shares {
		// Get expense details
		expense, err := s.expenseRepo.GetByID(ctx, share.ExpenseID)
		if err != nil {
			s.logger.Warnw("Failed to get expense details for share", "error", err, "expenseId", share.ExpenseID)
			// Continue with minimal expense info
			shareResponses[i] = model.ExpenseShareResponse{
				ID:          share.ID,
				ExpenseID:   share.ExpenseID,
				UserID:      share.UserID,
				Username:    user.Username,
				FullName:    fmt.Sprintf("%s %s", user.FirstName, user.LastName),
				ShareAmount: share.ShareAmount,
				IsPaid:      share.IsPaid,
				PaidAt:      share.PaidAt,
			}
			continue
		}

		shareResponses[i] = model.ExpenseShareResponse{
			ID:          share.ID,
			ExpenseID:   share.ExpenseID,
			ExpenseDesc: expense.Description,
			UserID:      share.UserID,
			Username:    user.Username,
			FullName:    fmt.Sprintf("%s %s", user.FirstName, user.LastName),
			ShareAmount: share.ShareAmount,
			Currency:    expense.Currency,
			IsPaid:      share.IsPaid,
			PaidAt:      share.PaidAt,
		}
	}

	return model.ExpenseSharesResponse{
		Shares: shareResponses,
		Total:  total,
	}, nil
}

// GetUserDebts retrieves a summary of a user's debts
func (s *shareService) GetUserDebts(ctx context.Context, userID uuid.UUID) (model.UserDebtSummary, error) {
	summary, err := s.shareRepo.GetUserDebts(ctx, userID)
	if err != nil {
		s.logger.Errorw("Failed to get user debt summary", "error", err, "userId", userID)
		return model.UserDebtSummary{}, fmt.Errorf("failed to get user debt summary: %w", err)
	}

	return summary, nil
}

// GetEventUserDebts retrieves a summary of a user's debts for a specific event
func (s *shareService) GetEventUserDebts(ctx context.Context, eventID, userID uuid.UUID) (model.EventDebtSummary, error) {
	summary, err := s.shareRepo.GetEventUserDebts(ctx, eventID, userID)
	if err != nil {
		s.logger.Errorw("Failed to get event user debt summary", "error", err, "eventId", eventID, "userId", userID)
		return model.EventDebtSummary{}, fmt.Errorf("failed to get event user debt summary: %w", err)
	}

	return summary, nil
}
