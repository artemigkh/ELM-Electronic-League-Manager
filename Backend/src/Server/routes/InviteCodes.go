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

/**
 * @api{GET} /api/inviteCodes/team/getInformation/:code Get Information about Team Manager Invite Code
 * @apiGroup InviteCodes
 * @apiDescription Get Information about a team manager invite code
 *
 * @apiParam {code} teamId The 16 character invite code
 *
 * @apiSuccess {string} code A 16 character string that is the invite code
 * @apiSuccess {int} creationTime The time the invite code was created
 * @apiSuccess {int} leagueId The id of the league to which the invite code applies
 * @apiSuccess {int} teamId The id of the team to which the invite code applies
 * @apiSuccess {bool} administrator Whether or not the invite gives admin privileges
 * @apiSuccess {bool} information Whether or not the invite gives edit information privileges
 * @apiSuccess {bool} players Whether or not the invite gives edit players privileges
 * @apiSuccess {bool} reportResults Whether or not the invite gives report game result privileges
 *
 * @apiError inviteCodeDoesNotExist The specified invite code does not exist
 */
func getTeamManagerInviteCodeInformation(ctx *gin.Context) {
	code := ctx.Param("code")
	codeInfo, err := InviteCodesDAO.GetTeamManagerInviteCodeInformation(code)
	if checkErr(ctx, err) {
		return
	}
	if codeInfo == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "inviteCodeDoesNotExist"})
		return
	}

	ctx.JSON(http.StatusOK, codeInfo)
}

/**
 * @api{POST} /api/inviteCodes/team/useCode/:code Use a Team Manager Invite Code
 * @apiGroup InviteCodes
 * @apiDescription Use a team manager invite code
 *
 * @apiParam {code} teamId The 16 character invite code
 *
 * @apiError notLoggedIn No user is logged in
 */
func useTeamManagerInviteCode(ctx *gin.Context) {
	//TODO: do not allow use of code if user is already manager of the team
	//TODO: join league as well if user is not part of current league
	err := InviteCodesDAO.UseTeamManagerInviteCode(ctx.GetInt("userId"), ctx.Param("code"))
	if checkErr(ctx, err) {
		return
	}
	ctx.Status(http.StatusOK)
}

func RegisterInviteCodeHandlers(g *gin.RouterGroup) {
	g.POST("/team/create", getActiveLeague(), authenticate(), createTeamManagerInviteCode)
	g.GET("/team/getInformation/:code", getTeamManagerInviteCodeInformation)
	g.POST("/team/useCode/:code", authenticate(), useTeamManagerInviteCode)
}
