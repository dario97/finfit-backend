package expense

import (
	"finfit-backend/internal/domain/models"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func NewRepositoryMock() *RepositoryMock {
	return &RepositoryMock{}
}

func (r *RepositoryMock) Add(expense *models.Expense) (*models.Expense, error) {
	args := r.Called(expense)

	savedExpense := args.Get(0)
	if savedExpense == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(*models.Expense), nil
	}
}

func (r *RepositoryMock) MockAdd(callArguments, returnArguments []interface{}, times int) {
	r.On("Add", callArguments...).Return(returnArguments...).Times(times)
}
