package models

import (
	"finfit-backend/pkg"
	"github.com/google/uuid"
	"time"
)

type Expense struct {
	id          uuid.UUID
	amount      float64
	currency    string
	expenseDate time.Time
	description string
	expenseType *ExpenseType
}

func NewExpense(amount float64, currency string, expenseDate time.Time, description string, expenseType *ExpenseType) *Expense {
	return &Expense{id: pkg.NewUUID(), amount: amount, currency: currency, expenseDate: expenseDate, description: description, expenseType: expenseType}
}

func NewExpenseWithId(id uuid.UUID, amount float64, currency string, expenseDate time.Time, description string, expenseType *ExpenseType) *Expense {
	return &Expense{id: id, amount: amount, currency: currency, expenseDate: expenseDate, description: description, expenseType: expenseType}
}

func (e Expense) Id() uuid.UUID {
	return e.id
}

func (e Expense) Amount() float64 {
	return e.amount
}

func (e Expense) Currency() string {
	return e.currency
}

func (e Expense) ExpenseDate() time.Time {
	return e.expenseDate
}

func (e Expense) Description() string {
	return e.description
}

func (e Expense) ExpenseType() *ExpenseType {
	return e.expenseType
}
