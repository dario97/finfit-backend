package add

import (
	"github.com/stretchr/testify/mock"
)

type serviceMock struct {
	mock.Mock
}

func NewServiceMock() *serviceMock {
	return &serviceMock{}
}

func (s *serviceMock) Add(command command) (*createdExpense, error) {
	args := s.Called(command)

	var expenseToReturn *createdExpense
	if args.Get(0) != nil {
		expenseToReturn = args.Get(0).(*createdExpense)
	}

	return expenseToReturn, args.Error(1)
}
