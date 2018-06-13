package routesTest

import (
	"testing"
	"github.com/gin-gonic/gin"
	"esports-league-manager/Backend/Server/routes"
	"bytes"
	"encoding/json"
	"errors"
)

type userProfile struct {
	Id int `json:"id"`
}

//set up mock session managers
type mockSessionManager struct {
	id int
	err error
}

func (s *mockSessionManager) AuthenticateAndGetUserID(ctx *gin.Context) (int, error) {
	return s.id, s.err
}

func (s *mockSessionManager) LogIn(ctx *gin.Context, userID int) error {
	return nil
}

func testGetProfileNotLoggedIn(t *testing.T) {
	routes.ElmSessions = &mockSessionManager{
		id: -1,
		err: nil,
	}
	responseCodeAndErrorJsonTest(t, new(bytes.Buffer), "notLoggedIn", "GET", 403)
}

func testGetProfileDatabaseError(t *testing.T) {
	routes.ElmSessions = &mockSessionManager{
		id: -1,
		err: errors.New("fake database error"),
	}
	responseCodeTest(t, new(bytes.Buffer), 500, "GET")
}

func testGetProfileCorrectly(t *testing.T) {
	routes.ElmSessions = &mockSessionManager{
		id: 1,
		err: nil,
	}
	resBodyB, _ := json.Marshal(userProfile{Id: 1})

	testResponseAndCode(t, new(bytes.Buffer), bytes.NewBuffer(resBodyB), "GET", 200)
}

func Test_GetProfile(t *testing.T) {
	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
	router = gin.New()
	router.GET("/", routes.Testing_Export_getProfile)

	t.Run("getProfileNotLoggedIn", testGetProfileNotLoggedIn)
	t.Run("getProfileDatabaseError", testGetProfileDatabaseError)
	t.Run("getProfileCorrectly", testGetProfileCorrectly)
}