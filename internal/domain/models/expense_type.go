package models

import (
	"finfit-backend/pkg"
	"github.com/google/uuid"
)

type ExpenseType struct {
	id   uuid.UUID
	name string
}

func NewExpenseType(name string) *ExpenseType {
	return &ExpenseType{id: pkg.NewUUID(), name: name}
}

func NewExpenseTypeWithId(id uuid.UUID, name string) *ExpenseType {
	return &ExpenseType{id: id, name: name}
}

func (e ExpenseType) Id() uuid.UUID {
	return e.id
}

func (e ExpenseType) Name() string {
	return e.name
}
