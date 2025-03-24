package task

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/PabloPerdolie/event-manager/core-service/internal/model"
	"github.com/PabloPerdolie/event-manager/core-service/internal/repository/task"
	"github.com/PabloPerdolie/event-manager/core-service/internal/repository/user"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// AssignmentService provides task assignment-related operations
type AssignmentService interface {
	Create(ctx context.Context, req model.TaskAssignmentCreateRequest) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.TaskAssignmentResponse, error)
	MarkComplete(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByTask(ctx context.Context, taskID uuid.UUID, page, size int) (model.TaskAssignmentsResponse, error)
	ListByUser(ctx context.Context, userID uuid.UUID, page, size int) (model.TaskAssignmentsResponse, error)
	GetUserSummary(ctx context.Context, userID uuid.UUID) (model.UserTaskSummary, error)
}

type assignmentService struct {
	repo     task.AssignmentRepository
	taskRepo task.Repository
	userRepo user.Repository
	logger   *zap.SugaredLogger
}

// NewAssignmentService creates a new task assignment service
func NewAssignmentService(repo task.AssignmentRepository, taskRepo task.Repository, userRepo user.Repository, logger *zap.SugaredLogger) AssignmentService {
	return &assignmentService{
		repo:     repo,
		taskRepo: taskRepo,
		userRepo: userRepo,
		logger:   logger,
	}
}

// Create creates a new task assignment
func (s *assignmentService) Create(ctx context.Context, req model.TaskAssignmentCreateRequest) (uuid.UUID, error) {
	// Verify that user exists
	_, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		s.logger.Errorw("Failed to get user for assignment creation", "error", err, "userId", req.UserID)
		return uuid.Nil, fmt.Errorf("invalid user: %w", err)
	}

	// Verify that task exists
	task, err := s.taskRepo.GetByID(ctx, req.TaskID)
	if err != nil {
		s.logger.Errorw("Failed to get task for assignment creation", "error", err, "taskId", req.TaskID)
		return uuid.Nil, fmt.Errorf("invalid task: %w", err)
	}

	// Check if assignment already exists
	_, err = s.repo.GetByTaskAndUser(ctx, req.TaskID, req.UserID)
	if err == nil {
		return uuid.Nil, errors.New("user is already assigned to this task")
	}

	assignment := model.TaskAssignment{
		TaskID:     req.TaskID,
		UserID:     req.UserID,
		AssignedAt: time.Now(),
	}

	// If the task is already completed, set the completed_at field
	if task.Status == model.TaskStatusCompleted {
		now := time.Now()
		assignment.CompletedAt = &now
	}

	id, err := s.repo.Create(ctx, assignment)
	if err != nil {
		s.logger.Errorw("Failed to create assignment", "error", err, "taskId", req.TaskID, "userId", req.UserID)
		return uuid.Nil, fmt.Errorf("failed to create assignment: %w", err)
	}

	return id, nil
}

// GetByID retrieves a task assignment by ID
func (s *assignmentService) GetByID(ctx context.Context, id uuid.UUID) (model.TaskAssignmentResponse, error) {
	assignment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to get assignment by ID", "error", err, "id", id)
		return model.TaskAssignmentResponse{}, fmt.Errorf("failed to get assignment: %w", err)
	}

	// Get task info
	task, err := s.taskRepo.GetByID(ctx, assignment.TaskID)
	if err != nil {
		s.logger.Warnw("Failed to get assignment task details", "error", err, "taskId", assignment.TaskID)
		// Continue even if we can't get task details
	}

	// Get user info
	user, err := s.userRepo.GetByID(ctx, assignment.UserID)
	if err != nil {
		s.logger.Warnw("Failed to get assignment user details", "error", err, "userId", assignment.UserID)
		// Continue even if we can't get user details
	}

	return model.TaskAssignmentResponse{
		ID:          assignment.ID,
		TaskID:      assignment.TaskID,
		TaskTitle:   task.Title,
		UserID:      assignment.UserID,
		Username:    user.Username,
		FullName:    fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		AssignedAt:  assignment.AssignedAt,
		CompletedAt: assignment.CompletedAt,
		IsCompleted: assignment.CompletedAt != nil,
	}, nil
}

// MarkComplete marks a task assignment as complete
func (s *assignmentService) MarkComplete(ctx context.Context, id uuid.UUID) error {
	assignment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to get assignment for completion", "error", err, "id", id)
		return fmt.Errorf("failed to get assignment: %w", err)
	}

	// If already completed, do nothing
	if assignment.CompletedAt != nil {
		return nil
	}

	now := time.Now()
	assignment.CompletedAt = &now

	if err := s.repo.Update(ctx, assignment); err != nil {
		s.logger.Errorw("Failed to mark assignment as complete", "error", err, "id", id)
		return fmt.Errorf("failed to update assignment: %w", err)
	}

	// Check if all assignments for this task are complete
	// If so, update the task status to completed
	assignments, _, err := s.repo.ListByTask(ctx, assignment.TaskID, 100, 0)
	if err != nil {
		s.logger.Warnw("Failed to list task assignments", "error", err, "taskId", assignment.TaskID)
		return nil // Don't fail the whole operation
	}

	allComplete := true
	for _, a := range assignments {
		if a.CompletedAt == nil {
			allComplete = false
			break
		}
	}

	if allComplete && len(assignments) > 0 {
		task, err := s.taskRepo.GetByID(ctx, assignment.TaskID)
		if err != nil {
			s.logger.Warnw("Failed to get task for status update", "error", err, "taskId", assignment.TaskID)
			return nil // Don't fail the whole operation
		}

		task.Status = model.TaskStatusCompleted
		if err := s.taskRepo.Update(ctx, task); err != nil {
			s.logger.Warnw("Failed to update task status", "error", err, "taskId", assignment.TaskID)
			return nil // Don't fail the whole operation
		}
	}

	return nil
}

// Delete deletes a task assignment
func (s *assignmentService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Errorw("Failed to delete assignment", "error", err, "id", id)
		return fmt.Errorf("failed to delete assignment: %w", err)
	}

	return nil
}

// ListByTask retrieves a list of assignments for a specific task with pagination
func (s *assignmentService) ListByTask(ctx context.Context, taskID uuid.UUID, page, size int) (model.TaskAssignmentsResponse, error) {
	// Set default pagination values if not provided
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	offset := (page - 1) * size
	assignments, total, err := s.repo.ListByTask(ctx, taskID, size, offset)
	if err != nil {
		s.logger.Errorw("Failed to list task assignments", "error", err, "taskId", taskID, "page", page, "size", size)
		return model.TaskAssignmentsResponse{}, fmt.Errorf("failed to list assignments: %w", err)
	}

	// Get task info
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		s.logger.Warnw("Failed to get task details for assignments", "error", err, "taskId", taskID)
		// Continue with minimal task info
	}

	// Convert to response objects
	assignmentResponses := make([]model.TaskAssignmentResponse, len(assignments))
	for i, assignment := range assignments {
		// Get user info
		user, err := s.userRepo.GetByID(ctx, assignment.UserID)
		if err != nil {
			s.logger.Warnw("Failed to get assignment user details", "error", err, "userId", assignment.UserID)
			// Continue with minimal user info
			assignmentResponses[i] = model.TaskAssignmentResponse{
				ID:          assignment.ID,
				TaskID:      assignment.TaskID,
				TaskTitle:   task.Title,
				UserID:      assignment.UserID,
				AssignedAt:  assignment.AssignedAt,
				CompletedAt: assignment.CompletedAt,
				IsCompleted: assignment.CompletedAt != nil,
			}
			continue
		}

		assignmentResponses[i] = model.TaskAssignmentResponse{
			ID:          assignment.ID,
			TaskID:      assignment.TaskID,
			TaskTitle:   task.Title,
			UserID:      assignment.UserID,
			Username:    user.Username,
			FullName:    fmt.Sprintf("%s %s", user.FirstName, user.LastName),
			AssignedAt:  assignment.AssignedAt,
			CompletedAt: assignment.CompletedAt,
			IsCompleted: assignment.CompletedAt != nil,
		}
	}

	return model.TaskAssignmentsResponse{
		Assignments: assignmentResponses,
		Total:       total,
	}, nil
}

// ListByUser retrieves a list of task assignments for a specific user with pagination
func (s *assignmentService) ListByUser(ctx context.Context, userID uuid.UUID, page, size int) (model.TaskAssignmentsResponse, error) {
	// Set default pagination values if not provided
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	offset := (page - 1) * size
	assignments, total, err := s.repo.ListByUser(ctx, userID, size, offset)
	if err != nil {
		s.logger.Errorw("Failed to list user assignments", "error", err, "userId", userID, "page", page, "size", size)
		return model.TaskAssignmentsResponse{}, fmt.Errorf("failed to list assignments: %w", err)
	}

	// Get user info
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.Warnw("Failed to get user details for assignments", "error", err, "userId", userID)
		// Continue with minimal info
	}

	// Convert to response objects
	assignmentResponses := make([]model.TaskAssignmentResponse, len(assignments))
	for i, assignment := range assignments {
		// Get task info
		task, err := s.taskRepo.GetByID(ctx, assignment.TaskID)
		if err != nil {
			s.logger.Warnw("Failed to get task details for assignment", "error", err, "taskId", assignment.TaskID)
			// Continue with minimal task info
			assignmentResponses[i] = model.TaskAssignmentResponse{
				ID:          assignment.ID,
				TaskID:      assignment.TaskID,
				UserID:      assignment.UserID,
				Username:    user.Username,
				FullName:    fmt.Sprintf("%s %s", user.FirstName, user.LastName),
				AssignedAt:  assignment.AssignedAt,
				CompletedAt: assignment.CompletedAt,
				IsCompleted: assignment.CompletedAt != nil,
			}
			continue
		}

		assignmentResponses[i] = model.TaskAssignmentResponse{
			ID:          assignment.ID,
			TaskID:      assignment.TaskID,
			TaskTitle:   task.Title,
			UserID:      assignment.UserID,
			Username:    user.Username,
			FullName:    fmt.Sprintf("%s %s", user.FirstName, user.LastName),
			AssignedAt:  assignment.AssignedAt,
			CompletedAt: assignment.CompletedAt,
			IsCompleted: assignment.CompletedAt != nil,
		}
	}

	return model.TaskAssignmentsResponse{
		Assignments: assignmentResponses,
		Total:       total,
	}, nil
}

// GetUserSummary retrieves a summary of a user's tasks
func (s *assignmentService) GetUserSummary(ctx context.Context, userID uuid.UUID) (model.UserTaskSummary, error) {
	var summary model.UserTaskSummary
	summary.UserID = userID

	// Get all assignments for the user
	assignments, _, err := s.repo.ListByUser(ctx, userID, 1000, 0)
	if err != nil {
		s.logger.Errorw("Failed to list user assignments for summary", "error", err, "userId", userID)
		return summary, fmt.Errorf("failed to get user assignments: %w", err)
	}

	// Count assignments by completion status
	totalAssignments := len(assignments)
	completedAssignments := 0

	for _, assignment := range assignments {
		if assignment.CompletedAt != nil {
			completedAssignments++
		}
	}

	pendingAssignments := totalAssignments - completedAssignments

	// Calculate completion percentage
	var completionPercentage float64
	if totalAssignments > 0 {
		completionPercentage = float64(completedAssignments) / float64(totalAssignments) * 100
	}

	summary.TotalAssignments = totalAssignments
	summary.CompletedAssignments = completedAssignments
	summary.PendingAssignments = pendingAssignments
	summary.CompletionPercentage = completionPercentage

	return summary, nil
}
