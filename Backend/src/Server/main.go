package main

import (
	"Server/config"
	"Server/databaseAccess"
	"Server/icons"
	"Server/lolApi"
	"Server/markdown"
	"Server/routes"
	"Server/sessionManager"
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
	routes.InviteCodesDAO = databaseAccess.CreateInviteCodesDAO()
	routes.ElmSessions = sessionManager.CreateCookieSessionManager(conf)
	routes.IconManager = icons.CreateGoIconManager(conf)
	routes.MarkdownManager = markdown.CreateGoMarkdownManager(conf)
	routes.LoLApi = lolApi.GetLolApiWrapper()

	routes.RegisterLoginHandlers(app.Group("/"))
	routes.RegisterUserHandlers(app.Group("/api/users"))
	routes.RegisterLeagueHandlers(app.Group("/api/leagues"))
	routes.RegisterTeamHandlers(app.Group("/api/teams"))
	routes.RegisterGameHandlers(app.Group("/api/games"))
	routes.RegisterSchedulingHandlers(app.Group("/api/scheduling"))
	routes.RegisterInviteCodeHandlers(app.Group("/api/inviteCodes"))

	routes.RegisterLeagueOfLegendsHandlers(app.Group("/api/league-of-legends"), conf)

	// should probably be replaced with apache or nginx in production
	app.Static("/icons", conf.GetIconsDir())

	return app
}

func main() {
	conf := config.GetConfig("Backend/conf.json")

	//start router/webapp
	app := NewApp(conf)
	app.Run(conf.GetPortString())
}
