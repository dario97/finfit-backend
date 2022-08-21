package add

import "time"

type expense struct {
	amount      float64
	expenseDate time.Time
	description string
	expenseType expenseType
}
