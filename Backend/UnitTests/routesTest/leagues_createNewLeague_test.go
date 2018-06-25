package routesTest

import (
	"bytes"
	"encoding/json"
	"esports-league-manager/Backend/Server/routes"
	"testing"
	"github.com/gin-gonic/gin"
	"esports-league-manager/Backend/Server/databaseAccess"
	"github.com/kataras/iris/core/errors"
)

func createLeagueRequestBody(name string, publicView, publicJoin bool) *bytes.Buffer {
	reqBody := routes.LeagueRequest{
		Name: name,
		PublicView: publicView,
		PublicJoin: publicJoin,
	}
	reqBodyB, _ := json.Marshal(&reqBody)
	return bytes.NewBuffer(reqBodyB)
}

type mockSessionsLeagueCreate struct {
	t *testing.T
	id int
	err error
}
func (s *mockSessionsLeagueCreate) AuthenticateAndGetUserID(ctx *gin.Context) (int, error) {
	return s.id, s.err
}
func (s *mockSessionsLeagueCreate) LogIn(ctx *gin.Context, userID int) error {
	return nil
}


type mockLeaguesDAO struct {
	t *testing.T
	id int
	err error
	inUse bool
	LeagueCreated bool
}
func (d *mockLeaguesDAO) CreateLeague(userID int, name string, publicView, publicJoin bool) (int, error) {
	if name != "testname" {
		d.t.Error("Did not get correct name from request. Should have gotten: testname")
	} else {
		d.LeagueCreated = true
	}
	return d.id, d.err
}
func (d *mockLeaguesDAO) IsNameInUse(name string) (bool, error) {
	return d.inUse, nil
}
func (d *mockLeaguesDAO) GetLeagueInformation(userID int) (*databaseAccess.LeagueInformation, error) {
	return nil, nil
}

func testCreateNewLeagueMalformedBody(t *testing.T) {
	responseCodeAndErrorJsonTest(t, new(bytes.Buffer), "malformedInput", "POST", 400)
}

func testCreateNewLeagueSessionError(t *testing.T) {
	routes.ElmSessions = &mockSessionsLeagueCreate{
		t: t,
		id: 1,
		err: errors.New("fake cookie error"),
	}
	responseCodeTest(t, createLeagueRequestBody("testname", true, true), 500, "POST")
}

func testCreateNewLeagueNotLoggedIn(t *testing.T) {
	routes.ElmSessions = &mockSessionsLeagueCreate{
		t: t,
		id: -1,
		err: nil,
	}
	responseCodeAndErrorJsonTest(t, createLeagueRequestBody("testname", true, true),
		"notLoggedIn","POST", 403)
}

func testCreateNewLeagueNameTooLong(t *testing.T) {
	routes.ElmSessions = &mockSessionsLeagueCreate{
		t: t,
		id: 1,
		err: nil,
	}
	responseCodeAndErrorJsonTest(t, createLeagueRequestBody("123456789012345678901234567890123456789012345678901",
		true, true), "nameTooLong", "POST", 400)
}

func testCreateNewLeagueNameInUse(t *testing.T) {
	routes.ElmSessions = &mockSessionsLeagueCreate{
		t: t,
		id: 1,
		err: nil,
	}
	routes.LeaguesDAO = &mockLeaguesDAO{
		t: t,
		id: 0,
		err: nil,
		inUse: true,
		LeagueCreated: false,
	}
	responseCodeAndErrorJsonTest(t, createLeagueRequestBody("12345678901234567890123456789012345678901234567890",
		true, true), "nameInUse", "POST", 400)
}

func testCreateNewLeagueDatabaseError(t *testing.T) {
	routes.ElmSessions = &mockSessionsLeagueCreate{
		t: t,
		id: 1,
		err: nil,
	}
	routes.LeaguesDAO = &mockLeaguesDAO{
		t: t,
		id: 0,
		err: errors.New("fake database error"),
		inUse: false,
		LeagueCreated: false,
	}
	responseCodeTest(t, createLeagueRequestBody("testname", true, true), 500, "POST")
}

func testCorrectLeagueCreation(t *testing.T) {
	routes.ElmSessions = &mockSessionsLeagueCreate{
		t: t,
		id: 1,
		err: nil,
	}
	mockDAO := &mockLeaguesDAO{
		t: t,
		id: 2,
		err: nil,
		inUse: false,
		LeagueCreated: false,
	}
	routes.LeaguesDAO = mockDAO
	res := responseCodeTest(t, createLeagueRequestBody("testname", true, true),
		200, "POST")

	if !mockDAO.LeagueCreated {
		t.Error("League creation DAO function was not called")
	}

	var idBody idResponse
	err := json.Unmarshal(res.Body.Bytes(), &idBody)
	if err != nil {
		t.Error("Response was not of a valid form")
	}

	if idBody.Id != 2 {
		t.Error("Did not receive correct ID in HTTP response")
	}
}

func Test_CreateNewLeague(t *testing.T) {
	//set up router and path to test
	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
	router = gin.New()
	router.POST("/", routes.Testing_Export_createNewLeague)

	t.Run("malformedBody", testCreateNewLeagueMalformedBody)
	t.Run("sessionsError", testCreateNewLeagueSessionError)
	t.Run("notLoggedIn", testCreateNewLeagueNotLoggedIn)
	t.Run("leagueNameTooLong", testCreateNewLeagueNameTooLong)
	t.Run("leagueNameInUse", testCreateNewLeagueNameInUse)
	t.Run("databaseError", testCreateNewLeagueDatabaseError)
	t.Run("correctLeagueCreation", testCorrectLeagueCreation)
}
