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
	routes.LeaguesDAO = databaseAccess.CreateLeaguesDAO()
	routes.TeamsDAO = databaseAccess.CreateTeamsDAO()
	routes.GamesDAO = databaseAccess.CreateGamesDAO()
	routes.ElmSessions = sessionManager.CreateCookieSessionManager()

	routes.RegisterLoginHandlers(app.Group("/"))
	routes.RegisterUserHandlers(app.Group("/api/users"))
	routes.RegisterLeagueHandlers(app.Group("/api/leagues"))
	routes.RegisterTeamHandlers(app.Group("/api/teams"))
	routes.RegisterGameHandlers(app.Group("/api/games"))

	return app
}

func main() {
	conf := config.GetConfig()

	//start router/webapp
	app := newApp(conf)
	app.Run(conf.GetPortString())
}
