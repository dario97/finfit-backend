package main

import (
	"finfit-backend/internal/application"
)

func main() {
	defer application.Close()
	application.Start()
}
