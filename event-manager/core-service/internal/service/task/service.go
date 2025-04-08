package task

import (
	"context"
	"github.com/PabloPerdolie/event-manager/core-service/internal/domain"
	"github.com/pkg/errors"
	"time"

	"github.com/PabloPerdolie/event-manager/core-service/internal/model"
	"go.uber.org/zap"
)

type Repository interface {
	Create(ctx context.Context, task model.Task) (int, error)
	GetById(ctx context.Context, id int) (model.Task, error)
	Update(ctx context.Context, task model.Task) error
	Delete(ctx context.Context, id int) error
	ListByEvent(ctx context.Context, eventId, limit, offset int) ([]model.Task, error)
	ListByUser(ctx context.Context, userId, limit, offset int) ([]model.Task, error)
	ListByStatus(ctx context.Context, eventId int, status string, limit, offset int) ([]model.Task, error)
}

type AssignmentRepository interface {
	Create(ctx context.Context, assignment model.TaskAssignment) (int, error)
	GetByTaskAndUser(ctx context.Context, taskId, userId int) (model.TaskAssignment, error)
	Update(ctx context.Context, assignment model.TaskAssignment) error
	Delete(ctx context.Context, id int) error
	ListByTask(ctx context.Context, taskId, limit, offset int) ([]model.TaskAssignment, error)
	ListByUser(ctx context.Context, userId, limit, offset int) ([]model.TaskAssignment, error)
}

type Service struct {
	taskRepo       Repository
	assignmentRepo AssignmentRepository
	logger         *zap.SugaredLogger
}

func NewService(taskRepo Repository, assignmentRepo AssignmentRepository, logger *zap.SugaredLogger) Service {
	return Service{
		taskRepo:       taskRepo,
		assignmentRepo: assignmentRepo,
		logger:         logger,
	}
}

func (s Service) Create(ctx context.Context, req domain.TaskCreateRequest) (*domain.TaskResponse, error) {
	task := model.Task{
		EventId:     req.EventId,
		ParentId:    req.ParentId,
		Title:       req.Title,
		Description: req.Description,
		StoryPoints: req.StoryPoints,
		Priority:    req.Priority,
		Status:      string(domain.TaskStatusPending),
		CreatedAt:   time.Now(),
	}

	id, err := s.taskRepo.Create(ctx, task)
	if err != nil {
		s.logger.Errorw("Failed to create task", "error", err, "eventId", req.EventId)
		return nil, errors.WithMessage(err, "create task")
	}

	if req.AssignedTo != nil {
		assignment := model.TaskAssignment{
			TaskId:     id,
			UserId:     *req.AssignedTo,
			AssignedAt: time.Now(),
		}

		_, err := s.assignmentRepo.Create(ctx, assignment)
		if err != nil {
			s.logger.Warnw("Failed to create task assignment", "error", err, "taskId", id, "userId", req.AssignedTo)
			// Continue even if one assignment fails
		}
	}

	return &domain.TaskResponse{
		Id:          id,
		EventId:     task.EventId,
		AssignedTo:  req.AssignedTo,
		ParentID:    task.ParentId,
		Title:       task.Title,
		Description: task.Description,
		StoryPoints: task.StoryPoints,
		Priority:    task.Priority,
		Status:      domain.TaskStatus(task.Status),
		CreatedAt:   task.CreatedAt,
	}, nil
}

func (s Service) Update(ctx context.Context, id int, req domain.TaskUpdateRequest) error {
	task, err := s.taskRepo.GetById(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to get task for update", "error", err, "id", id)
		return errors.WithMessage(err, "get task")
	}

	if req.Title != nil {
		task.Title = *req.Title
	}

	if req.Description != nil {
		task.Description = *req.Description
	}

	if req.Priority != nil {
		task.Priority = req.Priority
	}

	if req.Status != nil {
		task.Status = string(*req.Status)
	}

	if err := s.taskRepo.Update(ctx, task); err != nil {
		s.logger.Errorw("Failed to update task", "error", err, "id", id)
		return errors.WithMessage(err, "update task")
	}

	assignedTo := req.AssignedTo

	if assignedTo != nil {
		currentAssignments, err := s.assignmentRepo.ListByTask(ctx, id, 100, 0)
		if err != nil {
			s.logger.Warnw("Failed to get current task assignments", "error", err, "taskId", id)
			// Continue even if we can't get current assignments
		}

		currentAssigneeMap := make(map[int]model.TaskAssignment)
		for _, assignment := range currentAssignments {
			currentAssigneeMap[assignment.UserId] = assignment
		}

		//for _, assigneeId := range *req.AssignedTo { // todo

		if _, exists := currentAssigneeMap[*assignedTo]; !exists {
			assignment := model.TaskAssignment{
				TaskId:     id,
				UserId:     *assignedTo,
				AssignedAt: time.Now(),
			}

			_, err := s.assignmentRepo.Create(ctx, assignment)
			if err != nil {
				s.logger.Warnw("Failed to create task assignment", "error", err, "taskId", id, "userId", assignedTo)
				// Continue even if one assignment fails
			}
		}
		delete(currentAssigneeMap, *assignedTo)

		for userId, assignment := range currentAssigneeMap {
			err := s.assignmentRepo.Delete(ctx, assignment.UserId)
			if err != nil {
				s.logger.Warnw("Failed to delete task assignment", "error", err, "taskId", id, "userId", userId)
				// Continue even if one deletion fails
			}
		}
	}

	return nil
}

func (s Service) Delete(ctx context.Context, id int) error {
	if err := s.taskRepo.Delete(ctx, id); err != nil {
		s.logger.Errorw("Failed to delete task", "error", err, "id", id)
		return errors.WithMessage(err, "delete task")
	}

	return nil
}

func (s Service) ListByEvent(ctx context.Context, eventId int, page, size int) (*domain.TasksResponse, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	offset := (page - 1) * size
	tasks, err := s.taskRepo.ListByEvent(ctx, eventId, size, offset)
	if err != nil {
		s.logger.Errorw("Failed to list event tasks", "error", err, "eventId", eventId, "page", page, "size", size)
		return nil, errors.WithMessage(err, "list event tasks")
	}

	return s.convertToTasksResponse(ctx, tasks), nil
}

func (s Service) ListByUser(ctx context.Context, userId int, page, size int) (*domain.TasksResponse, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	offset := (page - 1) * size
	tasks, err := s.taskRepo.ListByUser(ctx, userId, size, offset)
	if err != nil {
		s.logger.Errorw("Failed to list user tasks", "error", err, "userId", userId, "page", page, "size", size)
		return nil, errors.WithMessage(err, "list user tasks")
	}

	return s.convertToTasksResponse(ctx, tasks), nil
}

func (s Service) UpdateStatus(ctx context.Context, id int, status domain.TaskStatus) error {
	task, err := s.taskRepo.GetById(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to get task for status update", "error", err, "id", id)
		return errors.WithMessage(err, "get task")
	}

	task.Status = string(status)

	if err := s.taskRepo.Update(ctx, task); err != nil {
		s.logger.Errorw("Failed to update task status", "error", err, "id", id, "status", status)
		return errors.WithMessage(err, "update task status")
	}

	if status == domain.TaskStatusCompleted {
		assignments, err := s.assignmentRepo.ListByTask(ctx, id, 100, 0)
		if err != nil {
			s.logger.Warnw("Failed to get task assignments for completion", "error", err, "taskId", id)
			// Continue even if we can't get assignments
		} else {
			now := time.Now()
			for _, assignment := range assignments {
				assignment.CompletedAt = &now
				err := s.assignmentRepo.Update(ctx, assignment)
				if err != nil {
					s.logger.Warnw("Failed to update assignment completion", "error", err, "assignmentId", assignment.TaskAssignmentID)
					// Continue even if one update fails
				}
			}
		}
	}

	return nil
}

func (s Service) convertToTasksResponse(ctx context.Context, tasks []model.Task) *domain.TasksResponse {
	taskResponses := make([]domain.TaskResponse, len(tasks))
	for i, task := range tasks {
		taskResponses[i] = domain.TaskResponse{
			Id:          task.TaskId,
			EventId:     task.EventId,
			Title:       task.Title,
			Description: task.Description,
			Priority:    task.Priority,
			Status:      domain.TaskStatus(task.Status),
			CreatedAt:   task.CreatedAt,
		}

		assignments, err := s.assignmentRepo.ListByTask(ctx, task.TaskId, 100, 0)
		if err != nil || len(assignments) == 0 {
			s.logger.Warnw("Failed to get task assignments", "error", err, "taskId", task.TaskId)
			continue
		}

		assigneeIds := make([]int, len(assignments))
		for j, assignment := range assignments {
			assigneeIds[j] = assignment.UserId
		}

		taskResponses[i].AssignedTo = &assigneeIds[0]
	}

	return &domain.TasksResponse{
		Tasks: taskResponses,
		Total: len(tasks),
	}
}
