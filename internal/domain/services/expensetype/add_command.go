package expensetype

import (
	"errors"
	"finfit-backend/pkg"
)

type AddCommand struct {
	name string
}

func NewAddCommand(name string) (*AddCommand, error) {
	if pkg.IsEmptyOrBlankString(name) {
		return nil, errors.New("invalid command")
	}
	return &AddCommand{name: name}, nil
}
