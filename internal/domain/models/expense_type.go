package models

import (
	"errors"
	"finfit-backend/pkg"
	"github.com/google/uuid"
)

type ExpenseType struct {
	id   uuid.UUID
	name string
}

func NewExpenseType(name string) (*ExpenseType, error) {
	id := pkg.NewUUID()
	err := validateExpenseType(id, name)
	if err != nil {
		return nil, err
	}
	return &ExpenseType{id: id, name: name}, nil
}

func NewExpenseTypeWithId(id uuid.UUID, name string) (*ExpenseType, error) {
	err := validateExpenseType(id, name)
	if err != nil {
		return nil, err
	}

	return &ExpenseType{id: id, name: name}, nil
}

func validateExpenseType(id uuid.UUID, name string) error {
	if id == uuid.Nil {
		return errors.New("invalid id, is must be a valid UUID")
	}

	if pkg.IsEmptyOrBlankString(name) {
		return errors.New("invalid name, cannot be empty")
	}
	return nil
}

func (e ExpenseType) Id() uuid.UUID {
	return e.id
}

func (e ExpenseType) Name() string {
	return e.name
}
