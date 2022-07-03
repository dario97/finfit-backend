package repository_mock

import (
	"finfit/finfit-backend/domain/entities"
	"github.com/stretchr/testify/mock"
)

type expenseTypeServiceMock struct {
	mock.Mock
}

func NewExpenseTypeServiceMock() *expenseTypeServiceMock {
	return &expenseTypeServiceMock{}
}

func (e *expenseTypeServiceMock) GetById(id int64) (*entities.ExpenseType, error) {
	args := e.Called(id)

	savedExpenseType := args.Get(0)
	err := args.Get(1)

	if savedExpenseType == nil && err == nil {
		return nil, nil
	}

	if savedExpenseType == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(*entities.ExpenseType), nil
	}
}
