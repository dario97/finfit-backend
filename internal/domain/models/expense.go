package models

import (
	"finfit-backend/pkg"
	"github.com/google/uuid"
	"time"
)

type Expense struct {
	Id          uuid.UUID
	Amount      float64
	Currency    string
	ExpenseDate time.Time
	Description string
	ExpenseType *ExpenseType
}

func NewExpense(amount float64, currency string, expenseDate time.Time, description string, expenseType *ExpenseType) *Expense {
	return &Expense{Id: pkg.NewUUID(), Amount: amount, Currency: currency, ExpenseDate: expenseDate, Description: description, ExpenseType: expenseType}
}
