package expensetype

import "time"

type ExpenseType struct {
	ID        string    `gorm:"primaryKey,column:id"`
	Name      string    `gorm:"column:name"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}
