package application

import "github.com/labstack/echo/v4"

func mapRoutes(e *echo.Echo) {
	v1Group := e.Group("/v1")
	v1Group.POST("/expenses", ExpenseHandler.Add)
	v1Group.POST("/expense-types", ExpenseTypeHandler.Add)
	v1Group.GET("/expenses", ExpenseHandler.SearchInPeriod)
}
