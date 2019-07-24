package routes

import (
	"Server/config"
	"Server/databaseAccess"
	"Server/icons"
	"Server/lolApi"
	"Server/markdown"
	"Server/sessionManager"
	"Server/validation"
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
	UserDAO = &databaseAccess.UserSqlDao{}
	LeagueDAO = &databaseAccess.LeagueSqlDao{}
	TeamDAO = &databaseAccess.TeamSqlDao{}
	GameDAO = &databaseAccess.GameSqlDao{}
	LeagueOfLegendsDAO = &databaseAccess.LeagueOfLegendsSqlDao{}

	Access = &validation.AccessChecker{}
	ElmSessions = sessionManager.CreateCookieSessionManager(conf)
	IconManager = icons.CreateGoIconManager(conf)
	MarkdownManager = markdown.CreateGoMarkdownManager(conf)
	LoLApi = lolApi.GetLolApiWrapper()

	RegisterLoginHandlers(app.Group("/"))
	RegisterUserHandlers(app.Group("/api/v1/users"))
	RegisterLeagueHandlers(app.Group("/api/v1/leagues"))
	RegisterTeamHandlers(app.Group("/api/v1"))
	RegisterGameHandlers(app.Group(""))
	RegisterSchedulingHandlers(app.Group("/api/v1"))
	//
	RegisterLeagueOfLegendsHandlers(app.Group("/api/v1/lol"))

	// should probably be replaced with apache or nginx in production
	app.Static("/icons", conf.GetIconsDir())

	return app
}
