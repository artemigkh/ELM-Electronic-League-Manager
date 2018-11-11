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

func createRemovePlayerRequestBody(teamId, playerId int) *bytes.Buffer {
	body := routes.PlayerRemoveInformation{
		TeamId:   teamId,
		PlayerId: playerId,
	}
	bodyB, _ := json.Marshal(&body)
	return bytes.NewBuffer(bodyB)
}

func testRemovePlayerFromTeamNoActiveLeague(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(-1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "DELETE", "/removePlayer", 403, testParams{Error: "noActiveLeague"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testRemovePlayerFromTeamSessionError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, errors.New("session error"))

	routes.ElmSessions = mockSession

	httpTest(t, nil, "DELETE", "/removePlayer", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testRemovePlayerFromTeamNotLoggedIn(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(1, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(-1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "DELETE", "/removePlayer", 403, testParams{Error: "notLoggedIn"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testRemovePlayerFromTeamMalformedBody(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(4, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "DELETE", "/removePlayer", 400, testParams{Error: "malformedInput"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testRemovePlayerFromTeamTeamDoesNotExist(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(4, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 5, 24).
		Return(false, nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, createRemovePlayerRequestBody(24, 31),
		"DELETE", "/removePlayer", 400, testParams{Error: "teamDoesNotExist"})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao)
}

func testRemovePlayerFromTeamCannotEditPlayersOnTeam(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(4, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 5, 24).
		Return(true, nil)
	mockTeamsDao.On("HasPlayerEditPermissions", 5, 24, 4).
		Return(false, nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, createRemovePlayerRequestBody(24, 31),
		"DELETE", "/removePlayer", 403, testParams{Error: "canNotEditPlayers"})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao)
}

func testRemovePlayerFromTeamPlayerDoesNotExist(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(4, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 5, 24).
		Return(true, nil)
	mockTeamsDao.On("HasPlayerEditPermissions", 5, 24, 4).
		Return(true, nil)
	mockTeamsDao.On("DoesPlayerExist", 24, 31).
		Return(false, nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, createRemovePlayerRequestBody(24, 31),
		"DELETE", "/removePlayer", 400, testParams{Error: "playerDoesNotExist"})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao)
}

func testRemovePlayerFromTeamDatabaseError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(4, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 5, 24).
		Return(true, nil)
	mockTeamsDao.On("HasPlayerEditPermissions", 5, 24, 4).
		Return(true, nil)
	mockTeamsDao.On("DoesPlayerExist", 24, 31).
		Return(true, nil)
	mockTeamsDao.On("RemovePlayer", 24, 31).
		Return(errors.New("fake db error"))

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, createRemovePlayerRequestBody(24, 31),
		"DELETE", "/removePlayer", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao)
}

func testRemovePlayerFromTeamCorrectRemovePlayer(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(4, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 5, 24).
		Return(true, nil)
	mockTeamsDao.On("HasPlayerEditPermissions", 5, 24, 4).
		Return(true, nil)
	mockTeamsDao.On("DoesPlayerExist", 24, 31).
		Return(true, nil)
	mockTeamsDao.On("RemovePlayer", 24, 31).
		Return(nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, createRemovePlayerRequestBody(24, 31),
		"DELETE", "/removePlayer", 200, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao)
}

func Test_RemovePlayerFromTeam(t *testing.T) {
	//set up router and path to test
	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
	router = gin.New()

	router.Use(routes.Testing_Export_getActiveLeague())
	router.DELETE("/removePlayer",
		routes.Testing_Export_authenticate(),
		routes.Testing_Export_removePlayerFromTeam)

	t.Run("NoActiveLeague", testRemovePlayerFromTeamNoActiveLeague)
	t.Run("SessionsError", testRemovePlayerFromTeamSessionError)
	t.Run("NotLoggedIn", testRemovePlayerFromTeamNotLoggedIn)
	t.Run("MalformedBody", testRemovePlayerFromTeamMalformedBody)
	t.Run("TeamDoesNotExist", testRemovePlayerFromTeamTeamDoesNotExist)
	t.Run("CannotEditPlayersOnTeam", testRemovePlayerFromTeamCannotEditPlayersOnTeam)
	t.Run("PlayerDoesNotExist", testRemovePlayerFromTeamPlayerDoesNotExist)
	t.Run("DatabaseError", testRemovePlayerFromTeamDatabaseError)
	t.Run("CorrectRemovePlayer", testRemovePlayerFromTeamCorrectRemovePlayer)
}
