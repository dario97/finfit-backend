package registry

import (
	repository2 "finfit-backend/src/domain/repository"
	"finfit-backend/src/infrastructure/repository"
	"github.com/jinzhu/gorm"
)

type RepositoryRegistry struct {
	db *gorm.DB
}

func NewRepositoryRegistry(db *gorm.DB) RepositoryRegistry {
	return RepositoryRegistry{
		db: db,
	}
}

func (receiver RepositoryRegistry) GetExpenseRepository() repository2.ExpenseRepository {
	return repository.NewExpenseRepository(receiver.db)
}

func (receiver RepositoryRegistry) GetExpenseTypeRepository() repository2.ExpenseTypeRepository {
	return repository.NewExpenseTypeRepository(receiver.db)
}
