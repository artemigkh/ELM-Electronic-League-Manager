package routesTest

import (
	"testing"
	"github.com/gin-gonic/gin"
	"esports-league-manager/Backend/Server/routes"
	"errors"
	"github.com/stretchr/testify/mock"
	"esports-league-manager/mocks"
	"bytes"
	"encoding/json"
)

type userProfile struct {
	Id int `json:"id"`
}

func createUserProfileBody(id int) *bytes.Buffer {
	body := userProfile{Id: id}
	reqBodyB, _ := json.Marshal(&body)
	return bytes.NewBuffer(reqBodyB)
}

func testGetProfileNotLoggedIn(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserID", mock.Anything).
		Return(-1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "GET", "/", 403, testParams{Error: "notLoggedIn"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testGetProfileSessionError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserID", mock.Anything).
		Return(-1, errors.New("fake session error"))

	routes.ElmSessions = mockSession

	httpTest(t, nil, "GET", "/", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testGetProfileCorrectly(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserID", mock.Anything).
		Return(14, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil, "GET", "/", 200, testParams{ResponseBody: createUserProfileBody(14)})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func Test_GetProfile(t *testing.T) {
	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
	router = gin.New()
	router.GET("/", routes.Testing_Export_authenticate(), routes.Testing_Export_getProfile)

	t.Run("getProfileNotLoggedIn", testGetProfileNotLoggedIn)
	t.Run("getProfileSessionError", testGetProfileSessionError)
	t.Run("getProfileCorrectly", testGetProfileCorrectly)
}