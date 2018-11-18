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

func createManagerCodeRequestBody(teamId int, administrator, information, players, reportResults bool) *bytes.Buffer {
	reqBody := routes.TeamManagerCodeRequest{
		TeamId:        teamId,
		Administrator: administrator,
		Information:   information,
		Players:       players,
		ReportResults: reportResults,
	}
	reqBodyB, _ := json.Marshal(&reqBody)
	return bytes.NewBuffer(reqBodyB)
}

type codeRes struct {
	Code string `json:"code"`
}

func createManagerCodeResponseBody(code string) *bytes.Buffer {
	resBody := codeRes{
		Code: code,
	}
	resBodyB, _ := json.Marshal(&resBody)
	return bytes.NewBuffer(resBodyB)
}

func testCreateManagerInviteCodeSessionsError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(5, errors.New("session error"))

	routes.ElmSessions = mockSession

	httpTest(t, nil,
		"POST", "/", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testCreateManagerInviteCodeNoActiveLeague(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(-1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil,
		"POST", "/", 403, testParams{Error: "noActiveLeague"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testCreateManagerInviteCodeNotLoggedIn(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(1, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(-1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil,
		"POST", "/", 403, testParams{Error: "notLoggedIn"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testCreateManagerInviteCodeMalformedBody(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(1, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil,
		"POST", "/", 400, testParams{Error: "malformedInput"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testCreateManagerInviteCodeTeamDoesNotExist(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(1, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 1, 2).Return(false, nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, createManagerCodeRequestBody(2, true, true, true, true),
		"POST", "/", 400, testParams{Error: "teamDoesNotExist"})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao)
}

func testCreateManagerInviteCodeNotAdmin(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(1, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(3, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 1, 2).Return(true, nil)
	mockTeamsDao.On("GetTeamPermissions", 2, 3).
		Return(TeamPermissions(false, true, true, true), nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 1, 3).
		Return(LeaguePermissions(false, false, false, false), nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, createManagerCodeRequestBody(2, true, true, true, true),
		"POST", "/", 403, testParams{Error: "notAdmin"})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao, mockLeaguesDao)
}

func testCreateManagerInviteCodeDatabaseError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(1, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(3, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 1, 2).Return(true, nil)
	mockTeamsDao.On("GetTeamPermissions", 2, 3).
		Return(TeamPermissions(true, true, true, true), nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 1, 3).
		Return(LeaguePermissions(false, false, false, false), nil)

	mockInviteCodesDao := new(mocks.InviteCodesDAO)
	mockInviteCodesDao.On("CreateTeamManagerInviteCode", 1, 2, true, true, true, true).
		Return("", errors.New("fake db error"))

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao
	routes.LeaguesDAO = mockLeaguesDao
	routes.InviteCodesDAO = mockInviteCodesDao

	httpTest(t, createManagerCodeRequestBody(2, true, true, true, true),
		"POST", "/", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao, mockLeaguesDao, mockInviteCodesDao)
}

func testCreateManagerInviteCodeCorrectCreate(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(1, nil)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(3, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 1, 2).Return(true, nil)
	mockTeamsDao.On("GetTeamPermissions", 2, 3).
		Return(TeamPermissions(true, true, true, true), nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("GetLeaguePermissions", 1, 3).
		Return(LeaguePermissions(false, false, false, false), nil)

	mockInviteCodesDao := new(mocks.InviteCodesDAO)
	mockInviteCodesDao.On("CreateTeamManagerInviteCode", 1, 2, true, true, true, true).
		Return("0123456789abcdef", nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao
	routes.LeaguesDAO = mockLeaguesDao
	routes.InviteCodesDAO = mockInviteCodesDao

	httpTest(t, createManagerCodeRequestBody(2, true, true, true, true),
		"POST", "/", 200, testParams{ResponseBody: createManagerCodeResponseBody("0123456789abcdef")})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao, mockLeaguesDao, mockInviteCodesDao)
}

func Test_CreateManagerInviteCode(t *testing.T) {
	//set up router and path to test
	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
	router = gin.New()
	router.Use(routes.Testing_Export_getActiveLeague())
	router.POST("/", routes.Testing_Export_authenticate(), routes.Testing_Export_createTeamManagerInviteCode)

	t.Run("SessionsError", testCreateManagerInviteCodeSessionsError)
	t.Run("NoActiveLeague", testCreateManagerInviteCodeNoActiveLeague)
	t.Run("NotLoggedIn", testCreateManagerInviteCodeNotLoggedIn)
	t.Run("MalformedBody", testCreateManagerInviteCodeMalformedBody)
	t.Run("TeamDoesNotExist", testCreateManagerInviteCodeTeamDoesNotExist)
	t.Run("NotAdmin", testCreateManagerInviteCodeNotAdmin)
	t.Run("DatabaseError", testCreateManagerInviteCodeDatabaseError)
	t.Run("CorrectCreate", testCreateManagerInviteCodeCorrectCreate)
}
