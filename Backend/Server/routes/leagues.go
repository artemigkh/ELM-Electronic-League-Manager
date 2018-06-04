package routes

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/kataras/iris"
	"database/sql"
	"github.com/kataras/iris/sessions"
)

type leagueRequest struct {
	Name string `json:"name"`
	PublicView bool `json:"publicView"`
	PublicJoin bool `json:"publicJoin"`
}

type idWrapper struct {
	Id int `json:"id"`
}

// /api/leagues
func RegisterLeagueHandlers(app iris.Party, db *sql.DB, sessions *sessions.Sessions) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
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
  * @apiError nameTooLong The league name has exceeded 50 characters
  * @apiError nameInUse The league name is currently in use
  */
	app.Post("/", func(ctx iris.Context) {
		//get params
		var lgRequest leagueRequest
		err := ctx.ReadJSON(&lgRequest)
		if checkErr(ctx, err) {return}

		session := sessions.Start(ctx)
		userID := authenticateAndGetCurrUserId(ctx, session)
		if userID == -1 {
			return
		}

		if failIfLeagueNameTooLong(lgRequest.Name, ctx)  {return}
		if failIfLeagueNameInUse(lgRequest.Name, ctx, psql, db)  {return}

		//create new league
		var leagueID int
		err = psql.Insert("leagues").Columns("name", "publicView", "publicJoin").
			Values(lgRequest.Name, lgRequest.PublicView, lgRequest.PublicJoin).Suffix("RETURNING \"id\"").
			RunWith(db).QueryRow().Scan(&leagueID)
		if checkErr(ctx, err) {return}

		//create permissions entry linking current user ID as the league creator
		_, err = psql.Insert("leaguePermissions").Columns("userID", "leagueID", "editPermissions", "editTeams",
					"editUsers", "editSchedule", "editResults").Values(userID, leagueID, true, true, true, true, true).
					RunWith(db).Exec()
		if checkErr(ctx, err) {return}

		ctx.JSON(idWrapper{Id: leagueID})
	})

/**
  * @api{POST} /api/leagues/setActiveLeague/:id Attempt to set the active league with id
  * @apiName setActiveLeague
  * @apiGroup leagues
  * @apiDescription Attempt to set the active league to :id
  * @apiParam {int} id the primary id of the league
  *
  * @apiError leagueDoesNotExist The league does not exist
  * @apiError 403 Forbidden
  */
	app.Post("/setActiveLeague/{id:int}", func(ctx iris.Context) {
		id, _ := ctx.Params().GetInt("id")

		//if public, set the session variable of current league to the requested league
		var publicview bool

		row := psql.Select("publicview").From("leagues").Where("id=?", id)
		err := row.Scan(&publicview)

		//check if league with specified id does not exist
		if err == sql.ErrNoRows {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(errorResponse{Error: "leagueDoesNotExist"})
			return
		}
		if checkErr(ctx, err) {return}

		//get the session
		session := sessions.Start(ctx)

		if publicview {
			//the league is publically viewable, so set the session variable to current league
			session.Set("activeLeague", id)
		} else {
			// the league is not public, so check if the logged in user (if logged in) has permissions to view
			userID := authenticateAndGetCurrUserId(ctx, session)
			if userID == -1 {
				ctx.StatusCode(iris.StatusForbidden)
				return
			}

			//if row exists, means that user is linked with this league and thus can view it as it is base privilege
			var userid int
			row := psql.Select("userid").From("leagues").Where("userid=? AND leagueid=?", userID, id)
			err := row.Scan(&userid)

			if err == sql.ErrNoRows {
				ctx.StatusCode(iris.StatusForbidden)
				return
			} else {
				session.Set("activeLeague", id)
			}
		}
	})


}