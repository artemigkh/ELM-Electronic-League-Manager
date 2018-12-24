package routesTest

import (
	"Server/routes"
	"errors"
	"github.com/Pallinder/go-randomdata"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"mocks"
	"testing"
)

func testUpateLeagueInfoSessionError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, errors.New("session error"))

	routes.ElmSessions = mockSession

	httpTest(t, createLeagueRequestBody("testname", "", true, true,
		1, 2, 3, 4),
		"PUT", "/", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testUpateLeagueInfoNotLoggedIn(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(-1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, createLeagueRequestBody("testname", "", true, true,
		1, 2, 3, 4),
		"PUT", "/", 403, testParams{Error: "notLoggedIn"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testUpdateLeagueInfoNoActiveLeague(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, nil)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(-1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "PUT", "/", 403, testParams{Error: "noActiveLeague"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testUpdateLeagueInfoNotAdmin(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, nil)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 2, 1).
		Return(LeaguePermissions(false, true, true, true), nil)
	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, nil, "PUT", "/", 403, testParams{Error: "notAdmin"})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)
}

func testUpateLeagueInfoMalformedBody(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, nil)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 2, 1).
		Return(LeaguePermissions(true, true, true, true), nil)
	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, nil, "PUT", "/", 400, testParams{Error: "malformedInput"})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)
}

func testUpateLeagueInfoDescriptionTooLong(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, nil)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 2, 1).
		Return(LeaguePermissions(true, true, true, true), nil)
	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, createLeagueRequestBody("123456789012345678901234567890123456789012345678901",
		randomdata.RandStringRunes(501), true, true, 1, 2, 3, 4),
		"PUT", "/", 400, testParams{Error: "descriptionTooLong"})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)
}

func testUpateLeagueInfoNameTooLong(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, nil)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 2, 1).
		Return(LeaguePermissions(true, true, true, true), nil)
	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, createLeagueRequestBody("123456789012345678901234567890123456789012345678901", "",
		true, true, 1, 2, 3, 4),
		"PUT", "/", 400, testParams{Error: "nameTooLong"})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)
}

func testUpateLeagueInfoNameInUse(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, nil)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 2, 1).
		Return(LeaguePermissions(true, true, true, true), nil)
	mockLeaguesDao.On("IsNameInUse", 2, "12345678901234567890123456789012345678901234567890").
		Return(true, nil)

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, createLeagueRequestBody("12345678901234567890123456789012345678901234567890",
		"", true, true, 1, 2, 3, 4),
		"PUT", "/", 400, testParams{Error: "nameInUse"})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)
}

func testUpateLeagueInfoDatabaseError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, nil)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 2, 1).
		Return(LeaguePermissions(true, true, true, true), nil)
	mockLeaguesDao.On("IsNameInUse", 2, "testName").
		Return(false, errors.New("fake db error"))

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao
	httpTest(t, createLeagueRequestBody("testName", "", true, true,
		1, 2, 3, 4),
		"PUT", "/", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)
}

func testCorrectLeagueInfoUpdate(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, nil)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 2, 1).
		Return(LeaguePermissions(true, true, true, true), nil)
	mockLeaguesDao.On("IsNameInUse", 2, "testName").
		Return(false, nil)
	mockLeaguesDao.On("UpdateLeague", 2, "testName", mock.AnythingOfType("string"), true, true, 1, 2, 3, 4).
		Return(nil)

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao
	httpTest(t, createLeagueRequestBody("testName", randomdata.RandStringRunes(500), true, true,
		1, 2, 3, 4),
		"PUT", "/", 200, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)
}

func Test_UpdateLeagueInfo(t *testing.T) {
	//set up router and path to test
	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
	router = gin.New()
	router.PUT("/",
		routes.Testing_Export_authenticate(),
		routes.Testing_Export_getActiveLeague(),
		routes.Testing_Export_failIfNotLeagueAdmin(),
		routes.Testing_Export_updateLeagueInfo)

	t.Run("sessionsError", testUpateLeagueInfoSessionError)
	t.Run("notLoggedIn", testUpateLeagueInfoNotLoggedIn)
	t.Run("NoActiveLeague", testUpdateLeagueInfoNoActiveLeague)
	t.Run("NotAdmin", testUpdateLeagueInfoNotAdmin)
	t.Run("malformedBody", testUpateLeagueInfoMalformedBody)
	t.Run("descriptionTooLong", testUpateLeagueInfoDescriptionTooLong)
	t.Run("leagueNameTooLong", testUpateLeagueInfoNameTooLong)
	t.Run("leagueNameInUse", testUpateLeagueInfoNameInUse)
	t.Run("databaseError", testUpateLeagueInfoDatabaseError)
	t.Run("correctLeagueInfoUpdate", testCorrectLeagueInfoUpdate)
}
