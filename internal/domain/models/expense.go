package models

import (
	"errors"
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

func NewExpense(amount *Money, expenseDate time.Time, description string, expenseType *ExpenseType) (*Expense, error) {
	id := pkg.NewUUID()
	err := validateExpense(id, amount, expenseDate, expenseType)
	if err != nil {
		return nil, err
	}

	return &Expense{id: id, amount: amount, expenseDate: expenseDate, description: description, expenseType: expenseType}, nil
}

func NewExpenseWithId(id uuid.UUID, amount *Money, expenseDate time.Time, description string, expenseType *ExpenseType) (*Expense, error) {
	err := validateExpense(id, amount, expenseDate, expenseType)
	if err != nil {
		return nil, err
	}

	return &Expense{id: id, amount: amount, expenseDate: expenseDate, description: description, expenseType: expenseType}, nil
}

func validateExpense(id uuid.UUID, amount *Money, expenseDate time.Time, expenseType *ExpenseType) error {
	if id == uuid.Nil {
		return errors.New("invalid id, is must be a valid UUID")
	}

	if amount.Amount() <= 0 {
		return errors.New("invalid expense amount, it cannot be lower than 1.0")
	}

	if expenseDate.IsZero() {
		return errors.New("invalid expense date, it cannot be zero")
	}

	if expenseType == nil {
		return errors.New("invalid expense type, it cannot be null")
	}
	return nil
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
