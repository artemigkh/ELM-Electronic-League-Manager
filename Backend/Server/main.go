package main

import (
	"github.com/kataras/iris"
	"esports-league-manager/Backend/Server/routes"
	"esports-league-manager/Backend/Server/errorCheck"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func newApp(db sql.DB) *iris.Application {
	app := iris.New()

	return app
}

func main() {
	//create database connection
	connStr := "user=postgres dbname=elmdb sslmode=verify-full"
	db, err := sql.Open("mysql", connStr)
	errorCheck.Check(err)
	defer db.Close()

	app := newApp(*db)
	routes.RegisterLeagueHandlers(app.Party("/leagues"))

	app.Run(iris.Addr(":8080"))
}
