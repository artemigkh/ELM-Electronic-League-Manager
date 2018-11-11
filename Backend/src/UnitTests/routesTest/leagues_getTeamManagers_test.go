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

func createTeamManagersBody(body []databaseAccess.TeamManagerInformation) *bytes.Buffer {
	bodyB, _ := json.Marshal(&body)
	return bytes.NewBuffer(bodyB)
}

func testGetTeamManagersSessionError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, errors.New("session error"))

	routes.ElmSessions = mockSession

	httpTest(t, nil, "POST", "/", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testGetTeamManagersNotLoggedIn(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(-1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "POST", "/", 403, testParams{Error: "notLoggedIn"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testGetTeamManagersNoActiveLeague(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, nil)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(-1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "POST", "/", 403, testParams{Error: "noActiveLeague"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testGetTeamManagersNotAdmin(t *testing.T) {
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

func testGetTeamManagersDatabaseError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, nil)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 2, 1).
		Return(LeaguePermissions(true, true, true, true), nil)
	mockLeaguesDao.On("GetTeamManagerInformation", 2).
		Return(nil, errors.New("fake db error"))

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, nil, "POST", "/", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)
}

func testCorrectGetManagers(t *testing.T) {
	teamManagerSummary := []databaseAccess.TeamManagerInformation{
		{
			TeamId:   15,
			TeamName: "team1",
			TeamTag:  "ABC",
			Managers: []databaseAccess.ManagerInformation{
				{
					UserId:          3,
					UserEmail:       "testEmail3@test.com",
					EditPermissions: true,
					EditTeamInfo:    true,
					EditPlayers:     true,
					ReportResult:    true,
				}, {
					UserId:          4,
					UserEmail:       "testEmail4@test.com",
					EditPermissions: false,
					EditTeamInfo:    false,
					EditPlayers:     false,
					ReportResult:    true,
				},
			},
		}, {
			TeamId:   16,
			TeamName: "team2",
			TeamTag:  "ABD",
			Managers: []databaseAccess.ManagerInformation{
				{
					UserId:          6,
					UserEmail:       "testEmail3@test.com",
					EditPermissions: false,
					EditTeamInfo:    true,
					EditPlayers:     true,
					ReportResult:    true,
				},
			},
		},
	}

	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, nil)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 2, 1).
		Return(LeaguePermissions(true, true, true, true), nil)
	mockLeaguesDao.On("GetTeamManagerInformation", 2).
		Return(teamManagerSummary, nil)

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, nil, "POST", "/", 200,
		testParams{ResponseBody: createTeamManagersBody(teamManagerSummary)})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)
}

func Test_GetTeamManagers(t *testing.T) {
	//set up router and path to test
	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
	router = gin.New()
	router.POST("/",
		routes.Testing_Export_authenticate(),
		routes.Testing_Export_getActiveLeague(),
		routes.Testing_Export_failIfNotLeagueAdmin(),
		routes.Testing_Export_getTeamManagers)

	t.Run("SessionsError", testGetTeamManagersSessionError)
	t.Run("NotLoggedIn", testGetTeamManagersNotLoggedIn)
	t.Run("NoActiveLeague", testGetTeamManagersNoActiveLeague)
	t.Run("NotAdmin", testGetTeamManagersNotAdmin)
	t.Run("DatabaseError", testGetTeamManagersDatabaseError)
	t.Run("CorrectGetManagers", testCorrectGetManagers)
}
