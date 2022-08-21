package service

import (
	"finfit-backend/src/domain/entities"
	"finfit-backend/src/domain/repository"
	"finfit-backend/src/domain/use_cases/custom_errors"
	"time"
)

const invalidExpenseTypeErrorMsg = "the expense type doesn't exists"

type CreateExpenseCommand struct {
	amount      float64
	expenseDate time.Time
	description string
	expenseType entities.ExpenseType
}

func NewCreateExpenseCommand(amount float64, expenseDate time.Time, description string, expenseType entities.ExpenseType) CreateExpenseCommand {
	return CreateExpenseCommand{
		amount:      amount,
		expenseDate: expenseDate,
		description: description,
		expenseType: expenseType,
	}
}

type SearchInPeriodCommand struct {
	startDate time.Time
	endDate   time.Time
}

func NewSearchInPeriodCommand(startDate time.Time, endDate time.Time) SearchInPeriodCommand {
	return SearchInPeriodCommand{
		startDate: startDate,
		endDate:   endDate,
	}
}

type ExpenseService interface {
	GetById(id int64) (*entities.Expense, error)
	SearchInPeriod(searchInPeriodCommand SearchInPeriodCommand) ([]*entities.Expense, error)
	Create(createExpenseCommand CreateExpenseCommand) (*entities.Expense, error)
	DeleteById(id int64) (*entities.Expense, error)
	Update(entity entities.Expense) (*entities.Expense, error)
}

type expenseService struct {
	expenseRepository  repository.ExpenseRepository
	expenseTypeService ExpenseTypeService
}

func NewExpenseService(repository repository.ExpenseRepository, expenseTypeService ExpenseTypeService) ExpenseService {
	return &expenseService{expenseRepository: repository, expenseTypeService: expenseTypeService}
}

func (e expenseService) GetById(id int64) (*entities.Expense, error) {
	//TODO implement me
	panic("implement me")
}

func (e expenseService) SearchInPeriod(command SearchInPeriodCommand) ([]*entities.Expense, error) {
	//TODO implement me
	panic("implement me")
}

func (e expenseService) Create(createExpenseCommand CreateExpenseCommand) (*entities.Expense, error) {
	expenseType, expenseTypServiceError := e.expenseTypeService.GetById(createExpenseCommand.expenseType.Id())

	if expenseTypServiceError != nil {
		return nil, custom_errors.UnexpectedError{Msg: expenseTypServiceError.Error()}
	}

	if expenseType == nil {
		return nil, custom_errors.InvalidExpenseTypeError{Msg: invalidExpenseTypeErrorMsg}
	}

	expenseToCreate := entities.NewExpense(createExpenseCommand.amount,
		createExpenseCommand.expenseDate,
		createExpenseCommand.description,
		createExpenseCommand.expenseType)

	createdExpense, expenseRepositoryError := e.expenseRepository.Save(expenseToCreate)

	if expenseRepositoryError != nil {
		return nil, custom_errors.UnexpectedError{Msg: expenseRepositoryError.Error()}
	}

	return createdExpense, nil
}

func (e expenseService) DeleteById(id int64) (*entities.Expense, error) {
	//TODO implement me
	panic("implement me")
}

func (e expenseService) Update(entity entities.Expense) (*entities.Expense, error) {
	//TODO implement me
	panic("implement me")
}
