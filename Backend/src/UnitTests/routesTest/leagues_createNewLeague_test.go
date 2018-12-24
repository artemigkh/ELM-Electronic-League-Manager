package routesTest

import (
	"Server/routes"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Pallinder/go-randomdata"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"mocks"
	"testing"
)

func createLeagueRequestBody(name, description string, publicView, publicJoin bool,
	signupStart, signupEnd, leagueStart, leagueEnd int) *bytes.Buffer {

	reqBody := routes.LeagueRequest{
		Name:        name,
		Description: description,
		PublicView:  publicView,
		PublicJoin:  publicJoin,
		SignupStart: signupStart,
		SignupEnd:   signupEnd,
		LeagueStart: leagueStart,
		LeagueEnd:   leagueEnd,
	}
	reqBodyB, _ := json.Marshal(&reqBody)
	return bytes.NewBuffer(reqBodyB)
}

type leagueRes struct {
	Id int `json:"id"`
}

func createLeagueResponseBody(id int) *bytes.Buffer {
	resBody := leagueRes{
		Id: id,
	}
	resBodyB, _ := json.Marshal(&resBody)
	return bytes.NewBuffer(resBodyB)
}

func testCreateNewLeagueMalformedBody(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "POST", "/", 400, testParams{Error: "malformedInput"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testCreateNewLeagueSessionError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, errors.New("session error"))

	routes.ElmSessions = mockSession

	httpTest(t, createLeagueRequestBody("testname", "", true, true,
		1, 2, 3, 4),
		"POST", "/", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testCreateNewLeagueNotLoggedIn(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(-1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, createLeagueRequestBody("testname", "", true, true,
		1, 2, 3, 4),
		"POST", "/", 403, testParams{Error: "notLoggedIn"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testCreateNewLeagueDescriptionTooLong(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, createLeagueRequestBody("123456789012345678901234567890123456789012345678901",
		randomdata.RandStringRunes(501), true, true,
		1, 2, 3, 4),
		"POST", "/", 400, testParams{Error: "descriptionTooLong"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testCreateNewLeagueNameTooLong(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, createLeagueRequestBody("123456789012345678901234567890123456789012345678901",
		"", true, true, 1, 2, 3, 4),
		"POST", "/", 400, testParams{Error: "nameTooLong"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testCreateNewLeagueNameTooShort(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, createLeagueRequestBody("A", "", true, true,
		1, 2, 3, 4),
		"POST", "/", 400, testParams{Error: "nameTooShort"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testCreateNewLeagueNameInUse(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("IsNameInUse", -1, "12345678901234567890123456789012345678901234567890").
		Return(true, nil)

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, createLeagueRequestBody("12345678901234567890123456789012345678901234567890",
		"", true, true, 1, 2, 3, 4),
		"POST", "/", 400, testParams{Error: "nameInUse"})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)
}

func testCreateNewLeagueDatabaseError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("IsNameInUse", -1, "testName").
		Return(false, nil)
	mockLeaguesDao.On("CreateLeague", 1, "testName", "", true, true, 1, 2, 3, 4).
		Return(-1, errors.New("fake db error"))

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, createLeagueRequestBody("testName", "", true, true,
		1, 2, 3, 4),
		"POST", "/", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)
}

func testCorrectLeagueCreation(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, nil)

	mockLeaguesDao := new(mocks.LeaguesDAO)
	mockLeaguesDao.On("IsNameInUse", -1, "testName").
		Return(false, nil)
	mockLeaguesDao.On("CreateLeague", 1, "testName", mock.AnythingOfType("string"), true, true, 1, 2, 3, 4).
		Return(3, nil)

	routes.ElmSessions = mockSession
	routes.LeaguesDAO = mockLeaguesDao

	httpTest(t, createLeagueRequestBody("testName", randomdata.RandStringRunes(500), true, true,
		1, 2, 3, 4),
		"POST", "/", 200, testParams{ResponseBody: createLeagueResponseBody(3)})

	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)
}

func Test_CreateNewLeague(t *testing.T) {
	//set up router and path to test
	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
	router = gin.New()
	router.POST("/", routes.Testing_Export_authenticate(), routes.Testing_Export_createNewLeague)

	t.Run("malformedBody", testCreateNewLeagueMalformedBody)
	t.Run("sessionsError", testCreateNewLeagueSessionError)
	t.Run("notLoggedIn", testCreateNewLeagueNotLoggedIn)
	t.Run("descriptionTooLong", testCreateNewLeagueDescriptionTooLong)
	t.Run("leagueNameTooLong", testCreateNewLeagueNameTooLong)
	t.Run("leagueNameTooShort", testCreateNewLeagueNameTooShort)
	t.Run("leagueNameInUse", testCreateNewLeagueNameInUse)
	t.Run("databaseError", testCreateNewLeagueDatabaseError)
	t.Run("correctLeagueCreation", testCorrectLeagueCreation)
}
