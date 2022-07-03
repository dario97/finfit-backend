package repository_mock

import (
	"finfit-backend/src/domain/entities"
	"finfit-backend/src/domain/use_cases/service"
	"github.com/stretchr/testify/mock"
)

type expenseServiceMock struct {
	mock.Mock
}

func NewExpenseServiceMock() *expenseServiceMock {
	return &expenseServiceMock{}
}

func (e *expenseServiceMock) GetById(id int64) (*entities.Expense, error) {
	//TODO implement me
	panic("implement me")
}

func (e *expenseServiceMock) Search() (*entities.Expense, error) {
	//TODO implement me
	panic("implement me")
}

func (e *expenseServiceMock) Create(createExpenseCommand service.CreateExpenseCommand) (*entities.Expense, error) {
	args := e.Called(createExpenseCommand)

	var expenseToReturn *entities.Expense
	if args.Get(0) != nil {
		expenseToReturn = args.Get(0).(*entities.Expense)
	}

	return expenseToReturn, args.Error(1)
}

func (e *expenseServiceMock) DeleteById(id int64) (*entities.Expense, error) {
	//TODO implement me
	panic("implement me")
}

func (e *expenseServiceMock) Update(entity entities.Expense) (*entities.Expense, error) {
	//TODO implement me
	panic("implement me")
}
