package application

import (
	"database/sql"
	expenseService "finfit-backend/internal/domain/services/expense"
	expenseTypeServ "finfit-backend/internal/domain/services/expensetype"
	expense2 "finfit-backend/internal/infrastructure/interfaces/handler/rest/expense"
	"finfit-backend/internal/infrastructure/repository/sql/expense"
	"finfit-backend/internal/infrastructure/repository/sql/expensetype"
	"finfit-backend/pkg/fieldvalidation"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var WireExpenseTypeRepository func()
var WireExpenseRepository func()
var WireExpenseTypeService func()
var WireExpenseService func()
var WireExpenseHandler func()
var WireDbConnection func()
var WireGenericFieldsValidator func()

func wireExpenseTypeRepository() {
	ExpenseTypeRepository = expensetype.NewRepository(Database)
}

func wireExpenseRepository() {
	ExpenseRepository = expense.NewRepository(Database)
}

func wireExpenseTypeService() {
	ExpenseTypeService = expenseTypeServ.NewService(ExpenseTypeRepository)
}

func wireExpenseService() {
	ExpenseService = expenseService.NewService(ExpenseRepository, ExpenseTypeService)
}

func wireExpenseHandler() {
	ExpenseHandler = expense2.NewHandler(ExpenseService, GenericFieldsValidator)
}

func wireDbConnection() {
	sqlDB, err := sql.Open("mysql", "gorm.db")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	SqlDbConnection = sqlDB
	Database = db
}

func wireGenericFieldsValidator() {
	GenericFieldsValidator, _ = fieldvalidation.RegisterFieldsValidator(nil, nil)
}
