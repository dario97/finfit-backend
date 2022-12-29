package expense

import (
	"finfit-backend/internal/domain/models"
	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

func NewServiceMock() *ServiceMock {
	return &ServiceMock{}
}

func (s *ServiceMock) Add(command *AddCommand) (*models.Expense, error) {
	args := s.Called(command)

	err := args.Error(1)
	expenseToReturn := args.Get(0)
	if err == nil && expenseToReturn == nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return expenseToReturn.(*models.Expense), nil
	}
}

func (s *ServiceMock) SearchInPeriod(command *SearchInPeriodCommand) ([]*models.Expense, error) {
	args := s.Called(command)

	err := args.Error(1)
	expenses := args.Get(0)
	if err == nil && expenses == nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return expenses.([]*models.Expense), nil
	}
}

func (s *ServiceMock) MockAdd(callArguments, returnArguments []interface{}, times int) {
	s.On("Add", callArguments...).Return(returnArguments...).Times(times)
}

func (s *ServiceMock) MockSearchInPeriod(callArguments, returnArguments []interface{}, times int) {
	s.On("SearchInPeriod", callArguments...).Return(returnArguments...).Times(times)
}
