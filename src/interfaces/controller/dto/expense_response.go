package dto

import (
	"finfit-backend/src/domain/entities"
	"time"
)

type ExpenseResponse struct {
	ID          int64       `json:"id"`
	Amount      float64     `json:"amount"`
	ExpenseDate time.Time   `json:"expense_date"`
	Description string      `json:"description"`
	ExpenseType ExpenseType `json:"expense_type"`
}

func NewExpenseResponseFromExpense(expense entities.Expense) ExpenseResponse {
	expenseType := ExpenseType{
		ID:   expense.ExpenseType().Id(),
		Name: expense.ExpenseType().Name(),
	}
	return ExpenseResponse{
		ID:          expense.Id(),
		Amount:      expense.Amount(),
		ExpenseDate: expense.ExpenseDate(),
		Description: expense.Description(),
		ExpenseType: expenseType,
	}
}
