package expense

import (
	"finfit-backend/internal/domain/models"
	"finfit-backend/internal/infrastructure/repository/sql/expensetype"
	"github.com/google/uuid"
	"time"
)

type Expense struct {
	ID            string `gorm:"primaryKey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Amount        float64
	ExpenseDate   time.Time
	Description   string
	ExpenseTypeID string
	ExpenseType   expensetype.ExpenseType
}

func (receiver Expense) MapToDomainExpense() *models.Expense {
	id, _ := uuid.Parse(receiver.ID)
	return &models.Expense{
		Id:          id,
		Amount:      receiver.Amount,
		ExpenseDate: receiver.ExpenseDate,
		Description: receiver.Description,
		ExpenseType: receiver.ExpenseType.MapToDomainExpenseType(),
	}
}
