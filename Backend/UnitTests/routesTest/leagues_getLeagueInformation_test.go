package routesTest

import (
	"github.com/gin-gonic/gin"
	"esports-league-manager/Backend/Server/routes"
	"testing"
	"esports-league-manager/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/kataras/iris/core/errors"
	"esports-league-manager/Backend/Server/databaseAccess"
	"bytes"
	"encoding/json"
)

func createLeagueInfoBody(id int) *bytes.Buffer {
	body := databaseAccess.LeagueInformation{
		Id: id,
	}
	bodyB, _ := json.Marshal(&body)
	return bytes.NewBuffer(bodyB)
}

func testGetLeagueDataSessionError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(1, errors.New("fake session error"))

	routes.ElmSessions = mockSession

	httpTest(t, nil, "GET", "/", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testGetLeagueDataNoActiveLeague(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(-1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "GET", "/", 403, testParams{Error: "noActiveLeague"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testGetLeagueDataDatabaseError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeagueInformation", 2).Return(&databaseAccess.LeagueInformation{
		Id: 2,
	}, errors.New("Fake db error"))

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, nil, "GET", "/", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)
}

func testCorrectGetLeagueData(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeagueInformation", 2).Return(&databaseAccess.LeagueInformation{
		Id: 2,
	}, nil)

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, nil, "GET", "/", 200, testParams{ResponseBody: createLeagueInfoBody(2)})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)
}

func Test_GetLeagueInformation(t *testing.T) {
	//set up router and path to test
	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
	router = gin.New()
	router.GET("/", routes.Testing_Export_getActiveLeague(), routes.Testing_Export_getActiveLeagueInformation)

	t.Run("sessionError", testGetLeagueDataSessionError)
	t.Run("noActiveLeague", testGetLeagueDataNoActiveLeague)
	t.Run("databaseError", testGetLeagueDataDatabaseError)
	t.Run("correctGetLeagueData", testCorrectGetLeagueData)
}