package router

import (
	"finfit/finfit-backend/src/interfaces/controller"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func NewRouter(e *echo.Echo, appController controller.AppController) *echo.Echo {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/expense", func(context echo.Context) error {
		return appController.Create(context)
	})

	return e
}
