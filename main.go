package main

import (
	"finfit-backend/src/infrastructure/datastore"
	"finfit-backend/src/infrastructure/registry"
	"finfit-backend/src/infrastructure/router"
	"fmt"
	"github.com/labstack/echo"
	"log"
	"net/http"
)

func main() {
	db := datastore.NewDB()
	db.LogMode(true)
	defer db.Close()

	repositoryRegistry := registry.NewRepositoryRegistry(db)
	serviceRegistry := registry.NewServiceRegistry(repositoryRegistry)
	controllerRegistry := registry.NewControllerRegistry(serviceRegistry)

	expenseController := controllerRegistry.GetExpenseController()

	e := echo.New()
	appRouter := router.NewRouter(e)

	appRouter.RegisterEndpoint(http.MethodPost, "/expense", func(context echo.Context) error {
		return expenseController.Create(context)
	})

	fmt.Println("Server listen at http://localhost" + ":8080")
	if err := appRouter.Start(":8080"); err != nil {
		log.Fatalln(err)
	}
}
