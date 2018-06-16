package main

import (
	"esports-league-manager/Backend/Server/databaseAccess"
	"esports-league-manager/Backend/Server/routes"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/stdlib"
	"esports-league-manager/Backend/Server/sessionManager"
    "esports-league-manager/Backend/Server/config"
)

type Configuration struct {
	DbUser string `json:"dbUser"`
	DbPass string `json:"dbPass"`
	DbName string `json:"dbName"`
	Port   string `json:"port"`
}

func newApp(conf config.Config) *gin.Engine {
	app := gin.Default()

	databaseAccess.Init(conf)
	routes.UsersDAO = databaseAccess.CreateUsersDao()
	routes.ElmSessions = sessionManager.CreateCookieSessionManager()

	//routesTest.RegisterLoginHandlers(app.Party("/login"), &db, elmSessions)
	routes.RegisterLoginHandlers(app.Group("/login"))
	routes.RegisterUserHandlers(app.Group("/api/users"))

	//routesTest.RegisterLeagueHandlers(app.Party("/api/leagues"), &db, elmSessions)

	return app
}

func main() {
	conf := config.GetConfig()

	//start router/webapp
	app := newApp(conf)
	app.Run(conf.GetPortString())
}
