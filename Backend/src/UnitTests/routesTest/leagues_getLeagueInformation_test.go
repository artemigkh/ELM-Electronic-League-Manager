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

func createLeagueInfoBody(id int, name, description string) *bytes.Buffer {
	body := databaseAccess.LeagueInformation{
		Id:          id,
		Name:        name,
		Description: description,
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
		Id:          2,
		Name:        "testName",
		Description: "testDescription",
	}, nil)

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, nil, "GET", "/", 200,
		testParams{ResponseBody: createLeagueInfoBody(2, "testName", "testDescription")})

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
