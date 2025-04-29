package task

import (
	"context"
	"github.com/PabloPerdolie/event-manager/core-service/internal/domain"
	"github.com/PabloPerdolie/event-manager/core-service/internal/model"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"testing"
	"time"
)

// Мок репозитория задач
type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(ctx context.Context, task model.Task) (int, error) {
	args := m.Called(ctx, task)
	return args.Int(0), args.Error(1)
}

func (m *MockTaskRepository) GetById(ctx context.Context, id int) (model.Task, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Task), args.Error(1)
}

func (m *MockTaskRepository) Update(ctx context.Context, task model.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTaskRepository) ListByEvent(ctx context.Context, eventId, limit, offset int) ([]model.Task, error) {
	args := m.Called(ctx, eventId, limit, offset)
	return args.Get(0).([]model.Task), args.Error(1)
}

func (m *MockTaskRepository) ListByUser(ctx context.Context, userId, limit, offset int) ([]model.Task, error) {
	args := m.Called(ctx, userId, limit, offset)
	return args.Get(0).([]model.Task), args.Error(1)
}

func (m *MockTaskRepository) ListByStatus(ctx context.Context, eventId int, status string, limit, offset int) ([]model.Task, error) {
	args := m.Called(ctx, eventId, status, limit, offset)
	return args.Get(0).([]model.Task), args.Error(1)
}

// Мок репозитория назначений задач
type MockAssignmentRepository struct {
	mock.Mock
}

func (m *MockAssignmentRepository) Create(ctx context.Context, assignment model.TaskAssignment) (int, error) {
	args := m.Called(ctx, assignment)
	return args.Int(0), args.Error(1)
}

func (m *MockAssignmentRepository) GetByTaskAndUser(ctx context.Context, taskId, userId int) (model.TaskAssignment, error) {
	args := m.Called(ctx, taskId, userId)
	return args.Get(0).(model.TaskAssignment), args.Error(1)
}

func (m *MockAssignmentRepository) Update(ctx context.Context, assignment model.TaskAssignment) error {
	args := m.Called(ctx, assignment)
	return args.Error(0)
}

func (m *MockAssignmentRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAssignmentRepository) ListByTask(ctx context.Context, taskId, limit, offset int) ([]model.TaskAssignment, error) {
	args := m.Called(ctx, taskId, limit, offset)
	return args.Get(0).([]model.TaskAssignment), args.Error(1)
}

func (m *MockAssignmentRepository) ListByUser(ctx context.Context, userId, limit, offset int) ([]model.TaskAssignment, error) {
	args := m.Called(ctx, userId, limit, offset)
	return args.Get(0).([]model.TaskAssignment), args.Error(1)
}

func setupTaskService() (*Service, *MockTaskRepository, *MockAssignmentRepository) {
	taskRepo := new(MockTaskRepository)
	assignmentRepo := new(MockAssignmentRepository)
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()

	service := NewService(taskRepo, assignmentRepo, sugar)

	return &service, taskRepo, assignmentRepo
}

// Тесты для метода Create

// Тест 1: Успешное создание задачи без назначений
func TestCreate_Success_WithoutAssignment(t *testing.T) {
	// Подготовка
	service, taskRepo, _ := setupTaskService()
	ctx := context.Background()

	taskID := 1

	storyPoints := 3
	priorityStr := "5"

	req := domain.TaskCreateRequest{
		EventId:     10,
		ParentId:    nil,
		Title:       "Test Task",
		Description: "Test Description",
		StoryPoints: &storyPoints,
		Priority:    &priorityStr,
		AssignedTo:  nil, // без назначения
	}

	// Настраиваем моки
	taskRepo.On("Create", ctx, mock.MatchedBy(func(task model.Task) bool {
		return task.EventId == req.EventId &&
			task.Title == req.Title &&
			task.Description == req.Description &&
			*task.StoryPoints == *req.StoryPoints &&
			*task.Priority == *req.Priority &&
			task.Status == string(domain.TaskStatusPending)
	})).Return(taskID, nil)

	// Действие
	resp, err := service.Create(ctx, req)

	// Проверка
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, taskID, resp.Id)
	assert.Equal(t, req.EventId, resp.EventId)
	assert.Equal(t, req.Title, resp.Title)
	assert.Equal(t, req.Description, resp.Description)
	assert.Equal(t, req.StoryPoints, resp.StoryPoints)
	assert.Equal(t, req.Priority, resp.Priority)
	assert.Equal(t, domain.TaskStatusPending, resp.Status)
	assert.Nil(t, resp.AssignedTo)

	taskRepo.AssertExpectations(t)
}

// Тест 2: Успешное создание задачи с назначением
func TestCreate_Success_WithAssignment(t *testing.T) {
	// Подготовка
	service, taskRepo, assignmentRepo := setupTaskService()
	ctx := context.Background()

	taskID := 1
	assigneeID := 5

	storyPoints := 3
	priorityStr := "5"

	req := domain.TaskCreateRequest{
		EventId:     10,
		ParentId:    nil,
		Title:       "Test Task",
		Description: "Test Description",
		StoryPoints: &storyPoints,
		Priority:    &priorityStr,
		AssignedTo:  &assigneeID,
	}

	// Настраиваем моки
	taskRepo.On("Create", ctx, mock.MatchedBy(func(task model.Task) bool {
		return task.EventId == req.EventId &&
			task.Title == req.Title &&
			task.Description == req.Description &&
			*task.StoryPoints == *req.StoryPoints &&
			*task.Priority == *req.Priority &&
			task.Status == string(domain.TaskStatusPending)
	})).Return(taskID, nil)

	assignmentRepo.On("Create", ctx, mock.MatchedBy(func(assignment model.TaskAssignment) bool {
		return assignment.TaskId == taskID && assignment.UserId == assigneeID
	})).Return(1, nil)

	// Действие
	resp, err := service.Create(ctx, req)

	// Проверка
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, taskID, resp.Id)
	assert.Equal(t, req.EventId, resp.EventId)
	assert.Equal(t, req.Title, resp.Title)
	assert.Equal(t, req.Description, resp.Description)
	assert.Equal(t, req.StoryPoints, resp.StoryPoints)
	assert.Equal(t, req.Priority, resp.Priority)
	assert.Equal(t, domain.TaskStatusPending, resp.Status)
	assert.NotNil(t, resp.AssignedTo)
	assert.Equal(t, assigneeID, *resp.AssignedTo)

	taskRepo.AssertExpectations(t)
	assignmentRepo.AssertExpectations(t)
}

// Тест 3: Ошибка при создании задачи
func TestCreate_Error_TaskCreation(t *testing.T) {
	// Подготовка
	service, taskRepo, _ := setupTaskService()
	ctx := context.Background()

	storyPoints := 3
	priorityStr := "5"

	req := domain.TaskCreateRequest{
		EventId:     10,
		Title:       "Test Task",
		Description: "Test Description",
		StoryPoints: &storyPoints,
		Priority:    &priorityStr,
	}

	// Настраиваем моки
	expectedError := errors.New("database error")
	taskRepo.On("Create", ctx, mock.AnythingOfType("model.Task")).Return(0, expectedError)

	// Действие
	resp, err := service.Create(ctx, req)

	// Проверка
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "create task")

	taskRepo.AssertExpectations(t)
}

// Тест 4: Успешное создание задачи с ошибкой при создании назначения
func TestCreate_Success_WithAssignmentError(t *testing.T) {
	// Подготовка
	service, taskRepo, assignmentRepo := setupTaskService()
	ctx := context.Background()

	taskID := 1
	assigneeID := 5

	storyPoints := 3
	priorityStr := "5"

	req := domain.TaskCreateRequest{
		EventId:     10,
		Title:       "Test Task",
		Description: "Test Description",
		StoryPoints: &storyPoints,
		Priority:    &priorityStr,
		AssignedTo:  &assigneeID,
	}

	// Настраиваем моки
	taskRepo.On("Create", ctx, mock.AnythingOfType("model.Task")).Return(taskID, nil)

	assignmentError := errors.New("assignment error")
	assignmentRepo.On("Create", ctx, mock.AnythingOfType("model.TaskAssignment")).Return(0, assignmentError)

	// Действие
	resp, err := service.Create(ctx, req)

	// Проверка - задача должна быть создана успешно, несмотря на ошибку назначения
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, taskID, resp.Id)

	taskRepo.AssertExpectations(t)
	assignmentRepo.AssertExpectations(t)
}

// Тесты для метода Update

// Тест 1: Успешное обновление задачи без изменения назначения
func TestUpdate_Success_NoAssignmentChange(t *testing.T) {
	// Подготовка
	service, taskRepo, _ := setupTaskService()
	ctx := context.Background()
	taskID := 1

	existingTask := model.Task{
		TaskId:      taskID,
		EventId:     10,
		Title:       "Old Title",
		Description: "Old Description",
		Priority:    func() *string { s := "3"; return &s }(),
		Status:      string(domain.TaskStatusPending),
	}

	newTitle := "New Title"
	newDescription := "New Description"
	newPriority := "5"

	req := domain.TaskUpdateRequest{
		Title:       &newTitle,
		Description: &newDescription,
		Priority:    &newPriority,
		AssignedTo:  nil, // не меняем назначение
	}

	// Настраиваем моки
	taskRepo.On("GetById", ctx, taskID).Return(existingTask, nil)

	taskRepo.On("Update", ctx, mock.MatchedBy(func(task model.Task) bool {
		return task.TaskId == taskID &&
			task.Title == newTitle &&
			task.Description == newDescription &&
			*task.Priority == newPriority &&
			task.Status == string(domain.TaskStatusPending)
	})).Return(nil)

	// Действие
	err := service.Update(ctx, taskID, req)

	// Проверка
	assert.NoError(t, err)
	taskRepo.AssertExpectations(t)
}

// Тест 2: Успешное обновление задачи с изменением назначения
func TestUpdate_Success_WithAssignmentChange(t *testing.T) {
	// Подготовка
	service, taskRepo, assignmentRepo := setupTaskService()
	ctx := context.Background()
	taskID := 1
	newAssigneeID := 2

	existingTask := model.Task{
		TaskId:      taskID,
		EventId:     10,
		Title:       "Old Title",
		Description: "Old Description",
		Priority:    func() *string { s := "3"; return &s }(),
		Status:      string(domain.TaskStatusPending),
	}

	existingAssignments := []model.TaskAssignment{
		{
			TaskAssignmentID: 100,
			TaskId:           taskID,
			UserId:           1, // старый исполнитель
			AssignedAt:       time.Now().Add(-24 * time.Hour),
		},
	}

	newTitle := "New Title"

	req := domain.TaskUpdateRequest{
		Title:      &newTitle,
		AssignedTo: &newAssigneeID, // меняем назначение
	}

	// Настраиваем моки
	taskRepo.On("GetById", ctx, taskID).Return(existingTask, nil)
	taskRepo.On("Update", ctx, mock.AnythingOfType("model.Task")).Return(nil)

	assignmentRepo.On("ListByTask", ctx, taskID, 100, 0).Return(existingAssignments, nil)

	// Создание нового назначения
	assignmentRepo.On("Create", ctx, mock.MatchedBy(func(assignment model.TaskAssignment) bool {
		return assignment.TaskId == taskID && assignment.UserId == newAssigneeID
	})).Return(101, nil)

	// Удаление старого назначения
	assignmentRepo.On("Delete", ctx, existingAssignments[0].UserId).Return(nil)

	// Действие
	err := service.Update(ctx, taskID, req)

	// Проверка
	assert.NoError(t, err)
	taskRepo.AssertExpectations(t)
	assignmentRepo.AssertExpectations(t)
}

// Тест 3: Ошибка при получении задачи для обновления
func TestUpdate_Error_GetTask(t *testing.T) {
	// Подготовка
	service, taskRepo, _ := setupTaskService()
	ctx := context.Background()
	taskID := 1

	req := domain.TaskUpdateRequest{
		Title: nil,
	}

	// Настраиваем моки
	expectedError := errors.New("task not found")
	taskRepo.On("GetById", ctx, taskID).Return(model.Task{}, expectedError)

	// Действие
	err := service.Update(ctx, taskID, req)

	// Проверка
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "get task")

	taskRepo.AssertExpectations(t)
}

// Тест 4: Ошибка при обновлении задачи
func TestUpdate_Error_UpdateTask(t *testing.T) {
	// Подготовка
	service, taskRepo, _ := setupTaskService()
	ctx := context.Background()
	taskID := 1

	existingTask := model.Task{
		TaskId:      taskID,
		EventId:     10,
		Title:       "Old Title",
		Description: "Old Description",
		Priority:    func() *string { s := "3"; return &s }(),
		Status:      string(domain.TaskStatusPending),
	}

	newTitle := "New Title"

	req := domain.TaskUpdateRequest{
		Title: &newTitle,
	}

	// Настраиваем моки
	taskRepo.On("GetById", ctx, taskID).Return(existingTask, nil)

	updateError := errors.New("update error")
	taskRepo.On("Update", ctx, mock.AnythingOfType("model.Task")).Return(updateError)

	// Действие
	err := service.Update(ctx, taskID, req)

	// Проверка
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update task")

	taskRepo.AssertExpectations(t)
}

// Тесты для методов получения списка задач

// Тест 1: Успешное получение списка задач по событию
func TestListByEvent_Success(t *testing.T) {
	// Подготовка
	service, taskRepo, assignmentRepo := setupTaskService()
	ctx := context.Background()
	eventID := 10
	page := 1
	size := 10

	tasks := []model.Task{
		{
			TaskId:      1,
			EventId:     eventID,
			Title:       "Task 1",
			Description: "Description 1",
			Status:      string(domain.TaskStatusPending),
			CreatedAt:   time.Now(),
		},
		{
			TaskId:      2,
			EventId:     eventID,
			Title:       "Task 2",
			Description: "Description 2",
			Status:      string(domain.TaskStatusInProgress),
			CreatedAt:   time.Now(),
		},
	}

	// Настраиваем моки для получения задач
	taskRepo.On("ListByEvent", ctx, eventID, size, 0).Return(tasks, nil)

	// Настраиваем моки для получения назначений
	assignmentRepo.On("ListByTask", ctx, tasks[0].TaskId, 100, 0).Return([]model.TaskAssignment{}, nil)
	assignmentRepo.On("ListByTask", ctx, tasks[1].TaskId, 100, 0).Return([]model.TaskAssignment{}, nil)

	// Действие
	resp, err := service.ListByEvent(ctx, eventID, page, size)

	// Проверка
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, len(tasks), resp.Total)
	assert.Equal(t, len(tasks), len(resp.Tasks))

	// Проверяем соответствие полей в ответе
	for i, task := range tasks {
		assert.Equal(t, task.TaskId, resp.Tasks[i].Id)
		assert.Equal(t, task.EventId, resp.Tasks[i].EventId)
		assert.Equal(t, task.Title, resp.Tasks[i].Title)
		assert.Equal(t, task.Description, resp.Tasks[i].Description)
		assert.Equal(t, domain.TaskStatus(task.Status), resp.Tasks[i].Status)
	}

	taskRepo.AssertExpectations(t)
	assignmentRepo.AssertExpectations(t)
}

// Тест 2: Ошибка при получении списка задач по событию
func TestListByEvent_Error(t *testing.T) {
	// Подготовка
	service, taskRepo, _ := setupTaskService()
	ctx := context.Background()
	eventID := 10
	page := 1
	size := 10

	// Настраиваем моки
	expectedError := errors.New("database error")
	taskRepo.On("ListByEvent", ctx, eventID, size, 0).Return([]model.Task{}, expectedError)

	// Действие
	resp, err := service.ListByEvent(ctx, eventID, page, size)

	// Проверка
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "list event tasks")

	taskRepo.AssertExpectations(t)
}

// Тест 3: Успешное получение списка задач по пользователю
func TestListByUser_Success(t *testing.T) {
	// Подготовка
	service, taskRepo, assignmentRepo := setupTaskService()
	ctx := context.Background()
	userID := 5
	page := 1
	size := 10

	tasks := []model.Task{
		{
			TaskId:      1,
			EventId:     10,
			Title:       "User Task 1",
			Description: "Description 1",
			Status:      string(domain.TaskStatusPending),
			CreatedAt:   time.Now(),
		},
		{
			TaskId:      2,
			EventId:     11,
			Title:       "User Task 2",
			Description: "Description 2",
			Status:      string(domain.TaskStatusCompleted),
			CreatedAt:   time.Now(),
		},
	}

	assignments := []model.TaskAssignment{
		{
			TaskAssignmentID: 101,
			TaskId:           1,
			UserId:           userID,
			AssignedAt:       time.Now(),
		},
	}

	// Настраиваем моки
	taskRepo.On("ListByUser", ctx, userID, size, 0).Return(tasks, nil)

	// Для первой задачи возвращаем назначение
	assignmentRepo.On("ListByTask", ctx, tasks[0].TaskId, 100, 0).Return(assignments, nil)

	// Для второй задачи возвращаем пустой список
	assignmentRepo.On("ListByTask", ctx, tasks[1].TaskId, 100, 0).Return([]model.TaskAssignment{}, nil)

	// Действие
	resp, err := service.ListByUser(ctx, userID, page, size)

	// Проверка
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, len(tasks), resp.Total)
	assert.Equal(t, len(tasks), len(resp.Tasks))

	// Проверяем соответствие полей в ответе
	for i, task := range tasks {
		assert.Equal(t, task.TaskId, resp.Tasks[i].Id)
		assert.Equal(t, task.EventId, resp.Tasks[i].EventId)
		assert.Equal(t, task.Title, resp.Tasks[i].Title)
		assert.Equal(t, task.Description, resp.Tasks[i].Description)
		assert.Equal(t, domain.TaskStatus(task.Status), resp.Tasks[i].Status)
	}

	// Для первой задачи должен быть установлен AssignedTo
	assert.NotNil(t, resp.Tasks[0].AssignedTo)
	assert.Equal(t, userID, *resp.Tasks[0].AssignedTo)

	taskRepo.AssertExpectations(t)
	assignmentRepo.AssertExpectations(t)
}

// Тест 4: Ошибка при получении списка задач по пользователю
func TestListByUser_Error(t *testing.T) {
	// Подготовка
	service, taskRepo, _ := setupTaskService()
	ctx := context.Background()
	userID := 5
	page := 1
	size := 10

	// Настраиваем моки
	expectedError := errors.New("database error")
	taskRepo.On("ListByUser", ctx, userID, size, 0).Return([]model.Task{}, expectedError)

	// Действие
	resp, err := service.ListByUser(ctx, userID, page, size)

	// Проверка
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "list user tasks")

	taskRepo.AssertExpectations(t)
}
