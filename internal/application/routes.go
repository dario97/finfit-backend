package application

import "github.com/labstack/echo/v4"

func mapRoutes(e *echo.Echo) {
	e.POST("/expense", ExpenseHandler.Add)
	e.GET("/expense", ExpenseHandler.SearchInPeriod)
}
