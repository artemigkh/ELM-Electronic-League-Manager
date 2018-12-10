package routesTest

import (
	"Server/routes"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"mocks"
	"testing"
)

func createSetLeaguePermissionRequestBody(id int, administrator, createTeams, editTeams, editGames bool) *bytes.Buffer {
	reqBody := routes.LeaguePermissionChange{
		Id:            id,
		Administrator: administrator,
		CreateTeams:   createTeams,
		EditTeams:     editTeams,
		EditGames:     editGames,
	}
	reqBodyB, _ := json.Marshal(&reqBody)
	return bytes.NewBuffer(reqBodyB)
}

func testSetLeaguePermissionsSessionError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, errors.New("session error"))

	routes.ElmSessions = mockSession

	httpTest(t, nil, "POST", "/", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testSetLeaguePermissionsNotLoggedIn(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(-1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "POST", "/", 403, testParams{Error: "notLoggedIn"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testSetLeaguePermissionsNoActiveLeague(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, nil)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(-1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "POST", "/", 403, testParams{Error: "noActiveLeague"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testSetLeaguePermissionsNotAdmin(t *testing.T) {
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

	httpTest(t, nil, "POST", "/", 403, testParams{Error: "notAdmin"})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)
}

func testSetLeaguePermissionsDatabaseError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, nil)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 2, 1).
		Return(LeaguePermissions(true, true, true, true), nil)
	mockLeaguesDao.On("SetLeaguePermissions", 2, 1, false, false, false, true).
		Return(errors.New("fake db error"))

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, createSetLeaguePermissionRequestBody(3, false, false, false, true),
		"POST", "/", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)
}

func testCorrectSetLeaguePermissions(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, nil)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 2, 1).
		Return(LeaguePermissions(true, true, true, true), nil)
	mockLeaguesDao.On("SetLeaguePermissions", 2, 1, false, false, false, true).
		Return(nil)

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, createSetLeaguePermissionRequestBody(3, false, false, false, true),
		"POST", "/", 200, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)
}

func Test_SetLeaguePermissions(t *testing.T) {
	//set up router and path to test
	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
	router = gin.New()
	router.POST("/",
		routes.Testing_Export_authenticate(),
		routes.Testing_Export_getActiveLeague(),
		routes.Testing_Export_failIfNotLeagueAdmin(),
		routes.Testing_Export_setLeaguePermissions)

	t.Run("SessionsError", testSetLeaguePermissionsSessionError)
	t.Run("NotLoggedIn", testSetLeaguePermissionsNotLoggedIn)
	t.Run("NoActiveLeague", testSetLeaguePermissionsNoActiveLeague)
	t.Run("NotAdmin", testSetLeaguePermissionsNotAdmin)
	t.Run("DatabaseError", testSetLeaguePermissionsDatabaseError)
	t.Run("CorrectSetLeaguePermissions", testCorrectSetLeaguePermissions)
}
