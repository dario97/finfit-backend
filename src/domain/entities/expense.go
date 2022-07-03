package entities

import (
	"time"
)

type Expense struct {
	id          int64
	amount      float64
	expenseDate time.Time
	description string
	expenseType ExpenseType
}

func NewExpense(amount float64, expenseDate time.Time, description string, expenseType ExpenseType) Expense {
	return Expense{
		amount:      amount,
		expenseDate: expenseDate,
		description: description,
		expenseType: expenseType,
	}
}

func NewExpenseWithId(id int64, amount float64, expenseDate time.Time, description string, expenseType ExpenseType) Expense {
	return Expense{
		id:          id,
		amount:      amount,
		expenseDate: expenseDate,
		description: description,
		expenseType: expenseType,
	}
}

func (e Expense) Id() int64 {
	return e.id
}

func (e Expense) Amount() float64 {
	return e.amount
}

func (e Expense) ExpenseDate() time.Time {
	return e.expenseDate
}

func (e Expense) Description() string {
	return e.description
}

func (e Expense) ExpenseType() ExpenseType {
	return e.expenseType
}
