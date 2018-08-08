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
 * @api{POST} /api/leagues/ Create a new league
 * @apiName createNewLeague
 * @apiGroup leagues
 * @apiDescription Register a new league in the database
 *
 * @apiParam {string} name the name of the league
 * @apiParam {boolean} publicView should the league be viewable by people not playing in the league?
 * @apiParam {boolean} publicJoin should the league be joinable by any team that has viewing rights?
 *
 * @apiSuccess {int} id the primary id of the created league
 *
 * @apiError notLoggedIn No user is logged in to create a league
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

	leagueID, err := LeaguesDAO.CreateLeague(ctx.GetInt("userID"), lgRequest.Name, lgRequest.PublicView, lgRequest.PublicJoin)
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": leagueID})
}

/**
 * @api{POST} /api/leagues/setActiveLeague/:id Attempt to set the active league with id
 * @apiName setActiveLeague
 * @apiGroup leagues
 * @apiDescription Attempt to set the active league to :id
 * @apiParam {int} id the primary id of the league
 *
 * @apiError 403 Forbidden
 */
//TODO: check if league exists
func setActiveLeague(ctx *gin.Context) {
	//get user ID (or -1 if not logged in)
	userID, err := ElmSessions.AuthenticateAndGetUserID(ctx)
	if checkErr(ctx, err) {
		return
	}

	viewable, err := LeaguesDAO.IsLeagueViewable(ctx.GetInt("urlId"), userID)
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

func getActiveLeagueInformation(ctx *gin.Context) {
	leagueInfo, err := LeaguesDAO.GetLeagueInformation(ctx.GetInt("leagueID"))
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, leagueInfo)
}

func getTeamSummary(ctx *gin.Context) {
	teamSummary, err := LeaguesDAO.GetTeamSummary(ctx.GetInt("leagueID"))
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, teamSummary)
}

func RegisterLeagueHandlers(g *gin.RouterGroup) {
	g.POST("/", authenticate(), createNewLeague)
	g.POST("/setActiveLeague/:id", getUrlId(), setActiveLeague)
	g.GET("/", getActiveLeague(), getActiveLeagueInformation)
	g.GET("/teamSummary", getActiveLeague(), getTeamSummary)
}
