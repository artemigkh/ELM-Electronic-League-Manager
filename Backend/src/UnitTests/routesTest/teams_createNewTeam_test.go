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

func createTeamRequestBody(name, tag string) *bytes.Buffer {
	body := routes.TeamInformation{
		Name: name,
		Tag:  tag,
	}
	bodyB, _ := json.Marshal(&body)
	return bytes.NewBuffer(bodyB)
}

type teamRes struct {
	Id int `json:"id"`
}

func createTeamResponseBody(id int) *bytes.Buffer {
	resBody := teamRes{
		Id: id,
	}
	resBodyB, _ := json.Marshal(&resBody)
	return bytes.NewBuffer(resBodyB)
}

func testCreateNewTeamMalformedBody(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(4, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 5, 4).
		Return(LeaguePermissions(false, true, false, false), nil)

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, nil, "POST", "/", 400, testParams{Error: "malformedInput"})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)

}

func testCreateNewTeamSessionError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, errors.New("session error"))

	routes.ElmSessions = mockSession

	httpTest(t, createTeamRequestBody("sampleName", "TAG"),
		"POST", "/", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testCreateNewTeamNotLoggedIn(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(1, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(-1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, createTeamRequestBody("sampleName", "TAG"),
		"POST", "/", 403, testParams{Error: "notLoggedIn"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testCreateNewTeamNoActiveLeague(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(-1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, createTeamRequestBody("sampleName", "TAG"),
		"POST", "/", 403, testParams{Error: "noActiveLeague"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testCreateNewTeamNoEditPermissions(t *testing.T) {
	mockSession := new(mocks.SessionManager)

	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(4, nil)
	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 5, 4).
		Return(LeaguePermissions(false, false, false, true), nil)

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, createTeamRequestBody("sampleName", "TAG"),
		"POST", "/", 403, testParams{Error: "noEditTeamPermissions"})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)
}

func testCreateNewTeamDbError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(4, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 5, 4).
		Return(nil, errors.New("fake db error"))

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, createTeamRequestBody("sampleName", "TAG"),
		"POST", "/", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)
}

func testCreateNewTeamNameTooLong(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(4, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 5, 4).
		Return(LeaguePermissions(false, true, false, false), nil)

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, createTeamRequestBody("123456789012345678901234567890123456789012345678901", "TAG"),
		"POST", "/", 400, testParams{Error: "nameTooLong"})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)
}

func testCreateNewTeamTagTooLong(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(4, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 5, 4).
		Return(LeaguePermissions(false, true, false, false), nil)

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, createTeamRequestBody("12345678901234567890123456789012345678901234567890", "123456"),
		"POST", "/", 400, testParams{Error: "tagTooLong"})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)
}

func testCreateNewTeamNameInUse(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(4, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 5, 4).
		Return(LeaguePermissions(false, true, false, false), nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("IsInfoInUse", 5, -1, "sampleName", "TAG").
		Return(true, "nameInUse", nil)

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, createTeamRequestBody("sampleName", "TAG"),
		"POST", "/", 400, testParams{Error: "nameInUse"})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao, mockTeamsDao)
}

func testCreateNewTeamTagInUse(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(4, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 5, 4).
		Return(LeaguePermissions(false, true, false, false), nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("IsInfoInUse", 5, -1, "sampleName", "TAG").
		Return(true, "tagInUse", nil)

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, createTeamRequestBody("sampleName", "TAG"),
		"POST", "/", 400, testParams{Error: "tagInUse"})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao, mockTeamsDao)
}

func testCorrectTeamCreation(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(4, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 5, 4).
		Return(LeaguePermissions(false, true, false, false), nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("IsInfoInUse", 5, -1, "sampleName", "TAG").
		Return(false, "", nil)
	mockTeamsDao.On("CreateTeam", 5, 4, "sampleName", "TAG").
		Return(6, nil)

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, createTeamRequestBody("sampleName", "TAG"),
		"POST", "/", 200, testParams{ResponseBody: createTeamResponseBody(6)})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao, mockTeamsDao)
}

func Test_CreateNewTeam(t *testing.T) {
	//set up router and path to test
	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
	router = gin.New()

	router.Use(routes.Testing_Export_getActiveLeague())
	router.POST("/",
		routes.Testing_Export_authenticate(),
		routes.Testing_Export_failIfNoTeamCreatePermissions(),
		routes.Testing_Export_createNewTeam)

	t.Run("malformedBody", testCreateNewTeamMalformedBody)
	t.Run("sessionsError", testCreateNewTeamSessionError)
	t.Run("notLoggedIn", testCreateNewTeamNotLoggedIn)
	t.Run("teamNameTooLong", testCreateNewTeamNameTooLong)
	t.Run("teamTagTooLong", testCreateNewTeamTagTooLong)
	t.Run("noActiveLeague", testCreateNewTeamNoActiveLeague)
	t.Run("noTeamEditPermissions", testCreateNewTeamNoEditPermissions)
	t.Run("dbError", testCreateNewTeamDbError)
	t.Run("nameInUse", testCreateNewTeamNameInUse)
	t.Run("tagInUse", testCreateNewTeamTagInUse)
	t.Run("correctTeamCreation", testCorrectTeamCreation)
}
