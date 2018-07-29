package routesTest

import (
	"testing"
	"github.com/gin-gonic/gin"
	"esports-league-manager/Backend/Server/routes"
	"bytes"
	"esports-league-manager/Backend/Server/databaseAccess"
	"encoding/json"
	"esports-league-manager/mocks"
	"github.com/stretchr/testify/mock"
	"errors"
)

func createTeamInfoBody(name, tag string, wins, losses int, members []databaseAccess.UserInformation) *bytes.Buffer {
	body := databaseAccess.TeamInformation{
		Name: name,
		Tag: tag,
		Wins: wins,
		Losses: losses,
		Members: members,
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
	mockTeamsDao.On("DoesTeamExist", 1, 2).Return(false, nil)

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
	mockTeamsDao.On("DoesTeamExist", 1, 2).Return(true, nil)
	mockTeamsDao.On("GetTeamInformation", 1, 2).
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
	mockTeamsDao.On("DoesTeamExist", 1, 2).Return(true, nil)
	mockTeamsDao.On("GetTeamInformation", 1, 2).
		Return(&databaseAccess.TeamInformation{
		Name: "sampleName",
		Tag: "TAG",
		Wins: 10,
		Losses: 2,
		Members: []databaseAccess.UserInformation{
			{
				Id: 1,
				Email: "test1@email.com",
			},
		},
	}, nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, nil, "GET", "/1", 200,
		testParams{ResponseBody: createTeamInfoBody("sampleName", "TAG", 10, 2,
			[]databaseAccess.UserInformation{
				{
					Id: 1,
					Email: "test1@email.com",
				},
		})})

	mock.AssertExpectationsForObjects(t, mockSession, mockTeamsDao)
}

func testCorrectGetTeamInformationManyMembers(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("GetActiveLeague", mock.Anything).
		Return(2, nil)

	mockTeamsDao := new(mocks.TeamsDAO)
	mockTeamsDao.On("DoesTeamExist", 1, 2).Return(true, nil)
	mockTeamsDao.On("GetTeamInformation", 1, 2).
		Return(&databaseAccess.TeamInformation{
		Name: "sampleName",
		Tag: "TAG",
		Wins: 10,
		Losses: 2,
		Members: []databaseAccess.UserInformation{
			{
				Id: 1,
				Email: "test1@email.com",
			},
			{
				Id: 5,
				Email: "test5@email.com",
			},
			{
				Id: 3,
				Email: "test3@email.com",
			},
		},
	}, nil)

	routes.ElmSessions = mockSession
	routes.TeamsDAO = mockTeamsDao

	httpTest(t, nil, "GET", "/1", 200,
		testParams{ResponseBody: createTeamInfoBody("sampleName", "TAG", 10, 2,
			[]databaseAccess.UserInformation{
				{
					Id: 1,
					Email: "test1@email.com",
				},
				{
					Id: 5,
					Email: "test5@email.com",
				},
				{
					Id: 3,
					Email: "test3@email.com",
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
	t.Run("getManyMembersTeam", testCorrectGetTeamInformationManyMembers)
}
