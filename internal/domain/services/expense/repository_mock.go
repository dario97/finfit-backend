package expense

import (
	"finfit-backend/internal/domain/models"
	"github.com/stretchr/testify/mock"
	"time"
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
	err := args.Error(1)
	if err == nil && savedExpense == nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return args.Get(0).(*models.Expense), nil
	}
}

func (r *RepositoryMock) SearchInPeriod(startDate time.Time, endDate time.Time) ([]*models.Expense, error) {
	args := r.Called(startDate, endDate)

	err := args.Error(1)
	if err != nil {
		return nil, err
	} else {
		return args.Get(0).([]*models.Expense), nil
	}
}

func (r *RepositoryMock) MockAdd(callArguments, returnArguments []interface{}, times int) {
	r.On("Add", callArguments...).Return(returnArguments...).Times(times)
}

func (r *RepositoryMock) MockSearchInPeriod(callArguments, returnArguments []interface{}, times int) {
	r.On("SearchInPeriod", callArguments...).Return(returnArguments...).Times(times)
}
