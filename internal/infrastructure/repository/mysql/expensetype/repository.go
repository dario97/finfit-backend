package expensetype

import (
	"finfit-backend/internal/domain/models"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r repository) GetByID(id uuid.UUID) (*models.ExpenseType, error) {
	var expenseType models.ExpenseType
	result := r.db.First(&DbModel{}, "id = ?", id.String())

	if err := result.Error; err != nil {
		return nil, err
	}

	return &expenseType, nil
}
