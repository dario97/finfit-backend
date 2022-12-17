package expense

import (
	"finfit-backend/internal/infrastructure/repository/mysql/expensetype"
	"time"
)

type dbModel struct {
	ID            string `gorm:"primaryKey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Amount        float64
	ExpenseDate   time.Time
	Description   string
	ExpenseTypeID string
	ExpenseType   expensetype.DbModel
}
