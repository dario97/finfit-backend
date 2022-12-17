package expense

import (
	"finfit-backend/internal/domain/models"
	"finfit-backend/internal/domain/services/expensetype"
)

const invalidExpenseTypeErrorMsg = "the expense type doesn't exists"

type Repository interface {
	Add(entity *models.Expense) (*models.Expense, error)
}
type Service interface {
	Add(command *AddCommand) (*models.Expense, error)
}

type service struct {
	repository         Repository
	expenseTypeService expensetype.Service
}

func NewService(expenseRepository Repository, expenseTypeService expensetype.Service) *service {
	return &service{repository: expenseRepository, expenseTypeService: expenseTypeService}
}

func (s service) Add(command *AddCommand) (*models.Expense, error) {
	expenseType, expenseTypeServiceError := s.expenseTypeService.GetById(command.expenseType.Id)

	if expenseTypeServiceError != nil {
		return nil, UnexpectedError{Msg: expenseTypeServiceError.Error()}
	}

	if expenseType == nil {
		return nil, InvalidExpenseTypeError{Msg: invalidExpenseTypeErrorMsg}
	}

	expenseToCreate := models.NewExpense(command.amount, command.expenseDate, command.description, command.expenseType)

	createdExpense, repoError := s.repository.Add(expenseToCreate)

	if repoError != nil {
		return nil, UnexpectedError{Msg: repoError.Error()}
	}

	return createdExpense, nil
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
