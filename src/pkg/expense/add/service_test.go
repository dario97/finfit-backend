package add

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGivenAnExpense_WhenAdd_ThenReturnCreatedExpense(t *testing.T) {
	repoMock := NewRepositoryMock()
	expenseService := NewService(repoMock)

	expenseToCreate := getExpense()

	expectedCreatedExpense := createdExpense{
		id:          1,
		amount:      expenseToCreate.amount,
		expenseDate: expenseToCreate.expenseDate,
		description: expenseToCreate.description,
		expenseType: expenseToCreate.expenseType,
	}

	repoMock.On("AddExpense", expenseToCreate).Return(&expectedCreatedExpense, nil)
	repoMock.On("GetExpenseTypeById", expenseToCreate.expenseType.id).Return(&expenseToCreate.expenseType, nil)

	actualCreatedExpense, err := expenseService.Add(buildCommandFromExpense(expenseToCreate))

	assert.Nil(t, err, "Error must to be nil")
	assertEqualsExpense(t, expectedCreatedExpense, *actualCreatedExpense)
}

func TestGivenThatFailToGetExpenseType_WhenAdd_ThenReturnError(t *testing.T) {
	repoMock := NewRepositoryMock()
	expenseService := NewService(repoMock)

	expenseToCreate := getExpense()

	repoError := errors.New("fail to ge expenseDBModel type")
	expectedError := UnexpectedError{
		Msg: repoError.Error(),
	}

	repoMock.On("GetExpenseTypeById", expenseToCreate.expenseType.id).Return(nil, repoError)

	actualCreatedExpense, err := expenseService.Add(buildCommandFromExpense(expenseToCreate))

	assert.Nil(t, actualCreatedExpense)
	assert.NotNil(t, err, "Error must not be nil")
	assert.Equal(t, expectedError, err)
}

func TestGivenThatExpenseTypeNotExists_WhenAdd_ThenReturnError(t *testing.T) {
	repoMock := NewRepositoryMock()
	expenseService := NewService(repoMock)

	expenseToCreate := getExpense()

	expectedError := InvalidExpenseTypeError{
		Msg: "the expenseDBModel type doesn't exists",
	}

	repoMock.On("GetExpenseTypeById", expenseToCreate.expenseType.id).Return(nil, nil)

	actualCreatedExpense, err := expenseService.Add(buildCommandFromExpense(expenseToCreate))

	assert.Nil(t, actualCreatedExpense)
	assert.NotNil(t, err, "Error must not be nil")
	assert.Equal(t, expectedError, err)
}

func TestGivenThatSaveExpenseIntoDatabaseFails_WhenAdd_ThenReturnError(t *testing.T) {
	repoMock := NewRepositoryMock()
	expenseService := NewService(repoMock)

	expenseToCreate := getExpense()

	repoError := errors.New("fail to save expenseDBModel")
	expectedError := UnexpectedError{
		Msg: repoError.Error(),
	}

	repoMock.On("GetExpenseTypeById", expenseToCreate.expenseType.id).Return(&expenseToCreate.expenseType, nil)
	repoMock.On("AddExpense", expenseToCreate).Return(nil, repoError)

	actualCreatedExpense, err := expenseService.Add(buildCommandFromExpense(expenseToCreate))

	assert.Nil(t, actualCreatedExpense)
	assert.NotNil(t, err, "Error must not be nil")
	assert.Equal(t, expectedError, err)
}

func getExpense() expense {
	return expense{
		amount:      100.50,
		expenseDate: time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC),
		description: "Lomitos",
		expenseType: getExpenseType(),
	}

}

func getExpenseType() expenseType {
	return expenseType{
		id:   1,
		name: "Delivery",
	}
}

func assertEqualsExpense(t *testing.T, expected createdExpense, actual createdExpense) {
	assert.Equal(t, expected.id, actual.id, "id are not equals")
	assert.Equal(t, expected.expenseType, actual.expenseType, "expenseTypes are not equals")
	assert.Equal(t, expected.amount, actual.amount, "amounts are not equals")
	assert.Equalf(t, expected.expenseDate, actual.expenseDate, "expenseDates are not equals")
	assert.Equalf(t, expected.description, actual.description, "descriptions are not equals")
}

func buildCommandFromExpense(expense expense) command {
	return command{
		amount:      expense.amount,
		expenseDate: expense.expenseDate,
		description: expense.description,
		expenseType: expense.expenseType,
	}
}
