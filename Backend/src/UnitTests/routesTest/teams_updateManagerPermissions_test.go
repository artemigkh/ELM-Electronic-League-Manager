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

func createManagerPermissionChangeRequestBody(teamId, userId int, administrator, information, players, reportResults bool) *bytes.Buffer {
	reqBody := routes.TeamManagerPermissionChange{
		TeamId:        teamId,
		UserId:        userId,
		Administrator: administrator,
		Information:   information,
		Players:       players,
		ReportResults: reportResults,
	}
	reqBodyB, _ := json.Marshal(&reqBody)
	return bytes.NewBuffer(reqBodyB)
}

func testChangeManagerPermissionsNoActiveLeague(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(-1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "POST", "/", 403, testParams{Error: "noActiveLeague"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testChangeManagerPermissionsSessionError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, errors.New("session error"))

	routes.ElmSessions = mockSession

	httpTest(t, nil, "POST", "/", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testChangeManagerPermissionsNotLoggedIn(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(1, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(-1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "POST", "/", 403, testParams{Error: "notLoggedIn"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testChangeManagerPermissionsMalformedBody(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(4, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "POST", "/", 400, testParams{Error: "malformedInput"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testChangeManagerPermissionsTeamDoesNotExist(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(4, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExistInLeague", 5, 12).
		Return(false, nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, createManagerPermissionChangeRequestBody(12, 15, false, true, true, true),
		"POST", "/", 400, testParams{Error: "teamDoesNotExist"})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao)
}

func testChangeManagerPermissionsManagerDoesNotExist(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(4, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExistInLeague", 5, 12).
		Return(true, nil)
	mockTeamsDao.On("GetTeamPermissions", 12, 15).
		Return(TeamPermissions(false, false, false, false), nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, createManagerPermissionChangeRequestBody(12, 15, false, true, true, true),
		"POST", "/", 400, testParams{Error: "managerDoesNotExist"})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao)
}

func testChangeManagerPermissionsNotAdmin(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(4, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExistInLeague", 5, 12).
		Return(true, nil)
	mockTeamsDao.On("GetTeamPermissions", 12, 15).
		Return(TeamPermissions(false, false, false, true), nil)
	mockTeamsDao.On("GetTeamPermissions", 12, 4).
		Return(TeamPermissions(false, true, true, true), nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 5, 4).
		Return(LeaguePermissions(false, false, false, false), nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, createManagerPermissionChangeRequestBody(12, 15, false, true, true, true),
		"POST", "/", 403, testParams{Error: "notAdmin"})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao, mockLeaguesDao)
}
func testChangeManagerPermissionsDatabaseError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(4, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExistInLeague", 5, 12).
		Return(true, nil)
	mockTeamsDao.On("GetTeamPermissions", 12, 15).
		Return(TeamPermissions(false, false, false, true), nil)
	mockTeamsDao.On("GetTeamPermissions", 12, 4).
		Return(TeamPermissions(true, true, true, true), nil)
	mockTeamsDao.On("ChangeManagerPermissions", 12, 15, false, true, true, true).
		Return(errors.New("fake db error"))

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 5, 4).
		Return(LeaguePermissions(false, false, false, false), nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, createManagerPermissionChangeRequestBody(12, 15, false, true, true, true),
		"POST", "/", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao, mockLeaguesDao)
}
func testChangeManagerPermissionsCorrectChangePermissions(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(4, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExistInLeague", 5, 12).
		Return(true, nil)
	mockTeamsDao.On("GetTeamPermissions", 12, 15).
		Return(TeamPermissions(false, false, false, true), nil)
	mockTeamsDao.On("GetTeamPermissions", 12, 4).
		Return(TeamPermissions(true, true, true, true), nil)
	mockTeamsDao.On("ChangeManagerPermissions", 12, 15, false, true, true, true).
		Return(nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 5, 4).
		Return(LeaguePermissions(false, false, false, false), nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, createManagerPermissionChangeRequestBody(12, 15, false, true, true, true),
		"POST", "/", 200, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao, mockLeaguesDao)
}

func Test_ChangeManagerPermissions(t *testing.T) {
	//set up router and path to test
	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
	router = gin.New()

	router.Use(routes.Testing_Export_getActiveLeague())
	router.POST("/",
		routes.Testing_Export_authenticate(),
		routes.Testing_Export_updateManagerPermissions)

	t.Run("NoActiveLeague", testChangeManagerPermissionsNoActiveLeague)
	t.Run("SessionsError", testChangeManagerPermissionsSessionError)
	t.Run("NotLoggedIn", testChangeManagerPermissionsNotLoggedIn)
	t.Run("MalformedBody", testChangeManagerPermissionsMalformedBody)
	t.Run("TeamDoesNotExist", testChangeManagerPermissionsTeamDoesNotExist)
	t.Run("ManagerDoesNotExist", testChangeManagerPermissionsManagerDoesNotExist)
	t.Run("NotAdmin", testChangeManagerPermissionsNotAdmin)
	t.Run("DatabaseError", testChangeManagerPermissionsDatabaseError)
	t.Run("CorrectChangePermissions", testChangeManagerPermissionsCorrectChangePermissions)
}
