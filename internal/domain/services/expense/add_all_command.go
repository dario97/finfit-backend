package expense

import "errors"

type AddAllCommand struct {
	addCommands []*AddCommand
}

func NewAddAllCommand(addCommands []*AddCommand) (*AddAllCommand, error) {
	if addCommands == nil {
		return nil, errors.New("invalid command")
	}
	return &AddAllCommand{addCommands: addCommands}, nil
}
