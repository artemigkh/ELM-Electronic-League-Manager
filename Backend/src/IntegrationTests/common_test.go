package IntegrationTests

import (
	"Server/config"
	"Server/databaseAccess"
	"Server/routes"
	"Server/sessionManager"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/publicsuffix"
	"log"
	"net/http"
	"net/http/cookiejar"
)

type errorResponse struct {
	Error string `json:"error"`
}

type idResponse struct {
	Id float64 `json:"id"`
}

var router *gin.Engine
var client *http.Client

const baseUrl = "http://localhost:8080/"

func createRouterAndHttpClient() {
	newClient()

	conf := config.GetConfig("../conf.json")

	//start router/webapp
	router = gin.Default()
	databaseAccess.Init(conf)
	routes.UsersDAO = databaseAccess.CreateUsersDao()
	routes.LeaguesDAO = databaseAccess.CreateLeaguesDAO()
	routes.TeamsDAO = databaseAccess.CreateTeamsDAO()
	routes.GamesDAO = databaseAccess.CreateGamesDAO()
	routes.ElmSessions = sessionManager.CreateCookieSessionManager()

	routes.RegisterLoginHandlers(router.Group("/"))
	routes.RegisterUserHandlers(router.Group("/api/users"))
	routes.RegisterLeagueHandlers(router.Group("/api/leagues"))
	routes.RegisterTeamHandlers(router.Group("/api/teams"))
	routes.RegisterGameHandlers(router.Group("/api/games"))

	go router.Run(conf.GetPortString())
}

func newClient() {
	cookieJar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal(err)
	}

	client = &http.Client{
		Jar: cookieJar,
	}
}
