package add

import (
	"github.com/jinzhu/gorm"
	"time"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return repository{db: db}
}

func (r repository) AddExpense(expense expense) (*createdExpense, error) {
	expenseDbModel := mapExpenseDBModelFromExpense(expense)
	result := r.db.Create(&expenseDbModel)

	if err := result.Error; err != nil {
		return nil, err
	}

	createdExpense := createdExpense{
		id:          expenseDbModel.ID,
		amount:      expenseDbModel.Amount,
		expenseDate: expenseDbModel.ExpenseDate,
		description: expenseDbModel.Description,
		expenseType: expenseDbModel.ExpenseType,
	}

	return &createdExpense, nil
}

func (r repository) GetExpenseTypeById(id uint64) (*expenseType, error) {
	var expenseType expenseType
	result := r.db.Find(expenseType)

	if err := result.Error; err != nil {
		return nil, err
	}

	return &expenseType, nil
}

func mapExpenseDBModelFromExpense(expenseToAdd expense) expenseDBModel {
	return expenseDBModel{
		Amount:        expenseToAdd.amount,
		ExpenseDate:   expenseToAdd.expenseDate,
		Description:   expenseToAdd.description,
		ExpenseTypeID: expenseToAdd.expenseType.id,
		ExpenseType:   expenseToAdd.expenseType,
	}
}

type expenseDBModel struct {
	ID            uint64 `gorm:"primary_key"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Amount        float64
	ExpenseDate   time.Time
	Description   string
	ExpenseTypeID uint64
	ExpenseType   expenseType
}
