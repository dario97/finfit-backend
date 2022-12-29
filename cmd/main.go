package main

import (
	"finfit-backend/internal/application"
)

func main() {
	defer application.Finish()
	application.Start()
}
