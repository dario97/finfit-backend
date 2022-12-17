package models

import (
	"finfit-backend/pkg"
	"github.com/google/uuid"
	"time"
)

type Expense struct {
	Id          uuid.UUID
	Amount      float64
	ExpenseDate time.Time
	Description string
	ExpenseType *ExpenseType
}

func NewExpense(amount float64, expenseDate time.Time, description string, expenseType *ExpenseType) *Expense {
	return &Expense{Id: pkg.NewUUID(), Amount: amount, ExpenseDate: expenseDate, Description: description, ExpenseType: expenseType}
}
