package expense

import (
	"finfit-backend/internal/domain/models"
	"finfit-backend/internal/domain/services/expensetype"
	"time"
)

const invalidExpenseTypeErrorMsg = "the expense type doesn't exists"

type Repository interface {
	Add(entity *models.Expense) (*models.Expense, error)
	SearchInPeriod(startDate time.Time, endDate time.Time) ([]*models.Expense, error)
}

type Service interface {
	Add(command *AddCommand) (*models.Expense, error)
	SearchInPeriod(command *SearchInPeriodCommand) ([]*models.Expense, error)
}

type service struct {
	repository         Repository
	expenseTypeService expensetype.Service
}

func NewService(expenseRepository Repository, expenseTypeService expensetype.Service) *service {
	return &service{repository: expenseRepository, expenseTypeService: expenseTypeService}
}

func (s service) Add(command *AddCommand) (*models.Expense, error) {
	expenseType, expenseTypeServiceError := s.expenseTypeService.GetById(command.expenseTypeId)

	if expenseTypeServiceError != nil {
		return nil, UnexpectedError{Msg: expenseTypeServiceError.Error()}
	}

	if expenseType == nil {
		return nil, InvalidExpenseTypeError{Msg: invalidExpenseTypeErrorMsg}
	}

	expenseToCreate := models.NewExpense(command.amount, command.expenseDate, command.description, expenseType)

	createdExpense, repoError := s.repository.Add(expenseToCreate)

	if repoError != nil {
		return nil, UnexpectedError{Msg: repoError.Error()}
	}

	return createdExpense, nil
}

func (s service) SearchInPeriod(command *SearchInPeriodCommand) ([]*models.Expense, error) {
	expenses, err := s.repository.SearchInPeriod(command.startDate, command.endDate)
	if err != nil {
		return nil, UnexpectedError{Msg: err.Error()}
	}
	return expenses, nil
}

type UnexpectedError struct {
	Msg string
}

func (receiver UnexpectedError) Error() string {
	return receiver.Msg
}

type InvalidExpenseTypeError struct {
	Msg string
}

func (receiver InvalidExpenseTypeError) Error() string {
	return receiver.Msg
}
