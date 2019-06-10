package routes

import (
	"Server/databaseAccess"
	"github.com/gin-gonic/gin"
	"net/http"
)

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/getLeagueGames
func getAllGamesInLeague(ctx *gin.Context) {
	games, err := GamesDAO.GetAllGamesInLeague(getLeagueId(ctx))
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, games)
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/createGame
func createNewGame(ctx *gin.Context) {
	hasPermissions, err := Access.Game(databaseAccess.Create, getLeagueId(ctx), 0, getUserId(ctx))
	if accessForbidden(ctx, hasPermissions, err) {
		return
	}

	var game databaseAccess.GameCreationInformation
	if bindAndCheckErr(ctx, &game) {
		return
	}

	if validator.DataInvalid(ctx, func() (bool, string, error) { return game.Validate(getLeagueId(ctx)) }) {
		return
	}

	gameId, err := GamesDAO.CreateGame(getLeagueId(ctx), game)
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"gameId": gameId})
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/getGameInfo
func getGameInformation(ctx *gin.Context) {
	hasPermissions, err := Access.Game(databaseAccess.View, getLeagueId(ctx), getGameId(ctx), getUserId(ctx))
	if accessForbidden(ctx, hasPermissions, err) {
		return
	}

	gameInfo, err := GamesDAO.GetGameInformation(getGameId(ctx))
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, gameInfo)
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/deleteGame
func deleteGame(ctx *gin.Context) {
	hasPermissions, err := Access.Game(databaseAccess.Delete, getLeagueId(ctx), getGameId(ctx), getUserId(ctx))
	if accessForbidden(ctx, hasPermissions, err) {
		return
	}

	err = GamesDAO.DeleteGame(getGameId(ctx))
	if checkErr(ctx, err) {
		return
	}

	ctx.Status(http.StatusOK)
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/rescheduleGame
func rescheduleGame(ctx *gin.Context) {
	hasPermissions, err := Access.Game(databaseAccess.Edit, getLeagueId(ctx), getGameId(ctx), getUserId(ctx))
	if accessForbidden(ctx, hasPermissions, err) {
		return
	}

	var gameTime databaseAccess.GameTime
	if bindAndCheckErr(ctx, &gameTime) {
		return
	}

	if validator.DataInvalid(ctx, func() (bool, string, error) { return gameTime.Validate(getLeagueId(ctx), getGameId(ctx)) }) {
		return
	}

	err = GamesDAO.RescheduleGame(getGameId(ctx), gameTime.GameTime)
	if checkErr(ctx, err) {
		return
	}

	ctx.Status(http.StatusOK)
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/reportGame
func reportGameResult(ctx *gin.Context) {
	hasPermissions, err := Access.Game(databaseAccess.Edit, getLeagueId(ctx), getGameId(ctx), getUserId(ctx))
	if accessForbidden(ctx, hasPermissions, err) {
		return
	}

	var gameResult databaseAccess.GameResult
	if bindAndCheckErr(ctx, &gameResult) {
		return
	}

	if validator.DataInvalid(ctx, func() (bool, string, error) { return gameResult.Validate(getLeagueId(ctx)) }) {
		return
	}

	err = GamesDAO.ReportGame(getGameId(ctx), gameResult)
	if checkErr(ctx, err) {
		return
	}

	ctx.Status(http.StatusOK)
}

func RegisterGameHandlers(g *gin.RouterGroup) {
	g.GET("/", getAllGamesInLeague)
	g.GET("/:gameId", storeGameId(), getGameInformation)
	g.POST("/", createNewGame)

	withId := g.Group("/:gameId").Use(storeGameId())
	withId.DELETE("", deleteGame)
	withId.POST("/reschedule", rescheduleGame)
	withId.POST("/report", reportGameResult)
}
