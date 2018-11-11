package routesTest

import (
	"Server/databaseAccess"
	"Server/routes"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"mocks"
	"testing"
)

func createPublicLeagueSummaryBody(body []databaseAccess.PublicLeagueInformation) *bytes.Buffer {
	bodyB, _ := json.Marshal(&body)
	return bytes.NewBuffer(bodyB)
}

func testGetPublicLeaguesDatabaseError(t *testing.T) {
	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetPublicLeagueList").Return([]databaseAccess.PublicLeagueInformation{},
		errors.New("fake db error"))

	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, nil, "GET", "/", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockLeaguesDao)
}
func testCorrectGetPublicLeagues(t *testing.T) {
	var leagueSummary []databaseAccess.PublicLeagueInformation
	leagueSummary = append(leagueSummary, databaseAccess.PublicLeagueInformation{
		Id:          1,
		Name:        "testleague1",
		Description: "test description 1",
		PublicJoin:  false,
	})
	leagueSummary = append(leagueSummary, databaseAccess.PublicLeagueInformation{
		Id:          2,
		Name:        "testleague2",
		Description: "test description 2",
		PublicJoin:  true,
	})

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetPublicLeagueList").Return(leagueSummary, nil)

	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, nil, "GET", "/", 200,
		testParams{ResponseBody: createPublicLeagueSummaryBody(leagueSummary)})

	mock.AssertExpectationsForObjects(t, mockLeaguesDao)
}

func Test_GetPublicLeagues(t *testing.T) {
	//set up router and path to test
	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
	router = gin.New()
	router.GET("/", routes.Testing_Export_getPublicLeagues)

	t.Run("databaseError", testGetPublicLeaguesDatabaseError)
	t.Run("correctGetPublicLeagues", testCorrectGetPublicLeagues)
}
