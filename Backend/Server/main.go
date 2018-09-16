package main

import (
	"esports-league-manager/Backend/Server/config"
	"esports-league-manager/Backend/Server/databaseAccess"
	"esports-league-manager/Backend/Server/routes"
	"esports-league-manager/Backend/Server/sessionManager"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/stdlib"
)

type Configuration struct {
	DbUser string `json:"dbUser"`
	DbPass string `json:"dbPass"`
	DbName string `json:"dbName"`
	Port   string `json:"port"`
}

func NewApp(conf config.Config) *gin.Engine {
	app := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:4200"}
	corsConfig.AllowCredentials = true

	app.Use(cors.New(corsConfig))
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
	conf := config.GetConfig("Backend/conf.json")

	//start router/webapp
	app := NewApp(conf)
	app.Run(conf.GetPortString())
}
