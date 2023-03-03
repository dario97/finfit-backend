package expense

import (
	"errors"
	"finfit-backend/internal/domain/models"
	"strings"
	"time"
)

type AddCommand struct {
	amount      float64
	expenseDate time.Time
	description string
	expenseType *models.ExpenseType
}

func NewAddCommand(amount float64, expenseDate time.Time, description string, expenseType *models.ExpenseType) (*AddCommand, error) {
	if amount <= 0 || expenseDate.IsZero() || expenseType == nil {
		return nil, errors.New("invalid command")
	}
	return &AddCommand{amount: amount, expenseDate: expenseDate, description: strings.TrimSpace(description), expenseType: expenseType}, nil
}
