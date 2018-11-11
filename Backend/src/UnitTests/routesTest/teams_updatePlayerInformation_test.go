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

func createUpdatePlayerRequestBody(teamId, playerId int, name, gameIdentifier string, mainRoster bool) *bytes.Buffer {
	body := routes.PlayerInformationChange{
		TeamId:         teamId,
		PlayerId:       playerId,
		Name:           name,
		GameIdentifier: gameIdentifier,
		MainRoster:     mainRoster,
	}
	bodyB, _ := json.Marshal(&body)
	return bytes.NewBuffer(bodyB)
}

func testUpdatePlayerNoActiveLeague(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(-1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "POST", "/updatePlayer", 403, testParams{Error: "noActiveLeague"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testUpdatePlayerSessionError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, errors.New("session error"))

	routes.ElmSessions = mockSession

	httpTest(t, nil, "POST", "/updatePlayer", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testUpdatePlayerNotLoggedIn(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(1, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(-1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "POST", "/updatePlayer", 403, testParams{Error: "notLoggedIn"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testUpdatePlayerMalformedBody(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(4, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "POST", "/updatePlayer", 400, testParams{Error: "malformedInput"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testUpdatePlayerTeamDoesNotExist(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(4, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 5, 1).
		Return(false, nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, createUpdatePlayerRequestBody(1, 2, "test name", "inGameTestName", true),
		"POST", "/updatePlayer", 400, testParams{Error: "teamDoesNotExist"})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao)
}

func testUpdatePlayerCannotEditPlayersOnTeam(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(4, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 5, 1).
		Return(true, nil)
	mockTeamsDao.On("HasPlayerEditPermissions", 5, 1, 4).
		Return(false, nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, createUpdatePlayerRequestBody(1, 2, "test name", "inGameTestName", true),
		"POST", "/updatePlayer", 403, testParams{Error: "canNotEditPlayers"})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao)
}

func testUpdatePlayerGameIdentifierTooLong(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(4, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 5, 1).
		Return(true, nil)
	mockTeamsDao.On("HasPlayerEditPermissions", 5, 1, 4).
		Return(true, nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, createUpdatePlayerRequestBody(1, 2, "12345678901234567890123456789012345678901234567890",
		"123456789012345678901234567890123456789012345678901", true),
		"POST", "/updatePlayer", 400, testParams{Error: "gameIdentifierTooLong"})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao)
}

func testUpdatePlayerNameTooLong(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(4, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 5, 1).
		Return(true, nil)
	mockTeamsDao.On("HasPlayerEditPermissions", 5, 1, 4).
		Return(true, nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, createUpdatePlayerRequestBody(1, 2, "123456789012345678901234567890123456789012345678901",
		"12345678901234567890123456789012345678901234567890", true),
		"POST", "/updatePlayer", 400, testParams{Error: "nameTooLong"})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao)
}

func testUpdatePlayerGameIdentifierInUse(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(4, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 5, 1).
		Return(true, nil)
	mockTeamsDao.On("HasPlayerEditPermissions", 5, 1, 4).
		Return(true, nil)
	mockTeamsDao.On("GetTeamInformation", 5, 1).
		Return(&databaseAccess.TeamInformation{
			Name:   "sampleName",
			Tag:    "TAG",
			Wins:   10,
			Losses: 2,
			Players: []databaseAccess.PlayerInformation{
				{
					Id:             11,
					Name:           "Test Player1",
					GameIdentifier: "21",
					MainRoster:     true,
				},
				{
					Id:             12,
					Name:           "Test Player2",
					GameIdentifier: "inGameTestName",
					MainRoster:     true,
				},
				{
					Id:             13,
					Name:           "Test Player3",
					GameIdentifier: "41",
					MainRoster:     false,
				},
			},
		}, nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, createUpdatePlayerRequestBody(1, 2, "Test Player1", "inGameTestName", true),
		"POST", "/updatePlayer", 400, testParams{Error: "gameIdentifierInUse"})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao)
}

func testUpdatePlayerDatabaseError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(4, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 5, 1).
		Return(true, nil)
	mockTeamsDao.On("HasPlayerEditPermissions", 5, 1, 4).
		Return(true, nil)
	mockTeamsDao.On("GetTeamInformation", 5, 1).
		Return(&databaseAccess.TeamInformation{
			Name:   "sampleName",
			Tag:    "TAG",
			Wins:   10,
			Losses: 2,
			Players: []databaseAccess.PlayerInformation{
				{
					Id:             1,
					Name:           "Test Player1",
					GameIdentifier: "21",
					MainRoster:     true,
				},
				{
					Id:             2,
					Name:           "Test Player2",
					GameIdentifier: "37",
					MainRoster:     true,
				},
				{
					Id:             3,
					Name:           "Test Player3",
					GameIdentifier: "41",
					MainRoster:     false,
				},
			},
		}, nil)
	mockTeamsDao.On("UpdatePlayer", 1, 2, "inGameTestName", "Test Player1", true).
		Return(errors.New("fake db error"))

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, createUpdatePlayerRequestBody(1, 2, "Test Player1", "inGameTestName", true),
		"POST", "/updatePlayer", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao)
}

func testUpdatePlayerCorrectUpdatePlayer(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(4, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 5, 1).
		Return(true, nil)
	mockTeamsDao.On("HasPlayerEditPermissions", 5, 1, 4).
		Return(true, nil)
	mockTeamsDao.On("GetTeamInformation", 5, 1).
		Return(&databaseAccess.TeamInformation{
			Name:   "sampleName",
			Tag:    "TAG",
			Wins:   10,
			Losses: 2,
			Players: []databaseAccess.PlayerInformation{
				{
					Id:             1,
					Name:           "Test Player1",
					GameIdentifier: "21",
					MainRoster:     true,
				},
				{
					Id:             2,
					Name:           "Test Player2",
					GameIdentifier: "37",
					MainRoster:     true,
				},
				{
					Id:             3,
					Name:           "Test Player3",
					GameIdentifier: "41",
					MainRoster:     false,
				},
			},
		}, nil)
	mockTeamsDao.On("UpdatePlayer", 1, 2, "inGameTestName", "Test Player1", true).
		Return(nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, createUpdatePlayerRequestBody(1, 2, "Test Player1", "inGameTestName", true),
		"POST", "/updatePlayer", 200, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao)
}

func Test_UpdatePlayer(t *testing.T) {
	//set up router and path to test
	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
	router = gin.New()

	router.Use(routes.Testing_Export_getActiveLeague())
	router.POST("/updatePlayer",
		routes.Testing_Export_authenticate(),
		routes.Testing_Export_updatePlayer)

	t.Run("NoActiveLeague", testUpdatePlayerNoActiveLeague)
	t.Run("SessionsError", testUpdatePlayerSessionError)
	t.Run("NotLoggedIn", testUpdatePlayerNotLoggedIn)
	t.Run("MalformedBody", testUpdatePlayerMalformedBody)
	t.Run("TeamDoesNotExist", testUpdatePlayerTeamDoesNotExist)
	t.Run("CannotEditPlayersOnTeam", testUpdatePlayerCannotEditPlayersOnTeam)
	t.Run("GameIdentifierTooLong", testUpdatePlayerGameIdentifierTooLong)
	t.Run("NameTooLong", testUpdatePlayerNameTooLong)
	t.Run("GameIdentifierInUse", testUpdatePlayerGameIdentifierInUse)
	t.Run("DatabaseError", testUpdatePlayerDatabaseError)
	t.Run("CorrectUpdatePlayer", testUpdatePlayerCorrectUpdatePlayer)
}
