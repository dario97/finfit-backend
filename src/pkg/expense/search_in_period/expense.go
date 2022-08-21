package search_in_period

import "time"

type expense struct {
	id          uint64
	amount      float64
	expenseDate time.Time
	description string
	expenseType expenseType
}
