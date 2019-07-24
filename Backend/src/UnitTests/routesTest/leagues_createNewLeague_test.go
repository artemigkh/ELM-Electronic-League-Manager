package routesTest

import (
	"Server/databaseAccess"
	"Server/routes"
	"bytes"
	"encoding/json"
	"github.com/Pallinder/go-randomdata"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"mocks"
	"testing"
)

func createLeagueRequestBody(name, description, game string, publicView, publicJoin bool,
	signupStart, signupEnd, leagueStart, leagueEnd int) *bytes.Buffer {

	reqBody := databaseAccess.LeagueDTO{
		Name:        name,
		Description: description,
		Game:        game,
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

type createNewLeague struct {
	suite.Suite
	mockSession    *mocks.SessionManager
	mockValidator  *mocks.DataValidator
	mockLeaguesDao *mocks.LeaguesDAO
}

func (s *createNewLeague) SetupSuite() {
	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
	router = gin.New()
	router.POST("/", routes.Testing_Export_authenticate(), routes.Testing_Export_createNewLeague)
}

func (s *createNewLeague) SetupTest() {
	s.mockSession = new(mocks.SessionManager)
	s.mockValidator = new(mocks.DataValidator)
	s.mockLeaguesDao = new(mocks.LeaguesDAO)

	routes.ElmSessions = s.mockSession
	routes.DataValidator = s.mockValidator
	routes.LeagueDAO = s.mockLeaguesDao
}
func (s *createNewLeague) SetDefaultMockValues() {
	s.mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(1, nil)
	s.mockValidator.On("ValidateLeagueDTO", mock.Anything).
		Return(true, "", nil)
	s.mockLeaguesDao.On("CreateLeague", 1, mock.Anything).
		Return(3, nil)
}

//func (s *createNewLeague) TearDownTest() {
//	mock.AssertExpectationsForObjects(s.T(), s.mockSession, s.mockValidator, s.mockLeaguesDao)
//}

func (s *createNewLeague) TestCreateNewLeagueMalformedBody() {
	httpTest(s.T(), nil, "POST", "/", 400, testParams{Error: "malformedInput"})
}

func (s *createNewLeague) TestCorrectLeagueCreation() {
	httpTest(s.T(), createLeagueRequestBody("testName", randomdata.RandStringRunes(500), "volleyball", true, true,
		1, 2, 3, 4),
		"POST", "/", 200, testParams{ResponseBody: createLeagueResponseBody(3)})
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(createNewLeague))
}

//
//func testCreateNewLeagueMalformedBody(t *testing.T) {
//	mockSession := new(mocks.SessionManager)
//	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
//		Return(1, nil)
//
//	routes.ElmSessions = mockSession
//
//	httpTest(t, nil, "POST", "/", 400, testParams{Error: "malformedInput"})
//
//	mock.AssertExpectationsForObjects(t, mockSession)
//}
//
//func testCreateNewLeagueSessionError(t *testing.T) {
//	mockSession := new(mocks.SessionManager)
//	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
//		Return(1, errors.New("session error"))
//
//	routes.ElmSessions = mockSession
//
//	httpTest(t, createLeagueRequestBody("testname", "", "genericsport", true, true,
//		1, 2, 3, 4),
//		"POST", "/", 500, testParams{})
//
//	mock.AssertExpectationsForObjects(t, mockSession)
//}
//
//func testCreateNewLeagueNotLoggedIn(t *testing.T) {
//	mockSession := new(mocks.SessionManager)
//	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
//		Return(-1, nil)
//
//	routes.ElmSessions = mockSession
//
//	httpTest(t, createLeagueRequestBody("testname", "", "genericsport", true, true,
//		1, 2, 3, 4),
//		"POST", "/", 403, testParams{Error: "notLoggedIn"})
//
//	mock.AssertExpectationsForObjects(t, mockSession)
//}
//
//func testCreateNewLeagueInvalidRequestData(t *testing.T) {
//	mockSession := new(mocks.SessionManager)
//	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
//		Return(1, nil)
//
//	mockValidator := new(mocks.DataValidator)
//	mockValidator.On("ValidateLeagueDTO", mock.Anything).
//		Return(false, "problemString", nil)
//
//	routes.ElmSessions = mockSession
//	routes.DataValidator = mockValidator
//
//	httpTest(t, createLeagueRequestBody("testname", "", "genericsport", true, true,
//		1, 2, 3, 4),
//		"POST", "/", 400, testParams{Error: "problemString"})
//
//	mock.AssertExpectationsForObjects(t, mockSession)
//}
//
//
//func testCreateNewLeagueDatabaseError(t *testing.T) {
//	mockSession := new(mocks.SessionManager)
//	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
//		Return(1, nil)
//
//	mockValidator := new(mocks.DataValidator)
//	mockValidator.On("ValidateLeagueDTO", mock.Anything).
//		Return(true, "", nil)
//
//	mockLeaguesDao := new(mocks.LeagueDAO)
//	mockLeaguesDao.On("CreateLeague", 1, mock.Anything).
//		Return(-1, errors.New("fake db error"))
//
//	routes.ElmSessions = mockSession
//	routes.DataValidator = mockValidator
//	routes.LeagueDAO = mockLeaguesDao
//
//	httpTest(t, createLeagueRequestBody("testName", "", "genericsport", true, true,
//		1, 2, 3, 4),
//		"POST", "/", 500, testParams{})
//
//	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)
//}
//
//func testCorrectLeagueCreation(t *testing.T) {
//	mockSession := new(mocks.SessionManager)
//	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
//		Return(1, nil)
//
//	mockValidator := new(mocks.DataValidator)
//	mockValidator.On("ValidateLeagueDTO", mock.Anything).
//		Return(true, "", nil)
//
//	mockLeaguesDao := new(mocks.LeagueDAO)
//	mockLeaguesDao.On("CreateLeague", 1, mock.Anything).
//		Return(3, nil)
//
//	routes.ElmSessions = mockSession
//	routes.DataValidator = mockValidator
//	routes.LeagueDAO = mockLeaguesDao
//
//	httpTest(t, createLeagueRequestBody("testName", randomdata.RandStringRunes(500), "volleyball", true, true,
//		1, 2, 3, 4),
//		"POST", "/", 200, testParams{ResponseBody: createLeagueResponseBody(3)})
//
//	mock.AssertExpectationsForObjects(t, mockSession, mockLeaguesDao)
//}
//
//func Test_CreateNewLeague(t *testing.T) {
//	//set up router and path to test
//	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
//	router = gin.New()
//	router.POST("/", routes.Testing_Export_authenticate(), routes.Testing_Export_createNewLeague)
//
//	t.Run("malformedBody", testCreateNewLeagueMalformedBody)
//	t.Run("sessionsError", testCreateNewLeagueSessionError)
//	t.Run("notLoggedIn", testCreateNewLeagueNotLoggedIn)
//	t.Run("invalidRequestData", testCreateNewLeagueInvalidRequestData)
//	t.Run("databaseError", testCreateNewLeagueDatabaseError)
//	t.Run("correctLeagueCreation", testCorrectLeagueCreation)
//}
