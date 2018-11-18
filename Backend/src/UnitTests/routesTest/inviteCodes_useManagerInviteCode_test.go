package routesTest

import (
	"Server/routes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"mocks"
	"testing"
)

func testUseManagerInviteCodeSessionsError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(5, errors.New("session error"))

	routes.ElmSessions = mockSession

	httpTest(t, nil,
		"POST", "/96b5ad2e42c08d1e", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testUseManagerInviteCodeNotLoggedIn(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(-1, nil)

	routes.ElmSessions = mockSession

	httpTest(t, nil,
		"POST", "/96b5ad2e42c08d1e", 403, testParams{Error: "notLoggedIn"})

	mock.AssertExpectationsForObjects(t, mockSession)
}

func testUseManagerInviteCodeNoCode(t *testing.T) {
	httpTest(t, nil,
		"POST", "/", 404, testParams{})
}

func testUseManagerInviteCodeDatabaseError(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(5, nil)

	mockInviteCodesDao := new(mocks.InviteCodesDAO)
	mockInviteCodesDao.On("UseTeamManagerInviteCode", 5, "96b5ad2e42c08d1e").
		Return(errors.New("fake db error"))

	routes.ElmSessions = mockSession
	routes.InviteCodesDAO = mockInviteCodesDao

	httpTest(t, nil,
		"POST", "/96b5ad2e42c08d1e", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockInviteCodesDao)
}

func testUseManagerInviteCodeCorrectCreate(t *testing.T) {
	mockSession := new(mocks.SessionManager)
	mockSession.On("AuthenticateAndGetUserId", mock.Anything).
		Return(5, nil)

	mockInviteCodesDao := new(mocks.InviteCodesDAO)
	mockInviteCodesDao.On("UseTeamManagerInviteCode", 5, "96b5ad2e42c08d1e").
		Return(nil)

	routes.ElmSessions = mockSession
	routes.InviteCodesDAO = mockInviteCodesDao

	httpTest(t, nil,
		"POST", "/96b5ad2e42c08d1e", 200, testParams{})

	mock.AssertExpectationsForObjects(t, mockSession, mockInviteCodesDao)
}

func Test_UseManagerInviteCode(t *testing.T) {
	//set up router and path to test
	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
	router = gin.New()
	router.POST("/:code", routes.Testing_Export_authenticate(), routes.Testing_Export_useTeamManagerInviteCode)

	t.Run("SessionsError", testUseManagerInviteCodeSessionsError)
	t.Run("NotLoggedIn", testUseManagerInviteCodeNotLoggedIn)
	t.Run("NoCode", testUseManagerInviteCodeNoCode)
	t.Run("DatabaseError", testUseManagerInviteCodeDatabaseError)
	t.Run("CorrectUse", testUseManagerInviteCodeCorrectCreate)
}
