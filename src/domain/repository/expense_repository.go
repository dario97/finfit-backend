package repository

import (
	"finfit-backend/src/domain/entities"
)

type ExpenseRepository interface {
	FindById(id int64) (*entities.Expense, error)
	Search() ([]*entities.Expense, error)
	Save(entity entities.Expense) (*entities.Expense, error)
	DeleteById(id int64) (*entities.Expense, error)
	Update(entity entities.Expense) (*entities.Expense, error)
}
