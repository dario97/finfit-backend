package search_in_period

import (
	"time"
)

type command struct {
	startDate time.Time
	endDate   time.Time
}

type Service interface {
	SearchInPeriod(command command) ([]*expense, error)
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

func (s service) SearchInPeriod(command command) ([]*expense, error) {
	//TODO implement me
	panic("implement me")
}
