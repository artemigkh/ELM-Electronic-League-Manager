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

/**
 * @api{POST} /api/teams Create New Team
 * @apiName createNewTeam
 * @apiGroup Teams
 * @apiDescription Register a new team
 *
 * @apiParam {string} name The name of the team to be created
 * @apiParam {string} tag The tag of the team to be created
 *
 * @apiSuccess {int} id the unique numerical identifier of the created team
 *
 * @apiError notLoggedIn No user is logged in
 * @apiError noEditTeamPermissions The currently logged in user does not have permissions to edit teams in this league
 * @apiError nameTooLong The team name has exceeded 50 characters
 * @apiError tagTooLong The team tag has exceeded 5 characters
 * @apiError nameInUse The team name is currently in use
 * @apiError tagInUse The team tag is currently in use
 */
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

/**
 * @api{GET} /api/teams/:id Get Team Information
 * @apiGroup Teams
 * @apiDescription Get information about the team with specified id
 *
 * @apiParam {int} id The unique numerical identifier of the team
 *
 * @apiSuccess {string} name The name of the team
 * @apiSuccess {string} tag The tag of the team
 * @apiSuccess {int} wins The number of wins this team has
 * @apiSuccess {int} losses The number of losses this team has
 * @apiSuccess {[]Object} players An array of json objects representing the players on the team
 * @apiSuccess {int} players.id The unique numerical identifier of the player
 * @apiSuccess {string} players.name The name of the player
 * @apiSuccess {string} players.gameIdentifier The in-game name identifier of the player (jersey number, ign, etc.)
 * @apiSuccess {bool} players.mainRoster If true, this player is on the main roster, otherwise is a substitute
 *
 * @apiError IdMustBeInteger The id in the url must be an integer value
 * @apiError noActiveLeague There is no active league selected
 * @apiError teamDoesNotExist The specified team does not exist
 */
func getTeamInformation(ctx *gin.Context) {
	if failIfTeamDoesNotExist(ctx, ctx.GetInt("urlId"), ctx.GetInt("leagueId")) {
		return
	}

	teamInfo, err := TeamsDAO.GetTeamInformation(ctx.GetInt("urlId"), ctx.GetInt("leagueId"))
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, teamInfo)
}

/**
 * @api{POST} /api/teams/addPlayer Add Player To Team
 * @apiGroup Teams
 * @apiDescription Create a new player and add him to the teams roster
 *
 * @apiParam {int} teamId The unique numerical identifier of the team the player is to be added to
 * @apiParam {string} name The name of the player (can be left blank)
 * @apiParam {string} gameIdentifier The in-game name identifier of the player (jersey number, ign, etc.)
 * @apiParam {bool} players.mainRoster If true, this player is on the main roster, otherwise is a substitute
 *
 * @apiSuccess {int} id the unique numerical identifier of the created player
 *
 * @apiError notLoggedIn No user is logged in
 * @apiError noActiveLeague There is no active league selected
 * @apiError teamDoesNotExist The specified team does not exist
 * @apiError canNotEditPlayers The currently logged in player does not have permission to edit the players on this team
 * @apiError gameIdentifierTooLong The game identifier exceeds 50 characters
 * @apiError nameTooLong The name exceeds 50 characters
 * @apiError gameIdentifierInUse This game identifier is already in use in this league
 */
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
