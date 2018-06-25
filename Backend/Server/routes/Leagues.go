package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//
//import (
//	sq "github.com/Masterminds/squirrel"
//	"github.com/kataras/iris"
//	"database/sql"
//	"github.com/kataras/iris/sessions"
//)
//

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
	//get string from path and convert to int
	idString := ctx.Param("id")
	if idString == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "noIdSpecified"})
	}

	leagueId, err := strconv.Atoi(idString)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "IdMustBeInteger"})
	}

	//get user ID (or -1 if not logged in)
	userID, err := ElmSessions.AuthenticateAndGetUserID(ctx)
	if checkErr(ctx, err) {
		return
	}

	viewable, err := LeaguesDAO.IsLeagueViewable(leagueId, userID)

	if !viewable {
		ctx.JSON(http.StatusForbidden, nil)
	} else {
		ElmSessions.SetActiveLeague(ctx, leagueId)
	}
}

//	app.Post("/setActiveLeague/{id:int}", func(ctx iris.Context) {
//		id, _ := ctx.Params().GetInt("id")
//
//		//if public, set the session variable of current league to the requested league
//		var publicview bool
//
//		row := psql.Select("publicview").From("leagues").Where("id=?", id).RunWith(db).QueryRow()
//		err := row.Scan(&publicview)
//
//		//check if league with specified id does not exist
//		if err == sql.ErrNoRows {
//			ctx.StatusCode(iris.StatusBadRequest)
//			ctx.JSON(errorResponse{Error: "leagueDoesNotExist"})
//			return
//		}
//		if checkErr(ctx, err) {return}
//
//		//get the session
//		session := sessions.Start(ctx)
//
//		if publicview {
//			//the league is publically viewable, so set the session variable to current league
//			session.Set("activeLeague", id)
//		} else {
//			// the league is not public, so check if the logged in user (if logged in) has permissions to view
//			userID := authenticateAndGetCurrUserId(ctx, session)
//			if userID == -1 {
//				ctx.StatusCode(iris.StatusForbidden)
//				return
//			}
//
//			//if row exists, means that user is linked with this league and thus can view it as it is base privilege
//			var userid int
//			row := psql.Select("userid").From("leagues").Where("userid=? AND leagueid=?", userID, id).
//					RunWith(db).QueryRow()
//			err := row.Scan(&userid)
//
//			if err == sql.ErrNoRows {
//				ctx.StatusCode(iris.StatusForbidden)
//				return
//			} else {
//				session.Set("activeLeague", id)
//			}
//		}
//	})
//

func RegisterLeagueHandlers(g *gin.RouterGroup) {
	g.POST("/", createNewLeague)
	g.POST("/setActiveLeague/:id")
}

//// /api/leagues
//func RegisterLeagueHandlers(app iris.Party, db *sql.DB, sessions *sessions.Sessions) {
//

//	app.Post("/", func(ctx iris.Context) {
//		//get params
//		var lgRequest LeagueRequest
//		err := ctx.ReadJSON(&lgRequest)
//		if checkErr(ctx, err) {return}
//
//		session := sessions.Start(ctx)
//		userID := authenticateAndGetCurrUserId(ctx, session)
//		if userID == -1 {
//			return
//		}
//
//		if failIfLeagueNameTooLong(lgRequest.Name, ctx)  {return}
//		if failIfLeagueNameInUse(lgRequest.Name, ctx, psql, db)  {return}
//
//		//create new league
//		var leagueID int
//		err = psql.Insert("leagues").Columns("name", "publicView", "publicJoin").
//			Values(lgRequest.Name, lgRequest.PublicView, lgRequest.PublicJoin).Suffix("RETURNING \"id\"").
//			RunWith(db).QueryRow().Scan(&leagueID)
//		if checkErr(ctx, err) {return}
//
//		//create permissions entry linking current user ID as the league creator
//		_, err = psql.Insert("leaguePermissions").Columns("userID", "leagueID", "editPermissions", "editTeams",
//					"editUsers", "editSchedule", "editResults").Values(userID, leagueID, true, true, true, true, true).
//					RunWith(db).Exec()
//		if checkErr(ctx, err) {return}
//
//		ctx.JSON(idWrapper{Id: leagueID})
//	})
//

///**
// * @api{GET} /api/leagues/ Get information about the current league
// * @apiName getLeagueInformation
// * @apiGroup leagues
// * @apiDescription get information about the current league
// *
// * @apiError noCurrentActiveLeague There is no currently active league
// *
// * @apiSuccess {int} id the id of the current league
// *
// */
//	app.Get("/", func(ctx iris.Context) {
//		session := sessions.Start(ctx)
//
//		//fail if active league does not exist
//		id, err := session.GetInt("activeLeague")
//		if err != nil {
//			ctx.JSON(errorResponse{Error: "noCurrentActiveLeague"})
//			ctx.StatusCode(iris.StatusBadRequest)
//			return
//		}
//
//		ctx.JSON(leagueInformation{Id: id})
//		ctx.StatusCode(iris.StatusOK)
//	})
//
//}
