package routesTest

import (
	"Server/databaseAccess"
	"Server/routes"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"mocks"
	"testing"
)

func createManagerCodeInfoResponseBody(code string, creationTime, leagueId, teamId int,
	administrator, information, players, reportResults bool) *bytes.Buffer {
	resBody := databaseAccess.TeamManagerInviteCode{
		Code:          code,
		CreationTime:  creationTime,
		LeagueId:      leagueId,
		TeamId:        teamId,
		Administrator: administrator,
		Information:   information,
		Players:       players,
		ReportResults: reportResults,
	}
	resBodyB, _ := json.Marshal(&resBody)
	return bytes.NewBuffer(resBodyB)
}

func testGetManagerInviteCodeInformationNoCode(t *testing.T) {
	httpTest(t, nil,
		"GET", "/", 404, testParams{})
}

func testGetManagerInviteCodeInformationCodeDoesNotExist(t *testing.T) {
	mockInviteCodesDao := new(mocks.InviteCodesDAO)
	mockInviteCodesDao.On("GetTeamManagerInviteCodeInformation", "96b5ad2e42c08d1e").
		Return(nil, nil)

	routes.InviteCodesDAO = mockInviteCodesDao

	httpTest(t, nil, "GET", "/96b5ad2e42c08d1e", 400, testParams{Error: "inviteCodeDoesNotExist"})

	mock.AssertExpectationsForObjects(t, mockInviteCodesDao)
}

func testGetManagerInviteCodeInformationDatabaseError(t *testing.T) {
	mockInviteCodesDao := new(mocks.InviteCodesDAO)
	mockInviteCodesDao.On("GetTeamManagerInviteCodeInformation", "96b5ad2e42c08d1e").
		Return(nil, errors.New("fake db error"))

	routes.InviteCodesDAO = mockInviteCodesDao

	httpTest(t, nil, "GET", "/96b5ad2e42c08d1e", 500, testParams{})

	mock.AssertExpectationsForObjects(t, mockInviteCodesDao)

}

func testGetManagerInviteCodeInformationCorrectGetCodeInfo(t *testing.T) {
	mockInviteCodesDao := new(mocks.InviteCodesDAO)
	mockInviteCodesDao.On("GetTeamManagerInviteCodeInformation", "96b5ad2e42c08d1e").
		Return(&databaseAccess.TeamManagerInviteCode{
			Code:          "96b5ad2e42c08d1e",
			CreationTime:  100,
			LeagueId:      1,
			TeamId:        2,
			Administrator: true,
			Information:   true,
			Players:       true,
			ReportResults: true,
		}, nil)

	routes.InviteCodesDAO = mockInviteCodesDao

	httpTest(t, nil, "GET", "/96b5ad2e42c08d1e", 200, testParams{ResponseBody: createManagerCodeInfoResponseBody("96b5ad2e42c08d1e", 100, 1, 2, true, true, true, true)})

	mock.AssertExpectationsForObjects(t, mockInviteCodesDao)
}

func Test_GetManagerInviteCodeInformation(t *testing.T) {
	//set up router and path to test
	gin.SetMode(gin.ReleaseMode) //opposite of gin.DebugMode to make tests faster by removing logging
	router = gin.New()
	router.GET("/:code", routes.Testing_Export_getTeamManagerInviteCodeInformation)

	t.Run("NoCode", testGetManagerInviteCodeInformationNoCode)
	t.Run("CodeDoesNotExist", testGetManagerInviteCodeInformationCodeDoesNotExist)
	t.Run("DatabaseError", testGetManagerInviteCodeInformationDatabaseError)
	t.Run("CorrectGetCodeInfo", testGetManagerInviteCodeInformationCorrectGetCodeInfo)
}
