package task

import (
	"context"
	"fmt"
	"time"

	"github.com/PabloPerdolie/event-manager/core-service/internal/model"
	"github.com/PabloPerdolie/event-manager/core-service/internal/repository/task"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Service provides task-related operations
type Service interface {
	Create(ctx context.Context, req model.TaskCreateRequest) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.TaskResponse, error)
	Update(ctx context.Context, id uuid.UUID, req model.TaskUpdateRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByEvent(ctx context.Context, eventID uuid.UUID, page, size int) (model.TasksResponse, error)
	ListByUser(ctx context.Context, userID uuid.UUID, page, size int) (model.TasksResponse, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status model.TaskStatus) error
	GetEventSummary(ctx context.Context, eventID uuid.UUID) (model.TaskEventSummary, error)
}

type service struct {
	taskRepo       task.Repository
	assignmentRepo task.AssignmentRepository
	logger         *zap.SugaredLogger
}

// NewService creates a new task service
func NewService(taskRepo task.Repository, assignmentRepo task.AssignmentRepository, logger *zap.SugaredLogger) Service {
	return &service{
		taskRepo:       taskRepo,
		assignmentRepo: assignmentRepo,
		logger:         logger,
	}
}

// Create creates a new task
func (s *service) Create(ctx context.Context, req model.TaskCreateRequest) (uuid.UUID, error) {
	task := model.Task{
		EventID:     req.EventID,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		Priority:    req.Priority,
		Status:      model.TaskStatusPending,
		CreatedBy:   req.CreatedBy,
	}

	id, err := s.taskRepo.Create(ctx, task)
	if err != nil {
		s.logger.Errorw("Failed to create task", "error", err, "eventId", req.EventID)
		return uuid.Nil, fmt.Errorf("failed to create task: %w", err)
	}

	// Create task assignments if assignees are provided
	if len(req.AssigneeIDs) > 0 {
		for _, assigneeID := range req.AssigneeIDs {
			assignment := model.TaskAssignment{
				TaskID:     id,
				UserID:     assigneeID,
				AssignedAt: time.Now(),
			}

			_, err := s.assignmentRepo.Create(ctx, assignment)
			if err != nil {
				s.logger.Warnw("Failed to create task assignment", "error", err, "taskId", id, "userId", assigneeID)
				// Continue even if one assignment fails
			}
		}
	}

	return id, nil
}

// GetByID retrieves a task by ID
func (s *service) GetByID(ctx context.Context, id uuid.UUID) (model.TaskResponse, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to get task by ID", "error", err, "id", id)
		return model.TaskResponse{}, fmt.Errorf("failed to get task: %w", err)
	}

	// Get task assignments
	assignments, _, err := s.assignmentRepo.ListByTask(ctx, id, 100, 0)
	if err != nil {
		s.logger.Warnw("Failed to get task assignments", "error", err, "taskId", id)
		// Continue even if we can't get assignments
	}

	// Extract assignee IDs
	assigneeIDs := make([]uuid.UUID, len(assignments))
	for i, assignment := range assignments {
		assigneeIDs[i] = assignment.UserID
	}

	return model.TaskResponse{
		ID:          task.ID,
		EventID:     task.EventID,
		Title:       task.Title,
		Description: task.Description,
		DueDate:     task.DueDate,
		Priority:    task.Priority,
		Status:      task.Status,
		CreatedBy:   task.CreatedBy,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		AssigneeIDs: assigneeIDs,
	}, nil
}

// Update updates a task
func (s *service) Update(ctx context.Context, id uuid.UUID, req model.TaskUpdateRequest) error {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to get task for update", "error", err, "id", id)
		return fmt.Errorf("failed to get task: %w", err)
	}

	// Update fields if provided
	if req.Title != nil {
		task.Title = *req.Title
	}

	if req.Description != nil {
		task.Description = *req.Description
	}

	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}

	if req.Priority != nil {
		task.Priority = *req.Priority
	}

	if req.Status != nil {
		task.Status = *req.Status
	}

	if err := s.taskRepo.Update(ctx, task); err != nil {
		s.logger.Errorw("Failed to update task", "error", err, "id", id)
		return fmt.Errorf("failed to update task: %w", err)
	}

	// Update assignees if provided
	if req.AssigneeIDs != nil {
		// Get current assignments
		currentAssignments, _, err := s.assignmentRepo.ListByTask(ctx, id, 100, 0)
		if err != nil {
			s.logger.Warnw("Failed to get current task assignments", "error", err, "taskId", id)
			// Continue even if we can't get current assignments
		}

		// Create a map of current assignee IDs for fast lookup
		currentAssigneeMap := make(map[uuid.UUID]model.TaskAssignment)
		for _, assignment := range currentAssignments {
			currentAssigneeMap[assignment.UserID] = assignment
		}

		// Create new assignments
		for _, assigneeID := range *req.AssigneeIDs {
			if _, exists := currentAssigneeMap[assigneeID]; !exists {
				// Create new assignment
				assignment := model.TaskAssignment{
					TaskID:     id,
					UserID:     assigneeID,
					AssignedAt: time.Now(),
				}

				_, err := s.assignmentRepo.Create(ctx, assignment)
				if err != nil {
					s.logger.Warnw("Failed to create task assignment", "error", err, "taskId", id, "userId", assigneeID)
					// Continue even if one assignment fails
				}
			}
			// Remove from the map to track which ones need to be deleted
			delete(currentAssigneeMap, assigneeID)
		}

		// Delete assignments that are no longer needed
		for userID, assignment := range currentAssigneeMap {
			err := s.assignmentRepo.Delete(ctx, assignment.ID)
			if err != nil {
				s.logger.Warnw("Failed to delete task assignment", "error", err, "taskId", id, "userId", userID)
				// Continue even if one deletion fails
			}
		}
	}

	return nil
}

// Delete deletes a task
func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	// First check if the task exists
	_, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to get task for deletion", "error", err, "id", id)
		return fmt.Errorf("failed to get task: %w", err)
	}

	if err := s.taskRepo.Delete(ctx, id); err != nil {
		s.logger.Errorw("Failed to delete task", "error", err, "id", id)
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

// ListByEvent retrieves a list of tasks for a specific event with pagination
func (s *service) ListByEvent(ctx context.Context, eventID uuid.UUID, page, size int) (model.TasksResponse, error) {
	// Set default pagination values if not provided
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	offset := (page - 1) * size
	tasks, total, err := s.taskRepo.ListByEvent(ctx, eventID, size, offset)
	if err != nil {
		s.logger.Errorw("Failed to list event tasks", "error", err, "eventId", eventID, "page", page, "size", size)
		return model.TasksResponse{}, fmt.Errorf("failed to list event tasks: %w", err)
	}

	// Convert to response objects
	taskResponses := make([]model.TaskResponse, len(tasks))
	for i, task := range tasks {
		// Get task assignments
		assignments, _, err := s.assignmentRepo.ListByTask(ctx, task.ID, 100, 0)
		if err != nil {
			s.logger.Warnw("Failed to get task assignments", "error", err, "taskId", task.ID)
			// Continue even if we can't get assignments
			taskResponses[i] = model.TaskResponse{
				ID:          task.ID,
				EventID:     task.EventID,
				Title:       task.Title,
				Description: task.Description,
				DueDate:     task.DueDate,
				Priority:    task.Priority,
				Status:      task.Status,
				CreatedBy:   task.CreatedBy,
				CreatedAt:   task.CreatedAt,
				UpdatedAt:   task.UpdatedAt,
			}
			continue
		}

		// Extract assignee IDs
		assigneeIDs := make([]uuid.UUID, len(assignments))
		for j, assignment := range assignments {
			assigneeIDs[j] = assignment.UserID
		}

		taskResponses[i] = model.TaskResponse{
			ID:          task.ID,
			EventID:     task.EventID,
			Title:       task.Title,
			Description: task.Description,
			DueDate:     task.DueDate,
			Priority:    task.Priority,
			Status:      task.Status,
			CreatedBy:   task.CreatedBy,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
			AssigneeIDs: assigneeIDs,
		}
	}

	return model.TasksResponse{
		Tasks: taskResponses,
		Total: total,
	}, nil
}

// ListByUser retrieves a list of tasks assigned to a specific user with pagination
func (s *service) ListByUser(ctx context.Context, userID uuid.UUID, page, size int) (model.TasksResponse, error) {
	// Set default pagination values if not provided
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	offset := (page - 1) * size
	tasks, total, err := s.taskRepo.ListByUser(ctx, userID, size, offset)
	if err != nil {
		s.logger.Errorw("Failed to list user tasks", "error", err, "userId", userID, "page", page, "size", size)
		return model.TasksResponse{}, fmt.Errorf("failed to list user tasks: %w", err)
	}

	// Convert to response objects
	taskResponses := make([]model.TaskResponse, len(tasks))
	for i, task := range tasks {
		// Get task assignments for this task
		assignments, _, err := s.assignmentRepo.ListByTask(ctx, task.ID, 100, 0)
		if err != nil {
			s.logger.Warnw("Failed to get task assignments", "error", err, "taskId", task.ID)
			// Continue even if we can't get assignments
		}

		// Extract assignee IDs
		assigneeIDs := make([]uuid.UUID, len(assignments))
		for j, assignment := range assignments {
			assigneeIDs[j] = assignment.UserID
		}

		taskResponses[i] = model.TaskResponse{
			ID:          task.ID,
			EventID:     task.EventID,
			Title:       task.Title,
			Description: task.Description,
			DueDate:     task.DueDate,
			Priority:    task.Priority,
			Status:      task.Status,
			CreatedBy:   task.CreatedBy,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
			AssigneeIDs: assigneeIDs,
		}
	}

	return model.TasksResponse{
		Tasks: taskResponses,
		Total: total,
	}, nil
}

// UpdateStatus updates the status of a task
func (s *service) UpdateStatus(ctx context.Context, id uuid.UUID, status model.TaskStatus) error {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to get task for status update", "error", err, "id", id)
		return fmt.Errorf("failed to get task: %w", err)
	}

	task.Status = status

	if err := s.taskRepo.Update(ctx, task); err != nil {
		s.logger.Errorw("Failed to update task status", "error", err, "id", id, "status", status)
		return fmt.Errorf("failed to update task status: %w", err)
	}

	// If task is completed, update all assignments' completed_at
	if status == model.TaskStatusCompleted {
		assignments, _, err := s.assignmentRepo.ListByTask(ctx, id, 100, 0)
		if err != nil {
			s.logger.Warnw("Failed to get task assignments for completion", "error", err, "taskId", id)
			// Continue even if we can't get assignments
		} else {
			now := time.Now()
			for _, assignment := range assignments {
				assignment.CompletedAt = &now
				err := s.assignmentRepo.Update(ctx, assignment)
				if err != nil {
					s.logger.Warnw("Failed to update assignment completion", "error", err, "assignmentId", assignment.ID)
					// Continue even if one update fails
				}
			}
		}
	}

	return nil
}

// GetEventSummary retrieves a summary of tasks for a specific event
func (s *service) GetEventSummary(ctx context.Context, eventID uuid.UUID) (model.TaskEventSummary, error) {
	var summary model.TaskEventSummary
	summary.EventID = eventID

	// Get all tasks for the event
	tasks, _, err := s.taskRepo.ListByEvent(ctx, eventID, 1000, 0)
	if err != nil {
		s.logger.Errorw("Failed to list event tasks for summary", "error", err, "eventId", eventID)
		return summary, fmt.Errorf("failed to get event tasks: %w", err)
	}

	// Count tasks by status
	totalTasks := len(tasks)
	completedTasks := 0
	inProgressTasks := 0
	pendingTasks := 0

	for _, task := range tasks {
		switch task.Status {
		case model.TaskStatusCompleted:
			completedTasks++
		case model.TaskStatusInProgress:
			inProgressTasks++
		case model.TaskStatusPending:
			pendingTasks++
		}
	}

	// Calculate completion percentage
	var completionPercentage float64
	if totalTasks > 0 {
		completionPercentage = float64(completedTasks) / float64(totalTasks) * 100
	}

	summary.TotalTasks = totalTasks
	summary.CompletedTasks = completedTasks
	summary.InProgressTasks = inProgressTasks
	summary.PendingTasks = pendingTasks
	summary.CompletionPercentage = completionPercentage

	return summary, nil
}
