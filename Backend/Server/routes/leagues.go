package routes

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
)

// /api/leagues
func RegisterLeagueHandlers(app iris.Party) {
	app.Post("/", func(ctx context.Context) {
		ctx.Writef("new league created")
	})

	app.Get("/", func(ctx context.Context) {
		ctx.Writef("list of all leagues")
	})

	app.Delete("/{id:int}", func(ctx context.Context) {
		id, _ := ctx.Params().GetInt("id")
		println(id)
	})
}