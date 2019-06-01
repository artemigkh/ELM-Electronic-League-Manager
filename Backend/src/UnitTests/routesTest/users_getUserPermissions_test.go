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
//func createUserPermissionsBody(body databaseAccess.UserPermissions) *bytes.Buffer {
//	bodyB, _ := json.Marshal(&body)
//	return bytes.NewBuffer(bodyB)
//}
//
//func testGetUserPermissionsSessionError(t *testing.T) {
//	mockSession := new(mocks.SessionManager)
//	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
//		Return(1, errors.New("session error"))
//
//	routes.ElmSessions = mockSession
//
//	httpTest(t, nil, "GET", "/", 500, testParams{})
//
//	mock.AssertExpectationsForObjects(t, mockSession)
//}
//
//func testGetUserPermissionsNotLoggedIn(t *testing.T) {
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
//func testGetUserPermissionsNoActiveLeague(t *testing.T) {
//	mockSession := new(mocks.SessionManager)
//	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
//		Return(1, nil)
//	mockSession.On("GetActiveLeague", mock.Anything).
//		Return(-1, nil)
//
//	routes.ElmSessions = mockSession
//
//	httpTest(t, nil, "GET", "/", 403, testParams{Error: "noActiveLeague"})
//
//	mock.AssertExpectationsForObjects(t, mockSession)
//}
//
//func testGetUserPermissionsDatabaseError(t *testing.T) {
//	mockSession := new(mocks.SessionManager)
//	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
//		Return(1, nil)
//	mockSession.On("GetActiveLeague", mock.Anything).
//		Return(2, nil)
//
//	mockUsersDao := new(mocks.UsersDAO)
//	mockUsersDao.On("GetPermissions", 2, 1).Return(nil, errors.New("fake db error"))
//
//	routes.ElmSessions = mockSession
//	routes.UsersDAO = mockUsersDao
//
//	httpTest(t, nil, "GET", "/", 500, testParams{})
//
//	mock.AssertExpectationsForObjects(t, mockSession)
//}
//
//func testCorrectGetUserPermissions(t *testing.T) {
//	userPermissions := databaseAccess.UserPermissions{
//		LeaguePermissions: databaseAccess.LeaguePermissions{
//			Administrator: false,
//			CreateTeams:   false,
//			EditTeams:     false,
//			EditGames:     true,
//		},
//		TeamPermissions: []databaseAccess.TeamPermissionsInformation{
//			{
//				Id:            220,
//				Administrator: false,
//				Information:   true,
//				Players:       true,
//				ReportResults: true,
//			},
//			{
//				Id:            217,
//				Administrator: true,
//				Information:   true,
//				Players:       true,
//				ReportResults: true,
//			},
//		},
//	}
//
//	mockSession := new(mocks.SessionManager)
//	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
//		Return(1, nil)
//	mockSession.On("GetActiveLeague", mock.Anything).
//		Return(2, nil)
//
//	mockUsersDao := new(mocks.UsersDAO)
//	mockUsersDao.On("GetPermissions", 2, 1).Return(&userPermissions, nil)
//
//	routes.ElmSessions = mockSession
//	routes.UsersDAO = mockUsersDao
//
//	httpTest(t, nil, "GET", "/", 200, testParams{
//		ResponseBody: createUserPermissionsBody(userPermissions)})
//
//	mock.AssertExpectationsForObjects(t, mockSession)
//}
//
//func Test_GetUserPermissions(t *testing.T) {
//	//set up router and path to test
//	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
//	router = gin.New()
//	router.GET("/",
//		routes.Testing_Export_authenticate(),
//		routes.Testing_Export_getActiveLeague(),
//		routes.Testing_Export_getUserPermissions)
//
//	t.Run("SessionsError", testGetUserPermissionsSessionError)
//	t.Run("NotLoggedIn", testGetUserPermissionsNotLoggedIn)
//	t.Run("NoActiveLeague", testGetUserPermissionsNoActiveLeague)
//	t.Run("DatabaseError", testGetUserPermissionsDatabaseError)
//	t.Run("CorrectGetManagers", testCorrectGetUserPermissions)
//}
