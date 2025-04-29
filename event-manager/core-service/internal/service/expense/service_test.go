package expense

import (
	"context"
	"github.com/PabloPerdolie/event-manager/core-service/internal/domain"
	"github.com/PabloPerdolie/event-manager/core-service/internal/model"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"testing"
)

// Mock репозитория расходов
type MockExpenseRepository struct {
	mock.Mock
}

func (m *MockExpenseRepository) CreateExpense(ctx context.Context, expense model.Expense) (int, error) {
	args := m.Called(ctx, expense)
	return args.Int(0), args.Error(1)
}

func (m *MockExpenseRepository) GetExpenseById(ctx context.Context, id int) (*model.Expense, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Expense), args.Error(1)
}

func (m *MockExpenseRepository) ListExpensesByEventId(ctx context.Context, eventId int, limit, offset int) ([]model.Expense, int, error) {
	args := m.Called(ctx, eventId, limit, offset)
	return args.Get(0).([]model.Expense), args.Int(1), args.Error(2)
}

func (m *MockExpenseRepository) GetEventTotalExpenses(ctx context.Context, eventId int) (float64, error) {
	args := m.Called(ctx, eventId)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockExpenseRepository) DeleteExpense(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Mock репозитория долей расходов
type MockExpenseShareRepository struct {
	mock.Mock
}

func (m *MockExpenseShareRepository) CreateExpenseShare(ctx context.Context, share model.ExpenseShare) (int, error) {
	args := m.Called(ctx, share)
	return args.Int(0), args.Error(1)
}

func (m *MockExpenseShareRepository) ListExpenseSharesByExpenseId(ctx context.Context, expenseId int) ([]model.ExpenseShare, error) {
	args := m.Called(ctx, expenseId)
	return args.Get(0).([]model.ExpenseShare), args.Error(1)
}

func (m *MockExpenseShareRepository) UpdateExpenseSharePaidStatus(ctx context.Context, shareId int, isPaid bool) error {
	args := m.Called(ctx, shareId, isPaid)
	return args.Error(0)
}

func (m *MockExpenseShareRepository) DeleteExpenseSharesByExpenseId(ctx context.Context, expenseId int) error {
	args := m.Called(ctx, expenseId)
	return args.Error(0)
}

func (m *MockExpenseShareRepository) GetEventBalanceReport(ctx context.Context, eventId int) ([]domain.UserBalance, error) {
	args := m.Called(ctx, eventId)
	return args.Get(0).([]domain.UserBalance), args.Error(1)
}

func setupService() (*Service, *MockExpenseRepository, *MockExpenseShareRepository) {
	expenseRepo := new(MockExpenseRepository)
	expenseShareRepo := new(MockExpenseShareRepository)
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()

	service := NewService(expenseRepo, expenseShareRepo, sugar)

	return &service, expenseRepo, expenseShareRepo
}

// Тест 1: Успешное создание долей расходов с методом разделения Equal
func TestCreateExpenseShares_Equal_Success(t *testing.T) {
	// Подготовка
	service, expenseRepo, expenseShareRepo := setupService()
	ctx := context.Background()
	expenseID := 1
	expense := &model.Expense{
		ExpenseID: expenseID,
		Amount:    100.0,
	}
	userIDs := []int{1, 2, 3, 4}

	// Настраиваем моки
	expenseRepo.On("GetExpenseById", ctx, expenseID).Return(expense, nil)
	expenseShareRepo.On("DeleteExpenseSharesByExpenseId", ctx, expenseID).Return(nil)

	// Для каждого пользователя настраиваем создание доли
	expectedAmount := 25.0 // 100 / 4 = 25
	for _, userID := range userIDs {
		expectedShare := model.ExpenseShare{
			ExpenseID: expenseID,
			UserID:    userID,
			Amount:    expectedAmount,
			IsPaid:    false,
		}
		expenseShareRepo.On("CreateExpenseShare", ctx, mock.MatchedBy(func(share model.ExpenseShare) bool {
			return share.ExpenseID == expectedShare.ExpenseID &&
				share.UserID == expectedShare.UserID &&
				share.Amount == expectedShare.Amount &&
				share.IsPaid == expectedShare.IsPaid
		})).Return(userID, nil).Once()
	}

	// Действие
	err := service.CreateExpenseShares(ctx, expenseID, userIDs, model.SplitMethodEqual)

	// Проверка
	assert.NoError(t, err)
	expenseRepo.AssertExpectations(t)
	expenseShareRepo.AssertExpectations(t)
}

// Тест 2: Успешное создание долей расходов с методом разделения Percent
func TestCreateExpenseShares_Percent_Success(t *testing.T) {
	// Подготовка
	service, expenseRepo, expenseShareRepo := setupService()
	ctx := context.Background()
	expenseID := 1
	expense := &model.Expense{
		ExpenseID: expenseID,
		Amount:    100.0,
	}
	userIDs := []int{1, 2}

	// Настраиваем моки
	expenseRepo.On("GetExpenseById", ctx, expenseID).Return(expense, nil)
	expenseShareRepo.On("DeleteExpenseSharesByExpenseId", ctx, expenseID).Return(nil)

	// Для каждого пользователя настраиваем создание доли
	expectedAmount := 50.0 // 100 / 2 = 50
	for _, userID := range userIDs {
		expectedShare := model.ExpenseShare{
			ExpenseID: expenseID,
			UserID:    userID,
			Amount:    expectedAmount,
			IsPaid:    false,
		}
		expenseShareRepo.On("CreateExpenseShare", ctx, mock.MatchedBy(func(share model.ExpenseShare) bool {
			return share.ExpenseID == expectedShare.ExpenseID &&
				share.UserID == expectedShare.UserID &&
				share.Amount == expectedShare.Amount &&
				share.IsPaid == expectedShare.IsPaid
		})).Return(userID, nil).Once()
	}

	// Действие
	err := service.CreateExpenseShares(ctx, expenseID, userIDs, model.SplitMethodPercent)

	// Проверка
	assert.NoError(t, err)
	expenseRepo.AssertExpectations(t)
	expenseShareRepo.AssertExpectations(t)
}

// Тест 3: Ошибка при использовании метода разделения Exact (пока не поддерживается полностью)
func TestCreateExpenseShares_Exact_NotSupported(t *testing.T) {
	// Подготовка
	service, expenseRepo, expenseShareRepo := setupService()
	ctx := context.Background()
	expenseID := 1
	expense := &model.Expense{
		ExpenseID: expenseID,
		Amount:    100.0,
	}
	userIDs := []int{1, 2, 3}

	// Настраиваем моки
	expenseRepo.On("GetExpenseById", ctx, expenseID).Return(expense, nil)
	expenseShareRepo.On("DeleteExpenseSharesByExpenseId", ctx, expenseID).Return(nil)

	// Действие
	err := service.CreateExpenseShares(ctx, expenseID, userIDs, model.SplitMethodExact)

	// Проверка
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exact split method requires additional data")
	expenseRepo.AssertExpectations(t)
	expenseShareRepo.AssertExpectations(t)
}

// Тест 4: Ошибка при попытке разделить расход без участников
func TestCreateExpenseShares_NoParticipants(t *testing.T) {
	// Подготовка
	service, expenseRepo, expenseShareRepo := setupService()
	ctx := context.Background()
	expenseID := 1
	expense := &model.Expense{
		ExpenseID: expenseID,
		Amount:    100.0,
	}
	userIDs := []int{} // Пустой список пользователей

	// Настраиваем моки
	expenseRepo.On("GetExpenseById", ctx, expenseID).Return(expense, nil)
	expenseShareRepo.On("DeleteExpenseSharesByExpenseId", ctx, expenseID).Return(nil)

	// Действие
	err := service.CreateExpenseShares(ctx, expenseID, userIDs, model.SplitMethodEqual)

	// Проверка
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot split expense with no participants")
	expenseRepo.AssertExpectations(t)
	expenseShareRepo.AssertExpectations(t)
}

// Тест 1: Успешное получение отчета о балансе с непустым списком пользователей
func TestGetEventBalanceReport_Success(t *testing.T) {
	// Подготовка
	service, expenseRepo, expenseShareRepo := setupService()
	ctx := context.Background()
	eventID := 1

	// Настраиваем моки
	totalAmount := 300.0
	expenseRepo.On("GetEventTotalExpenses", ctx, eventID).Return(totalAmount, nil)

	userBalances := []domain.UserBalance{
		{UserID: 1, Username: "user1", Balance: 100.0},
		{UserID: 2, Username: "user2", Balance: -50.0},
		{UserID: 3, Username: "user3", Balance: -50.0},
	}
	expenseShareRepo.On("GetEventBalanceReport", ctx, eventID).Return(userBalances, nil)

	// Действие
	report, err := service.GetEventBalanceReport(ctx, eventID)

	// Проверка
	assert.NoError(t, err)
	assert.Equal(t, eventID, report.EventID)
	assert.Equal(t, totalAmount, report.TotalAmount)
	assert.Equal(t, len(userBalances), len(report.UserBalances))

	for i, balance := range userBalances {
		assert.Equal(t, balance.UserID, report.UserBalances[i].UserID)
		assert.Equal(t, balance.Username, report.UserBalances[i].Username)
		assert.Equal(t, balance.Balance, report.UserBalances[i].Balance)
	}

	expenseRepo.AssertExpectations(t)
	expenseShareRepo.AssertExpectations(t)
}

// Тест 2: Обработка ошибки при получении общей суммы расходов
func TestGetEventBalanceReport_ErrorOnGetTotalExpenses(t *testing.T) {
	// Подготовка
	service, expenseRepo, expenseShareRepo := setupService()
	ctx := context.Background()
	eventID := 1

	// Настраиваем мок для ошибки при получении общей суммы
	expenseRepo.On("GetEventTotalExpenses", ctx, eventID).Return(0.0, errors.New("database error"))

	// Действие
	_, err := service.GetEventBalanceReport(ctx, eventID)

	// Проверка
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get event total expenses")

	expenseRepo.AssertExpectations(t)
	// Проверяем, что второй репозиторий не вызывался
	expenseShareRepo.AssertNotCalled(t, "GetEventBalanceReport")
}

// Тест 3: Обработка ошибки при получении баланса пользователей
func TestGetEventBalanceReport_ErrorOnGetUserBalances(t *testing.T) {
	// Подготовка
	service, expenseRepo, expenseShareRepo := setupService()
	ctx := context.Background()
	eventID := 1

	// Настраиваем моки
	totalAmount := 300.0
	expenseRepo.On("GetEventTotalExpenses", ctx, eventID).Return(totalAmount, nil)
	expenseShareRepo.On("GetEventBalanceReport", ctx, eventID).Return([]domain.UserBalance{}, errors.New("user balance error"))

	// Действие
	_, err := service.GetEventBalanceReport(ctx, eventID)

	// Проверка
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get event balance report")

	expenseRepo.AssertExpectations(t)
	expenseShareRepo.AssertExpectations(t)
}

// Тест 4: Успешное получение отчета о балансе с пустым списком пользователей
func TestGetEventBalanceReport_EmptyUserBalances(t *testing.T) {
	// Подготовка
	service, expenseRepo, expenseShareRepo := setupService()
	ctx := context.Background()
	eventID := 1

	// Настраиваем моки
	totalAmount := 0.0
	expenseRepo.On("GetEventTotalExpenses", ctx, eventID).Return(totalAmount, nil)

	// Пустой список балансов пользователей
	userBalances := []domain.UserBalance{}
	expenseShareRepo.On("GetEventBalanceReport", ctx, eventID).Return(userBalances, nil)

	// Действие
	report, err := service.GetEventBalanceReport(ctx, eventID)

	// Проверка
	assert.NoError(t, err)
	assert.Equal(t, eventID, report.EventID)
	assert.Equal(t, totalAmount, report.TotalAmount)
	assert.Empty(t, report.UserBalances)

	expenseRepo.AssertExpectations(t)
	expenseShareRepo.AssertExpectations(t)
}
