package cmd

import (
	"finfit-backend/src/pkg/expense/add"
	"finfit-backend/src/pkg/fieldvalidation"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en2 "github.com/go-playground/validator/v10/translations/en"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
)

var (
	AddExpenseHandler      add.Handler
	Database               *gorm.DB
	GenericFieldsValidator fieldvalidation.FieldsValidator
)

func init() {
	wireDbConnection()
	wireGenericFieldsValidator()
	wireAddExpenseHandler()
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
	validate := validator.New()
	english := en.New()
	uni := ut.New(english, english)
	translator, _ := uni.GetTranslator("en")
	_ = en2.RegisterDefaultTranslations(validate, translator)

	GenericFieldsValidator = fieldvalidation.NewGenericFieldsValidator(validate, translator)
}

func wireAddExpenseHandler() {
	addExpenseRepository := add.NewRepository(Database)
	addExpenseService := add.NewService(addExpenseRepository)
	addExpenseHandler := add.NewHandler(addExpenseService, GenericFieldsValidator)

	AddExpenseHandler = addExpenseHandler
}
