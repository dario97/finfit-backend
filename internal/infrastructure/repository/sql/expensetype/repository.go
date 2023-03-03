package expensetype

import (
	"errors"
	"finfit-backend/internal/domain/models"
	"finfit-backend/internal/infrastructure/repository/sql"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type repository struct {
	table string
	db    sql.Database
}

func NewRepository(db sql.Database, table string) *repository {
	return &repository{db: db, table: table}
}

func (r repository) GetByID(id uuid.UUID) (*models.ExpenseType, error) {
	var expenseType models.ExpenseType
	result := r.db.Table(r.table).First(&ExpenseType{}, "id = ?", id.String())

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err := result.Error; err != nil {
		return nil, err
	}

	return &expenseType, nil
}
