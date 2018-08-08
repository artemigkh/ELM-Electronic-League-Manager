package routesTest

import (
	"errors"
	"esports-league-manager/Backend/Server/routes"
	"esports-league-manager/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"testing"
)

func testSetLeagueNoId(t *testing.T) {
	httpTest(t, nil, "POST", "/", 404, testParams{})
}

func testSetLeagueIdNotInt(t *testing.T) {
	httpTest(t, nil, "POST", "/a", 400, testParams{Error: "IdMustBeInteger"})
}

func testSetLeagueNotViewable(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserID", mock.Anything).
		Return(1, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("IsLeagueViewable", 2, 1).Return(false, nil)

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, nil, "POST", "/2", 403, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)
}

func testSetLeagueSessionInternalError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserID", mock.Anything).
		Return(1, errors.New("session error"))

	routes.ElmSessions = mockSession

	httpTest(t, nil, "POST", "/2", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testSetLeagueDatabaseError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserID", mock.Anything).
		Return(1, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("IsLeagueViewable", 2, 1).
		Return(false, errors.New("database error"))

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, nil, "POST", "/2", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)
}

func testSetLeagueSetSessionError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserID", mock.Anything).
		Return(1, nil)
	mockSession.On("SetActiveLeague", mock.Anything, 2).Return(errors.New("set session error"))

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("IsLeagueViewable", 2, 1).Return(true, nil)

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, nil, "POST", "/2", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)
}

func testCorrectSetLeague(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserID", mock.Anything).
		Return(1, nil)
	mockSession.On("SetActiveLeague", mock.Anything, 2).Return(nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("IsLeagueViewable", 2, 1).Return(true, nil)

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, nil, "POST", "/2", 200, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)
}

func Test_CreateSetLeague(t *testing.T) {
	//set up router and path to test
	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
	router = gin.New()
	router.POST("/:id", routes.Testing_Export_getUrlId(), routes.Testing_Export_setActiveLeague)

	t.Run("noId", testSetLeagueNoId)
	t.Run("IdNotInt", testSetLeagueIdNotInt)
	t.Run("notViewable", testSetLeagueNotViewable)
	t.Run("sessionError", testSetLeagueSessionInternalError)
	t.Run("databaseError", testSetLeagueDatabaseError)
	t.Run("setSessionError", testSetLeagueSetSessionError)
	t.Run("correctSetLeague", testCorrectSetLeague)
}
