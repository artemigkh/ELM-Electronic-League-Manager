package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type TeamInformation struct {
	Name string `json:"name"`
	Tag  string `json:"tag"`
}

type PlayerInformation struct {
	TeamId         int    `json:"teamId"`
	Name           string `json:"name"`
	GameIdentifier string `json:"gameIdentifier"` // Jersey Number, IGN, etc.
	MainRoster     bool   `json:"mainRoster"`
}

func createNewTeam(ctx *gin.Context) {
	//get parameters
	var teamInfo TeamInformation
	err := ctx.ShouldBindJSON(&teamInfo)
	if checkJsonErr(ctx, err) {
		return
	}

	if failIfTeamNameTooLong(ctx, teamInfo.Name) {
		return
	}
	if failIfTeamTagTooLong(ctx, teamInfo.Tag) {
		return
	}
	if failIfTeamInfoInUse(ctx, teamInfo.Name, teamInfo.Tag, ctx.GetInt("leagueId")) {
		return
	}

	teamId, err := TeamsDAO.CreateTeam(ctx.GetInt("leagueId"), ctx.GetInt("userId"), teamInfo.Name, teamInfo.Tag)
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": teamId})
}

func getTeamInformation(ctx *gin.Context) {
	println("got here 1")
	if failIfTeamDoesNotExist(ctx, ctx.GetInt("urlId"), ctx.GetInt("leagueId")) {
		return
	}
	println("got here 2")
	teamInfo, err := TeamsDAO.GetTeamInformation(ctx.GetInt("urlId"), ctx.GetInt("leagueId"))
	if checkErr(ctx, err) {
		return
	}
	println("got here 3")
	ctx.JSON(http.StatusOK, teamInfo)
}

func addPlayerToTeam(ctx *gin.Context) {
	//get parameters
	var playerInfo PlayerInformation
	err := ctx.ShouldBindJSON(&playerInfo)
	if checkJsonErr(ctx, err) {
		return
	}

	if failIfTeamDoesNotExist(ctx, playerInfo.TeamId, ctx.GetInt("leagueId")) {
		return
	}
	if failIfCannotEditPlayersOnTeam(ctx, ctx.GetInt("userId"), playerInfo.TeamId, ctx.GetInt("leagueId")) {
		return
	}
	if failIfGameIdentifierTooLong(ctx, playerInfo.GameIdentifier) {
		return
	}
	if failIfNameTooLong(ctx, playerInfo.Name) {
		return
	}
	if failIfGameIdentifierInUse(ctx, playerInfo.GameIdentifier, playerInfo.TeamId, ctx.GetInt("leagueId")) {
		return
	}

	playerId, err := TeamsDAO.AddNewPlayer(playerInfo.TeamId, playerInfo.GameIdentifier,
		playerInfo.Name, playerInfo.MainRoster)
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": playerId})
}

func RegisterTeamHandlers(g *gin.RouterGroup) {
	g.Use(getActiveLeague())

	g.POST("/", authenticate(), getTeamEditPermissions(), createNewTeam)
	g.POST("/addPlayer", authenticate(), addPlayerToTeam)
	g.GET("/:id", getUrlId(), getTeamInformation)
}
