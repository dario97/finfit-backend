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
	Currency      string
	ExpenseDate   time.Time
	Description   string
	ExpenseTypeID string
	ExpenseType   expensetype.ExpenseType
}

func (receiver Expense) MapToDomainExpense() (*models.Expense, error) {
	id, _ := uuid.Parse(receiver.ID)
	money, err := models.NewMoney(receiver.Amount, receiver.Currency)
	if err != nil {
		return nil, err
	}
	expenseType, err := receiver.ExpenseType.MapToDomainExpenseType()
	if err != nil {
		return nil, err
	}

	return models.NewExpenseWithId(id, money, receiver.ExpenseDate, receiver.Description, expenseType)
}
