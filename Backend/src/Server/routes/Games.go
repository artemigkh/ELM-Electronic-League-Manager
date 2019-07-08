package routes

import (
	"Server/databaseAccess"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/getLeagueGames
func getAllGamesInLeague(ctx *gin.Context) {
	games, err := GamesDAO.GetAllGamesInLeague(getLeagueId(ctx))
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, games)
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/getSortedGames
func getSortedGames() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		teamIdString := ctx.DefaultQuery("teamId", "0")
		teamId, err := strconv.Atoi(teamIdString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "IdMustBeInteger"})
		}

		limitString := ctx.DefaultQuery("limit", "0")
		limit, err := strconv.Atoi(limitString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "limitMustBeInteger"})
		}

		games, err := GamesDAO.GetSortedGames(getLeagueId(ctx), teamId, limit)
		if checkErr(ctx, err) {
			return
		}

		ctx.JSON(http.StatusOK, games)
	}
}

func getGamesByWeek() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		timeZoneOffsetString := ctx.DefaultQuery("timeZoneOffset", "0")
		timeZoneOffset, err := strconv.Atoi(timeZoneOffsetString)
		fmt.Printf("timezone offset is %v", timeZoneOffset)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "IdMustBeInteger"})
		}
		games, err := GamesDAO.GetGamesByWeek(getLeagueId(ctx), timeZoneOffset)
		if checkErr(ctx, err) {
			return
		}

		ctx.JSON(http.StatusOK, games)
	}
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/getGame
func getGameInformation() gin.HandlerFunc {
	return endpoint{
		Entity:     Game,
		AccessType: View,
		Core:       func(ctx *gin.Context) (interface{}, error) { return GamesDAO.GetGameInformation(getGameId(ctx)) },
	}.createEndpointHandler()
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/createGame
func createNewGame() gin.HandlerFunc {
	var game databaseAccess.GameCreationInformation
	return endpoint{
		Entity:        Game,
		AccessType:    Create,
		BindData:      func(ctx *gin.Context) bool { return bindAndCheckErr(ctx, &game) },
		IsDataInvalid: func(ctx *gin.Context) (bool, string, error) { return game.Validate(getLeagueId(ctx)) },
		Core: func(ctx *gin.Context) (interface{}, error) {
			gameId, err := GamesDAO.CreateGame(getLeagueId(ctx), game)
			return gin.H{"gameId": gameId}, err
		},
	}.createEndpointHandler()
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/deleteGame
func deleteGame() gin.HandlerFunc {
	return endpoint{
		Entity:     Game,
		AccessType: Delete,
		Core: func(ctx *gin.Context) (interface{}, error) {
			return nil, GamesDAO.DeleteGame(getGameId(ctx))
		},
	}.createEndpointHandler()
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/rescheduleGame
func rescheduleGame() gin.HandlerFunc {
	var gameTime databaseAccess.GameTime
	return endpoint{
		Entity:     Game,
		AccessType: Edit,
		BindData:   func(ctx *gin.Context) bool { return bindAndCheckErr(ctx, &gameTime) },
		IsDataInvalid: func(ctx *gin.Context) (bool, string, error) {
			return gameTime.Validate(getLeagueId(ctx), getGameId(ctx))
		},
		Core: func(ctx *gin.Context) (interface{}, error) {
			return nil, GamesDAO.RescheduleGame(getGameId(ctx), gameTime.GameTime)
		},
	}.createEndpointHandler()
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/reportGame
func reportGameResult() gin.HandlerFunc {
	var gameResult databaseAccess.GameResult
	return endpoint{
		Entity:        Game,
		AccessType:    Edit,
		BindData:      func(ctx *gin.Context) bool { return bindAndCheckErr(ctx, &gameResult) },
		IsDataInvalid: func(ctx *gin.Context) (bool, string, error) { return gameResult.Validate(getLeagueId(ctx)) },
		Core: func(ctx *gin.Context) (interface{}, error) {
			fmt.Printf("%+v\n", gameResult)
			return nil, GamesDAO.ReportGame(getGameId(ctx), gameResult)
		},
	}.createEndpointHandler()
}

func RegisterGameHandlers(g *gin.RouterGroup) {
	g.GET("/api/v1/sortedGames", getSortedGames())
	g.GET("/api/v1/gamesByWeek", getGamesByWeek())
	games := g.Group("/api/v1/games")
	games.GET("", getAllGamesInLeague)
	games.GET("/:gameId", storeGameId(), getGameInformation())
	games.POST("", createNewGame())

	withId := games.Group("/:gameId", storeGameId())
	withId.DELETE("", deleteGame())
	withId.POST("/reschedule", rescheduleGame())
	withId.POST("/report", reportGameResult())
}
