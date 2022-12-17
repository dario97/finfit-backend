package models

import (
	"finfit-backend/pkg"
	"github.com/google/uuid"
)

type ExpenseType struct {
	Id   uuid.UUID
	Name string
}

func NewExpenseType(name string) *ExpenseType {
	return &ExpenseType{Id: pkg.NewUUID(), Name: name}
}
