package application

import (
	"github.com/labstack/echo/v4"
)

type Application interface {
	LoadDependencyConfiguration()
	Start() error
	Finish()
}

type application struct {
	echo *echo.Echo
}

func NewApplication(echo *echo.Echo) *application {
	return &application{echo: echo}
}

func (a application) LoadDependencyConfiguration() {
	WireExpenseTypeRepository = wireExpenseTypeRepository
	WireExpenseRepository = wireExpenseRepository
	WireExpenseTypeService = wireExpenseTypeService
	WireExpenseService = wireExpenseService
	WireExpenseHandler = wireExpenseHandler
	WireDbConnection = wireDbConnection
	WireGenericFieldsValidator = wireGenericFieldsValidator
	WireConfigurations = wireConfigurations
}

func (a application) Start() error {
	injectDependencies()
	mapRoutes(a.echo)
	return a.echo.Start("8080")
}

func (a application) Finish() {
	_ = SqlDbConnection.Close()
}
