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

type GameRescheduleInformation struct {
	Id       int `json:"id"`
	GameTime int `json:"gameTime"`
}

/**
 * @api{POST} /api/games/ Create New Game
 * @apiGroup Games
 * @apiDescription Schedule a new game
 *
 * @apiParam {int} team1Id The unique numerical identifier of the team in position 1
 * @apiParam {int} team2Id The unique numerical identifier of the team in position 2
 * @apiParam {int} gameTime The unix time of when the game is scheduled for
 *
 * @apiSuccess {int} id The primary id of the created game
 *
 * @apiError notLoggedIn No user is logged in
 * @apiError teamsAreSame Team 1 and Team 2 are the same team
 * @apiError noActiveLeague There is no active league selected
 * @apiError teamDoesNotExist One of the teams specified does not exist
 * @apiError conflictExists One of the teams already has a game scheduled at this time
 */
func createNewGame(ctx *gin.Context) {
	//get parameters
	var gameInfo GameInformation
	err := ctx.ShouldBindJSON(&gameInfo)
	if checkJsonErr(ctx, err) {
		return
	}

	if gameInfo.Team1Id == gameInfo.Team2Id {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "teamsAreSame"})
		return
	}
	if failIfTeamDoesNotExist(ctx, ctx.GetInt("leagueId"), gameInfo.Team1Id) {
		return
	}
	if failIfTeamDoesNotExist(ctx, ctx.GetInt("leagueId"), gameInfo.Team2Id) {
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

/**
 * @api{GET} /api/games/:id Get Game Information
 * @apiGroup Games
 * @apiDescription Get the information about the game with specified id
 *
 * @apiParam {int} id The unique numerical identifier of the game
 *
 * @apiSuccess {int} id The unique numerical identifier of the game
 * @apiSuccess {int} leagueId The unique numerical identifier of the league the game is in
 * @apiSuccess {int} team1Id The unique numerical identifier of the team in position 1
 * @apiSuccess {int} team2Id The unique numerical identifier of the team in position 2
 * @apiSuccess {int} gameTime The unix time of when the game is scheduled for
 * @apiSuccess {bool} complete Whether or not the game has been played and recorded
 * @apiSuccess {int} winnerId The unique numerical identifier of the team that won; -1 if game is not complete
 * @apiSuccess {int} scoreTeam1 The score of the team in position 1
 * @apiSuccess {int} scoreTeam2 The score of the team in position 2
 *
 * @apiError IdMustBeInteger The id in the url must be an integer value
 * @apiError noActiveLeague There is no active league selected
 * @apiError gameDoesNotExist The game with specified id does not exist
 */
func getGameInformation(ctx *gin.Context) {
	gameInformation, err := GamesDAO.GetGameInformation(ctx.GetInt("leagueId"), ctx.GetInt("urlId"))
	if checkErr(ctx, err) {
		return
	}

	if gameInformation == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "gameDoesNotExist"})
		return
	}

	ctx.JSON(http.StatusOK, gameInformation)
}

/**
 * @api{POST} /api/games/report/:id Report Game Result
 * @apiGroup Games
 * @apiDescription Report the result of a scheduled game
 *
 * @apiParam {int} id The unique numerical identifier of the game
 * @apiParam {int} winnerId The unique numerical identifier of the winning team
 * @apiParam {int} scoreTeam1 The score of the team in position 1
 * @apiParam {int} scoreTeam2 The score of the team in position 2
 *
 * @apiError notLoggedIn No user is logged in
 * @apiError IdMustBeInteger The id in the url must be an integer value
 * @apiError noActiveLeague There is no active league selected
 * @apiError gameDoesNotExist The game with specified id does not exist
 * @apiError gameDoesNotContainWinner The game does not contain the specified winner ID
 * @apiError noReportResultPermissions The currently logged in user does not have permissions to report results for this team
 */
func reportGameResult(ctx *gin.Context) {
	//get parameters
	var gameInfo GameReportInformation
	err := ctx.ShouldBindJSON(&gameInfo)
	if checkJsonErr(ctx, err) {
		return
	}

	if failIfGameDoesNotExist(ctx, ctx.GetInt("leagueId"), ctx.GetInt("urlId")) {
		return
	}
	if failIfGameDoesNotContainWinner(ctx, ctx.GetInt("leagueId"), ctx.GetInt("urlId"), gameInfo.WinnerId) {
		return
	}

	//report the result
	err = GamesDAO.ReportGame(ctx.GetInt("leagueId"), ctx.GetInt("urlId"),
		gameInfo.WinnerId, gameInfo.ScoreTeam1, gameInfo.ScoreTeam2)
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

/**
 * @api{DELETE} /api/games/:id Remove Game
 * @apiGroup Games
 * @apiDescription Unschedule a game
 *
 * @apiParam {int} id The unique numerical identifier of game
 *
 * @apiError notLoggedIn No user is logged in
 * @apiError IdMustBeInteger The specified Id in the URL must be an integer
 * @apiError noActiveLeague There is no active league selected
 * @apiError noEditSchedulePermissions The currently logged in user does not have permissions to edit the schedule
 * @apiError gameDoesNotExist The game with specified id does not exist in this league
 */
func deleteGame(ctx *gin.Context) {
	if failIfGameDoesNotExist(ctx, ctx.GetInt("leagueId"), ctx.GetInt("urlId")) {
		return
	}

	err := GamesDAO.DeleteGame(ctx.GetInt("leagueId"), ctx.GetInt("urlId"))
	if checkErr(ctx, err) {
		return
	}
}

/**
 * @api{PUT} /api/games/ Reschedule Game
 * @apiGroup Games
 * @apiDescription Reschedule a game
 *
 * @apiParam {int} id The unique numerical identifier of game
 * @apiParam {int} gameTime The unix time of when the game is scheduled for
 *
 * @apiError notLoggedIn No user is logged in
 * @apiError noActiveLeague There is no active league selected
 * @apiError noEditSchedulePermissions The currently logged in user does not have permissions to edit the schedule
 * @apiError gameDoesNotExist The game with specified id does not exist in this league
 * @apiError gameIsComplete The game with specified id is already complete
 * @apiError conflictExists One of the teams already has a game scheduled at this time
 */
func rescheduleGame(ctx *gin.Context) {
	//get parameters
	var gameInfo GameRescheduleInformation
	err := ctx.ShouldBindJSON(&gameInfo)
	if checkJsonErr(ctx, err) {
		return
	}

	if failIfGameDoesNotExist(ctx, ctx.GetInt("leagueId"), gameInfo.Id) {
		return
	}

	// get the game information to get the two team Ids to check for conflicts
	gameInformation, err := GamesDAO.GetGameInformation(ctx.GetInt("leagueId"), gameInfo.Id)
	if checkErr(ctx, err) {
		return
	}

	if gameInformation.Complete {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "gameIsComplete"})
		return
	}

	if failIfConflictExists(ctx, gameInformation.Team1Id, gameInformation.Team2Id, gameInfo.GameTime) {
		return
	}

	err = GamesDAO.RescheduleGame(ctx.GetInt("leagueId"), gameInfo.Id, gameInfo.GameTime)
	if checkErr(ctx, err) {
		return
	}

	ctx.Status(http.StatusOK)
}

func RegisterGameHandlers(g *gin.RouterGroup) {
	g.Use(getActiveLeague())

	g.POST("/", authenticate(), failIfNoEditSchedulePermissions(), createNewGame)
	g.PUT("/", authenticate(), failIfNoEditSchedulePermissions(), rescheduleGame)
	g.POST("/report/:id", authenticate(), getUrlId(), failIfNoReportResultPermissions(), reportGameResult)
	g.GET("/:id", getUrlId(), getGameInformation)
	g.DELETE("/:id", authenticate(), getUrlId(), failIfNoEditSchedulePermissions(), deleteGame)
}
