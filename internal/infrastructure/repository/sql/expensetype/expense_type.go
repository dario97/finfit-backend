package expensetype

import "time"

type ExpenseType struct {
	ID        string `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
}
