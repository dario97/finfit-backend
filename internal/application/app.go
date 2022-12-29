package application

import (
	"github.com/labstack/echo/v4"
)

func Start() {
	e := echo.New()
	injectDependencies()
	mapRoutes(e)
	_ = e.Start(":8080")
}

func Finish() {
	_ = Database.Close()
}
