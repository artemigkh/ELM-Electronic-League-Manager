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

func createGameSummaryBody(body []databaseAccess.GameSummaryInformation) *bytes.Buffer {
	bodyB, _ := json.Marshal(&body)
	return bytes.NewBuffer(bodyB)
}

func testGetGameSummarySessionError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(1, errors.New("fake session error"))

	routes.ElmSessions = mockSession

	httpTest(t, nil, "GET", "/", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testGetGameSummaryNoActiveLeague(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(-1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "GET", "/", 403, testParams{Error: "noActiveLeague"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testGetGameSummaryDatabaseError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetGameSummary", 2).Return([]databaseAccess.GameSummaryInformation{},
		errors.New("fake db error"))

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, nil, "GET", "/", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)
}

func testCorrectGetGameSummary(t *testing.T) {
	var gameSummary []databaseAccess.GameSummaryInformation
	gameSummary = append(gameSummary, databaseAccess.GameSummaryInformation{
		Id:         1,
		Team1Id:    5,
		Team2Id:    6,
		GameTime:   1532913359,
		Complete:   true,
		WinnerId:   5,
		ScoreTeam1: 3,
		ScoreTeam2: 2,
	})
	gameSummary = append(gameSummary, databaseAccess.GameSummaryInformation{
		Id:         2,
		Team1Id:    5,
		Team2Id:    7,
		GameTime:   1532912359,
		Complete:   false,
		WinnerId:   -1,
		ScoreTeam1: 0,
		ScoreTeam2: 0,
	})
	gameSummary = append(gameSummary, databaseAccess.GameSummaryInformation{
		Id:         3,
		Team1Id:    6,
		Team2Id:    7,
		GameTime:   1532911359,
		Complete:   false,
		WinnerId:   -1,
		ScoreTeam1: 0,
		ScoreTeam2: 0,
	})

	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetGameSummary", 2).Return(gameSummary, nil)

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, nil, "GET", "/", 200,
		testParams{ResponseBody: createGameSummaryBody(gameSummary)})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)
}

func Test_GetGameSummary(t *testing.T) {
	//set up router and path to test
	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
	router = gin.New()
	router.GET("/", routes.Testing_Export_getActiveLeague(), routes.Testing_Export_getGameSummary)

	t.Run("sessionError", testGetGameSummarySessionError)
	t.Run("noActiveLeague", testGetGameSummaryNoActiveLeague)
	t.Run("databaseError", testGetGameSummaryDatabaseError)
	t.Run("correctGetGameSummary", testCorrectGetGameSummary)
}
