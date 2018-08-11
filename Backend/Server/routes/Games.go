package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type GameInformation struct {
	Team1Id  int `json:"team1Id"`
	Team2Id  int `json:"team2Id"`
	GameTime int `json:"gameTime"`
}

type GameReportInformation struct {
	WinnerId   int `json:"winnerId"`
	ScoreTeam1 int `json:"scoreTeam1"`
	ScoreTeam2 int `json:"scoreTeam2"`
}

//TODO: check if there exists the same player on both teams
func createNewGame(ctx *gin.Context) {
	//get parameters
	var gameInfo GameInformation
	err := ctx.ShouldBindJSON(&gameInfo)
	if checkJsonErr(ctx, err) {
		return
	}

	if failIfTeamDoesNotExist(ctx, gameInfo.Team1Id, ctx.GetInt("leagueId")) {
		return
	}
	if failIfTeamDoesNotExist(ctx, gameInfo.Team2Id, ctx.GetInt("leagueId")) {
		return
	}
	if failIfConflictExists(ctx, gameInfo.Team1Id, gameInfo.Team2Id, gameInfo.GameTime) {
		return
	}

	gameId, err := GamesDAO.CreateGame(ctx.GetInt("leagueId"), gameInfo.Team1Id, gameInfo.Team2Id, gameInfo.GameTime)
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": gameId})
}

func getGameInformation(ctx *gin.Context) {
	gameInformation, err := GamesDAO.GetGameInformation(ctx.GetInt("urlId"), ctx.GetInt("leagueId"))
	if checkErr(ctx, err) {
		return
	}

	if gameInformation == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "gameDoesNotExist"})
		return
	}

	ctx.JSON(http.StatusOK, gameInformation)
}

//TODO: check if the winner Id is one of the two team Ids in the game
func reportGameResult(ctx *gin.Context) {
	//get parameters
	var gameInfo GameReportInformation
	err := ctx.ShouldBindJSON(&gameInfo)
	if checkJsonErr(ctx, err) {
		return
	}

	if failIfGameDoesNotExist(ctx) {
		return
	}

	//report the result
	err = GamesDAO.ReportGame(ctx.GetInt("urlId"), ctx.GetInt("leagueId"),
		gameInfo.WinnerId, gameInfo.ScoreTeam1, gameInfo.ScoreTeam2)
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func RegisterGameHandlers(g *gin.RouterGroup) {
	g.Use(getActiveLeague())

	g.POST("/", authenticate(), createNewGame)
	g.POST("/report/:id", authenticate(), getUrlId(), getReportResultPermissions(), reportGameResult)
	g.GET("/:id", getUrlId(), getGameInformation)
}
