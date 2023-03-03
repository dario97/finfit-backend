package expense

import (
	"finfit-backend/internal/domain/models"
	"finfit-backend/internal/infrastructure/repository/sql"
	"time"
)

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

func (r repository) SearchInPeriod(startDate time.Time, endDate time.Time) ([]*models.Expense, error) {
	//TODO implement me
	panic("implement me")
}

func (r repository) mapExpenseDBModelFromExpense(expenseToAdd *models.Expense) Expense {
	return Expense{
		ID:            expenseToAdd.Id.String(),
		Amount:        expenseToAdd.Amount,
		ExpenseDate:   expenseToAdd.ExpenseDate,
		Description:   expenseToAdd.Description,
		ExpenseTypeID: expenseToAdd.ExpenseType.Id.String(),
	}
}
