package routesTest

//
//import (
//	"Server/routes"
//	"bytes"
//	"encoding/json"
//	"errors"
//	"github.com/gin-gonic/gin"
//	"github.com/stretchr/testify/mock"
//	"mocks"
//	"testing"
//)
//
//func createGamesRequestBody(team1Id, team2Id, gameTime int) *bytes.Buffer {
//	reqBody := routes.GameInformation{
//		Team1Id:  team1Id,
//		Team2Id:  team2Id,
//		GameTime: gameTime,
//	}
//	reqBodyB, _ := json.Marshal(&reqBody)
//	return bytes.NewBuffer(reqBodyB)
//}
//
//type gamesRes struct {
//	Id int `json:"id"`
//}
//
//func createGamesResponseBody(id int) *bytes.Buffer {
//	resBody := gamesRes{
//		Id: id,
//	}
//	resBodyB, _ := json.Marshal(&resBody)
//	return bytes.NewBuffer(resBodyB)
//}
//
//func testCreateNewGameNotLoggedIn(t *testing.T) {
//	mockSession := new(mocks.SessionManager)
//	mockSession.On("GetActiveLeague", mock.Anything).
//		Return(2, nil)
//	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
//		Return(-1, nil)
//
//	routes.ElmSessions = mockSession
//
//	httpTest(t, createGamesRequestBody(1, 2, 1532913359),
//		"POST", "/", 403, testParams{Error: "notLoggedIn"})
//
//	mock.AssertExpectationsForObjects(t, mockSession)
//}
//
//func testCreateNewGameNoActiveLeague(t *testing.T) {
//	mockSession := new(mocks.SessionManager)
//	mockSession.On("GetActiveLeague", mock.Anything).
//		Return(-1, nil)
//
//	routes.ElmSessions = mockSession
//
//	httpTest(t, createGamesRequestBody(1, 2, 1532913359),
//		"POST", "/", 403, testParams{Error: "noActiveLeague"})
//
//	mock.AssertExpectationsForObjects(t, mockSession)
//}
//
//func testCreateNewGameNoEditSchedulePermissions(t *testing.T) {
//	mockSession := new(mocks.SessionManager)
//	mockSession.On("GetActiveLeague", mock.Anything).
//		Return(2, nil)
//	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
//		Return(1, nil)
//
//	mockLeaguesDao := new(mocks.LeaguesDAO)
//	mockLeaguesDao.On("GetLeaguePermissions", 2, 1).
//		Return(LeaguePermissions(false, false, false, false), nil)
//
//	routes.ElmSessions = mockSession
//	routes.LeaguesDAO = mockLeaguesDao
//
//	httpTest(t, createGamesRequestBody(1, 2, 1532913359),
//		"POST", "/", 403, testParams{Error: "noEditSchedulePermissions"})
//
//	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)
//}
//
//func testCreateNewGameTeamsAreSame(t *testing.T) {
//	mockSession := new(mocks.SessionManager)
//	mockSession.On("GetActiveLeague", mock.Anything).
//		Return(2, nil)
//	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
//		Return(1, nil)
//
//	routes.ElmSessions = mockSession
//
//	httpTest(t, createGamesRequestBody(5, 5, 1532913359),
//		"POST", "/", 400, testParams{Error: "teamsAreSame"})
//
//	mock.AssertExpectationsForObjects(t, mockSession)
//}
//
//func testCreateNewGameTeam1DoesNotExist(t *testing.T) {
//	mockSession := new(mocks.SessionManager)
//	mockSession.On("GetActiveLeague", mock.Anything).
//		Return(2, nil)
//	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
//		Return(1, nil)
//
//	mockTeamsDao := new(mocks.TeamsDAO)
//	mockTeamsDao.On("DoesTeamExistInLeague", 2, 1).Return(false, nil)
//
//	routes.ElmSessions = mockSession
//	routes.TeamsDAO = mockTeamsDao
//
//	httpTest(t, createGamesRequestBody(1, 2, 1532913359),
//		"POST", "/", 400, testParams{Error: "teamDoesNotExist"})
//
//	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao)
//}
//
//func testCreateNewGameTeam2DoesNotExist(t *testing.T) {
//	mockSession := new(mocks.SessionManager)
//	mockSession.On("GetActiveLeague", mock.Anything).
//		Return(2, nil)
//	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
//		Return(1, nil)
//
//	mockTeamsDao := new(mocks.TeamsDAO)
//	mockTeamsDao.On("DoesTeamExistInLeague", 2, 1).Return(true, nil)
//	mockTeamsDao.On("DoesTeamExistInLeague", 2, 3).Return(false, nil)
//
//	routes.ElmSessions = mockSession
//	routes.TeamsDAO = mockTeamsDao
//
//	httpTest(t, createGamesRequestBody(1, 3, 1532913359),
//		"POST", "/", 400, testParams{Error: "teamDoesNotExist"})
//
//	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao)
//}
//
//func testCreateNewGameSessionError(t *testing.T) {
//	mockSession := new(mocks.SessionManager)
//	mockSession.On("GetActiveLeague", mock.Anything).
//		Return(2, nil)
//	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
//		Return(1, errors.New("Fake Cookie Error"))
//
//	routes.ElmSessions = mockSession
//
//	httpTest(t, createGamesRequestBody(1, 3, 1532913359),
//		"POST", "/", 500, testParams{})
//
//	mock.AssertExpectationsForObjects(t, mockSession)
//}
//
//func testCreateNewGameDatabaseError(t *testing.T) {
//	mockSession := new(mocks.SessionManager)
//	mockSession.On("GetActiveLeague", mock.Anything).
//		Return(2, nil)
//	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
//		Return(1, nil)
//
//	mockTeamsDao := new(mocks.TeamsDAO)
//	mockTeamsDao.On("DoesTeamExistInLeague", 2, 1).Return(true, nil)
//	mockTeamsDao.On("DoesTeamExistInLeague", 2, 3).Return(true, errors.New("Fake database error"))
//
//	routes.ElmSessions = mockSession
//	routes.TeamsDAO = mockTeamsDao
//
//	httpTest(t, createGamesRequestBody(1, 3, 1532913359),
//		"POST", "/", 500, testParams{})
//
//	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao)
//}
//
//func testCreateNewGameConflictExists(t *testing.T) {
//	mockSession := new(mocks.SessionManager)
//	mockSession.On("GetActiveLeague", mock.Anything).
//		Return(2, nil)
//	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
//		Return(1, nil)
//
//	mockTeamsDao := new(mocks.TeamsDAO)
//	mockTeamsDao.On("DoesTeamExistInLeague", 2, 1).Return(true, nil)
//	mockTeamsDao.On("DoesTeamExistInLeague", 2, 3).Return(true, nil)
//
//	mockGamesDao := new(mocks.GamesDAO)
//	mockGamesDao.On("DoesExistConflict", 1, 3, 1532913359).Return(true, nil)
//
//	routes.ElmSessions = mockSession
//	routes.TeamsDAO = mockTeamsDao
//	routes.GamesDAO = mockGamesDao
//
//	httpTest(t, createGamesRequestBody(1, 3, 1532913359),
//		"POST", "/", 400, testParams{Error: "conflictExists"})
//
//	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao, mockGamesDao)
//}
//
//func testCreateNewGameCorrectDatabaseError(t *testing.T) {
//	mockSession := new(mocks.SessionManager)
//	mockSession.On("GetActiveLeague", mock.Anything).
//		Return(2, nil)
//	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
//		Return(1, nil)
//
//	mockTeamsDao := new(mocks.TeamsDAO)
//	mockTeamsDao.On("DoesTeamExistInLeague", 2, 1).Return(true, nil)
//	mockTeamsDao.On("DoesTeamExistInLeague", 2, 3).Return(true, nil)
//
//	mockGamesDao := new(mocks.GamesDAO)
//	mockGamesDao.On("DoesExistConflict", 1, 3, 1532913359).Return(false, nil)
//	mockGamesDao.On("CreateGame", 2, 1, 3, 1532913359, "").Return(14, errors.New("Fake db error"))
//
//	routes.ElmSessions = mockSession
//	routes.TeamsDAO = mockTeamsDao
//	routes.GamesDAO = mockGamesDao
//
//	httpTest(t, createGamesRequestBody(1, 3, 1532913359),
//		"POST", "/", 500, testParams{})
//
//	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao, mockGamesDao)
//}
//
//func testCreateNewGameCorrectCreation(t *testing.T) {
//	mockSession := new(mocks.SessionManager)
//	mockSession.On("GetActiveLeague", mock.Anything).
//		Return(2, nil)
//	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
//		Return(1, nil)
//
//	mockTeamsDao := new(mocks.TeamsDAO)
//	mockTeamsDao.On("DoesTeamExistInLeague", 2, 1).Return(true, nil)
//	mockTeamsDao.On("DoesTeamExistInLeague", 2, 3).Return(true, nil)
//
//	mockGamesDao := new(mocks.GamesDAO)
//	mockGamesDao.On("DoesExistConflict", 1, 3, 1532913359).Return(false, nil)
//	mockGamesDao.On("CreateGame", 2, 1, 3, 1532913359, "").Return(14, nil)
//
//	routes.ElmSessions = mockSession
//	routes.TeamsDAO = mockTeamsDao
//	routes.GamesDAO = mockGamesDao
//
//	httpTest(t, createGamesRequestBody(1, 3, 1532913359),
//		"POST", "/", 200, testParams{ResponseBody: createGamesResponseBody(14)})
//
//	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao, mockGamesDao)
//}
//
//func Test_CreateNewGame(t *testing.T) {
//	//set up router and path to test
//	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
//	router = gin.New()
//
//	router.Use(routes.Testing_Export_getActiveLeague())
//	router.POST("/",
//		routes.Testing_Export_authenticate(),
//		routes.Testing_Export_failIfNoEditSchedulePermissions(),
//		routes.Testing_Export_createNewGame)
//
//	t.Run("notLoggedIn", testCreateNewGameNotLoggedIn)
//	t.Run("noActiveLeague", testCreateNewGameNoActiveLeague)
//	t.Run("NoEditSchedulePermissions", testCreateNewGameNoEditSchedulePermissions)
//	mockLeaguesDao := new(mocks.LeaguesDAO)
//	mockLeaguesDao.On("GetLeaguePermissions", 2, 1).
//		Return(LeaguePermissions(false, false, false, true), nil)
//	routes.LeaguesDAO = mockLeaguesDao
//	t.Run("teamsAreSame", testCreateNewGameTeamsAreSame)
//	t.Run("team1DoesNotExist", testCreateNewGameTeam1DoesNotExist)
//	t.Run("team2DoesNotExist", testCreateNewGameTeam2DoesNotExist)
//	t.Run("sessionError", testCreateNewGameSessionError)
//	t.Run("databaseError", testCreateNewGameDatabaseError)
//	t.Run("conflictExists", testCreateNewGameConflictExists)
//	t.Run("correctDatabaseError", testCreateNewGameCorrectDatabaseError)
//	t.Run("correctCreation", testCreateNewGameCorrectCreation)
//}
