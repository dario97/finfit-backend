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
	var storedExpenseType ExpenseType
	result := r.db.Table(r.table).First(&storedExpenseType, "id = ?", id.String())

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err := result.Error; err != nil {
		return nil, err
	}

	return storedExpenseType.MapToDomainExpenseType()
}

func (r repository) GetByName(name string) (*models.ExpenseType, error) {
	var storedExpenseType ExpenseType
	result := r.db.Table(r.table).First(&storedExpenseType, "name = ?", name)
	result.Row()
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err := result.Error; err != nil {
		return nil, err
	}

	return storedExpenseType.MapToDomainExpenseType()
}

func (r repository) Add(expenseType *models.ExpenseType) (*models.ExpenseType, error) {
	expenseDbModel := r.mapExpenseTypeDBModelFromExpenseType(expenseType)
	result := r.db.Table(r.table).Create(&expenseDbModel)

	if err := result.Error; err != nil {
		return nil, err
	}

	return expenseType, nil
}

func (r repository) GetAll() ([]*models.ExpenseType, error) {
	//TODO implement me
	panic("implement me")
}

func (r repository) mapExpenseTypeDBModelFromExpenseType(expenseType *models.ExpenseType) ExpenseType {
	return ExpenseType{
		ID:   expenseType.Id().String(),
		Name: expenseType.Name(),
	}
}
