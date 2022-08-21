package add

import (
	"github.com/stretchr/testify/mock"
)

type repositoryMock struct {
	mock.Mock
}

func NewRepositoryMock() *repositoryMock {
	return &repositoryMock{}
}

func (r *repositoryMock) AddExpense(expense expense) (*createdExpense, error) {
	args := r.Called(expense)

	savedExpense := args.Get(0)
	if savedExpense == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(*createdExpense), nil
	}
}

func (r *repositoryMock) GetExpenseTypeById(id uint64) (*expenseType, error) {
	args := r.Called(id)

	savedExpenseType := args.Get(0)
	err := args.Get(1)

	if savedExpenseType == nil && err == nil {
		return nil, nil
	}

	if savedExpenseType == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(*expenseType), nil
	}
}
