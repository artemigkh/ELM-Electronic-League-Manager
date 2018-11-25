package routesTest

import (
	"Server/routes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"mocks"
	"testing"
)

func testUpdateTeamInformationNoId(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "PUT", "/", 404, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testUpdateTeamInformationNotInt(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "PUT", "/a", 400, testParams{Error: "IdMustBeInteger"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testUpdateTeamSessionError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(1, errors.New("fake session error"))

	routes.ElmSessions = mockSession

	httpTest(t, nil, "PUT", "/1", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testUpdateTeamNotLoggedIn(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "PUT", "/a", 400, testParams{Error: "IdMustBeInteger"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testUpdateTeamNoEditPermissions(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(2, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("GetTeamPermissions", 7, 2).
		Return(TeamPermissions(false, false, true, true), nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 5, 2).
		Return(LeaguePermissions(false, false, false, true), nil)

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, nil, "PUT", "/7", 403, testParams{Error: "noEditTeamInformationPermissions"})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao, mockTeamsDao)
}

func testUpdateTeamInformationTeamDoesNotExist(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(2, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 5, 7).Return(false, nil)
	mockTeamsDao.On("GetTeamPermissions", 7, 2).
		Return(TeamPermissions(false, false, false, false), nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 5, 2).
		Return(LeaguePermissions(false, false, true, false), nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, createTeamRequestBody("sampleName", "TAG"),
		"PUT", "/7", 400, testParams{Error: "teamDoesNotExist"})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao, mockLeaguesDao)
}

func testUpdateTeamMalformedBody(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(2, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("GetTeamPermissions", 7, 2).
		Return(TeamPermissions(false, false, false, false), nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 5, 2).
		Return(LeaguePermissions(false, false, true, false), nil)

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, nil, "PUT", "/7", 400, testParams{Error: "malformedInput"})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao, mockTeamsDao)
}

func testUpdateTeamNameTooLong(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(2, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 5, 7).Return(true, nil)
	mockTeamsDao.On("GetTeamPermissions", 7, 2).
		Return(TeamPermissions(false, false, false, false), nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 5, 2).
		Return(LeaguePermissions(false, false, true, false), nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, createTeamRequestBody("123456789012345678901234567890123456789012345678901", "TAG"),
		"PUT", "/7", 400, testParams{Error: "nameTooLong"})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao, mockLeaguesDao)
}

func testUpdateTeamTagTooLong(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(2, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 5, 7).Return(true, nil)
	mockTeamsDao.On("GetTeamPermissions", 7, 2).
		Return(TeamPermissions(false, false, false, false), nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 5, 2).
		Return(LeaguePermissions(false, false, true, false), nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, createTeamRequestBody("12345678901234567890123456789012345678901234567890", "123456"),
		"PUT", "/7", 400, testParams{Error: "tagTooLong"})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao, mockLeaguesDao)
}

func testUpdateTeamNoActiveLeague(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(-1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "PUT", "/7", 403, testParams{Error: "noActiveLeague"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testUpdateTeamDbError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(2, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 5, 7).Return(true, nil)
	mockTeamsDao.On("IsInfoInUse", 5, 7, "sampleName", "TAG").
		Return(false, "", nil)
	mockTeamsDao.On("UpdateTeam", 5, 7, "sampleName", "TAG", "").Return(errors.New("fake db error"))
	mockTeamsDao.On("GetTeamPermissions", 7, 2).
		Return(TeamPermissions(false, false, false, false), nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 5, 2).
		Return(LeaguePermissions(false, false, true, false), nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, createTeamRequestBody("sampleName", "TAG"),
		"PUT", "/7", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao, mockLeaguesDao)
}

func testUpdateTeamNameInUse(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(2, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 5, 7).Return(true, nil)
	mockTeamsDao.On("IsInfoInUse", 5, 7, "sampleName", "TAG").
		Return(true, "nameInUse", nil)
	mockTeamsDao.On("GetTeamPermissions", 7, 2).
		Return(TeamPermissions(false, false, false, false), nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 5, 2).
		Return(LeaguePermissions(false, false, true, false), nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, createTeamRequestBody("sampleName", "TAG"),
		"PUT", "/7", 400, testParams{Error: "nameInUse"})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao, mockLeaguesDao)
}
func testUpdateTeamTagInUse(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(2, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 5, 7).Return(true, nil)
	mockTeamsDao.On("IsInfoInUse", 5, 7, "sampleName", "TAG").
		Return(true, "tagInUse", nil)
	mockTeamsDao.On("GetTeamPermissions", 7, 2).
		Return(TeamPermissions(false, false, false, false), nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 5, 2).
		Return(LeaguePermissions(false, false, true, false), nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, createTeamRequestBody("sampleName", "TAG"),
		"PUT", "/7", 400, testParams{Error: "tagInUse"})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao, mockLeaguesDao)
}

func testUpdateTeam(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(2, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 5, 7).Return(true, nil)
	mockTeamsDao.On("IsInfoInUse", 5, 7, "sampleName", "TAG").
		Return(false, "", nil)
	mockTeamsDao.On("UpdateTeam", 5, 7, "sampleName", "TAG", "").Return(nil)
	mockTeamsDao.On("GetTeamPermissions", 7, 2).
		Return(TeamPermissions(false, false, false, false), nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 5, 2).
		Return(LeaguePermissions(false, false, true, false), nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, createTeamRequestBody("sampleName", "TAG"),
		"PUT", "/7", 200, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao, mockLeaguesDao)
}

func Test_UpdateTeam(t *testing.T) {
	//set up router and path to test
	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
	router = gin.New()

	router.Use(routes.Testing_Export_getActiveLeague())
	router.PUT("/:id",
		routes.Testing_Export_getUrlId(),
		routes.Testing_Export_authenticate(),
		routes.Testing_Export_failIfCanNotEditTeamInformation(),
		routes.Testing_Export_updateTeam)

	t.Run("NoId", testUpdateTeamInformationNoId)
	t.Run("IdNotInt", testUpdateTeamInformationNotInt)
	t.Run("sessionsError", testUpdateTeamSessionError)
	t.Run("notLoggedIn", testUpdateTeamNotLoggedIn)
	t.Run("noTeamEditPermissions", testUpdateTeamNoEditPermissions)
	t.Run("TeamDoesNotExist", testUpdateTeamInformationTeamDoesNotExist)
	t.Run("malformedBody", testUpdateTeamMalformedBody)
	t.Run("teamNameTooLong", testUpdateTeamNameTooLong)
	t.Run("teamTagTooLong", testUpdateTeamTagTooLong)
	t.Run("noActiveLeague", testUpdateTeamNoActiveLeague)
	t.Run("dbError", testUpdateTeamDbError)
	t.Run("nameInUse", testUpdateTeamNameInUse)
	t.Run("tagInUse", testUpdateTeamTagInUse)
	t.Run("correctUpdate", testUpdateTeam)
}
