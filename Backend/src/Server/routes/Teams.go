package routes

import (
	"Server/databaseAccess"
	"github.com/gin-gonic/gin"
	"net/http"
)

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/createTeam
func createNewTeam() gin.HandlerFunc {
	var team databaseAccess.TeamCore
	return endpoint{
		Entity:        Team,
		AccessType:    Create,
		BindData:      func(ctx *gin.Context) bool { return bindAndCheckErr(ctx, &team) },
		IsDataInvalid: func(ctx *gin.Context) (bool, string, error) { return team.ValidateNew(getLeagueId(ctx)) },
		Core: func(ctx *gin.Context) (interface{}, error) {
			teamId, err := TeamsDAO.CreateTeam(getLeagueId(ctx), getUserId(ctx), team)
			return gin.H{"teamId": teamId}, err
		},
	}.createEndpointHandler()
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/getLeagueTeams
func getAllTeams() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		games, err := TeamsDAO.GetAllTeamsInLeague(getLeagueId(ctx))
		if checkErr(ctx, err) {
			return
		}

		ctx.JSON(http.StatusOK, games)
	}
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/getTeam
func getTeamInfo() gin.HandlerFunc {
	return endpoint{
		Entity:     Team,
		AccessType: View,
		Core: func(ctx *gin.Context) (interface{}, error) {
			return TeamsDAO.GetTeamInformation(getTeamId(ctx))
		},
	}.createEndpointHandler()
}

//
//func editTeam() gin.HandlerFunc {
//
//}
//
//func deleteTeam() gin.HandlerFunc {
//
//}
//
//func editTeamManagerPermissions() gin.HandlerFunc {
//
//}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/createPlayer
func createNewPlayer() gin.HandlerFunc {
	var player databaseAccess.PlayerCore
	return endpoint{
		Entity:     Player,
		AccessType: Create,
		BindData:   func(ctx *gin.Context) bool { return bindAndCheckErr(ctx, &player) },
		IsDataInvalid: func(ctx *gin.Context) (bool, string, error) {
			return player.ValidateNew(getLeagueId(ctx), getTeamId(ctx))
		},
		Core: func(ctx *gin.Context) (interface{}, error) {
			playerId, err := TeamsDAO.CreatePlayer(getLeagueId(ctx), getTeamId(ctx), player)
			return gin.H{"playerId": playerId}, err
		},
	}.createEndpointHandler()
}

//func editPlayer() gin.HandlerFunc {
//
//}
//
//func deletePlayer() gin.HandlerFunc {
//
//}

//
///**
// * @api{POST} /api/teams/withIcon Create New Team With Icon
// * @apiName createNewTeamWithIcon
// * @apiGroup Teams
// * @apiDescription Register a new team with icon
// *
// * @apiParam {string} name The name of the team to be created in form
// * @apiParam {string} tag The tag of the team to be created in form
// * @apiParam {string} description The description of the team to be created in form
// * @apiParam {File} icon The icon png as multipart/form-data
// *
// * @apiSuccess {int} id the unique numerical identifier of the created team
// *
// * @apiError notLoggedIn No user is logged in
// * @apiError noActiveLeague There is no active league selected
// * @apiError noEditTeamPermissions The currently logged in user does not have permissions to edit teams in this league
// * @apiError nameTooLong The team name has exceeded 50 characters
// * @apiError tagTooLong The team tag has exceeded 5 characters
// * @apiError nameTooShort The name must be at least 2 characters in length
// * @apiError tagTooShort The tag must be at least 2 characters in length
// * @apiError nameInUse The team name is currently in use
// * @apiError tagInUse The team tag is currently in use
// * @apiError iconError There was an error while processing the icon image png file
// */
//func createNewTeamWithIcon(ctx *gin.Context) {
//	//get parameters
//	var teamInfo databaseAccess.TeamDTO
//	teamInfo.Name = ctx.PostForm("name")
//	teamInfo.Tag = ctx.PostForm("tag")
//	teamInfo.Description = ctx.PostForm("description")
//
//	if teamInfo.Name == "" || teamInfo.Tag == "" {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": "malformedInput"})
//		return
//	}
//
//	if failIfNameTooLong(ctx, teamInfo.Name) {
//		return
//	}
//	if failIfTeamTagTooLong(ctx, teamInfo.Tag) {
//		return
//	}
//	if failIfNameTooShort(ctx, teamInfo.Name) {
//		return
//	}
//	if failIfTagTooShort(ctx, teamInfo.Tag) {
//		return
//	}
//	if failIfDescriptionTooLong(ctx, teamInfo.Description) {
//		return
//	}
//	if failIfTeamInfoInUse(ctx, ctx.GetInt("getLeagueId"), -1, teamInfo.Name, teamInfo.Tag) {
//		return
//	}
//
//	_, err := ctx.FormFile("icon")
//	if err == nil {
//		smallIcon, largeIcon, err := IconManager.StoreNewIcon(ctx)
//		if err != nil {
//			ctx.JSON(http.StatusBadRequest, gin.H{"error": "iconError"})
//			return
//		}
//
//		teamInfo.IconSmall = smallIcon
//		teamInfo.IconLarge = largeIcon
//
//		getTeamId, err := TeamsDAO.CreateTeamWithIcon(ctx.GetInt("getLeagueId"), ctx.GetInt("getUserId"),
//			teamInfo)
//		if checkErr(ctx, err) {
//			return
//		}
//
//		ctx.JSON(http.StatusOK, gin.H{"id": getTeamId})
//	} else {
//		getTeamId, err := TeamsDAO.CreateTeam(
//			ctx.GetInt("getLeagueId"), ctx.GetInt("getUserId"), teamInfo)
//		if checkErr(ctx, err) {
//			return
//		}
//
//		ctx.JSON(http.StatusOK, gin.H{"id": getTeamId})
//	}
//}
//
///**
// * @api{PUT} /api/teams/updateTeamWithIcon/:id Update Team Information
// * @apiName updateTeam
// * @apiGroup Teams
// * @apiDescription Change Team Information
// *
// * @apiParam {int} id The unique numerical identifier of the team
// * @apiParam {string} name The updated name of the team
// * @apiParam {string} tag The updated tag of the team
// * @apiParam {string} description The description of the team to be created
// * @apiParam {File} icon The icon png as multipart/form-data
// *
// * @apiError notLoggedIn No user is logged in
// * @apiError noActiveLeague There is no active league selected
// * @apiError IdMustBeInteger The id in the url must be an integer value
// * @apiError teamDoesNotExist The specified team does not exist
// * @apiError noEditTeamInformationPermissions The currently logged in user does not have permissions to edit this team information
// * @apiError nameTooLong The team name has exceeded 50 characters
// * @apiError tagTooLong The team tag has exceeded 5 characters
// * @apiError nameInUse The team name is currently in use
// * @apiError tagInUse The team tag is currently in use
// * @apiError iconError There was an error while processing the icon image png file
// */
//func updateTeamWithIcon(ctx *gin.Context) {
//	//get parameters
//	var teamInfo databaseAccess.TeamDTO
//	teamInfo.Name = ctx.PostForm("name")
//	teamInfo.Tag = ctx.PostForm("tag")
//	teamInfo.Description = ctx.PostForm("description")
//
//	if teamInfo.Name == "" || teamInfo.Tag == "" {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": "malformedInput"})
//		return
//	}
//
//	if failIfTeamDoesNotExist(ctx, ctx.GetInt("getLeagueId"), ctx.GetInt("urlId")) {
//		return
//	}
//	if failIfNameTooLong(ctx, teamInfo.Name) {
//		return
//	}
//	if failIfTeamTagTooLong(ctx, teamInfo.Tag) {
//		return
//	}
//	if failIfNameTooShort(ctx, teamInfo.Name) {
//		return
//	}
//	if failIfTagTooShort(ctx, teamInfo.Tag) {
//		return
//	}
//	if failIfDescriptionTooLong(ctx, teamInfo.Description) {
//		return
//	}
//	if failIfTeamInfoInUse(ctx, ctx.GetInt("getLeagueId"), ctx.GetInt("urlId"), teamInfo.Name, teamInfo.Tag) {
//		return
//	}
//
//	_, err := ctx.FormFile("icon")
//	if err == nil {
//		smallIcon, largeIcon, err := IconManager.StoreNewIcon(ctx)
//		if err != nil {
//			ctx.JSON(http.StatusBadRequest, gin.H{"error": "iconError"})
//			return
//		}
//
//		err = TeamsDAO.UpdateTeamIcon(ctx.GetInt("urlId"), smallIcon, largeIcon)
//		if checkErr(ctx, err) {
//			return
//		}
//	}
//
//	err = TeamsDAO.UpdateTeam(teamInfo)
//	if checkErr(ctx, err) {
//		return
//	}
//
//	ctx.Status(http.StatusOK)
//}
//
///**
// * @api{PUT} /api/teams/updateTeam/:id Update Team Information
// * @apiName updateTeam
// * @apiGroup Teams
// * @apiDescription Change Team Information
// *
// * @apiParam {int} id The unique numerical identifier of the team
// * @apiParam {string} name The updated name of the team
// * @apiParam {string} tag The updated tag of the team
// * @apiParam {string} description The description of the team to be created
// *
// * @apiError notLoggedIn No user is logged in
// * @apiError noActiveLeague There is no active league selected
// * @apiError IdMustBeInteger The id in the url must be an integer value
// * @apiError teamDoesNotExist The specified team does not exist
// * @apiError noEditTeamInformationPermissions The currently logged in user does not have permissions to edit this team information
// * @apiError nameTooLong The team name has exceeded 50 characters
// * @apiError tagTooLong The team tag has exceeded 5 characters
// * @apiError nameInUse The team name is currently in use
// * @apiError tagInUse The team tag is currently in use
// */
//func updateTeam(ctx *gin.Context) {
//	//get parameters
//	var teamInfo databaseAccess.TeamDTO
//	err := ctx.ShouldBindJSON(&teamInfo)
//	if checkJsonErr(ctx, err) {
//		return
//	}
//
//	if failIfTeamDoesNotExist(ctx, ctx.GetInt("getLeagueId"), ctx.GetInt("urlId")) {
//		return
//	}
//	if failIfNameTooLong(ctx, teamInfo.Name) {
//		return
//	}
//	if failIfTeamTagTooLong(ctx, teamInfo.Tag) {
//		return
//	}
//	if failIfNameTooShort(ctx, teamInfo.Name) {
//		return
//	}
//	if failIfTagTooShort(ctx, teamInfo.Tag) {
//		return
//	}
//	if failIfDescriptionTooLong(ctx, teamInfo.Description) {
//		return
//	}
//	if failIfTeamInfoInUse(ctx, ctx.GetInt("getLeagueId"), ctx.GetInt("urlId"), teamInfo.Name, teamInfo.Tag) {
//		return
//	}
//
//	err = TeamsDAO.UpdateTeam(teamInfo)
//	if checkErr(ctx, err) {
//		return
//	}
//
//	ctx.Status(http.StatusOK)
//}
//
///**
// * @api{PUT} /api/teams/updateTeamIcon/:id Update Team Icon
// * @apiName updateTeamIcon
// * @apiGroup Teams
// * @apiDescription Change Team Icon
// *
// * @apiParam {int} id The unique numerical identifier of the team
// * @apiParam {File} icon The icon png as multipart/form-data
// *
// * @apiError notLoggedIn No user is logged in
// * @apiError noActiveLeague There is no active league selected
// * @apiError IdMustBeInteger The id in the url must be an integer value
// * @apiError teamDoesNotExist The specified team does not exist
// * @apiError noEditTeamInformationPermissions The currently logged in user does not have permissions to edit this team information
// * @apiError iconError There was an error while processing the icon image png file
// */
//func updateTeamIcon(ctx *gin.Context) {
//	if failIfTeamDoesNotExist(ctx, ctx.GetInt("getLeagueId"), ctx.GetInt("urlId")) {
//		return
//	}
//
//	smallIcon, largeIcon, err := IconManager.StoreNewIcon(ctx)
//	if err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": "iconError"})
//		return
//	}
//
//	err = TeamsDAO.UpdateTeamIcon(ctx.GetInt("urlId"), smallIcon, largeIcon)
//	if checkErr(ctx, err) {
//		return
//	}
//
//	ctx.Status(http.StatusOK)
//}
//
///**
// * @api{DELETE} /api/teams/removeTeam/:id Delete Team
// * @apiName deleteTeam
// * @apiGroup Teams
// * @apiDescription Delete a team from the current league and its players
// *
// * @apiParam {int} id The unique numerical identifier of the team
// *
// * @apiError IdMustBeInteger The id in the url must be an integer value
// * @apiError notLoggedIn No user is logged in
// * @apiError noActiveLeague There is no active league selected
// * @apiError teamDoesNotExist The specified team does not exist
// * @apiError teamIsActive This team cannot be deleted because it has played games in this league
// * @apiError notTeamAdmin The currently logged in user does not have permissions to edit teams in this league
// */
//func deleteTeam(ctx *gin.Context) {
//	if failIfTeamDoesNotExist(ctx, ctx.GetInt("getLeagueId"), ctx.GetInt("urlId")) {
//		return
//	}
//
//	err := TeamsDAO.DeleteTeam(ctx.GetInt("urlId"))
//	if checkErr(ctx, err) {
//		return
//	}
//
//	ctx.Status(http.StatusOK)
//}
//
///**
// * @api{GET} /api/teams/:id Get Team Information
// * @apiGroup Teams
// * @apiDescription Get information about the team with specified id
// *
// * @apiParam {int} id The unique numerical identifier of the team
// *
// * @apiSuccess {string} name The name of the team
// * @apiSuccess {string} tag The tag of the team
// * @apiSuccess {string} description The team description
// * @apiSuccess {int} wins The number of wins this team has
// * @apiSuccess {int} losses The number of losses this team has
// * @apiSuccess {string} iconSmall The small icon filename
// * @apiSuccess {string} iconLarge The large icon filename
// * @apiSuccess {[]Object} players An array of json objects representing the players on the team
// * @apiSuccess {int} players.id The unique numerical identifier of the player
// * @apiSuccess {string} players.name The name of the player
// * @apiSuccess {string} players.gameIdentifier The in-game name identifier of the player (jersey number, ign, etc.)
// * @apiSuccess {bool} players.mainRoster If true, this player is on the main roster, otherwise is a substitute
// *
// * @apiError IdMustBeInteger The id in the url must be an integer value
// * @apiError noActiveLeague There is no active league selected
// * @apiError teamDoesNotExist The specified team does not exist
// */
//func getTeamInformation(ctx *gin.Context) {
//	if failIfTeamDoesNotExist(ctx, ctx.GetInt("getLeagueId"), ctx.GetInt("urlId")) {
//		return
//	}
//
//	teamInfo, err := TeamsDAO.GetTeamInformation(ctx.GetInt("urlId"))
//	if checkErr(ctx, err) {
//		return
//	}
//
//	ctx.JSON(http.StatusOK, teamInfo)
//}
//
///**
// * @api{POST} /api/teams/addPlayer Add Player To Team
// * @apiGroup Teams
// * @apiDescription Create a new player and add him to the teams roster
// *
// * @apiParam {int} getTeamId The unique numerical identifier of the team the player is to be added to
// * @apiParam {string} name The name of the player (can be left blank)
// * @apiParam {string} gameIdentifier The in-game name identifier of the player (jersey number, ign, etc.)
// * @apiParam {bool} mainRoster If true, this player is on the main roster, otherwise is a substitute
// *
// * @apiSuccess {int} id the unique numerical identifier of the created player
// *
// * @apiError notLoggedIn No user is logged in
// * @apiError noActiveLeague There is no active league selected
// * @apiError teamDoesNotExist The specified team does not exist
// * @apiError canNotEditPlayers The currently logged in player does not have permission to edit the players on this team
// * @apiError gameIdentifierTooLong The game identifier exceeds 50 characters
// * @apiError nameTooLong The name exceeds 50 characters
// * @apiError gameIdentifierTooShort The game identifier is smaller than 2 characters
// * @apiError gameIdentifierInUse This game identifier is already in use in this league
// */
////TODO: add length check for position
//func addPlayerToTeam(ctx *gin.Context) {
//	//get parameters
//	var playerInfo databaseAccess.PlayerDTO
//	err := ctx.ShouldBindBodyWith(&playerInfo, binding.JSON)
//	if checkJsonErr(ctx, err) {
//		return
//	}
//
//	if failIfTeamDoesNotExist(ctx, ctx.GetInt("getLeagueId"), playerInfo.TeamId) {
//		return
//	}
//	if failIfCannotEditPlayersOnTeam(ctx, ctx.GetInt("getLeagueId"), playerInfo.TeamId, ctx.GetInt("getUserId")) {
//		return
//	}
//	if failIfGameIdentifierTooLong(ctx, playerInfo.GameIdentifier) {
//		return
//	}
//	if failIfNameTooLong(ctx, playerInfo.Name) {
//		return
//	}
//	if failIfGameIdentifierTooShort(ctx, playerInfo.GameIdentifier) {
//		return
//	}
//	//TODO: check about game identifier league wide
//	if failIfGameIdentifierInUse(ctx, playerInfo.TeamId, -1, playerInfo.GameIdentifier) {
//		return
//	}
//
//	getPlayerId, err := TeamsDAO.CreatePlayer(playerInfo)
//	if checkErr(ctx, err) {
//		return
//	}
//
//	ctx.JSON(http.StatusOK, gin.H{"id": getPlayerId})
//}
//
///**
// * @api{DELETE} /api/teams/removePlayer Remove Player From Team
// * @apiGroup Teams
// * @apiDescription Remove a player from a teams roster
// *
// * @apiParam {int} getTeamId The unique numerical identifier of the team the player is to be added to
// * @apiParam {int} getPlayerId The unique numerical identifier of the player to be removed
// *
// * @apiError notLoggedIn No user is logged in
// * @apiError noActiveLeague There is no active league selected
// * @apiError teamDoesNotExist The specified team does not exist
// * @apiError canNotEditPlayers The currently logged in player does not have permission to edit the players on this team
// * @apiError playerDoesNotExist The specified player does not exist on this team
// */
//func removePlayerFromTeam(ctx *gin.Context) {
//	//get parameters
//	var playerRemoveInfo PlayerRemoveInformation
//	err := ctx.ShouldBindJSON(&playerRemoveInfo)
//	if checkJsonErr(ctx, err) {
//		return
//	}
//	if failIfTeamDoesNotExist(ctx, ctx.GetInt("getLeagueId"), playerRemoveInfo.TeamId) {
//		return
//	}
//	if failIfCannotEditPlayersOnTeam(ctx, ctx.GetInt("getLeagueId"), playerRemoveInfo.TeamId, ctx.GetInt("getUserId")) {
//		return
//	}
//	if failIfPlayerDoesNotExist(ctx, playerRemoveInfo.TeamId, playerRemoveInfo.PlayerId) {
//		return
//	}
//	//TODO: check if player on team
//	err = TeamsDAO.DeletePlayer(playerRemoveInfo.PlayerId)
//	if checkErr(ctx, err) {
//		return
//	}
//
//	ctx.Status(http.StatusOK)
//}
//
///**
// * @api{put} /api/teams/updatePlayer Update Player Information
// * @apiGroup Teams
// * @apiDescription Change a players information
// *
// * @apiParam {int} getTeamId The unique numerical identifier of the team the player is to be added to
// * @apiParam {int} getPlayerId The unique numerical identifier of the player
// * @apiParam {string} name The updated name of the player (can be left blank)
// * @apiParam {string} gameIdentifier The updated in-game name identifier of the player (jersey number, ign, etc.)
// * @apiParam {bool} mainRoster If true, this player is on the main roster, otherwise is a substitute
// *
// * @apiError notLoggedIn No user is logged in
// * @apiError noActiveLeague There is no active league selected
// * @apiError teamDoesNotExist The specified team does not exist
// * @apiError canNotEditPlayers The currently logged in player does not have permission to edit the players on this team
// * @apiError gameIdentifierTooLong The game identifier exceeds 50 characters
// * @apiError nameTooLong The name exceeds 50 characters
// * @apiError gameIdentifierTooShort The game identifier is smaller than 2 characters
// * @apiError gameIdentifierInUse This game identifier is already in use in this league
// */
////TODO: add length check for position
//func updatePlayer(ctx *gin.Context) {
//	//get parameters
//	var player databaseAccess.PlayerDTO
//	err := ctx.ShouldBindBodyWith(&player, binding.JSON)
//	if checkJsonErr(ctx, err) {
//		return
//	}
//
//	allowedAccess, err := Access.Team(
//		databaseAccess.Edit,
//		ctx.GetInt("getLeagueId"),
//		player.TeamId,
//		ctx.GetInt("getUserId"))
//	if checkErr(ctx, err) {
//		return
//	}
//	if !allowedAccess {
//		ctx.Status(http.StatusForbidden)
//		return
//	}
//
//	valid, problem, err := DataValidator.ValidatePlayerDTO(ctx.GetInt("getLeagueId"), player)
//	if checkErr(ctx, err) {
//		return
//	}
//	if !valid {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": problem})
//		return
//	}
//
//	err = TeamsDAO.UpdatePlayer(player)
//	if checkErr(ctx, err) {
//		return
//	}
//
//	ctx.Status(http.StatusOK)
//}
//
///**
// * @api{PUT} /api/teams/updatePermissions Update ManagerPermissions
// * @apiGroup Teams
// * @apiDescription Change a managers permissions. Must be either team or league admin.
// *
// * @apiParam {int} getTeamId The unique numerical identifier of the team
// * @apiParam {int} getUserId The unique numerical identifier of the manager
// * @apiParam {boolean} administrator Whether the manager can edit permissions of other managers on the team (admin)
// * @apiParam {boolean} information Whether the manager can edit team information
// * @apiParam {boolean} players Whether the manager can edit players on the team
// * @apiParam {boolean} reportResults Whether the manager can report game results of this team
// *
// * @apiError notLoggedIn No user is logged in
// * @apiError noActiveLeague There is no active league selected
// * @apiError teamDoesNotExist The specified team does not exist
// * @apiError managerDoesNotExist The specified manager does not exist on this team
// * @apiError notAdmin The currently logged in player does not have permission to edit the permissions on this team
// */
//func updateManagerPermissions(ctx *gin.Context) {
//	//get parameters
//	var permissionChange TeamManagerPermissionChange
//	err := ctx.ShouldBindJSON(&permissionChange)
//	if checkJsonErr(ctx, err) {
//		return
//	}
//
//	if failIfTeamDoesNotExist(ctx, ctx.GetInt("getLeagueId"), permissionChange.TeamId) {
//		return
//	}
//	if failIfManagerDoesNotExist(ctx, permissionChange.TeamId, permissionChange.UserId) {
//		return
//	}
//	if failIfNotTeamAdmin(ctx, ctx.GetInt("getLeagueId"), permissionChange.TeamId, ctx.GetInt("getUserId")) {
//		return
//	}
//
//	err = TeamsDAO.ChangeManagerPermissions(permissionChange.TeamId, permissionChange.UserId,
//		databaseAccess.TeamPermissionsDTO{
//			Administrator: permissionChange.Administrator,
//			Information:   permissionChange.Information,
//			Players:       permissionChange.Players,
//			ReportResults: permissionChange.ReportResults,
//		})
//	if checkErr(ctx, err) {
//		return
//	}
//
//	ctx.Status(http.StatusOK)
//}
//
func RegisterTeamHandlers(g *gin.RouterGroup) {

	g.POST("", createNewTeam())
	g.GET("", getAllTeams())

	withTeamId := g.Group("/:teamId", storeTeamId())
	withTeamId.GET("", getTeamInfo())
	//withTeamId.PUT("", editTeam())
	//withTeamId.DELETE("", deleteTeam())
	//
	//withTeamId.PUT("/permissions/:userId", storeTargetUserId(), editTeamManagerPermissions())
	//
	withTeamId.POST("/players", createNewPlayer())
	//withPlayerId := withTeamId.Group("/players/:playerId", storePlayerId())
	//withPlayerId.PUT("", editPlayer())
	//withPlayerId.DELETE("", deletePlayer())
}
