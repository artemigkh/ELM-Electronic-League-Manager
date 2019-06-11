package routes

import (
	"Server/databaseAccess"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"mocks"
	"net/http"
)

type leagues struct {
	suite.Suite
	mockSession    *mocks.SessionManager
	mockValidator  *mocks.Validator
	mockAccess     *mocks.Access
	mockLeaguesDao *mocks.LeaguesDAO

	leagueCreationInformation databaseAccess.LeagueCore
	markdown                  databaseAccess.Markdown
	league                    *databaseAccess.League
	leagues                   []*databaseAccess.League
	permissions               databaseAccess.LeaguePermissionsCore
}

func (s *leagues) SetupSuite() {
	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
	router = gin.New()
	router.Use(Authenticate())
	RegisterGameHandlers(router.Group("/"))

	name := "league name"
	description := "league description"
	game := "basketball"
	signupStart := 3501
	signupEnd := 3502
	leagueStart := 3701
	leagueEnd := 3702

	s.leagueCreationInformation = databaseAccess.LeagueCore{
		Name:        name,
		Description: description,
		Game:        game,
		PublicView:  true,
		PublicJoin:  true,
		SignupStart: signupStart,
		SignupEnd:   signupEnd,
		LeagueStart: leagueStart,
		LeagueEnd:   leagueEnd,
	}
}

// Mock Management
func (s *leagues) setSessionSuccess() {
	s.mockSession = new(mocks.SessionManager)
	s.mockSession.On("AuthenticateAndGetUserId", mock.Anything).Return(testUserId, nil)
	s.mockSession.On("GetActiveLeague", mock.Anything).Return(testLeagueId, nil)
	ElmSessions = s.mockSession
}

func (s *leagues) setSessionFail() {
	s.mockSession = new(mocks.SessionManager)
	s.mockSession.On("AuthenticateAndGetUserId", mock.Anything).Return(testLeagueId, testSessionErr)
	ElmSessions = s.mockSession
}

func (s *leagues) setDataValid() {
	s.mockValidator = new(mocks.Validator)
	s.mockValidator.On("DataInvalid", mock.Anything, mock.Anything).Return(false)
	validator = s.mockValidator
}

func (s *leagues) setDataInvalid() {
	s.mockValidator = new(mocks.Validator)
	s.mockValidator.On("DataInvalid", mock.Anything, mock.Anything).Return(true).
		Run(func(args mock.Arguments) {
			ctx := args.Get(0).(*gin.Context)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": testValidatorErrMsg})
		})
	validator = s.mockValidator
}

func (s *leagues) setLeagueCreateAccessTrue() {
	s.mockAccess = new(mocks.Access)
	s.mockAccess.On("League", databaseAccess.Create, 0, testUserId).Return(true, nil)
	Access = s.mockAccess
}

func (s *leagues) setLeagueEditAccessTrue() {
	s.mockAccess = new(mocks.Access)
	s.mockAccess.On("League", databaseAccess.Edit, testLeagueId, testUserId).Return(true, nil)
	Access = s.mockAccess
}

func (s *leagues) setLeagueViewAccessTrue() {
	s.mockAccess = new(mocks.Access)
	s.mockAccess.On("League", databaseAccess.View, testLeagueId, testUserId).Return(true, nil)
	Access = s.mockAccess
}

func (s *leagues) setLeagueDeleteAccessTrue() {
	s.mockAccess = new(mocks.Access)
	s.mockAccess.On("League", databaseAccess.Delete, testLeagueId, testUserId).Return(true, nil)
	Access = s.mockAccess
}

func (s *leagues) setLeagueAccessFalse() {
	s.mockAccess = new(mocks.Access)
	s.mockAccess.On("League", mock.Anything, mock.AnythingOfType("int"), testUserId).Return(false, nil)
	Access = s.mockAccess
}

func (s *leagues) setLeagueAccessFail() {
	s.mockAccess = new(mocks.Access)
	s.mockAccess.On("League", mock.Anything, mock.AnythingOfType("int"), testUserId).Return(false, testPermissionsErr)
	Access = s.mockAccess
}
