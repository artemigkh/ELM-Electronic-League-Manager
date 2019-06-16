package routesTest

//
//import (
//	"Server/databaseAccess"
//	"Server/routes"
//	"bytes"
//	"encoding/json"
//	"errors"
//	"github.com/gin-gonic/gin"
//	"github.com/stretchr/testify/mock"
//	"mocks"
//	"testing"
//)
//
//func createUserProfileBody(email string) *bytes.Buffer {
//	body := databaseAccess.UserInformation{Email: email}
//	reqBodyB, _ := json.Marshal(&body)
//	return bytes.NewBuffer(reqBodyB)
//}
//
//func testGetProfileNotLoggedIn(t *testing.T) {
//	mockSession := new(mocks.SessionManager)
//	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
//		Return(-1, nil)
//
//	routes.ElmSessions = mockSession
//
//	httpTest(t, nil, "GET", "/", 403, testParams{Error: "notLoggedIn"})
//
//	mock.AssertExpectationsForObjects(t, mockSession)
//}
//
//func testGetProfileSessionError(t *testing.T) {
//	mockSession := new(mocks.SessionManager)
//	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
//		Return(-1, errors.New("fake session error"))
//
//	routes.ElmSessions = mockSession
//
//	httpTest(t, nil, "GET", "/", 500, testParams{})
//
//	mock.AssertExpectationsForObjects(t, mockSession)
//}
//
//func testGetProfileDatabaseError(t *testing.T) {
//	mockSession := new(mocks.SessionManager)
//	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
//		Return(14, nil)
//
//	mockUsersDao := new(mocks.UsersDAO)
//	mockUsersDao.On("GetUserProfile", 14).
//		Return(&databaseAccess.UserInformation{Email: "test3@test.com"}, errors.New("fake db error"))
//
//	routes.ElmSessions = mockSession
//	routes.UsersDAO = mockUsersDao
//
//	httpTest(t, nil, "GET", "/", 500, testParams{})
//
//	mock.AssertExpectationsForObjects(t, mockSession, mockUsersDao)
//}
//
//func testGetProfileCorrectly(t *testing.T) {
//	mockSession := new(mocks.SessionManager)
//	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
//		Return(14, nil)
//
//	mockUsersDao := new(mocks.UsersDAO)
//	mockUsersDao.On("GetUserProfile", 14).
//		Return(&databaseAccess.UserInformation{Email: "test3@test.com"}, nil)
//
//	routes.ElmSessions = mockSession
//	routes.UsersDAO = mockUsersDao
//
//	httpTest(t, nil, "GET", "/", 200,
//		testParams{ResponseBody: createUserProfileBody("test3@test.com")})
//
//	mock.AssertExpectationsForObjects(t, mockSession, mockUsersDao)
//}
//
//func Test_GetProfile(t *testing.T) {
//	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
//	router = gin.New()
//	router.GET("/", routes.Testing_Export_authenticate(), routes.Testing_Export_getProfile)
//
//	t.Run("getProfileNotLoggedIn", testGetProfileNotLoggedIn)
//	t.Run("getProfileSessionError", testGetProfileSessionError)
//	t.Run("getProfileDatabaseError", testGetProfileDatabaseError)
//	t.Run("getProfileCorrectly", testGetProfileCorrectly)
//}
