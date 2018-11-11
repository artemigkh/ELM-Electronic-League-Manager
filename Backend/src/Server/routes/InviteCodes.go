package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type TeamManagerCodeRequest struct {
	TeamId        int  `json:"teamId"`
	Administrator bool `json:"administrator"`
	Information   bool `json:"information"`
	Players       bool `json:"players"`
	ReportResults bool `json:"reportResults"`
}

/**
 * @api{POST} /api/inviteCodes/team/create Create Team Manager Invite Code
 * @apiGroup InviteCodes
 * @apiDescription Create Team Manager Invite Code. Logged in user must be a team admin or league admin
 * (has editPermissions in either). Code expires after 24h
 *
 * @apiParam {number} teamId The unique numerical identifier of the team
 * @apiParam {boolean} administrator Whether the invited user can edit permissions of other managers on the team (admin)
 * @apiParam {boolean} information Whether the invited user can edit team information
 * @apiParam {boolean} players Whether the invited user can edit players on the team
 * @apiParam {boolean} reportResults Whether the invited user can report game results of this team
 *
 * @apiSuccess {string} code A 16 character string that is the invite code
 *
 * @apiError notLoggedIn No user is logged in
 * @apiError noActiveLeague There is no active league selected
 * @apiError teamDoesNotExist The specified team does not exist
 * @apiError notAdmin The user creating the invite code is not an admin of the team or league
 */
func createTeamManagerInviteCode(ctx *gin.Context) {
	var codeRequest TeamManagerCodeRequest
	err := ctx.ShouldBindJSON(&codeRequest)
	if checkJsonErr(ctx, err) {
		return
	}

	if failIfTeamDoesNotExist(ctx, ctx.GetInt("leagueId"), codeRequest.TeamId) {
		return
	}
	if failIfNotTeamAdmin(ctx, ctx.GetInt("leagueId"), codeRequest.TeamId, ctx.GetInt("userId")) {
		return
	}

	code, err := InviteCodesDAO.CreateTeamManagerInviteCode(ctx.GetInt("leagueId"), codeRequest.TeamId,
		codeRequest.Administrator, codeRequest.Information, codeRequest.Players, codeRequest.ReportResults)
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": code})
}

func RegisterInviteCodeHandlers(g *gin.RouterGroup) {
	g.Use(getActiveLeague())

	g.POST("/team/create", authenticate(), createTeamManagerInviteCode)
}
