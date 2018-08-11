package routesTest

import (
	"bytes"
	"encoding/json"
	"errors"
	"esports-league-manager/Backend/Server/routes"
	"esports-league-manager/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"testing"
	"esports-league-manager/Backend/Server/databaseAccess"
)

func createAddPlayerRequestBody(teamId int, name, gameIdentifier string, mainRoster bool) *bytes.Buffer {
	body := routes.PlayerInformation{
		TeamId: teamId,
		Name: name,
		GameIdentifier: gameIdentifier,
		MainRoster: mainRoster,
	}
	bodyB, _ := json.Marshal(&body)
	return bytes.NewBuffer(bodyB)
}

type addPlayerRes struct {
	Id int `json:"id"`
}

func createAddPlayerResponseBody(id int) *bytes.Buffer {
	resBody := addPlayerRes{
		Id: id,
	}
	resBodyB, _ := json.Marshal(&resBody)
	return bytes.NewBuffer(resBodyB)
}

func testAddPlayerToTeamNoActiveLeague(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(-1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "POST", "/addPlayer", 403, testParams{Error: "noActiveLeague"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testAddPlayerToTeamSessionError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, errors.New("session error"))

	routes.ElmSessions = mockSession

	httpTest(t, nil, "POST", "/addPlayer", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testAddPlayerToTeamNotLoggedIn(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(1, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(-1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "POST", "/addPlayer", 403, testParams{Error: "notLoggedIn"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testAddPlayerToTeamMalformedBody(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(4, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "POST", "/addPlayer", 400, testParams{Error: "malformedInput"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testAddPlayerToTeamTeamDoesNotExist(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(4, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 1, 5).
		Return(false, nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, createAddPlayerRequestBody(1, "test name", "inGameTestName", true),
		"POST", "/addPlayer", 400, testParams{Error: "teamDoesNotExist"})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao)
}

func testAddPlayerToTeamCannotEditPlayersOnTeam(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(4, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 1, 5).
		Return(true, nil)
	mockTeamsDao.On("HasPlayerEditPermissions", 1, 4, 5).
		Return(false, nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, createAddPlayerRequestBody(1, "test name", "inGameTestName", true),
		"POST", "/addPlayer", 403, testParams{Error: "canNotEditPlayers"})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao)
}

func testAddPlayerToTeamGameIdentifierTooLong(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(4, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 1, 5).
		Return(true, nil)
	mockTeamsDao.On("HasPlayerEditPermissions", 1, 4, 5).
		Return(true, nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, createAddPlayerRequestBody(1, "12345678901234567890123456789012345678901234567890",
		"123456789012345678901234567890123456789012345678901", true),
		"POST", "/addPlayer", 400, testParams{Error: "gameIdentifierTooLong"})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao)
}

func testAddPlayerToTeamNameTooLong(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(4, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 1, 5).
		Return(true, nil)
	mockTeamsDao.On("HasPlayerEditPermissions", 1, 4, 5).
		Return(true, nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, createAddPlayerRequestBody(1, "123456789012345678901234567890123456789012345678901",
		"12345678901234567890123456789012345678901234567890", true),
		"POST", "/addPlayer", 400, testParams{Error: "nameTooLong"})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao)
}

func testAddPlayerToTeamGameIdentifierInUse(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(4, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 1, 5).
		Return(true, nil)
	mockTeamsDao.On("HasPlayerEditPermissions", 1, 4, 5).
		Return(true, nil)
	mockTeamsDao.On("GetTeamInformation", 1, 5).
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
				GameIdentifier: "inGameTestName",
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

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, createAddPlayerRequestBody(1, "Test Player1", "inGameTestName", true),
		"POST", "/addPlayer", 400, testParams{Error: "gameIdentifierInUse"})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao)
}

func testAddPlayerToTeamDatabaseError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(4, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 1, 5).
		Return(true, nil)
	mockTeamsDao.On("HasPlayerEditPermissions", 1, 4, 5).
		Return(true, nil)
	mockTeamsDao.On("GetTeamInformation", 1, 5).
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
	mockTeamsDao.On("AddNewPlayer", 1, "inGameTestName", "Test Player1", true).
		Return(7, errors.New("fake db error"))

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, createAddPlayerRequestBody(1, "Test Player1", "inGameTestName", true),
		"POST", "/addPlayer", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao)
}

func testAddPlayerToTeamCorrectAddPlayer(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(4, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 1, 5).
		Return(true, nil)
	mockTeamsDao.On("HasPlayerEditPermissions", 1, 4, 5).
		Return(true, nil)
	mockTeamsDao.On("GetTeamInformation", 1, 5).
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
	mockTeamsDao.On("AddNewPlayer", 1, "inGameTestName", "Test Player1", true).
		Return(7, nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, createAddPlayerRequestBody(1, "Test Player1", "inGameTestName", true),
		"POST", "/addPlayer", 200, testParams{ResponseBody: createAddPlayerResponseBody(7)})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao)
}

func Test_AddPlayerToTeam(t *testing.T) {
	//set up router and path to test
	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
	router = gin.New()

	router.Use(routes.Testing_Export_getActiveLeague())
	router.POST("/addPlayer",
		routes.Testing_Export_authenticate(),
		routes.Testing_Export_addPlayerToTeam)

	t.Run("NoActiveLeague", testAddPlayerToTeamNoActiveLeague)
	t.Run("SessionsError", testAddPlayerToTeamSessionError)
	t.Run("NotLoggedIn", testAddPlayerToTeamNotLoggedIn)
	t.Run("MalformedBody", testAddPlayerToTeamMalformedBody)
	t.Run("TeamDoesNotExist", testAddPlayerToTeamTeamDoesNotExist)
	t.Run("CannotEditPlayersOnTeam", testAddPlayerToTeamCannotEditPlayersOnTeam)
	t.Run("GameIdentifierTooLong", testAddPlayerToTeamGameIdentifierTooLong)
	t.Run("NameTooLong", testAddPlayerToTeamNameTooLong)
	t.Run("GameIdentifierInUse", testAddPlayerToTeamGameIdentifierInUse)
	t.Run("DatabaseError", testAddPlayerToTeamDatabaseError)
	t.Run("CorrectAddPlayer", testAddPlayerToTeamCorrectAddPlayer)
}

