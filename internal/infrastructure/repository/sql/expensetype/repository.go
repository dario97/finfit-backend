package expensetype

import (
	"finfit-backend/internal/domain/models"
	"finfit-backend/internal/infrastructure/repository/sql"
	"github.com/google/uuid"
)

type repository struct {
	db sql.Database
}

func NewRepository(db sql.Database) *repository {
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
