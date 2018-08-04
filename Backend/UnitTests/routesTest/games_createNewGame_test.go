package routesTest

import (
	"bytes"
	"esports-league-manager/Backend/Server/routes"
	"encoding/json"
	"testing"
	"github.com/gin-gonic/gin"
	"esports-league-manager/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/kataras/iris/core/errors"
)

func createGamesRequestBody(team1ID, team2ID, gameTime int) *bytes.Buffer {
	reqBody := routes.GameInformation{
		Team1ID: team1ID,
		Team2ID: team2ID,
		GameTime: gameTime,
	}
	reqBodyB, _ := json.Marshal(&reqBody)
	return bytes.NewBuffer(reqBodyB)
}

type gamesRes struct {
	Id int `json:"id"`
}

func createGamesResponseBody(id int) *bytes.Buffer {
	resBody := gamesRes{
		Id: id,
	}
	resBodyB, _ := json.Marshal(&resBody)
	return bytes.NewBuffer(resBodyB)
}

func testCreateNewGameNotLoggedIn(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)
	mockSession.On("AuthenticateAndGetUserID", mock.Anything).
		Return(-1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, createGamesRequestBody(1, 2, 1532913359),
		"POST", "/", 403, testParams{Error: "notLoggedIn"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testCreateNewGameNoActiveLeague(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(-1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, createGamesRequestBody(1, 2, 1532913359),
		"POST", "/", 403, testParams{Error: "noActiveLeague"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testCreateNewGameTeam1DoesNotExist(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)
	mockSession.On("AuthenticateAndGetUserID", mock.Anything).
		Return(1, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 1, 2).Return(false, nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, createGamesRequestBody(1, 2, 1532913359),
		"POST", "/", 400, testParams{Error: "teamDoesNotExist"})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao)
}

func testCreateNewGameTeam2DoesNotExist(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)
	mockSession.On("AuthenticateAndGetUserID", mock.Anything).
		Return(1, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 1, 2).Return(true, nil)
	mockTeamsDao.On("DoesTeamExist", 3, 2).Return(false, nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, createGamesRequestBody(1, 3, 1532913359),
		"POST", "/", 400, testParams{Error: "teamDoesNotExist"})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao)
}

func testCreateNewGameSessionError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)
	mockSession.On("AuthenticateAndGetUserID", mock.Anything).
		Return(1, errors.New("Fake Cookie Error"))

	routes.ElmSessions = mockSession

	httpTest(t, createGamesRequestBody(1, 3, 1532913359),
		"POST", "/", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testCreateNewGameDatabaseError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)
	mockSession.On("AuthenticateAndGetUserID", mock.Anything).
		Return(1, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 1, 2).Return(true, nil)
	mockTeamsDao.On("DoesTeamExist", 3, 2).Return(true, errors.New("Fake database error"))

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, createGamesRequestBody(1, 3, 1532913359),
		"POST", "/", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao)
}

func testCreateNewGameConflictExists(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)
	mockSession.On("AuthenticateAndGetUserID", mock.Anything).
		Return(1, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 1, 2).Return(true, nil)
	mockTeamsDao.On("DoesTeamExist", 3, 2).Return(true, nil)

	mockGamesDao := new(mocks.GamesDAO)
	mockGamesDao.On("DoesExistConflict", 1, 3, 1532913359).Return(true, nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao
	routes.GamesDAO = mockGamesDao

	httpTest(t, createGamesRequestBody(1, 3, 1532913359),
		"POST", "/", 400, testParams{Error: "conflictExists"})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao, mockGamesDao)
}

func testCreateNewGameCorrectDatabaseError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)
	mockSession.On("AuthenticateAndGetUserID", mock.Anything).
		Return(1, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 1, 2).Return(true, nil)
	mockTeamsDao.On("DoesTeamExist", 3, 2).Return(true, nil)

	mockGamesDao := new(mocks.GamesDAO)
	mockGamesDao.On("DoesExistConflict", 1, 3, 1532913359).Return(false, nil)
	mockGamesDao.On("CreateGame", 2, 1, 3, 1532913359).Return(14, errors.New("Fake db error"))

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao
	routes.GamesDAO = mockGamesDao

	httpTest(t, createGamesRequestBody(1, 3, 1532913359),
		"POST", "/", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao, mockGamesDao)
}

func testCreateNewGameCorrectCreation(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)
	mockSession.On("AuthenticateAndGetUserID", mock.Anything).
		Return(1, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 1, 2).Return(true, nil)
	mockTeamsDao.On("DoesTeamExist", 3, 2).Return(true, nil)

	mockGamesDao := new(mocks.GamesDAO)
	mockGamesDao.On("DoesExistConflict", 1, 3, 1532913359).Return(false, nil)
	mockGamesDao.On("CreateGame", 2, 1, 3, 1532913359).Return(14, nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao
	routes.GamesDAO = mockGamesDao

	httpTest(t, createGamesRequestBody(1, 3, 1532913359),
		"POST", "/", 200, testParams{ResponseBody:createGamesResponseBody(14)})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao, mockGamesDao)
}

func Test_CreateNewGame(t *testing.T) {
	//set up router and path to test
	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
	router = gin.New()

	router.Use(routes.Testing_Export_getActiveLeague())
	router.POST("/",
		routes.Testing_Export_authenticate(),
		routes.Testing_Export_createNewGame)

	t.Run("notLoggedIn", testCreateNewGameNotLoggedIn)
	t.Run("noActiveLeague", testCreateNewGameNoActiveLeague)
	t.Run("team1DoesNotExist", testCreateNewGameTeam1DoesNotExist)
	t.Run("team2DoesNotExist", testCreateNewGameTeam2DoesNotExist)
	t.Run("sessionError", testCreateNewGameSessionError)
	t.Run("databaseError", testCreateNewGameDatabaseError)
	t.Run("conflictExists", testCreateNewGameConflictExists)
	t.Run("correctDatabaseError", testCreateNewGameCorrectDatabaseError)
	t.Run("correctCreation", testCreateNewGameCorrectCreation)
}
