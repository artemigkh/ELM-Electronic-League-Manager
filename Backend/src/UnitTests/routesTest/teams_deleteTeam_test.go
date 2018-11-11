package routesTest

import (
	"Server/routes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"mocks"
	"testing"
)

func testDeleteTeamInformationNoId(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "DELETE", "/", 404, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testDeleteTeamInformationNotInt(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "DELETE", "/a", 400, testParams{Error: "IdMustBeInteger"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testDeleteTeamInformationSessionError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(1, errors.New("fake session error"))

	routes.ElmSessions = mockSession

	httpTest(t, nil, "DELETE", "/1", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testDeleteTeamInformationNoActiveLeague(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(-1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "DELETE", "/1", 403, testParams{Error: "noActiveLeague"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testDeleteTeamInformationNotLoggedIn(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(1, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(-1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "DELETE", "/1", 403, testParams{Error: "notLoggedIn"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testDeleteTeamInformationTeamIsActive(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(2, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("IsTeamActive", 5, 7).Return(true, nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, nil, "DELETE", "/7", 400, testParams{Error: "teamIsActive"})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao)
}

func testDeleteTeamInformationNoTeamEditPermissions(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(2, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("IsTeamActive", 5, 7).Return(false, nil)
	mockTeamsDao.On("GetTeamPermissions", 7, 2).
		Return(TeamPermissions(false, true, true, true), nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 5, 2).
		Return(LeaguePermissions(false, false, false, false), nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, nil, "DELETE", "/7", 403, testParams{Error: "notTeamAdmin"})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao, mockLeaguesDao)
}

func testDeleteTeamInformationTeamDoesNotExist(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(2, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("IsTeamActive", 5, 7).Return(false, nil)
	mockTeamsDao.On("DoesTeamExist", 5, 7).Return(false, nil)
	mockTeamsDao.On("GetTeamPermissions", 7, 2).
		Return(TeamPermissions(true, true, true, true), nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 5, 2).
		Return(LeaguePermissions(false, false, false, false), nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, nil, "DELETE", "/7", 400, testParams{Error: "teamDoesNotExist"})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao, mockLeaguesDao)
}

func testDeleteTeamInformationDbError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(2, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("IsTeamActive", 5, 7).Return(false, nil)
	mockTeamsDao.On("DoesTeamExist", 5, 7).Return(true, nil)
	mockTeamsDao.On("DeleteTeam", 5, 7).Return(errors.New("fake db error"))
	mockTeamsDao.On("GetTeamPermissions", 7, 2).
		Return(TeamPermissions(true, true, true, true), nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 5, 2).
		Return(LeaguePermissions(false, false, false, false), nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, nil, "DELETE", "/7", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao, mockLeaguesDao)
}

func testDeleteTeamInformationCorrectDeleteTeam(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(2, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("IsTeamActive", 5, 7).Return(false, nil)
	mockTeamsDao.On("DoesTeamExist", 5, 7).Return(true, nil)
	mockTeamsDao.On("DeleteTeam", 5, 7).Return(nil)
	mockTeamsDao.On("GetTeamPermissions", 7, 2).
		Return(TeamPermissions(true, true, true, true), nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 5, 2).
		Return(LeaguePermissions(false, false, false, false), nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, nil, "DELETE", "/7", 200, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao, mockLeaguesDao)
}

func Test_DeleteTeam(t *testing.T) {
	//set up router and path to test
	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
	router = gin.New()

	router.Use(routes.Testing_Export_getActiveLeague())
	router.DELETE("/:id",
		routes.Testing_Export_getUrlId(),
		routes.Testing_Export_authenticate(),
		routes.Testing_Export_failIfTeamActive(),
		routes.Testing_Export_failIfNotTeamAdministrator(),
		routes.Testing_Export_deleteTeam)

	t.Run("NoId", testDeleteTeamInformationNoId)
	t.Run("IdNotInt", testDeleteTeamInformationNotInt)
	t.Run("SessionError", testDeleteTeamInformationSessionError)
	t.Run("NoActiveLeague", testDeleteTeamInformationNoActiveLeague)
	t.Run("NotLoggedIn", testDeleteTeamInformationNotLoggedIn)
	t.Run("NoTeamEditPermissions", testDeleteTeamInformationNoTeamEditPermissions)
	t.Run("TeamDoesNotExist", testDeleteTeamInformationTeamDoesNotExist)
	t.Run("TeamIsActive", testDeleteTeamInformationTeamIsActive)
	t.Run("DbError", testDeleteTeamInformationDbError)
	t.Run("CorrectDeleteTeam", testDeleteTeamInformationCorrectDeleteTeam)
}
