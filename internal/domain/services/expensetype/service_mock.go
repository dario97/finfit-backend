package expensetype

import (
	"finfit-backend/internal/domain/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

func NewServiceMock() *ServiceMock {
	return &ServiceMock{}
}

func (s *ServiceMock) GetById(id uuid.UUID) (*models.ExpenseType, error) {
	args := s.Called(id)

	err := args.Error(1)
	expenseType := args.Get(0)
	if err == nil && expenseType == nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return args.Get(0).(*models.ExpenseType), nil
	}
}

func (s *ServiceMock) MockGetByID(callArguments, returnArguments []interface{}, times int) {
	s.On("GetById", callArguments...).Return(returnArguments...).Times(times)
}
