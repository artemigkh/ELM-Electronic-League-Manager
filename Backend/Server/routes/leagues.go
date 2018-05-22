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
  * @apiError nameTooLong The league name has exceeded 50 characters
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

		if failIfLeagueNameTooLong(lgRequest.Name, ctx) {return}

		//create new league
		res, err := psql.Insert("leagues").Columns("name", "publicView", "publicJoin").
			Values(lgRequest.Name, lgRequest.PublicView, lgRequest.PublicJoin).RunWith(db).Exec()
		if checkErr(ctx, err) {return}

		id, err := res.LastInsertId()
		if checkErr(ctx, err) {return}
		println(id)

	})

	app.Get("/", func(ctx iris.Context) {
		ctx.Writef("list of all leagues")
	})

	app.Delete("/{id:int}", func(ctx iris.Context) {
		id, _ := ctx.Params().GetInt("id")
		println(id)
	})
}