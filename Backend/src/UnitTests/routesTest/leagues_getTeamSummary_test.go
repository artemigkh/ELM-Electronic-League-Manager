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

func createTeamSummaryBody(body []databaseAccess.TeamSummaryInformation) *bytes.Buffer {
	bodyB, _ := json.Marshal(&body)
	return bytes.NewBuffer(bodyB)
}

func testGetTeamSummarySessionError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(1, errors.New("fake session error"))

	routes.ElmSessions = mockSession

	httpTest(t, nil, "GET", "/teamSummary", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testGetTeamSummaryNoActiveLeague(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(-1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "GET", "/teamSummary", 403, testParams{Error: "noActiveLeague"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testGetTeamSummaryDatabaseError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetTeamSummary", 2).Return([]databaseAccess.TeamSummaryInformation{},
		errors.New("Fake db error"))

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, nil, "GET", "/teamSummary", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)
}

func testCorrectGetTeamSummary(t *testing.T) {
	var teamSummary []databaseAccess.TeamSummaryInformation
	teamSummary = append(teamSummary, databaseAccess.TeamSummaryInformation{
		Id:     1,
		Name:   "team1",
		Tag:    "T1",
		Wins:   2,
		Losses: 0,
	})
	teamSummary = append(teamSummary, databaseAccess.TeamSummaryInformation{
		Id:     2,
		Name:   "team2",
		Tag:    "T2",
		Wins:   1,
		Losses: 1,
	})
	teamSummary = append(teamSummary, databaseAccess.TeamSummaryInformation{
		Id:     3,
		Name:   "team3",
		Tag:    "T3",
		Wins:   0,
		Losses: 2,
	})

	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetTeamSummary", 2).Return(teamSummary, nil)

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, nil, "GET", "/teamSummary", 200, testParams{ResponseBody: createTeamSummaryBody(teamSummary)})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)
}

func Test_GetTeamSummary(t *testing.T) {
	//set up router and path to test
	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
	router = gin.New()
	router.GET("/teamSummary", routes.Testing_Export_getActiveLeague(), routes.Testing_Export_getTeamSummary)

	t.Run("sessionError", testGetTeamSummarySessionError)
	t.Run("noActiveLeague", testGetTeamSummaryNoActiveLeague)
	t.Run("databaseError", testGetTeamSummaryDatabaseError)
	t.Run("correctGetTeamSummary", testCorrectGetTeamSummary)
}
