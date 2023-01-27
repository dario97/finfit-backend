package main

import (
	"finfit-backend/internal/application"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	log.Info("starting application...")
	defer application.Finish()
	e := echo.New()
	application.LoadDependencyConfiguration()
	err := application.Start(e)
	if err != nil {
		panic(err)
	}

	log.Info("application started on port 8080")
}
