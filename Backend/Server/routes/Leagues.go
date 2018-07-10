package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

	userID, err := ElmSessions.AuthenticateAndGetUserID(ctx)
	if checkErr(ctx, err) {
		return
	}

	if userID == -1 {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "notLoggedIn"})
		return
	}

	if failIfLeagueNameTooLong(ctx, lgRequest.Name) {
		return
	}
	if failIfLeagueNameInUse(ctx, lgRequest.Name) {
		return
	}

	leagueID, err := LeaguesDAO.CreateLeague(userID, lgRequest.Name, lgRequest.PublicView, lgRequest.PublicJoin)
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
func setActiveLeague(ctx *gin.Context) {
	leagueId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "IdMustBeInteger"})
		return
	}

	//get user ID (or -1 if not logged in)
	userID, err := ElmSessions.AuthenticateAndGetUserID(ctx)
	if checkErr(ctx, err) {
		return
	}

	viewable, err := LeaguesDAO.IsLeagueViewable(leagueId, userID)
	if checkErr(ctx, err) {
		return
	}

	if !viewable {
		ctx.JSON(http.StatusForbidden, nil)
	} else {
		err := ElmSessions.SetActiveLeague(ctx, leagueId)
		if checkErr(ctx, err) {
			return
		}
	}
}

func RegisterLeagueHandlers(g *gin.RouterGroup) {
	g.POST("/", createNewLeague)
	g.POST("/setActiveLeague/:id", setActiveLeague)
}