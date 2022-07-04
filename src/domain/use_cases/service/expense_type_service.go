package service

import (
	"finfit-backend/src/domain/entities"
	"finfit-backend/src/domain/repository"
	"github.com/go-playground/validator/v10"
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
