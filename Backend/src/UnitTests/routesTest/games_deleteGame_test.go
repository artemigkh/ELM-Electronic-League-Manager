package routesTest

import (
	"Server/databaseAccess"
	"Server/routes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"mocks"
	"testing"
)

func testDeleteGameNoId(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "DELETE", "/", 404, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testDeleteGameSessionError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(1, errors.New("fake session error"))

	routes.ElmSessions = mockSession

	httpTest(t, nil, "DELETE", "/1", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testDeleteGameNoActiveLeague(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(-1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "DELETE", "/1", 403, testParams{Error: "noActiveLeague"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testDeleteGameNotLoggedIn(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(-1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "DELETE", "/1", 403, testParams{Error: "notLoggedIn"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testDeleteGameIdNotInt(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(5, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "DELETE", "/a", 400, testParams{Error: "IdMustBeInteger"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testDeleteGameNoEditSchedulePermissions(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(5, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 2, 5).
		Return(LeaguePermissions(false, false, false, false), nil)

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, nil, "DELETE", "/1", 403, testParams{Error: "noEditSchedulePermissions"})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)
}

func testDeleteGameGameDoesNotExist(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(5, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 2, 5).
		Return(LeaguePermissions(false, false, false, true), nil)

	mockGamesDao := new(mocks.GamesDAO)
	mockGamesDao.On("GetGameInformation", 2, 16).
		Return(nil, nil)

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao
	routes.GamesDAO = mockGamesDao

	httpTest(t, nil, "DELETE", "/16", 400, testParams{Error: "gameDoesNotExist"})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao, mockGamesDao)
}

func testDeleteGameDbError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(5, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 2, 5).
		Return(LeaguePermissions(false, false, false, true), nil)

	mockGamesDao := new(mocks.GamesDAO)
	mockGamesDao.On("GetGameInformation", 2, 16).
		Return(&databaseAccess.GameDTO{}, nil)
	mockGamesDao.On("DeleteGame", 2, 16).
		Return(errors.New("fake db error"))

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao
	routes.GamesDAO = mockGamesDao

	httpTest(t, nil, "DELETE", "/16", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao, mockGamesDao)
}

func testDeleteGameCorrectDeleteGame(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(5, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 2, 5).
		Return(LeaguePermissions(false, false, false, true), nil)

	mockGamesDao := new(mocks.GamesDAO)
	mockGamesDao.On("GetGameInformation", 2, 16).
		Return(&databaseAccess.GameDTO{}, nil)
	mockGamesDao.On("DeleteGame", 2, 16).
		Return(nil)

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao
	routes.GamesDAO = mockGamesDao

	httpTest(t, nil, "DELETE", "/16", 200, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao, mockGamesDao)
}

func Test_DeleteGame(t *testing.T) {
	//set up router and path to test
	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
	router = gin.New()

	router.Use(routes.Testing_Export_getActiveLeague())
	router.DELETE("/:id",
		routes.Testing_Export_authenticate(),
		routes.Testing_Export_getUrlId(),
		routes.Testing_Export_failIfNoEditSchedulePermissions(),
		routes.Testing_Export_deleteGame)

	t.Run("NoId", testDeleteGameNoId)
	t.Run("SessionError", testDeleteGameSessionError)
	t.Run("NoActiveLeague", testDeleteGameNoActiveLeague)
	t.Run("NotLoggedIn", testDeleteGameNotLoggedIn)
	t.Run("IdNotInt", testDeleteGameIdNotInt)
	t.Run("NoEditSchedulePermissions", testDeleteGameNoEditSchedulePermissions)
	t.Run("GameDoesNotExist", testDeleteGameGameDoesNotExist)
	t.Run("DbError", testDeleteGameDbError)
	t.Run("CorrectDeleteGame", testDeleteGameCorrectDeleteGame)
}
