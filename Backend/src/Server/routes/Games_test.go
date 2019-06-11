package routes

import (
	"Server/databaseAccess"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"mocks"
	"net/http"
	"strconv"
	"testing"
)

type games struct {
	suite.Suite
	mockSession   *mocks.SessionManager
	mockValidator *mocks.Validator
	mockAccess    *mocks.Access
	mockGamesDao  *mocks.GamesDAO

	gameCreationInformation databaseAccess.GameCreationInformation
	game                    *databaseAccess.Game
	games                   []*databaseAccess.Game
	gameId                  gameIdResponse

	rescheduleInformation databaseAccess.GameTime
	reportInformation     databaseAccess.GameResult
}

func (s *games) SetupSuite() {
	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
	router = gin.New()
	router.Use(Authenticate())
	RegisterGameHandlers(router.Group("/"))

	team1Id := 1001
	team2Id := 1002
	gameTime := 1000000

	s.gameCreationInformation = databaseAccess.GameCreationInformation{
		Team1Id:  team1Id,
		Team2Id:  team2Id,
		GameTime: gameTime,
	}

	s.gameId = gameIdResponse{
		GameId: testGameId,
	}

	s.game = &databaseAccess.Game{
		GameId:   testGameId,
		GameTime: gameTime,
		Team1: databaseAccess.TeamDisplay{
			TeamId:    1001,
			Name:      "team1",
			Tag:       "team1",
			IconSmall: "icon1.png",
		},
		Team2: databaseAccess.TeamDisplay{
			TeamId:    1002,
			Name:      "team2",
			Tag:       "team2",
			IconSmall: "icon2.png",
		},
		WinnerId:   0,
		LoserId:    0,
		ScoreTeam1: 0,
		ScoreTeam2: 0,
		Complete:   false,
	}
	s.games = append(s.games, s.game)

	s.rescheduleInformation = databaseAccess.GameTime{
		GameTime: 14,
	}

	s.reportInformation = databaseAccess.GameResult{
		WinnerId:   1001,
		LoserId:    1002,
		ScoreTeam1: 2,
		ScoreTeam2: 1,
	}
}

// Mock Management
func (s *games) setSessionSuccess() {
	s.mockSession = new(mocks.SessionManager)
	s.mockSession.On("AuthenticateAndGetUserId", mock.Anything).Return(testUserId, nil)
	s.mockSession.On("GetActiveLeague", mock.Anything).Return(testLeagueId, nil)
	ElmSessions = s.mockSession
}

func (s *games) setSessionFail() {
	s.mockSession = new(mocks.SessionManager)
	s.mockSession.On("AuthenticateAndGetUserId", mock.Anything).Return(testLeagueId, testSessionErr)
	ElmSessions = s.mockSession
}

func (s *games) setDataValid() {
	s.mockValidator = new(mocks.Validator)
	s.mockValidator.On("DataInvalid", mock.Anything, mock.Anything).Return(false)
	validator = s.mockValidator
}

func (s *games) setDataInvalid() {
	s.mockValidator = new(mocks.Validator)
	s.mockValidator.On("DataInvalid", mock.Anything, mock.Anything).Return(true).
		Run(func(args mock.Arguments) {
			ctx := args.Get(0).(*gin.Context)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": testValidatorErrMsg})
		})
	validator = s.mockValidator
}

func (s *games) setGameCreateAccessTrue() {
	s.mockAccess = new(mocks.Access)
	s.mockAccess.On("Game", databaseAccess.Create, testLeagueId, 0, testUserId).Return(true, nil)
	Access = s.mockAccess
}

func (s *games) setGameEditAccessTrue() {
	s.mockAccess = new(mocks.Access)
	s.mockAccess.On("Game", databaseAccess.Edit, testLeagueId, testGameId, testUserId).Return(true, nil)
	Access = s.mockAccess
}

func (s *games) setGameViewAccessTrue() {
	s.mockAccess = new(mocks.Access)
	s.mockAccess.On("Game", databaseAccess.View, testLeagueId, testGameId, testUserId).Return(true, nil)
	Access = s.mockAccess
}

func (s *games) setGameDeleteAccessTrue() {
	s.mockAccess = new(mocks.Access)
	s.mockAccess.On("Game", databaseAccess.Delete, testLeagueId, testGameId, testUserId).Return(true, nil)
	Access = s.mockAccess
}

func (s *games) setGameAccessFalse() {
	s.mockAccess = new(mocks.Access)
	s.mockAccess.On("Game", mock.Anything, testLeagueId, mock.AnythingOfType("int"), testUserId).Return(false, nil)
	Access = s.mockAccess
}

func (s *games) setGameAccessFail() {
	s.mockAccess = new(mocks.Access)
	s.mockAccess.On("Game", mock.Anything, testLeagueId, mock.AnythingOfType("int"), testUserId).Return(false, testPermissionsErr)
	Access = s.mockAccess
}

func (s *games) setCreateGameSuccess() {
	s.mockGamesDao = new(mocks.GamesDAO)
	s.mockGamesDao.On("CreateGame", testLeagueId, s.gameCreationInformation).Return(testGameId, nil)
	GamesDAO = s.mockGamesDao
}

func (s *games) setCreateGameFail() {
	s.mockGamesDao = new(mocks.GamesDAO)
	s.mockGamesDao.On("CreateGame", testLeagueId, s.gameCreationInformation).Return(-1, testDaoErr)
	GamesDAO = s.mockGamesDao
}

func (s *games) setGameInformationSuccess() {
	s.mockGamesDao = new(mocks.GamesDAO)
	s.mockGamesDao.On("GetGameInformation", testGameId).Return(s.game, nil)
	GamesDAO = s.mockGamesDao
}

func (s *games) setGameInformationFail() {
	s.mockGamesDao = new(mocks.GamesDAO)
	s.mockGamesDao.On("GetGameInformation", testGameId).Return(nil, testDaoErr)
	GamesDAO = s.mockGamesDao
}

func (s *games) setDeleteGameSuccess() {
	s.mockGamesDao = new(mocks.GamesDAO)
	s.mockGamesDao.On("DeleteGame", testGameId).Return(nil)
	GamesDAO = s.mockGamesDao
}

func (s *games) setDeleteGameFail() {
	s.mockGamesDao = new(mocks.GamesDAO)
	s.mockGamesDao.On("DeleteGame", testGameId).Return(testDaoErr)
	GamesDAO = s.mockGamesDao
}

func (s *games) setRescheduleGameSuccess() {
	s.mockGamesDao = new(mocks.GamesDAO)
	s.mockGamesDao.On("RescheduleGame", testGameId, s.rescheduleInformation.GameTime).Return(nil)
	GamesDAO = s.mockGamesDao
}

func (s *games) setRescheduleGameFail() {
	s.mockGamesDao = new(mocks.GamesDAO)
	s.mockGamesDao.On("RescheduleGame", testGameId, s.rescheduleInformation.GameTime).Return(testDaoErr)
	GamesDAO = s.mockGamesDao
}

func (s *games) setReportGameSuccess() {
	s.mockGamesDao = new(mocks.GamesDAO)
	s.mockGamesDao.On("ReportGame", testGameId, s.reportInformation).Return(nil)
	GamesDAO = s.mockGamesDao
}

func (s *games) setReportGameFail() {
	s.mockGamesDao = new(mocks.GamesDAO)
	s.mockGamesDao.On("ReportGame", testGameId, s.reportInformation).Return(testDaoErr)
	GamesDAO = s.mockGamesDao
}

func (s *games) setAllGamesSuccess() {
	s.mockGamesDao = new(mocks.GamesDAO)
	s.mockGamesDao.On("GetAllGamesInLeague", testLeagueId).Return(s.games, nil)
	GamesDAO = s.mockGamesDao
}

func (s *games) setAllGamesFail() {
	s.mockGamesDao = new(mocks.GamesDAO)
	s.mockGamesDao.On("GetAllGamesInLeague", testLeagueId).Return(nil, testDaoErr)
	GamesDAO = s.mockGamesDao
}

// Abstract Negative Test Functions
func (s *games) testSessionFail(baseName string, body interface{}, reqType, url string) {
	s.T().Run(baseName+": Session Fail", func(t *testing.T) {
		s.setSessionFail()
		httpTest{
			T:            t,
			RequestData:  body,
			Type:         reqType,
			Url:          url,
			ResponseCode: http.StatusInternalServerError,
		}.RunHttpTest()
	})
}

func (s *games) testUrlIdInvalid(baseName string, body interface{}, reqType, url string) {
	s.T().Run(baseName+": Invalid URL Parameter", func(t *testing.T) {
		httpTest{
			T:            t,
			RequestData:  body,
			Type:         reqType,
			Url:          url + "z",
			ResponseCode: http.StatusBadRequest,
			Error:        testInvalidUrlMsg,
		}.RunHttpTest()
	})
}

func (s *games) testAccessFalse(baseName string, body interface{}, reqType, url string) {
	s.T().Run(baseName+": Access Not Allowed", func(t *testing.T) {
		s.setGameAccessFalse()
		httpTest{
			T:            t,
			RequestData:  body,
			Type:         reqType,
			Url:          url,
			ResponseCode: http.StatusForbidden,
		}.RunHttpTest()
	})
}

func (s *games) testAccessFail(baseName string, body interface{}, reqType, url string) {
	s.T().Run(baseName+": Access Check Failed", func(t *testing.T) {
		s.setGameAccessFail()
		httpTest{
			T:            t,
			RequestData:  body,
			Type:         reqType,
			Url:          url,
			ResponseCode: http.StatusInternalServerError,
		}.RunHttpTest()
	})
}

func (s *games) testBindJsonFail(baseName string, _ interface{}, reqType, url string) {
	s.T().Run(baseName+": Bind Json Fail", func(t *testing.T) {
		httpTest{
			T:            t,
			RequestData:  "_}",
			Type:         reqType,
			Url:          url,
			ResponseCode: http.StatusBadRequest,
			Error:        testBindErrMsg,
		}.RunHttpTest()
	})
}

func (s *games) testDataInvalid(baseName string, body interface{}, reqType, url string) {
	s.T().Run(baseName+": Data Invalid", func(t *testing.T) {
		s.setDataInvalid()
		httpTest{
			T:            t,
			RequestData:  body,
			Type:         reqType,
			Url:          url,
			ResponseCode: http.StatusBadRequest,
			Error:        testValidatorErrMsg,
		}.RunHttpTest()
	})
}

// Test Cases
func (s *games) TestCreateNewGame() {
	baseName := "Create New Game"
	method := "POST"
	path := "/"
	setDefaults := func() {
		s.setSessionSuccess()
		s.setGameCreateAccessTrue()
		s.setDataValid()
		s.setCreateGameSuccess()
	}

	for _, negativeTest := range []func(baseName string, body interface{}, reqType, url string){
		s.testSessionFail,
		s.testAccessFalse,
		s.testAccessFail,
		s.testBindJsonFail,
		s.testDataInvalid,
	} {
		setDefaults()
		negativeTest(baseName, httpBody(s.gameCreationInformation), method, path)
	}

	s.T().Run(baseName+": Create Game Failure", func(t *testing.T) {
		setDefaults()
		s.setCreateGameFail()
		httpTest{
			T:            t,
			RequestData:  s.gameCreationInformation,
			Type:         method,
			Url:          path,
			ResponseCode: http.StatusInternalServerError,
		}.RunHttpTest()
	})

	s.T().Run(baseName+" Correctly", func(t *testing.T) {
		setDefaults()
		httpTest{
			T:            t,
			RequestData:  s.gameCreationInformation,
			ResponseData: s.gameId,
			Type:         method,
			Url:          path,
			ResponseCode: http.StatusCreated,
		}.RunHttpTest()
	})
}

func (s *games) TestGetGameInformation() {
	baseName := "Get Game Information"
	method := "GET"
	path := "/" + strconv.Itoa(testGameId)

	setDefaults := func() {
		s.setSessionSuccess()
		s.setGameViewAccessTrue()
		s.setGameInformationSuccess()
	}

	for _, negativeTest := range []func(baseName string, body interface{}, reqType, url string){
		s.testSessionFail,
		s.testUrlIdInvalid,
		s.testAccessFalse,
		s.testAccessFail,
	} {
		setDefaults()
		negativeTest(baseName, nil, method, path)
	}

	s.T().Run(baseName+": Get Game Information Failure", func(t *testing.T) {
		setDefaults()
		s.setGameInformationFail()
		httpTest{
			T:            t,
			Type:         method,
			Url:          path,
			ResponseCode: http.StatusInternalServerError,
		}.RunHttpTest()
	})

	s.T().Run(baseName+" Correctly", func(t *testing.T) {
		setDefaults()
		httpTest{
			T:            t,
			ResponseData: s.game,
			Type:         method,
			Url:          path,
			ResponseCode: http.StatusOK,
		}.RunHttpTest()
	})
}

func (s *games) TestDeleteGame() {
	baseName := "Delete Game"
	method := "DELETE"
	path := "/" + strconv.Itoa(testGameId)

	setDefaults := func() {
		s.setSessionSuccess()
		s.setGameDeleteAccessTrue()
		s.setDeleteGameSuccess()
	}

	for _, negativeTest := range []func(baseName string, body interface{}, reqType, url string){
		s.testSessionFail,
		s.testUrlIdInvalid,
		s.testAccessFalse,
		s.testAccessFail,
	} {
		setDefaults()
		negativeTest(baseName, nil, method, path)
	}

	s.T().Run(baseName+": Delete Game Failure", func(t *testing.T) {
		setDefaults()
		s.setDeleteGameFail()
		httpTest{
			T:            t,
			Type:         method,
			Url:          path,
			ResponseCode: http.StatusInternalServerError,
		}.RunHttpTest()
	})

	s.T().Run(baseName+" Correctly", func(t *testing.T) {
		setDefaults()
		httpTest{
			T:            t,
			Type:         method,
			Url:          path,
			ResponseCode: http.StatusOK,
		}.RunHttpTest()
	})
}

func (s *games) TestRescheduleGame() {
	baseName := "Reschedule Game"
	method := "POST"
	path := fmt.Sprintf("/%v/reschedule", strconv.Itoa(testGameId))

	setDefaults := func() {
		s.setSessionSuccess()
		s.setGameEditAccessTrue()
		s.setDataValid()
		s.setRescheduleGameSuccess()
	}

	for _, negativeTest := range []func(baseName string, body interface{}, reqType, url string){
		s.testSessionFail,
		//s.testUrlIdInvalid,
		s.testAccessFalse,
		s.testAccessFail,
		s.testBindJsonFail,
		s.testDataInvalid,
	} {
		setDefaults()
		negativeTest(baseName, httpBody(s.rescheduleInformation), method, path)
	}

	s.T().Run(baseName+": Reschedule Game Failure", func(t *testing.T) {
		setDefaults()
		s.setRescheduleGameFail()
		httpTest{
			T:            t,
			RequestData:  s.rescheduleInformation,
			Type:         method,
			Url:          path,
			ResponseCode: http.StatusInternalServerError,
		}.RunHttpTest()
	})

	s.T().Run(baseName+" Correctly", func(t *testing.T) {
		setDefaults()
		httpTest{
			T:            t,
			RequestData:  s.rescheduleInformation,
			Type:         method,
			Url:          path,
			ResponseCode: http.StatusOK,
		}.RunHttpTest()
	})
}

func (s *games) TestReportGame() {
	baseName := "Report Game"
	method := "POST"
	path := fmt.Sprintf("/%v/report", strconv.Itoa(testGameId))

	setDefaults := func() {
		s.setSessionSuccess()
		s.setGameEditAccessTrue()
		s.setDataValid()
		s.setReportGameSuccess()
	}

	for _, negativeTest := range []func(baseName string, body interface{}, reqType, url string){
		s.testSessionFail,
		//s.testUrlIdInvalid,
		s.testAccessFalse,
		s.testAccessFail,
		s.testBindJsonFail,
		s.testDataInvalid,
	} {
		setDefaults()
		negativeTest(baseName, httpBody(s.reportInformation), method, path)
	}

	s.T().Run(baseName+": Report Game Failure", func(t *testing.T) {
		setDefaults()
		s.setReportGameFail()
		httpTest{
			T:            t,
			RequestData:  s.reportInformation,
			Type:         method,
			Url:          path,
			ResponseCode: http.StatusInternalServerError,
		}.RunHttpTest()
	})

	s.T().Run(baseName+" Correctly", func(t *testing.T) {
		setDefaults()
		httpTest{
			T:            t,
			RequestData:  s.reportInformation,
			Type:         method,
			Url:          path,
			ResponseCode: http.StatusOK,
		}.RunHttpTest()
	})
}

func (s *games) TestGetAllGamesInLeague() {
	baseName := "Get All Games in League"
	method := "GET"
	path := "/"

	setDefaults := func() {
		s.setSessionSuccess()
		s.setAllGamesSuccess()
	}

	for _, negativeTest := range []func(baseName string, body interface{}, reqType, url string){
		s.testSessionFail,
	} {
		setDefaults()
		negativeTest(baseName, nil, method, path)
	}

	s.T().Run(baseName+": Get All Games Failure", func(t *testing.T) {
		setDefaults()
		s.setAllGamesFail()
		httpTest{
			T:            t,
			Type:         method,
			Url:          path,
			ResponseCode: http.StatusInternalServerError,
		}.RunHttpTest()
	})

	s.T().Run(baseName+" Correctly", func(t *testing.T) {
		setDefaults()
		httpTest{
			T:            t,
			ResponseData: s.games,
			Type:         method,
			Url:          path,
			ResponseCode: http.StatusOK,
		}.RunHttpTest()
	})
}

func TestGamesSuite(t *testing.T) {
	suite.Run(t, new(games))
}
