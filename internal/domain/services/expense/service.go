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
	AddAll(expenses []*models.Expense) ([]*models.Expense, error)
}

type Service interface {
	Add(command *AddCommand) (*models.Expense, error)
	AddAll(command *AddAllCommand) ([]*models.Expense, error)
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
	expenseType, expenseTypeServiceError := s.checkIfExpenseTypeExists(command)

	if expenseTypeServiceError != nil {
		return nil, UnexpectedError{Msg: expenseTypeServiceError.Error()}
	}

	if expenseType == nil {
		return nil, InvalidExpenseTypeError{Msg: invalidExpenseTypeErrorMsg}
	}

	expenseToCreate, err := s.mapAddCommandToExpense(command, expenseType)
	if err != nil {
		return nil, InvalidDomainModelError{Msg: err.Error()}
	}

	createdExpense, repoError := s.repository.Add(expenseToCreate)

	if repoError != nil {
		return nil, UnexpectedError{Msg: repoError.Error()}
	}

	return createdExpense, nil
}

func (s service) checkIfExpenseTypeExists(command *AddCommand) (*models.ExpenseType, error) {
	return s.expenseTypeService.GetById(command.expenseTypeId)
}

func (s service) SearchInPeriod(command *SearchInPeriodCommand) ([]*models.Expense, error) {
	expenses, err := s.repository.SearchInPeriod(command.startDate, command.endDate)
	if err != nil {
		return nil, UnexpectedError{Msg: err.Error()}
	}
	return expenses, nil
}

func (s service) AddAll(command *AddAllCommand) ([]*models.Expense, error) {
	money, _ := models.NewMoney(10.3, "ARS")
	expenseType, _ := models.NewExpenseType("Delivery")
	expense1, _ := models.NewExpense(money, time.Date(2022, 5, 28, 0, 0, 0, 0, time.Local), "Lomitos", expenseType)
	expense2, _ := models.NewExpense(money, time.Date(2022, 7, 28, 0, 0, 0, 0, time.Local), "Lomitos", expenseType)

	return []*models.Expense{expense1, expense2}, nil
}

func (s service) mapAddCommandToExpense(command *AddCommand, expenseType *models.ExpenseType) (*models.Expense, error) {
	money, err := models.NewMoney(command.amount, command.currency)
	if err != nil {
		return nil, err
	}
	expenseToCreate, err := models.NewExpense(money, command.expenseDate, command.description, expenseType)
	if err != nil {
		return nil, err
	}

	return expenseToCreate, err
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

type InvalidDomainModelError struct {
	Msg string
}

func (receiver InvalidDomainModelError) Error() string {
	return receiver.Msg
}
