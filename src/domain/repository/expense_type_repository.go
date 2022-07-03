package repository

import "finfit-backend/src/domain/entities"

type ExpenseTypeRepository interface {
	FindById(id int64) (*entities.ExpenseType, error)
	Search() ([]*entities.ExpenseType, error)
	Save(entity *entities.ExpenseType) (*entities.ExpenseType, error)
	DeleteById(id int64) (*entities.ExpenseType, error)
	Update(entity *entities.ExpenseType) (*entities.ExpenseType, error)
}
