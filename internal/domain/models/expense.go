package models

import (
	"finfit-backend/pkg"
	"github.com/google/uuid"
	"time"
)

type Expense struct {
	id          uuid.UUID
	amount      *Money
	expenseDate time.Time
	description string
	expenseType *ExpenseType
}

func NewExpense(amount *Money, expenseDate time.Time, description string, expenseType *ExpenseType) *Expense {
	return &Expense{id: pkg.NewUUID(), amount: amount, expenseDate: expenseDate, description: description, expenseType: expenseType}
}

func NewExpenseWithId(id uuid.UUID, amount *Money, expenseDate time.Time, description string, expenseType *ExpenseType) *Expense {
	return &Expense{id: id, amount: amount, expenseDate: expenseDate, description: description, expenseType: expenseType}
}

func (e Expense) Id() uuid.UUID {
	return e.id
}

func (e Expense) Amount() *Money {
	return e.amount
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
