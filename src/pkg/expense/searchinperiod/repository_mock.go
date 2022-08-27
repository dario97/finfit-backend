package searchinperiod

import (
	"github.com/stretchr/testify/mock"
	"time"
)

type repositoryMock struct {
	mock.Mock
}

func newRepositoryMock() *repositoryMock {
	return &repositoryMock{}
}

func (r *repositoryMock) SearchExpensesInPeriod(startDate time.Time, endDate time.Time) ([]*expense, error) {
	args := r.Called(startDate, endDate)

	savedExpenses := args.Get(0)
	err := args.Get(1)

	if savedExpenses == nil && err == nil {
		return nil, nil
	}

	if savedExpenses == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).([]*expense), nil
	}

}
