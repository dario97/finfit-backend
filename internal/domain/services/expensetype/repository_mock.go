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

func (r *RepositoryMock) GetByName(name string) (*models.ExpenseType, error) {
	args := r.Called(name)

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

func (r *RepositoryMock) GetAll() ([]*models.ExpenseType, error) {
	args := r.Called()

	err := args.Error(1)
	expenseType := args.Get(0)
	if err == nil && expenseType == nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return args.Get(0).([]*models.ExpenseType), nil
	}
}

func (r *RepositoryMock) Add(expenseType *models.ExpenseType) (*models.ExpenseType, error) {
	args := r.Called(expenseType)

	savedExpense := args.Get(0)
	err := args.Error(1)
	if err == nil && savedExpense == nil {
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

func (r *RepositoryMock) MockGetByName(callArguments, returnArguments []interface{}, times int) {
	r.On("GetByName", callArguments...).Return(returnArguments...).Times(times)
}

func (r *RepositoryMock) MockAdd(callArguments, returnArguments []interface{}, times int) {
	r.On("Add", callArguments...).Return(returnArguments...).Times(times)
}

func (r *RepositoryMock) MockGetAll(callArguments, returnArguments []interface{}, times int) {
	r.On("GetAll", callArguments...).Return(returnArguments...).Times(times)
}
