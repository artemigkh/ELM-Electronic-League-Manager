package routes

import (
	"github.com/kataras/iris"
)

// /api/leagues
func RegisterLeagueHandlers(app iris.Party) {
	app.Post("/", func(ctx iris.Context) {
		ctx.Writef("new league created")
	})

	app.Get("/", func(ctx iris.Context) {
		ctx.Writef("list of all leagues")
	})

	app.Delete("/{id:int}", func(ctx iris.Context) {
		id, _ := ctx.Params().GetInt("id")
		println(id)
	})
}