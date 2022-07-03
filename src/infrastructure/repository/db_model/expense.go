package db_model

import (
	"finfit/finfit-backend/src/domain/entities"
	"time"
)

type Expense struct {
	ID            int64 `gorm:"primary_key"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Amount        float64
	ExpenseDate   time.Time
	Description   string
	ExpenseTypeID int64
	ExpenseType   entities.ExpenseType
}

func NewExpenseDbModelFromExpense(expense entities.Expense) Expense {
	return Expense{
		ID:            expense.Id(),
		Amount:        expense.Amount(),
		ExpenseDate:   expense.ExpenseDate(),
		Description:   expense.Description(),
		ExpenseTypeID: expense.ExpenseType().Id(),
		ExpenseType:   expense.ExpenseType(),
	}
}
