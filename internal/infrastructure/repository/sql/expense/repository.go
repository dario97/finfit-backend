package expense

import (
	"errors"
	"finfit-backend/internal/domain/models"
	"finfit-backend/internal/infrastructure/repository/sql"
	"gorm.io/gorm"
	"time"
)

const dateFormat = "2006-01-02"

type repository struct {
	table string
	db    sql.Database
}

func NewRepository(db sql.Database, table string) *repository {
	return &repository{db: db, table: table}
}

func (r repository) Add(expense *models.Expense) (*models.Expense, error) {
	expenseDbModel := r.mapExpenseDBModelFromExpense(expense)
	result := r.db.Table(r.table).Create(&expenseDbModel)

	if err := result.Error; err != nil {
		return nil, err
	}

	return expense, nil
}

// TODO: no me gusta que el nombre de las tablas este atado a como lo resuelve GORM
func (r repository) SearchInPeriod(startDate time.Time, endDate time.Time) ([]*models.Expense, error) {
	storedExpenses := []Expense{}
	result := r.db.Table(r.table).
		Joins("ExpenseType").
		Find(&storedExpenses, "expense_date >= ?  AND expense_date <= ?", startDate.Format(dateFormat), endDate.Format(dateFormat))

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err := result.Error; err != nil {
		return nil, err
	}

	expenses := []*models.Expense{}
	for _, expense := range storedExpenses {
		domainExpense, err := expense.MapToDomainExpense()
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, domainExpense)
	}

	return expenses, nil
}

func (r repository) mapExpenseDBModelFromExpense(expenseToAdd *models.Expense) Expense {
	return Expense{
		ID:            expenseToAdd.Id().String(),
		Amount:        expenseToAdd.Amount().Amount(),
		Currency:      expenseToAdd.Amount().Currency(),
		ExpenseDate:   expenseToAdd.ExpenseDate(),
		Description:   expenseToAdd.Description(),
		ExpenseTypeID: expenseToAdd.ExpenseType().Id().String(),
	}
}
