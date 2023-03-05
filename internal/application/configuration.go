package application

import (
	"database/sql"
	expenseService "finfit-backend/internal/domain/services/expense"
	expenseTypeServ "finfit-backend/internal/domain/services/expensetype"
	expense2 "finfit-backend/internal/infrastructure/interfaces/handler/rest/expense"
	expensetype2 "finfit-backend/internal/infrastructure/interfaces/handler/rest/expensetype"
	"finfit-backend/internal/infrastructure/repository/sql/expense"
	"finfit-backend/internal/infrastructure/repository/sql/expensetype"
	"finfit-backend/pkg/fieldvalidation"
	"fmt"
	"github.com/labstack/gommon/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"
	"strings"
)

var WireExpenseTypeRepository func()
var WireExpenseRepository func()
var WireExpenseTypeService func()
var WireExpenseService func()
var WireExpenseHandler func()
var WireExpenseTypeHandler func()
var WireDbConnection func()
var WireGenericFieldsValidator func()
var WireConfigurations func()

const (
	databaseHostConfigKey     = "DATABASE_HOST"
	databasePortConfigKey     = "DATABASE_PORT"
	databaseUserConfigKey     = "DATABASE_USER"
	databasePasswordConfigKey = "DATABASE_PASSWORD"
	databaseNameConfigKey     = "DATABASE_NAME"
	databaseDriverConfigKey   = "DATABASE_DRIVER"
)

type Configurations interface {
	Load()
	GetString(key string) string
}
type configurations struct {
	configs map[string]string
}

func (c *configurations) Load() {
	c.configs = map[string]string{}
	environmentVariables := os.Environ()
	for _, variable := range environmentVariables {
		keyValuePair := strings.Split(variable, "=")
		c.configs[keyValuePair[0]] = keyValuePair[1]
	}
}

func (c *configurations) GetString(key string) string {
	return c.configs[key]
}

// TODO: el nombre de las tablas tiene que venir por config
func wireExpenseTypeRepository() {
	ExpenseTypeRepository = expensetype.NewRepository(Database, "expense_type")
}

func wireExpenseRepository() {
	ExpenseRepository = expense.NewRepository(Database, "expense")
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

func wireExpenseTypeHandler() {
	ExpenseTypeHandler = expensetype2.NewHandler(ExpenseTypeService, GenericFieldsValidator)
}

// TODO: el nombre del schema tiene que venir por config
func wireDbConnection() {
	log.Info("starting database connection...")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		Configs.GetString(databaseHostConfigKey),
		Configs.GetString(databasePortConfigKey),
		Configs.GetString(databaseUserConfigKey),
		Configs.GetString(databasePasswordConfigKey),
		Configs.GetString(databaseNameConfigKey))
	sqlDB, err := sql.Open(Configs.GetString(databaseDriverConfigKey), dsn)
	log.Info("DSN: " + dsn)
	if err != nil {
		log.Panic(err)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{NamingStrategy: schema.NamingStrategy{TablePrefix: "public.", SingularTable: true}})

	if err != nil {
		log.Panic(err)
	}

	log.Info("database connection started")
	SqlDbConnection = sqlDB
	Database = db
}

func wireGenericFieldsValidator() {
	GenericFieldsValidator, _ = fieldvalidation.RegisterFieldsValidator(nil, nil)
}

func wireConfigurations() {
	Configs = &configurations{}
	Configs.Load()
}
