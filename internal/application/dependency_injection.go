package application

import (
	expenseService "finfit-backend/internal/domain/services/expense"
	expenseTypeService "finfit-backend/internal/domain/services/expensetype"
	"finfit-backend/internal/infrastructure/interfaces/handler/rest/expense"
	expenseRepository "finfit-backend/internal/infrastructure/repository/mysql/expense"
	"finfit-backend/internal/infrastructure/repository/mysql/expensetype"
	"finfit-backend/pkg/fieldvalidation"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
)

var (
	ExpenseHandler         expense.Handler
	Database               *gorm.DB
	GenericFieldsValidator fieldvalidation.FieldsValidator
	ExpenseRepository      expenseService.Repository
	ExpenseTypeRepository  expenseTypeService.Repository
	ExpenseService         expenseService.Service
	ExpenseTypeService     expenseTypeService.Service
)

func injectDependencies() {
	wireDbConnection()
	wireGenericFieldsValidator()
	wireRepositories()
	wireServices()
	wireHandlers()
}

func wireRepositories() {
	wireExpenseTypeRepository()
	wireExpenseRepository()
}

func wireExpenseTypeRepository() {
	ExpenseTypeRepository = expensetype.NewRepository(Database)
}

func wireExpenseRepository() {
	ExpenseRepository = expenseRepository.NewRepository(Database)
}

func wireServices() {
	wireExpenseTypeService()
	wireExpenseService()
}

func wireExpenseTypeService() {
	ExpenseTypeService = expenseTypeService.NewService(ExpenseTypeRepository)
}

func wireExpenseService() {
	ExpenseService = expenseService.NewService(ExpenseRepository, ExpenseTypeService)
}

func wireHandlers() {
	wireExpenseHandler()
}

func wireExpenseHandler() {
	ExpenseHandler = expense.NewHandler(ExpenseService, GenericFieldsValidator)
}

func wireDbConnection() {
	DBMS := "mysql"
	mySqlConfig := &mysql.Config{
		User:                 "config.C.Database.User",
		Passwd:               "config.C.Database.Password",
		Net:                  "config.C.Database.Net",
		Addr:                 "config.C.Database.Addr",
		DBName:               "config.C.Database.DBName",
		AllowNativePasswords: false,
		Params: map[string]string{
			"parseTime": "config.C.Database.Params.ParseTime",
		},
	}

	db, err := gorm.Open(DBMS, mySqlConfig.FormatDSN())

	if err != nil {
		log.Fatalln(err)
	}
	db.LogMode(true)

	Database = db
}

func wireGenericFieldsValidator() {
	GenericFieldsValidator, _ = fieldvalidation.RegisterFieldsValidator(nil, nil)
}
