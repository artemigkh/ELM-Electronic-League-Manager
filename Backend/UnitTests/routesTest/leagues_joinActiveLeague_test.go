package routesTest

import (
	"errors"
	"esports-league-manager/Backend/Server/routes"
	"esports-league-manager/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"testing"
)

func testJoinActiveLeagueSessionError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, errors.New("session error"))

	routes.ElmSessions = mockSession

	httpTest(t, nil, "POST", "/", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testJoinActiveLeagueNotLoggedIn(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(-1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "POST", "/", 403, testParams{Error: "notLoggedIn"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testJoinActiveLeagueNoActiveLeague(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, nil)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(-1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "POST", "/", 403, testParams{Error: "noActiveLeague"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testJoinActiveLeagueCanNotJoin(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, nil)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("CanJoinLeague", 2, 1).Return(false, nil)

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, nil, "POST", "/", 400, testParams{Error: "canNotJoin"})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)
}

func testJoinActiveLeagueDatabaseError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, nil)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("CanJoinLeague", 2, 1).Return(true, nil)
	mockLeaguesDao.On("JoinLeague", 2, 1).Return(errors.New("fake db error"))

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, nil, "POST", "/", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)
}

func testCorrectJoinLeague(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, nil)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("CanJoinLeague", 2, 1).Return(true, nil)
	mockLeaguesDao.On("JoinLeague", 2, 1).Return(nil)

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, nil, "POST", "/", 200, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)
}

func Test_JoinActiveLeague(t *testing.T) {
	//set up router and path to test
	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
	router = gin.New()
	router.POST("/",
		routes.Testing_Export_authenticate(),
		routes.Testing_Export_getActiveLeague(),
		routes.Testing_Export_joinActiveLeague)

	t.Run("SessionsError", testJoinActiveLeagueSessionError)
	t.Run("NotLoggedIn", testJoinActiveLeagueNotLoggedIn)
	t.Run("NoActiveLeague", testJoinActiveLeagueNoActiveLeague)
	t.Run("CanNotJoin", testJoinActiveLeagueCanNotJoin)
	t.Run("DatabaseError", testJoinActiveLeagueDatabaseError)
	t.Run("CorrectJoinLeague", testCorrectJoinLeague)
}
