package expense_test

import (
	"finfit-backend/internal/application"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

type HandlerIntegrationTestSuite struct {
	suite.Suite
}

func TestHandlerIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerIntegrationTestSuite))
}
func (suite *HandlerIntegrationTestSuite) SetupSuite() {
	application.LoadConfigurations()
	application.WireDbConnection = func() {
		db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
		application.Database = db
	}
	application.Start()
}

func (suite *HandlerIntegrationTestSuite) TestHola() {
	database := application.Database
	print(database)
}
