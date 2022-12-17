package expense

import (
	"finfit-backend/internal/domain/models"
	"github.com/jinzhu/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r repository) Add(expense *models.Expense) (*models.Expense, error) {
	expenseDbModel := r.mapExpenseDBModelFromExpense(expense)
	result := r.db.Create(&expenseDbModel)

	if err := result.Error; err != nil {
		return nil, err
	}

	return expense, nil
}

func (r repository) mapExpenseDBModelFromExpense(expenseToAdd *models.Expense) dbModel {
	return dbModel{
		ID:            expenseToAdd.Id.String(),
		Amount:        expenseToAdd.Amount,
		ExpenseDate:   expenseToAdd.ExpenseDate,
		Description:   expenseToAdd.Description,
		ExpenseTypeID: expenseToAdd.ExpenseType.Id.String(),
	}
}
