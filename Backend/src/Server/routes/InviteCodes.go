package routes

import "github.com/gin-gonic/gin"

/**
 * @api{POST} /api/inviteCodes/team Create Team Manager Invite Code
 * @apiGroup InviteCodes
 * @apiDescription Create Team Manager Invite Code. Logged in user must be a team admin or league admin
 * (has editPermissions in either). Code expires after 24h
 *
 * @apiParam {number} teamId The unique numerical identifier of the team
 * @apiParam {boolean} editPermissions Whether the invited user can edit permissions of other managers on the team (admin)
 * @apiParam {boolean} editTeamInfo Whether the invited user can edit team information
 * @apiParam {boolean} editPlayers Whether the invited user can edit players on the team
 * @apiParam {boolean} reportResult Whether the invited user can report game results of this team
 *
 * @apiSuccess {string} A 16 character string that is the invite code
 *
 * @apiError notLoggedIn No user is logged in
 * @apiError noActiveLeague There is no active league selected
 * @apiError teamDoesNotExist The specified team does not exist
 * @apiError notAdmin The user creating the invite code is not an admin of the team or league
 */
func createTeamManagerInviteCode(ctx *gin.Context) {

}

func RegisterInviteCodeHandlers(g *gin.RouterGroup) {
	g.Use(getActiveLeague())

	g.POST("/team", authenticate(), createTeamManagerInviteCode)
}
