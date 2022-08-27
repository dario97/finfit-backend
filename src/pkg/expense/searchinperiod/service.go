package searchinperiod

import (
	"time"
)

type command struct {
	startDate time.Time
	endDate   time.Time
}

func newCommand(startDate time.Time, endDate time.Time) (*command, error) {
	if startDate.After(endDate) || startDate.Equal(endDate) {
		return nil, invalidArgumentsError{msg: "start date must be before end date"}
	}

	return &command{
		startDate: startDate,
		endDate:   endDate,
	}, nil
}

type Service interface {
	SearchInPeriod(command *command) ([]*expense, error)
}

type Repository interface {
	SearchExpensesInPeriod(startDate time.Time, endDate time.Time) ([]*expense, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return service{
		repository: repository,
	}
}

func (s service) SearchInPeriod(command *command) ([]*expense, error) {
	expenses, err := s.repository.SearchExpensesInPeriod(command.startDate, command.endDate)
	if err != nil {
		return nil, unexpectedError{msg: err.Error()}
	}

	return expenses, nil
}
