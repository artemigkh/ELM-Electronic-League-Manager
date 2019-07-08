package routes

import (
	"Server/config"
	"Server/databaseAccess"
	"Server/icons"
	"Server/lolApi"
	"Server/markdown"
	"Server/sessionManager"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRoutes(conf config.Config) *gin.Engine {
	app := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:4200"}
	corsConfig.AllowCredentials = true

	app.Use(cors.New(corsConfig), Authenticate())
	databaseAccess.Init(conf)
	UsersDAO = databaseAccess.CreateUsersDao()
	LeaguesDAO = databaseAccess.CreateLeaguesDAO()
	TeamsDAO = databaseAccess.CreateTeamsDAO()
	GamesDAO = databaseAccess.CreateGamesDAO()
	Access = &databaseAccess.AccessChecker{}
	ElmSessions = sessionManager.CreateCookieSessionManager(conf)
	IconManager = icons.CreateGoIconManager(conf)
	MarkdownManager = markdown.CreateGoMarkdownManager(conf)
	LoLApi = lolApi.GetLolApiWrapper()
	LeagueOfLegendsDAO = databaseAccess.CreateLeagueOfLegendsDAO()

	RegisterLoginHandlers(app.Group("/"))
	RegisterUserHandlers(app.Group("/api/v1/users"))
	RegisterLeagueHandlers(app.Group("/api/v1/leagues"))
	RegisterTeamHandlers(app.Group("/api/v1"))
	RegisterGameHandlers(app.Group(""))
	RegisterSchedulingHandlers(app.Group("/api/v1"))
	//
	//routes.RegisterLeagueOfLegendsHandlers(app.Group("/api/v1/league-of-legends"), conf)

	// should probably be replaced with apache or nginx in production
	app.Static("/icons", conf.GetIconsDir())

	return app
}
