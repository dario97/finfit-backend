package main

import (
	"finfit-backend/internal/application"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	log.Info("starting application...")
	e := echo.New()
	app := application.NewApplication(e)
	defer app.Finish()
	app.LoadDependencyConfiguration()
	err := app.Start()
	if err != nil {
		panic(err)
	}

	log.Info("application started on port 8080")
}
