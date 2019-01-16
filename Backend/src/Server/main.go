package main

import (
	"Server/config"
	"Server/databaseAccess"
	"Server/icons"
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

	routes.RegisterLoginHandlers(app.Group("/"))
	routes.RegisterUserHandlers(app.Group("/api/users"))
	routes.RegisterLeagueHandlers(app.Group("/api/leagues"))
	routes.RegisterTeamHandlers(app.Group("/api/teams"))
	routes.RegisterGameHandlers(app.Group("/api/games"))
	routes.RegisterInviteCodeHandlers(app.Group("/api/inviteCodes"))

	// should probably be replaced with apache or nginx in production
	app.Static("/icons", conf.GetIconsDir())

	// only for testing, will be used by routes later
	//app.POST("/uploadIcon", func(c *gin.Context) {
	//	iconManager := icons.GoIconManager{
	//		OutPath: conf.GetIconsDir(),
	//	}
	//	tmpFileLoc := randomdata.RandStringRunes(10)
	//	file, err := c.FormFile("file")
	//	if err != nil {
	//		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
	//		return
	//	}
	//
	//	if err := c.SaveUploadedFile(file, "tmp/"+tmpFileLoc); err != nil {
	//		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
	//		return
	//	}
	//
	//	smallLoc, largeLoc, err := iconManager.StoreNewIcon("tmp/" + tmpFileLoc)
	//
	//	if err != nil {
	//		print(err.Error())
	//		c.Status(http.StatusBadRequest)
	//		return
	//	}
	//
	//	println(smallLoc)
	//	println(largeLoc)
	//
	//	c.Status(http.StatusOK)
	//})
	return app
}

func main() {
	conf := config.GetConfig("Backend/conf.json")

	//start router/webapp
	app := NewApp(conf)
	app.Run(conf.GetPortString())
}
