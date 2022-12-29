package expense_test

import (
	"finfit-backend/internal/application"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
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
	e := echo.New()
	application.Start(e)

	req := httptest.NewRequest(http.MethodGet,
		fmt.Sprintf("/expense/search?%s",
			"start_date=2006-01-02"),
		nil)

	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)
}

func (suite *HandlerIntegrationTestSuite) TestHola() {
	database := application.Database
	print(database)
}
