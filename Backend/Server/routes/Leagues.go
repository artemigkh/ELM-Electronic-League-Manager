package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type LeagueRequest struct {
	Name       string `json:"name"`
	PublicView bool   `json:"publicView"`
	PublicJoin bool   `json:"publicJoin"`
}

/**
 * @api{POST} /api/leagues/ Create New League
 * @apiName createNewLeague
 * @apiGroup Leagues
 * @apiDescription Register a new league
 *
 * @apiParam {string} name the name of the league
 * @apiParam {boolean} publicView should the league be viewable by people not playing in the league?
 * @apiParam {boolean} publicJoin should the league be joinable by any team that has viewing rights?
 *
 * @apiSuccess {int} id the primary id of the created league
 *
 * @apiError notLoggedIn No user is logged in
 * @apiError nameTooLong The league name has exceeded 50 characters
 * @apiError nameInUse The league name is currently in use
 */
func createNewLeague(ctx *gin.Context) {
	var lgRequest LeagueRequest
	err := ctx.ShouldBindJSON(&lgRequest)
	if checkJsonErr(ctx, err) {
		return
	}

	if failIfNameTooLong(ctx, lgRequest.Name) {
		return
	}
	if failIfLeagueNameInUse(ctx, lgRequest.Name) {
		return
	}

	leagueId, err := LeaguesDAO.CreateLeague(ctx.GetInt("userId"), lgRequest.Name, lgRequest.PublicView, lgRequest.PublicJoin)
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": leagueId})
}

/**
 * @api{POST} /api/leagues/setActiveLeague/:id Set Active League
 * @apiName setActiveLeague
 * @apiGroup Leagues
 * @apiDescription Attempt to set the active league to :id
 * @apiParam {int} id the primary id of the league
 *
 * @apiError 403 Forbidden
 */
func setActiveLeague(ctx *gin.Context) {
	//TODO: check if league exists
	//get user Id (or -1 if not logged in)
	userId, err := ElmSessions.AuthenticateAndGetUserId(ctx)
	if checkErr(ctx, err) {
		return
	}

	viewable, err := LeaguesDAO.IsLeagueViewable(ctx.GetInt("urlId"), userId)
	if checkErr(ctx, err) {
		return
	}

	if !viewable {
		ctx.JSON(http.StatusForbidden, nil)
	} else {
		err := ElmSessions.SetActiveLeague(ctx, ctx.GetInt("urlId"))
		if checkErr(ctx, err) {
			return
		}
	}
}

/**
 * @api{GET} /api/leagues Get Active League Information
 * @apiGroup Leagues
 * @apiDescription Get information about the currently selected league
 *
 * @apiSuccess {int} id The unique numerical identifier of the league
 *
 * @apiError noActiveLeague There is no active league selected
 */
func getActiveLeagueInformation(ctx *gin.Context) {
	leagueInfo, err := LeaguesDAO.GetLeagueInformation(ctx.GetInt("leagueId"))
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, leagueInfo)
}

/**
 * @api{GET} /api/leagues/teamSummary Get Team Summary
 * @apiGroup Leagues
 * @apiDescription Get the team summary of the current league, sorted by standings
 *
 * @apiSuccess {jsonArray} _ An array of JSON objects, each representing a team
 * @apiSuccess {int} _.id The unique numerical identifier of the team
 * @apiSuccess {int} _.name The name of the team
 * @apiSuccess {int} _.tag The tag of the team
 * @apiSuccess {int} _.wins The number of wins of the team
 * @apiSuccess {int} _.losses The number of losses of the team
 *
 * @apiError noActiveLeague There is no active league selected
 */
func getTeamSummary(ctx *gin.Context) {
	teamSummary, err := LeaguesDAO.GetTeamSummary(ctx.GetInt("leagueId"))
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, teamSummary)
}

/**
 * @api{POST} /api/leagues/join Join Active League
 * @apiGroup Leagues
 * @apiDescription Join the currently selected league as a manager
 *
 * @apiError notLoggedIn No user is logged in
 * @apiError noActiveLeague There is no active league selected
 * @apiError canNotJoin The active league is not accepting new members
 */
func joinActiveLeague(ctx *gin.Context) {
	err := LeaguesDAO.JoinLeague(ctx.GetInt("leagueId"), ctx.GetInt("userId"))
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

/**
 * @api{GET} /api/leagues/teamManagers Get Team Managers
 * @apiGroup Leagues
 * @apiDescription If logged in as a league administrator, see all users that have permissions to manage teams in this league
 *
 * @apiSuccess {jsonArray} _ An array of JSON objects, each representing a team
 * @apiSuccess {int} _.teamId The unique numerical identifier of the team
 * @apiSuccess {string} _.teamName The name of the team
 * @apiSuccess {string} _.teamTag The tag of the team
 * @apiSuccess {[]Object} _.managers The users on this team that have management permissions
 * @apiSuccess {int} _.managers.userId The unique numerical identifier of the user/manager
 * @apiSuccess {string} _.managers.userEmail The email of the user/manager
 * @apiSuccess {bool} _.managers.editPermissions True if this user can manage permissions of other users on the team
 * @apiSuccess {bool} _.managers.editTeamInfo True if this user can edit information about the team
 * @apiSuccess {bool} _.managers.editPlayers True if this user can edit players on this team
 * @apiSuccess {bool} _.managers.reportResult True if this user can report results for this team
 *
 * @apiError notLoggedIn No user is logged in
 * @apiError noActiveLeague There is no active league selected
 * @apiError notAdmin The currently logged in user is not a league administrator
 */
func getTeamManagers(ctx *gin.Context) {
	teamManagerInfo, err := LeaguesDAO.GetTeamManagerInformation(ctx.GetInt("leagueId"))
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, teamManagerInfo)
}

/**
 * @api{GET} /api/leagues/gameSummary Get Game Summary
 * @apiGroup Leagues
 * @apiDescription Get the game summary of the current league, in chronological order
 *
 * @apiSuccess {jsonArray} _ An array of JSON objects, each representing a game
 * @apiSuccess {int} _.id The unique numerical identifier of the game
 * @apiSuccess {int} _.team1Id The unique numerical identifier of the team in position 1
 * @apiSuccess {int} _.team2Id The unique numerical identifier of the team in position 2
 * @apiSuccess {int} _.gameTime The unix epoch time in seconds when the game is played
 * @apiSuccess {bool} _.complete A boolean indicating if the game is finished or not
 * @apiSuccess {int} _.winnerId The Id of the winning team, or -1 if the game is not complete
 * @apiSuccess {int} _.scoreTeam1 The score of the team in position 1
 * @apiSuccess {int} _.scoreTeam2 The score of the team in position 2
 *
 * @apiError noActiveLeague There is no active league selected
 */
func getGameSummary(ctx *gin.Context) {
	gameSummary, err := LeaguesDAO.GetGameSummary(ctx.GetInt("leagueId"))
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, gameSummary)
}

//TODO: create endpoint to get list of all publicly visible leagues

func RegisterLeagueHandlers(g *gin.RouterGroup) {
	g.POST("/", authenticate(), createNewLeague)
	g.POST("/setActiveLeague/:id", getUrlId(), setActiveLeague)
	g.POST("/join", authenticate(), getActiveLeague(), failIfCannotJoinLeague(), joinActiveLeague)
	g.GET("/", getActiveLeague(), getActiveLeagueInformation)
	g.GET("/teamSummary", getActiveLeague(), getTeamSummary)
	g.GET("/gameSummary", getActiveLeague(), getGameSummary)
	g.GET("/teamManagers", authenticate(), getActiveLeague(), failIfNotLeagueAdmin(), getTeamManagers)
}
