package routesTest

import (
	"bytes"
	"encoding/json"
	"errors"
	"esports-league-manager/Backend/Server/databaseAccess"
	"esports-league-manager/Backend/Server/routes"
	"esports-league-manager/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"testing"
)

func createRescheduleGameRequestBody(id, gameTime int) *bytes.Buffer {
	reqBody := routes.GameRescheduleInformation{
		Id:       id,
		GameTime: gameTime,
	}
	reqBodyB, _ := json.Marshal(&reqBody)
	return bytes.NewBuffer(reqBodyB)
}

func testRescheduleGameSessionError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(1, errors.New("fake session error"))

	routes.ElmSessions = mockSession

	httpTest(t, nil, "PUT", "/", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testRescheduleGameNoActiveLeague(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(-1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "PUT", "/", 403, testParams{Error: "noActiveLeague"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testRescheduleGameNotLoggedIn(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(-1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "PUT", "/", 403, testParams{Error: "notLoggedIn"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testRescheduleGameNoEditSchedulePermissions(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(5, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("HasEditSchedulePermission", 2, 5).Return(false, nil)

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, nil, "PUT", "/", 403, testParams{Error: "noEditSchedulePermissions"})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)

}

func testRescheduleGameGameDoesNotExist(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(5, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("HasEditSchedulePermission", 2, 5).Return(true, nil)

	mockGamesDao := new(mocks.GamesDAO)
	mockGamesDao.On("GetGameInformation", 2, 16).
		Return(nil, nil)

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao
	routes.GamesDAO = mockGamesDao

	httpTest(t, createRescheduleGameRequestBody(16, 1532913359), "PUT", "/", 400,
		testParams{Error: "gameDoesNotExist"})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao, mockGamesDao)
}

func testRescheduleGameGameComplete(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(5, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("HasEditSchedulePermission", 2, 5).Return(true, nil)

	mockGamesDao := new(mocks.GamesDAO)
	mockGamesDao.On("GetGameInformation", 2, 16).
		Return(&databaseAccess.GameInformation{Team1Id: 7, Team2Id: 8, Complete: true}, nil)

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao
	routes.GamesDAO = mockGamesDao

	httpTest(t, createRescheduleGameRequestBody(16, 1532913359), "PUT", "/", 400,
		testParams{Error: "gameIsComplete"})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao, mockGamesDao)
}

func testRescheduleGameConflictExists(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(5, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("HasEditSchedulePermission", 2, 5).Return(true, nil)

	mockGamesDao := new(mocks.GamesDAO)
	mockGamesDao.On("GetGameInformation", 2, 16).
		Return(&databaseAccess.GameInformation{Team1Id: 7, Team2Id: 8, Complete: false}, nil)
	mockGamesDao.On("DoesExistConflict", 7, 8, 1532913359).
		Return(true, nil)

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao
	routes.GamesDAO = mockGamesDao

	httpTest(t, createRescheduleGameRequestBody(16, 1532913359), "PUT", "/", 400,
		testParams{Error: "conflictExists"})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao, mockGamesDao)
}

func testRescheduleGameDbError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(5, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("HasEditSchedulePermission", 2, 5).Return(true, nil)

	mockGamesDao := new(mocks.GamesDAO)
	mockGamesDao.On("GetGameInformation", 2, 16).
		Return(&databaseAccess.GameInformation{Team1Id: 7, Team2Id: 8, Complete: false}, nil)
	mockGamesDao.On("DoesExistConflict", 7, 8, 1532913359).
		Return(false, nil)
	mockGamesDao.On("RescheduleGame", 2, 16, 1532913359).
		Return(errors.New("fake db error"))

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao
	routes.GamesDAO = mockGamesDao

	httpTest(t, createRescheduleGameRequestBody(16, 1532913359), "PUT", "/", 500,
		testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao, mockGamesDao)
}

func testRescheduleGameCorrectReschedule(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(5, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("HasEditSchedulePermission", 2, 5).Return(true, nil)

	mockGamesDao := new(mocks.GamesDAO)
	mockGamesDao.On("GetGameInformation", 2, 16).
		Return(&databaseAccess.GameInformation{Team1Id: 7, Team2Id: 8, Complete: false}, nil)
	mockGamesDao.On("DoesExistConflict", 7, 8, 1532913359).
		Return(false, nil)
	mockGamesDao.On("RescheduleGame", 2, 16, 1532913359).Return(nil)

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao
	routes.GamesDAO = mockGamesDao

	httpTest(t, createRescheduleGameRequestBody(16, 1532913359), "PUT", "/", 200,
		testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao, mockGamesDao)
}

func Test_RescheduleGame(t *testing.T) {
	//set up router and path to test
	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
	router = gin.New()

	router.Use(routes.Testing_Export_getActiveLeague())
	router.PUT("/",
		routes.Testing_Export_authenticate(),
		routes.Testing_Export_failIfNoEditSchedulePermissions(),
		routes.Testing_Export_rescheduleGame)

	t.Run("SessionError", testRescheduleGameSessionError)
	t.Run("NoActiveLeague", testRescheduleGameNoActiveLeague)
	t.Run("NotLoggedIn", testRescheduleGameNotLoggedIn)
	t.Run("NoEditSchedulePermissions", testRescheduleGameNoEditSchedulePermissions)
	t.Run("GameDoesNotExist", testRescheduleGameGameDoesNotExist)
	t.Run("GameComplete", testRescheduleGameGameComplete)
	t.Run("ConflictExists", testRescheduleGameConflictExists)
	t.Run("DbError", testRescheduleGameDbError)
	t.Run("CorrectReschedule", testRescheduleGameCorrectReschedule)
}
