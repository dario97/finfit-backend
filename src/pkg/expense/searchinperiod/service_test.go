package searchinperiod

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGivenAnStartDateAndEndDate_whenSearchInPeriod_thenReturnExpensesInPeriod(t *testing.T) {
	repositoryMock := newRepositoryMock()
	service := NewService(repositoryMock)
	startDate := time.Date(2022, 7, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2022, 9, 1, 0, 0, 0, 0, time.UTC)
	command, _ := newCommand(startDate, endDate)
	expectedExpenses := []*expense{
		getExpense(),
		getAnotherExpense(),
	}
	repositoryMock.On("SearchExpensesInPeriod", startDate, endDate).Return(expectedExpenses, nil)

	actualExpenses, err := service.SearchInPeriod(command)

	assert.Nil(t, err)
	assert.Equal(t, expectedExpenses, actualExpenses)
}

func Test02(t *testing.T) {
	repositoryMock := newRepositoryMock()
	service := NewService(repositoryMock)
	startDate := time.Date(2022, 7, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2022, 9, 1, 0, 0, 0, 0, time.UTC)
	command, _ := newCommand(startDate, endDate)
	repoError := errors.New("fail to ge expenseDBModel type")
	expectedError := unexpectedError{msg: repoError.Error()}
	repositoryMock.On("SearchExpensesInPeriod", startDate, endDate).Return(nil, repoError)

	actualExpenses, err := service.SearchInPeriod(command)

	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, actualExpenses)
}

func getExpense() *expense {
	return &expense{
		id:          1,
		amount:      10.3,
		expenseDate: time.Date(2022, 8, 4, 0, 0, 0, 0, time.UTC),
		description: "playstation 5",
		expenseType: expenseType{
			id:   1,
			name: "Tecnologia",
		},
	}
}

func getAnotherExpense() *expense {
	return &expense{
		id:          2,
		amount:      101.2,
		expenseDate: time.Date(2022, 7, 2, 0, 0, 0, 0, time.UTC),
		description: "playstation 1",
		expenseType: expenseType{
			id:   1,
			name: "Antiguedades",
		},
	}
}
