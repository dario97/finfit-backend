package main

import (
	"finfit-backend/src/cmd"
	"finfit-backend/src/infrastructure/router"
	"fmt"
	"github.com/labstack/echo"
	"log"
	"net/http"
)

func main() {
	defer cmd.Database.Close()

	e := echo.New()
	appRouter := router.NewRouter(e)

	appRouter.RegisterEndpoint(http.MethodPost, "/expense", func(context echo.Context) error {
		return cmd.AddExpenseHandler.Add(context)
	})

	fmt.Println("Server listen at http://localhost" + ":8080")
	if err := appRouter.Start(":8080"); err != nil {
		log.Fatalln(err)
	}
}
