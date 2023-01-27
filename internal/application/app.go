package application

import (
	"github.com/labstack/echo/v4"
)

func LoadDependencyConfiguration() {
	WireExpenseTypeRepository = wireExpenseTypeRepository
	WireExpenseRepository = wireExpenseRepository
	WireExpenseTypeService = wireExpenseTypeService
	WireExpenseService = wireExpenseService
	WireExpenseHandler = wireExpenseHandler
	WireDbConnection = wireDbConnection
	WireGenericFieldsValidator = wireGenericFieldsValidator
	WireConfigurations = wireConfigurations
}
func Start(echo *echo.Echo) error {
	injectDependencies()
	mapRoutes(echo)
	return echo.Start("8080")
}

func Finish() {
	_ = SqlDbConnection.Close()
}
