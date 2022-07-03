package test

import (
	"errors"
	"finfit/finfit-backend/domain/entities"
	"finfit/finfit-backend/domain/use_cases/custom_errors"
	expense2 "finfit/finfit-backend/domain/use_cases/service"
	"finfit/finfit-backend/test/mock/repository_mock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGivenAExpenseToCreate_WhenCreate_ThenReturnCreatedExpense(t *testing.T) {
	expenseRepositoryMock := repository_mock.NewExpenseRepositoryMock()
	expenseTypeServiceMock := repository_mock.NewExpenseTypeServiceMock()
	expenseService := expense2.NewExpenseService(expenseRepositoryMock, expenseTypeServiceMock)

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

	actualCreatedExpense, err := expenseService.Create(getCreateExpenseCommandFromExpense(expenseToCreate))

	assert.Nil(t, err, "Error must to be nil")
	assertEqualsExpense(t, expectedCreatedExpense, *actualCreatedExpense)
}

func TestGivenThatExpenseTypeServiceFails_WhenCreate_ThenReturnInternalError(t *testing.T) {
	expenseRepositoryMock := repository_mock.NewExpenseRepositoryMock()
	expenseTypeServiceMock := repository_mock.NewExpenseTypeServiceMock()
	expenseService := expense2.NewExpenseService(expenseRepositoryMock, expenseTypeServiceMock)

	expenseType := entities.NewExpenseTypeWithId(1, "Delivery")
	expenseToCreate := entities.NewExpense(100.50,
		time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC),
		"Lomitos",
		expenseType)

	expenseTypeServiceError := custom_errors.InternalError{Msg: "fail to get expense type"}
	expectedError := custom_errors.InternalError{
		Msg: expenseTypeServiceError.Error(),
	}
	expenseTypeServiceMock.On("GetById", expenseType.Id()).Return(nil, expenseTypeServiceError)

	actualCreatedExpense, err := expenseService.Create(getCreateExpenseCommandFromExpense(expenseToCreate))

	assert.Nil(t, actualCreatedExpense)
	assert.NotNil(t, err, "Error must not be nil")
	assert.Equal(t, expectedError, err)
}

func TestGivenThatExpenseTypeNotExists_WhenCreate_ThenReturnInvalidExpenseTypeError(t *testing.T) {
	expenseRepositoryMock := repository_mock.NewExpenseRepositoryMock()
	expenseTypeServiceMock := repository_mock.NewExpenseTypeServiceMock()
	expenseService := expense2.NewExpenseService(expenseRepositoryMock, expenseTypeServiceMock)

	expenseType := entities.NewExpenseTypeWithId(1, "Delivery")
	expenseToCreate := entities.NewExpense(100.50,
		time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC),
		"Lomitos",
		expenseType)

	expectedError := custom_errors.InvalidExpenseTypeError{
		Msg: "the expense type doesn't exists",
	}

	expenseTypeServiceMock.On("GetById", expenseType.Id()).Return(nil, nil)

	actualCreatedExpense, err := expenseService.Create(getCreateExpenseCommandFromExpense(expenseToCreate))

	assert.Nil(t, actualCreatedExpense)
	assert.NotNil(t, err, "Error must not be nil")
	assert.Equal(t, expectedError, err)
}

func TestGivenThatExpenseRepositoryFails_WhenCreate_ThenReturnInternalError(t *testing.T) {
	expenseRepositoryMock := repository_mock.NewExpenseRepositoryMock()
	expenseTypeServiceMock := repository_mock.NewExpenseTypeServiceMock()
	expenseService := expense2.NewExpenseService(expenseRepositoryMock, expenseTypeServiceMock)

	expenseType := entities.NewExpenseTypeWithId(1, "Delivery")
	expenseToCreate := entities.NewExpense(100.50,
		time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC),
		"Lomitos",
		expenseType)

	expenseRepositoryError := errors.New("fail to save expense")
	expectedError := custom_errors.InternalError{
		Msg: expenseRepositoryError.Error(),
	}

	expenseTypeServiceMock.On("GetById", expenseType.Id()).Return(&expenseType, nil)
	expenseRepositoryMock.On("Save", expenseToCreate).Return(nil, expenseRepositoryError)

	actualCreatedExpense, err := expenseService.Create(getCreateExpenseCommandFromExpense(expenseToCreate))

	assert.Nil(t, actualCreatedExpense)
	assert.NotNil(t, err, "Error must not be nil")
	assert.Equal(t, expectedError, err)
}

func assertEqualsExpense(t *testing.T, expected entities.Expense, actual entities.Expense) {
	assert.Equal(t, expected.Id(), actual.Id(), "id are not equals")
	assert.Equal(t, expected.ExpenseType(), actual.ExpenseType(), "expenseTypes are not equals")
	assert.Equal(t, expected.Amount(), actual.Amount(), "amounts are not equals")
	assert.Equalf(t, expected.ExpenseDate(), actual.ExpenseDate(), "expenseDates are not equals")
	assert.Equalf(t, expected.Description(), actual.Description(), "descriptions are not equals")
}

func getCreateExpenseCommandFromExpense(expense entities.Expense) expense2.CreateExpenseCommand {
	return expense2.NewCreateExpenseCommand(expense.Amount(), expense.ExpenseDate(), expense.Description(), expense.ExpenseType())
}
