package add

import (
	"time"
)

type command struct {
	amount      float64
	expenseDate time.Time
	description string
	expenseType expenseType
}
