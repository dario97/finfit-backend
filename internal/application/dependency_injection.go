package application

import (
	"database/sql"
	expenseService "finfit-backend/internal/domain/services/expense"
	expenseTypeService "finfit-backend/internal/domain/services/expensetype"
	"finfit-backend/internal/infrastructure/interfaces/handler/rest/expense"
	"finfit-backend/internal/infrastructure/interfaces/handler/rest/expensetype"
	"finfit-backend/pkg/fieldvalidation"
	"gorm.io/gorm"
)

var (
	ExpenseHandler         expense.Handler
	ExpenseTypeHandler     expensetype.Handler
	Database               *gorm.DB
	GenericFieldsValidator fieldvalidation.FieldsValidator
	ExpenseRepository      expenseService.Repository
	ExpenseTypeRepository  expenseTypeService.Repository
	ExpenseService         expenseService.Service
	ExpenseTypeService     expenseTypeService.Service
	SqlDbConnection        *sql.DB
	Configs                Configurations
)

func injectDependencies() {
	WireConfigurations()
	WireDbConnection()
	WireGenericFieldsValidator()
	wireRepositories()
	wireServices()
	wireHandlers()
}

func wireRepositories() {
	WireExpenseTypeRepository()
	WireExpenseRepository()
}

func wireServices() {
	WireExpenseTypeService()
	WireExpenseService()
}

func wireHandlers() {
	WireExpenseHandler()
	WireExpenseTypeHandler()
}
