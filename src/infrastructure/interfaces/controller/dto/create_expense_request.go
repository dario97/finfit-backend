package dto

import (
	"time"
)

type CreateExpenseRequest struct {
	Amount      float64      `json:"amount" validate:"required,gt=0"`
	ExpenseDate time.Time    `json:"expense_date" validate:"required,datetime=2006-01-02"`
	Description string       `json:"description"`
	ExpenseType *ExpenseType `json:"expense_type" validate:"required"`
}

type ExpenseType struct {
	ID   int64  `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}
