package expense

import (
	"errors"
	"time"
)

type SearchInPeriodCommand struct {
	startDate time.Time
	endDate   time.Time
}

func NewSearchInPeriodCommand(startDate time.Time, endDate time.Time) (*SearchInPeriodCommand, error) {
	if startDate.IsZero() || endDate.IsZero() || startDate.After(endDate) {
		return nil, errors.New("invalid command")
	}
	return &SearchInPeriodCommand{startDate: startDate, endDate: endDate}, nil
}
