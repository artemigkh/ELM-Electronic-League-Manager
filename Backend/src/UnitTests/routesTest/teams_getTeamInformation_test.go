package routesTest

import (
	"Server/databaseAccess"
	"Server/routes"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"mocks"
	"testing"
)

func createTeamInfoBody(name, tag string, wins, losses int, players []databaseAccess.PlayerInformation) *bytes.Buffer {
	body := databaseAccess.TeamInformation{
		Name:    name,
		Tag:     tag,
		Wins:    wins,
		Losses:  losses,
		Players: players,
	}
	bodyB, _ := json.Marshal(&body)
	return bytes.NewBuffer(bodyB)
}

func testGetTeamInformationNoId(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "GET", "/", 404, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testGetTeamInformationNotInt(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "GET", "/a", 400, testParams{Error: "IdMustBeInteger"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testGetTeamInformationSessionError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(1, errors.New("fake session error"))

	routes.ElmSessions = mockSession

	httpTest(t, nil, "GET", "/1", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testGetTeamInformationNoActiveLeague(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(-1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "GET", "/1", 403, testParams{Error: "noActiveLeague"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testGetTeamInformationTeamDoesNotExist(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 2, 1).Return(false, nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, nil, "GET", "/1", 400, testParams{Error: "teamDoesNotExist"})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao)
}

func testGetTeamInformationDbError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 2, 1).Return(true, nil)
	mockTeamsDao.On("GetTeamInformation", 2, 1).
		Return(nil, errors.New("fake db error"))

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, nil, "GET", "/1", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao)
}

func testCorrectGetTeamInformationOneMember(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 2, 1).Return(true, nil)
	mockTeamsDao.On("GetTeamInformation", 2, 1).
		Return(&databaseAccess.TeamInformation{
			Name:   "sampleName",
			Tag:    "TAG",
			Wins:   10,
			Losses: 2,
			Players: []databaseAccess.PlayerInformation{
				{
					Id:             1,
					Name:           "Test Player1",
					GameIdentifier: "21",
					MainRoster:     true,
				},
			},
		}, nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, nil, "GET", "/1", 200,
		testParams{ResponseBody: createTeamInfoBody("sampleName", "TAG", 10, 2,
			[]databaseAccess.PlayerInformation{
				{
					Id:             1,
					Name:           "Test Player1",
					GameIdentifier: "21",
					MainRoster:     true,
				},
			})})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao)
}

func testCorrectGetTeamInformationManyPlayers(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 2, 1).Return(true, nil)
	mockTeamsDao.On("GetTeamInformation", 2, 1).
		Return(&databaseAccess.TeamInformation{
			Name:   "sampleName",
			Tag:    "TAG",
			Wins:   10,
			Losses: 2,
			Players: []databaseAccess.PlayerInformation{
				{
					Id:             1,
					Name:           "Test Player1",
					GameIdentifier: "21",
					MainRoster:     true,
				},
				{
					Id:             2,
					Name:           "Test Player2",
					GameIdentifier: "32",
					MainRoster:     true,
				},
				{
					Id:             3,
					Name:           "Test Player3",
					GameIdentifier: "41",
					MainRoster:     false,
				},
			},
		}, nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, nil, "GET", "/1", 200,
		testParams{ResponseBody: createTeamInfoBody("sampleName", "TAG", 10, 2,
			[]databaseAccess.PlayerInformation{
				{
					Id:             1,
					Name:           "Test Player1",
					GameIdentifier: "21",
					MainRoster:     true,
				},
				{
					Id:             2,
					Name:           "Test Player2",
					GameIdentifier: "32",
					MainRoster:     true,
				},
				{
					Id:             3,
					Name:           "Test Player3",
					GameIdentifier: "41",
					MainRoster:     false,
				},
			})})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao)
}

func Test_GetTeamInformation(t *testing.T) {
	//set up router and path to test
	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
	router = gin.New()

	router.Use(routes.Testing_Export_getActiveLeague())
	router.GET("/:id", routes.Testing_Export_getUrlId(), routes.Testing_Export_getTeamInformation)

	t.Run("noId", testGetTeamInformationNoId)
	t.Run("IdNotInt", testGetTeamInformationNotInt)
	t.Run("sessionError", testGetTeamInformationSessionError)
	t.Run("noActiveLeague", testGetTeamInformationNoActiveLeague)
	t.Run("teamDoesNotExist", testGetTeamInformationTeamDoesNotExist)
	t.Run("dbError", testGetTeamInformationDbError)
	t.Run("getOneMemberTeam", testCorrectGetTeamInformationOneMember)
	t.Run("getManyPlayersTeam", testCorrectGetTeamInformationManyPlayers)
}
