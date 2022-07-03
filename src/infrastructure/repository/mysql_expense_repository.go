package repository

import (
	"finfit-backend/src/domain/entities"
	"finfit-backend/src/domain/repository"
	"finfit-backend/src/infrastructure/repository/db_model"
	"github.com/jinzhu/gorm"
)

type mySqlExpenseRepository struct {
	db *gorm.DB
}

func NewExpenseRepository(db *gorm.DB) repository.ExpenseRepository {
	return &mySqlExpenseRepository{db}
}

func (e mySqlExpenseRepository) FindById(id int64) (*entities.Expense, error) {
	//TODO implement me
	panic("implement me")
}

func (e mySqlExpenseRepository) Search() ([]*entities.Expense, error) {
	//TODO implement me
	panic("implement me")
}

func (e mySqlExpenseRepository) Save(expenseToSave entities.Expense) (*entities.Expense, error) {
	expenseDbModel := db_model.NewExpenseDbModelFromExpense(expenseToSave)
	result := e.db.Create(&expenseDbModel)

	if err := result.Error; err != nil {
		return nil, err
	}

	createdExpense := entities.NewExpenseWithId(expenseDbModel.ID,
		expenseDbModel.Amount,
		expenseDbModel.ExpenseDate,
		expenseDbModel.Description,
		expenseDbModel.ExpenseType)

	return &createdExpense, nil
}

func (e mySqlExpenseRepository) DeleteById(id int64) (*entities.Expense, error) {
	//TODO implement me
	panic("implement me")
}

func (e mySqlExpenseRepository) Update(entity entities.Expense) (*entities.Expense, error) {
	//TODO implement me
	panic("implement me")
}
