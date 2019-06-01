package main

import (
	"Server/config"
	"Server/databaseAccess"
	"Server/icons"
	"Server/lolApi"
	"Server/markdown"
	"Server/routes"
	"Server/sessionManager"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/stdlib"
	"net/http"
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
	routes.DataValidator = &databaseAccess.DTOValidator{}
	routes.Access = &databaseAccess.AccessChecker{}
	routes.ElmSessions = sessionManager.CreateCookieSessionManager(conf)
	routes.IconManager = icons.CreateGoIconManager(conf)
	routes.MarkdownManager = markdown.CreateGoMarkdownManager(conf)
	routes.LoLApi = lolApi.GetLolApiWrapper()
	routes.LeagueOfLegendsDAO = databaseAccess.CreateLeagueOfLegendsDAO()

	routes.RegisterLoginHandlers(app.Group("/"))
	routes.RegisterUserHandlers(app.Group("/api/v1/users"))
	routes.RegisterLeagueHandlers(app.Group("/api/v1/leagues"))
	routes.RegisterTeamHandlers(app.Group("/api/v1/teams"))
	routes.RegisterGameHandlers(app.Group("/api/v1/games"))
	routes.RegisterSchedulingHandlers(app.Group("/api/v1/scheduling"))

	routes.RegisterLeagueOfLegendsHandlers(app.Group("/api/v1/league-of-legends"), conf)

	app.GET("/test", func(ctx *gin.Context) {
		var err2 error
		var gameInfo databaseAccess.GameDTO
		err := ctx.ShouldBindJSON(&gameInfo)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		fmt.Printf("%+v\n", gameInfo)
		fmt.Printf("%+v\n", err2)
	})

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
