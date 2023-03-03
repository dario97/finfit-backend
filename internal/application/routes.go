package application

import "github.com/labstack/echo/v4"

func mapRoutes(e *echo.Echo) {
	e.POST("/expense", ExpenseHandler.Add)
	e.POST("/expense-type", ExpenseTypeHandler.Add)
	e.GET("/expense/search", ExpenseHandler.SearchInPeriod)
}
