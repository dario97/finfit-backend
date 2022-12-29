package main

import (
	"finfit-backend/internal/application"
	"github.com/labstack/echo/v4"
)

func main() {
	defer application.Finish()
	e := echo.New()
	application.LoadConfigurations()
	application.Start(e)
	err := e.Start(":8090")
	if err != nil {
		panic(err)
	}
}
