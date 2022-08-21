package add

import "time"

type createdExpense struct {
	id          uint64
	amount      float64
	expenseDate time.Time
	description string
	expenseType expenseType
}
