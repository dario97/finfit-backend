package searchinperiod

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test01(t *testing.T) {
	repositoryMock := newRepositoryMock()
	service := NewService(repositoryMock)
	startDate := time.Date(2022, 7, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2022, 9, 1, 0, 0, 0, 0, time.UTC)
	command, _ := newCommand(startDate, endDate)

	actualExpenses, err := service.SearchInPeriod(command)

	expectedExpenses := []*expense{
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
	}

	assert.Nil(t, err)
	assert.Equal(t, expectedExpenses, actualExpenses)
}
