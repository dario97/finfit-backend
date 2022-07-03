package repository

import (
	"finfit/finfit-backend/domain/entities"
	"finfit/finfit-backend/domain/repository"
	"github.com/jinzhu/gorm"
)

type mySqlExpenseTypeRepository struct {
	db *gorm.DB
}

func NewExpenseTypeRepository(db *gorm.DB) repository.ExpenseTypeRepository {
	return &mySqlExpenseTypeRepository{db}
}

func (e mySqlExpenseTypeRepository) FindById(id int64) (*entities.ExpenseType, error) {
	var expenseType entities.ExpenseType
	result := e.db.Find(expenseType)

	if err := result.Error; err != nil {
		return nil, err
	}

	return &expenseType, nil
}

func (e mySqlExpenseTypeRepository) Save(expenseType *entities.ExpenseType) (*entities.ExpenseType, error) {
	result := e.db.Create(&expenseType)

	if err := result.Error; err != nil {
		return nil, err
	}

	return expenseType, nil
}

func (e mySqlExpenseTypeRepository) Search() ([]*entities.ExpenseType, error) {
	//TODO implement me
	panic("implement me")
}

func (e mySqlExpenseTypeRepository) DeleteById(id int64) (*entities.ExpenseType, error) {
	//TODO implement me
	panic("implement me")
}

func (e mySqlExpenseTypeRepository) Update(expenseType *entities.ExpenseType) (*entities.ExpenseType, error) {
	panic("ho")
}
