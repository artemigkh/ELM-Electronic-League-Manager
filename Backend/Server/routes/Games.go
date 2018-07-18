package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type GameInformation struct {
	Team1ID int `json:"team1Id"`
	Team2ID int `json:"team2Id"`
	GameTime int `json:"gameTime"`
}

func createNewGame(ctx *gin.Context) {
	//get parameters
	var gameInfo GameInformation
	err := ctx.ShouldBindJSON(&gameInfo)
	if checkJsonErr(ctx, err) {
		return
	}

	//must be logged in to create a game
	userID, err := ElmSessions.AuthenticateAndGetUserID(ctx)
	if checkErr(ctx, err) {
		return
	}
	if userID == -1 {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "notLoggedIn"})
		return
	}

	//must have an active game to create a team in it
	leagueId, err := ElmSessions.GetActiveLeague(ctx)
	if checkErr(ctx, err) {
		return
	}
	if leagueId == -1 {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "noActiveLeague"})
		return
	}
}

func RegisterGameHandlers(g *gin.RouterGroup) {
	g.POST("/", createNewGame)
}
