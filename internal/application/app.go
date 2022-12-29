package application

import (
	"github.com/labstack/echo/v4"
)

func LoadConfigurations() {
	WireExpenseTypeRepository = wireExpenseTypeRepository
	WireExpenseRepository = wireExpenseRepository
	WireExpenseTypeService = wireExpenseTypeService
	WireExpenseService = wireExpenseService
	WireExpenseHandler = wireExpenseHandler
	WireDbConnection = wireDbConnection
	WireGenericFieldsValidator = wireGenericFieldsValidator
}
func Start(echo *echo.Echo) {
	injectDependencies()
	mapRoutes(echo)
}

func Finish() {
	_ = SqlDbConnection.Close()
}
