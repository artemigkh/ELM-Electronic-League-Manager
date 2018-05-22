package main

import (
	"database/sql"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/kataras/iris"
	"esports-league-manager/Backend/Server/routes"
	"github.com/gorilla/securecookie"
	"github.com/kataras/iris/sessions"
	"log"
)

func newApp(db sql.DB) *iris.Application {
	app := iris.New()

	//create session manager with secure cookies
	cookieName := "elmsession"
	hashKey := securecookie.GenerateRandomKey(24)
	blockKey := securecookie.GenerateRandomKey(24)
	secureCookie := securecookie.New(hashKey, blockKey)

	elmSessions := sessions.New(sessions.Config{
		Cookie: cookieName,
		Encode: secureCookie.Encode,
		Decode: secureCookie.Decode,
	})
	routes.RegisterLoginHandlers(app.Party("/login"), &db, elmSessions)
	routes.RegisterUserHandlers(app.Party("/api/users"), &db, elmSessions)
	routes.RegisterLeagueHandlers(app.Party("/api/leagues"), &db, elmSessions)

	return app
}

func main() {
	//create database connection
	connStr := "user=postgres password=123 dbname=elmdb sslmode=disable"
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	app := newApp(*db)
	app.Run(iris.Addr(":8080"))
}
