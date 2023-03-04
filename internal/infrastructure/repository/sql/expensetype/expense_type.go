package expensetype

import (
	"finfit-backend/internal/domain/models"
	"github.com/google/uuid"
	"time"
)

type ExpenseType struct {
	ID        string    `gorm:"primaryKey,column:id"`
	Name      string    `gorm:"column:name"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (receiver ExpenseType) MapToDomainExpenseType() *models.ExpenseType {
	id, _ := uuid.Parse(receiver.ID)
	return &models.ExpenseType{
		Id:   id,
		Name: receiver.Name,
	}
}
