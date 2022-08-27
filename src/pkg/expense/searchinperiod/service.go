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
	return []*expense{
		&expense{
			id:          1,
			amount:      10.3,
			expenseDate: time.Date(2022, 8, 4, 0, 0, 0, 0, time.UTC),
			description: "playstation 5",
			expenseType: expenseType{
				id:   1,
				name: "Tecnologia",
			},
		},
		&expense{
			id:          2,
			amount:      101.2,
			expenseDate: time.Date(2022, 7, 2, 0, 0, 0, 0, time.UTC),
			description: "playstation 1",
			expenseType: expenseType{
				id:   1,
				name: "Antiguedades",
			},
		},
	}, nil
}
