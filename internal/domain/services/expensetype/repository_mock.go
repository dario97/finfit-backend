package expensetype

import (
	"finfit-backend/internal/domain/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func NewRepositoryMock() *RepositoryMock {
	return &RepositoryMock{}
}

func (r *RepositoryMock) GetByID(id uuid.UUID) (*models.ExpenseType, error) {
	args := r.Called(id)

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

func (r *RepositoryMock) MockGetByID(callArguments, returnArguments []interface{}, times int) {
	r.On("GetByID", callArguments...).Return(returnArguments...).Times(times)
}
