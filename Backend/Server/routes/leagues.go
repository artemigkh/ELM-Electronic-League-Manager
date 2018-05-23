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

	app.Get("/", func(ctx iris.Context) {
		ctx.Writef("list of all leagues")
	})

	app.Delete("/{id:int}", func(ctx iris.Context) {
		id, _ := ctx.Params().GetInt("id")
		println(id)
	})
}