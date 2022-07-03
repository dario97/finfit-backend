package service

import (
	"finfit/finfit-backend/domain/entities"
	"finfit/finfit-backend/domain/repository"
	"github.com/go-playground/validator"
)

type ExpenseTypeService interface {
	GetById(id int64) (*entities.ExpenseType, error)
}

type expenseTypeService struct {
	expenseRepository repository.ExpenseTypeRepository
	validator         *validator.Validate
}

func NewExpenseTypeService(repository repository.ExpenseTypeRepository) ExpenseTypeService {
	return &expenseTypeService{expenseRepository: repository, validator: validator.New()}
}

func (e expenseTypeService) GetById(id int64) (*entities.ExpenseType, error) {
	//TODO implement me
	panic("implement me")
}
