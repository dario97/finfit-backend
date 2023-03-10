package expense

import (
	"errors"
	"github.com/google/uuid"
	"strings"
	"time"
)

type AddCommand struct {
	amount        float64
	expenseDate   time.Time
	description   string
	expenseTypeId uuid.UUID
}

func NewAddCommand(amount float64, expenseDate time.Time, description string, expenseTypeId uuid.UUID) (*AddCommand, error) {
	if amount <= 0 || expenseDate.IsZero() || expenseTypeId == uuid.Nil {
		return nil, errors.New("invalid command")
	}
	return &AddCommand{amount: amount, expenseDate: expenseDate, description: strings.TrimSpace(description), expenseTypeId: expenseTypeId}, nil
}
