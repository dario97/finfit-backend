package main

import (
	"finfit/finfit-backend/infrastructure/datastore"
	"finfit/finfit-backend/infrastructure/registry"
	"finfit/finfit-backend/infrastructure/router"
	"fmt"
	"github.com/labstack/echo"
	"log"
)

func main() {
	db := datastore.NewDB()
	db.LogMode(true)
	defer db.Close()

	r := registry.NewRegistry(db)
	e := echo.New()
	e = router.NewRouter(e, r.NewAppController())

	fmt.Println("Server listen at http://localhost" + ":8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatalln(err)
	}
}
