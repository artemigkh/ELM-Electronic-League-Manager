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

type PlayerInformationChange struct {
	TeamId         int    `json:"teamId"`
	PlayerId       int    `json:"playerId"`
	Name           string `json:"name"`
	GameIdentifier string `json:"gameIdentifier"` // Jersey Number, IGN, etc.
	MainRoster     bool   `json:"mainRoster"`
}

type PlayerRemoveInformation struct {
	TeamId   int `json:"teamId"`
	PlayerId int `json:"playerId"`
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
//TODO: Add minimum name and tag length (1 char)
func createNewTeam(ctx *gin.Context) {
	//get parameters
	var teamInfo TeamInformation
	err := ctx.ShouldBindJSON(&teamInfo)
	if checkJsonErr(ctx, err) {
		return
	}

	if failIfNameTooLong(ctx, teamInfo.Name) {
		return
	}
	if failIfTeamTagTooLong(ctx, teamInfo.Tag) {
		return
	}
	if failIfTeamInfoInUse(ctx, ctx.GetInt("leagueId"), teamInfo.Name, teamInfo.Tag) {
		return
	}

	teamId, err := TeamsDAO.CreateTeam(ctx.GetInt("leagueId"), ctx.GetInt("userId"), teamInfo.Name, teamInfo.Tag)
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": teamId})
}

/**
 * @api{DELETE} /api/teams/removeTeam/:id Delete Team
 * @apiName deleteTeam
 * @apiGroup Teams
 * @apiDescription Delete a team from the current league and its players
 *
 * @apiParam {int} id The unique numerical identifier of the team
 *
 * @apiError IdMustBeInteger The id in the url must be an integer value
 * @apiError notLoggedIn No user is logged in
 * @apiError noActiveLeague There is no active league selected
 * @apiError teamDoesNotExist The specified team does not exist
 * @apiError teamIsActive This team cannot be deleted because it has played games in this league
 * @apiError noEditTeamPermissions The currently logged in user does not have permissions to edit teams in this league
 */
func deleteTeam(ctx *gin.Context) {
	if failIfTeamDoesNotExist(ctx, ctx.GetInt("leagueId"), ctx.GetInt("urlId")) {
		return
	}

	err := TeamsDAO.DeleteTeam(ctx.GetInt("leagueId"), ctx.GetInt("urlId"))
	if checkErr(ctx, err) {
		return
	}

	ctx.Status(http.StatusOK)
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
	if failIfTeamDoesNotExist(ctx, ctx.GetInt("leagueId"), ctx.GetInt("urlId")) {
		return
	}

	teamInfo, err := TeamsDAO.GetTeamInformation(ctx.GetInt("leagueId"), ctx.GetInt("urlId"))
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
 * @apiParam {bool} mainRoster If true, this player is on the main roster, otherwise is a substitute
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
//TODO: Add minimum name lengths (1 char each)
func addPlayerToTeam(ctx *gin.Context) {
	//get parameters
	var playerInfo PlayerInformation
	err := ctx.ShouldBindJSON(&playerInfo)
	if checkJsonErr(ctx, err) {
		return
	}

	if failIfTeamDoesNotExist(ctx, ctx.GetInt("leagueId"), playerInfo.TeamId) {
		return
	}
	if failIfCannotEditPlayersOnTeam(ctx, ctx.GetInt("leagueId"), playerInfo.TeamId, ctx.GetInt("userId")) {
		return
	}
	if failIfGameIdentifierTooLong(ctx, playerInfo.GameIdentifier) {
		return
	}
	if failIfNameTooLong(ctx, playerInfo.Name) {
		return
	}
	if failIfGameIdentifierInUse(ctx, ctx.GetInt("leagueId"), playerInfo.TeamId, playerInfo.GameIdentifier) {
		return
	}

	playerId, err := TeamsDAO.AddNewPlayer(playerInfo.TeamId, playerInfo.GameIdentifier,
		playerInfo.Name, playerInfo.MainRoster)
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": playerId})
}

/**
 * @api{DELETE} /api/teams/removePlayer Remove Player From Team
 * @apiGroup Teams
 * @apiDescription Remove a player from a teams roster
 *
 * @apiParam {int} teamId The unique numerical identifier of the team the player is to be added to
 * @apiParam {int} playerId The unique numerical identifier of the player to be removed
 *
 * @apiError notLoggedIn No user is logged in
 * @apiError noActiveLeague There is no active league selected
 * @apiError teamDoesNotExist The specified team does not exist
 * @apiError canNotEditPlayers The currently logged in player does not have permission to edit the players on this team
 * @apiError playerDoesNotExist The specified player does not exist on this team
 */
func removePlayerFromTeam(ctx *gin.Context) {
	//get parameters
	var playerRemoveInfo PlayerRemoveInformation
	err := ctx.ShouldBindJSON(&playerRemoveInfo)
	if checkJsonErr(ctx, err) {
		return
	}
	if failIfTeamDoesNotExist(ctx, ctx.GetInt("leagueId"), playerRemoveInfo.TeamId) {
		return
	}
	if failIfCannotEditPlayersOnTeam(ctx, ctx.GetInt("leagueId"), playerRemoveInfo.TeamId, ctx.GetInt("userId")) {
		return
	}
	if failIfPlayerDoesNotExist(ctx, playerRemoveInfo.TeamId, playerRemoveInfo.PlayerId) {
		return
	}
	err = TeamsDAO.RemovePlayer(playerRemoveInfo.TeamId, playerRemoveInfo.PlayerId)
	if checkErr(ctx, err) {
		return
	}

	ctx.Status(http.StatusOK)
}

/**
 * @api{put} /api/teams/updatePlayer Update Player Information
 * @apiGroup Teams
 * @apiDescription Change a players information
 *
 * @apiParam {int} teamId The unique numerical identifier of the team the player is to be added to
 * @apiParam {int} playerId The unique numerical identifier of the player
 * @apiParam {string} name The updated name of the player (can be left blank)
 * @apiParam {string} gameIdentifier The updated in-game name identifier of the player (jersey number, ign, etc.)
 * @apiParam {bool} mainRoster If true, this player is on the main roster, otherwise is a substitute
 *
 * @apiError notLoggedIn No user is logged in
 * @apiError noActiveLeague There is no active league selected
 * @apiError teamDoesNotExist The specified team does not exist
 * @apiError canNotEditPlayers The currently logged in player does not have permission to edit the players on this team
 * @apiError gameIdentifierTooLong The game identifier exceeds 50 characters
 * @apiError nameTooLong The name exceeds 50 characters
 * @apiError gameIdentifierInUse This game identifier is already in use in this league
 */
func updatePlayer(ctx *gin.Context) {
	//get parameters
	var playerInfoChange PlayerInformationChange
	err := ctx.ShouldBindJSON(&playerInfoChange)
	if checkJsonErr(ctx, err) {
		return
	}

	if failIfTeamDoesNotExist(ctx, ctx.GetInt("leagueId"), playerInfoChange.TeamId) {
		return
	}
	if failIfCannotEditPlayersOnTeam(ctx, ctx.GetInt("leagueId"), playerInfoChange.TeamId, ctx.GetInt("userId")) {
		return
	}
	if failIfGameIdentifierTooLong(ctx, playerInfoChange.GameIdentifier) {
		return
	}
	if failIfNameTooLong(ctx, playerInfoChange.Name) {
		return
	}
	if failIfGameIdentifierInUse(ctx, ctx.GetInt("leagueId"), playerInfoChange.TeamId, playerInfoChange.GameIdentifier) {
		return
	}

	err = TeamsDAO.UpdatePlayer(playerInfoChange.TeamId, playerInfoChange.PlayerId, playerInfoChange.GameIdentifier,
		playerInfoChange.Name, playerInfoChange.MainRoster)
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func RegisterTeamHandlers(g *gin.RouterGroup) {
	g.Use(getActiveLeague())

	g.POST("/", authenticate(), failIfNoTeamCreatePermissions(), createNewTeam)
	g.POST("/addPlayer", authenticate(), addPlayerToTeam)
	g.DELETE("/removePlayer", authenticate(), removePlayerFromTeam)
	g.PUT("/updatePlayer", authenticate(), updatePlayer)
	g.GET("/:id", getUrlId(), getTeamInformation)
	g.DELETE("/removeTeam/:id", getUrlId(), authenticate(), failIfTeamActive(), failIfCannotEditTeam(), deleteTeam)
}
