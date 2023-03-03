package expensetype

import (
	"errors"
	"finfit-backend/pkg"
)

type AddCommand struct {
	name string
}

func NewAddCommand(name string) (*AddCommand, error) {
	if pkg.IsEmptyOrBlankString(name) || !pkg.HasMin(name, 3) || pkg.ExceedsMax(name, 32) {
		return nil, errors.New("invalid command")
	}
	return &AddCommand{name: name}, nil
}
