package repository_mock

import (
	"finfit/finfit-backend/src/domain/entities"
	"github.com/stretchr/testify/mock"
)

type expenseRepositoryMock struct {
	mock.Mock
}

func NewExpenseRepositoryMock() *expenseRepositoryMock {
	return &expenseRepositoryMock{}
}

func (e expenseRepositoryMock) Search() ([]*entities.Expense, error) {
	//TODO implement me
	panic("implement me")
}

func (e *expenseRepositoryMock) Save(entity entities.Expense) (*entities.Expense, error) {
	args := e.Called(entity)

	savedExpense := args.Get(0)
	if savedExpense == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(*entities.Expense), nil
	}
}

func (e expenseRepositoryMock) DeleteById(id int64) (*entities.Expense, error) {
	//TODO implement me
	panic("implement me")
}

func (e expenseRepositoryMock) Update(entity entities.Expense) (*entities.Expense, error) {
	//TODO implement me
	panic("implement me")
}

func (e expenseRepositoryMock) FindById(id int64) (*entities.Expense, error) {
	//TODO implement me
	panic("implement me")
}
