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

	if failIfLeagueNameTooLong(ctx, lgRequest.Name) {
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
	if failIfCannotJoinLeague(ctx, ctx.GetInt("userId"), ctx.GetInt("leagueId")) {
		return
	}
	err := LeaguesDAO.JoinLeague(ctx.GetInt("userId"), ctx.GetInt("leagueId"))
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

//TODO: make "get managers of active league endpoint"

func RegisterLeagueHandlers(g *gin.RouterGroup) {
	g.POST("/", authenticate(), createNewLeague)
	g.POST("/setActiveLeague/:id", getUrlId(), setActiveLeague)
	g.POST("/join", authenticate(), getActiveLeague(), joinActiveLeague)
	g.GET("/", getActiveLeague(), getActiveLeagueInformation)
	g.GET("/teamSummary", getActiveLeague(), getTeamSummary)
}
