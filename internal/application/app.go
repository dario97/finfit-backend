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
func Start() {
	e := echo.New()
	injectDependencies()
	mapRoutes(e)
	//err := e.Start(":8090")
	//if err != nil {
	//	panic(err)
	//}
}

func Finish() {
	_ = SqlDbConnection.Close()
}
