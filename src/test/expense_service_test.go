package test

import (
	"errors"
	"finfit-backend/src/domain/entities"
	"finfit-backend/src/domain/use_cases/custom_errors"
	"finfit-backend/src/domain/use_cases/service"
	"finfit-backend/src/test/mock/repository_mock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGivenAExpenseToCreate_WhenCreate_ThenReturnCreatedExpense(t *testing.T) {
	expenseRepositoryMock := repository_mock.NewExpenseRepositoryMock()
	expenseTypeServiceMock := repository_mock.NewExpenseTypeServiceMock()
	expenseService := service.NewExpenseService(expenseRepositoryMock, expenseTypeServiceMock)

	expenseType := entities.NewExpenseTypeWithId(1, "Delivery")
	expenseToCreate := entities.NewExpense(100.50,
		time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC),
		"Lomitos",
		expenseType)

	expectedCreatedExpense := entities.NewExpenseWithId(1, 100.50,
		time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC),
		"Lomitos",
		entities.NewExpenseTypeWithId(1, "Delivery"))

	expenseRepositoryMock.On("Save", expenseToCreate).Return(&expectedCreatedExpense, nil)
	expenseTypeServiceMock.On("GetById", expenseType.Id()).Return(&expenseType, nil)

	actualCreatedExpense, err := expenseService.Create(buildCreateExpenseCommandFromExpense(expenseToCreate))

	assert.Nil(t, err, "Error must to be nil")
	assertEqualsExpense(t, expectedCreatedExpense, *actualCreatedExpense)
}

func TestGivenThatExpenseTypeServiceFails_WhenCreate_ThenReturnInternalError(t *testing.T) {
	expenseRepositoryMock := repository_mock.NewExpenseRepositoryMock()
	expenseTypeServiceMock := repository_mock.NewExpenseTypeServiceMock()
	expenseService := service.NewExpenseService(expenseRepositoryMock, expenseTypeServiceMock)

	expenseType := entities.NewExpenseTypeWithId(1, "Delivery")
	expenseToCreate := entities.NewExpense(100.50,
		time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC),
		"Lomitos",
		expenseType)

	expenseTypeServiceError := custom_errors.UnexpectedError{Msg: "fail to get expense type"}
	expectedError := custom_errors.UnexpectedError{
		Msg: expenseTypeServiceError.Error(),
	}
	expenseTypeServiceMock.On("GetById", expenseType.Id()).Return(nil, expenseTypeServiceError)

	actualCreatedExpense, err := expenseService.Create(buildCreateExpenseCommandFromExpense(expenseToCreate))

	assert.Nil(t, actualCreatedExpense)
	assert.NotNil(t, err, "Error must not be nil")
	assert.Equal(t, expectedError, err)
}

func TestGivenThatExpenseTypeNotExists_WhenCreate_ThenReturnInvalidExpenseTypeError(t *testing.T) {
	expenseRepositoryMock := repository_mock.NewExpenseRepositoryMock()
	expenseTypeServiceMock := repository_mock.NewExpenseTypeServiceMock()
	expenseService := service.NewExpenseService(expenseRepositoryMock, expenseTypeServiceMock)

	expenseType := entities.NewExpenseTypeWithId(1, "Delivery")
	expenseToCreate := entities.NewExpense(100.50,
		time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC),
		"Lomitos",
		expenseType)

	expectedError := custom_errors.InvalidExpenseTypeError{
		Msg: "the expense type doesn't exists",
	}

	expenseTypeServiceMock.On("GetById", expenseType.Id()).Return(nil, nil)

	actualCreatedExpense, err := expenseService.Create(buildCreateExpenseCommandFromExpense(expenseToCreate))

	assert.Nil(t, actualCreatedExpense)
	assert.NotNil(t, err, "Error must not be nil")
	assert.Equal(t, expectedError, err)
}

func TestGivenThatExpenseRepositoryFails_WhenCreate_ThenReturnInternalError(t *testing.T) {
	expenseRepositoryMock := repository_mock.NewExpenseRepositoryMock()
	expenseTypeServiceMock := repository_mock.NewExpenseTypeServiceMock()
	expenseService := service.NewExpenseService(expenseRepositoryMock, expenseTypeServiceMock)

	expenseType := entities.NewExpenseTypeWithId(1, "Delivery")
	expenseToCreate := entities.NewExpense(100.50,
		time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC),
		"Lomitos",
		expenseType)

	expenseRepositoryError := errors.New("fail to save expense")
	expectedError := custom_errors.UnexpectedError{
		Msg: expenseRepositoryError.Error(),
	}

	expenseTypeServiceMock.On("GetById", expenseType.Id()).Return(&expenseType, nil)
	expenseRepositoryMock.On("Save", expenseToCreate).Return(nil, expenseRepositoryError)

	actualCreatedExpense, err := expenseService.Create(buildCreateExpenseCommandFromExpense(expenseToCreate))

	assert.Nil(t, actualCreatedExpense)
	assert.NotNil(t, err, "Error must not be nil")
	assert.Equal(t, expectedError, err)
}

func TestGivenAStartDateAndEndDate_WhenSearchInPeriod_ThenReturnExpensesInPeriod(t *testing.T) {
	expenseRepositoryMock := repository_mock.NewExpenseRepositoryMock()
	expenseTypeServiceMock := repository_mock.NewExpenseTypeServiceMock()
	expenseService := service.NewExpenseService(expenseRepositoryMock, expenseTypeServiceMock)

	startDate := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2022, 7, 1, 0, 0, 0, 0, time.UTC)

	var expectedExpensesInPeriod []*entities.Expense
	actualExpensesInPeriod, err := expenseService.SearchInPeriod(service.NewSearchInPeriodCommand(startDate, endDate))

	assert.Nil(t, err, "Error must to be nil")

	assert.Equal(t, len(expectedExpensesInPeriod), len(actualExpensesInPeriod))
	for i, expectedExpense := range expectedExpensesInPeriod {
		assertEqualsExpense(t, *expectedExpense, *actualExpensesInPeriod[i])
	}
}

func assertEqualsExpense(t *testing.T, expected entities.Expense, actual entities.Expense) {
	assert.Equal(t, expected.Id(), actual.Id(), "id are not equals")
	assert.Equal(t, expected.ExpenseType(), actual.ExpenseType(), "expenseTypes are not equals")
	assert.Equal(t, expected.Amount(), actual.Amount(), "amounts are not equals")
	assert.Equalf(t, expected.ExpenseDate(), actual.ExpenseDate(), "expenseDates are not equals")
	assert.Equalf(t, expected.Description(), actual.Description(), "descriptions are not equals")
}

func buildCreateExpenseCommandFromExpense(expense entities.Expense) service.CreateExpenseCommand {
	return service.NewCreateExpenseCommand(expense.Amount(), expense.ExpenseDate(), expense.Description(), expense.ExpenseType())
}

func buildSearchInPeriodCommand(startDate time.Time, endDate time.Time) service.SearchInPeriodCommand {
	return service.NewSearchInPeriodCommand(startDate, endDate)
}
