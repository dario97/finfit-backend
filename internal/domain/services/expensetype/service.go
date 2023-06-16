package expensetype

import (
	"finfit-backend/internal/domain/models"
	"github.com/google/uuid"
)

type Repository interface {
	GetByID(id uuid.UUID) (*models.ExpenseType, error)
	GetByName(name string) (*models.ExpenseType, error)
	GetAll() ([]*models.ExpenseType, error)
	Add(expense *models.ExpenseType) (*models.ExpenseType, error)
}
type Service interface {
	GetById(id uuid.UUID) (*models.ExpenseType, error)
	Add(command *AddCommand) (*models.ExpenseType, error)
	GetAll() ([]*models.ExpenseType, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s service) GetById(id uuid.UUID) (*models.ExpenseType, error) {
	expenseType, err := s.repo.GetByID(id)
	if err != nil {
		return nil, UnexpectedError{Msg: err.Error()}
	}

	return expenseType, nil
}

func (s service) Add(command *AddCommand) (*models.ExpenseType, error) {
	storedExpenseType, err := s.repo.GetByName(command.name)

	if err != nil {
		return nil, UnexpectedError{Msg: err.Error()}
	}

	if storedExpenseType != nil {
		return storedExpenseType, nil
	}

	expenseTypeToAdd, err := mapExpenseTypeFromAddCommand(command)
	if err != nil {
		return nil, InvalidDomainModelError{Msg: err.Error()}
	}

	addedExpenseType, err := s.repo.Add(expenseTypeToAdd)
	if err != nil {
		return nil, UnexpectedError{Msg: err.Error()}
	}

	return addedExpenseType, nil
}

func (s service) GetAll() ([]*models.ExpenseType, error) {
	expenseTypes, err := s.repo.GetAll()
	if err != nil {
		return nil, UnexpectedError{Msg: err.Error()}
	}

	return expenseTypes, nil
}

func mapExpenseTypeFromAddCommand(command *AddCommand) (*models.ExpenseType, error) {
	return models.NewExpenseType(command.name)
}

type UnexpectedError struct {
	Msg string
}

func (receiver UnexpectedError) Error() string {
	return receiver.Msg
}

type InvalidDomainModelError struct {
	Msg string
}

func (receiver InvalidDomainModelError) Error() string {
	return receiver.Msg
}
