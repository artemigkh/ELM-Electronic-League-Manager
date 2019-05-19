package routes

import (
	"Server/databaseAccess"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LeaguePermissionChange struct {
	Id            int  `json:"id"`
	Administrator bool `json:"administrator"`
	CreateTeams   bool `json:"createTeams"`
	EditTeams     bool `json:"editTeams"`
	EditGames     bool `json:"editGames"`
}

type LeagueMarkdown struct {
	Markdown string `json:"markdown"`
}

/**
* @api{POST} /api/leagues/ Create New League
* @apiName createNewLeague
* @apiGroup Leagues
* @apiDescription Register a new league
*
* @apiParam {string} name the name of the league
* @apiParam {string} description A brief (<500) char description of the league
* @apiParam {string} game The type of game. Acceptable values:
  "genericsport", "basketball", "curling", "football", "hockey", "rugby", "soccer", "volleyball", "waterpolo",
  "genericesport", "csgo", "leagueoflegends", "overwatch"
* @apiParam {boolean} publicView should the league be viewable by people not playing in the league?
* @apiParam {boolean} publicJoin should the league be joinable by any team that has viewing rights?
* @apiParam {number} signupStart The unix timestamp of the start of the signup period
* @apiParam {number} signupEnd The unix timestamp of the end of the signup period
* @apiParam {number} leagueStart The unix timestamp of the start of the competition period
* @apiParam {number} leagueEnd The unix timestamp of the end of the competition period
*
* @apiSuccess {int} id the primary id of the created league
*
* @apiError notLoggedIn No user is logged in
* @apiError nameTooLong The league name has exceeded 50 characters
* @apiError nameTooShort The league name is shorter than 2 characters
* @apiError descriptionTooLong The description has exceeded 500 characters
* @apiError gameStringNotValid The game string is not one of the allowed values
* @apiError nameInUse The league name is currently in use
*/
func createNewLeague(ctx *gin.Context) {
	//TODO: here and in update, check that competition period after signup period
	var lgRequest databaseAccess.LeagueDTO
	err := ctx.ShouldBindJSON(&lgRequest)
	if checkJsonErr(ctx, err) {
		return
	}
	if failIfGameStringtNotValid(ctx, lgRequest.Game) {
		return
	}
	if failIfDescriptionTooLong(ctx, lgRequest.Description) {
		return
	}
	if failIfNameTooLong(ctx, lgRequest.Name) {
		return
	}
	if failIfNameTooShort(ctx, lgRequest.Name) {
		return
	}
	if failIfLeagueNameInUse(ctx, -1, lgRequest.Name) {
		return
	}

	leagueId, err := LeaguesDAO.CreateLeague(
		ctx.GetInt("userId"), lgRequest)
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": leagueId})
}

/**
* @api{PUT} /api/leagues/ Update League Information
* @apiName updateLeagueInformation
* @apiGroup Leagues
* @apiDescription Update currently active league information
*
* @apiParam {string} name the name of the league
* @apiParam {string} description A brief (<500) char description of the league
* @apiParam {string} game The type of game. Acceptable values:
  "genericsport", "basketball", "curling", "football", "hockey", "rugby", "soccer", "volleyball", "waterpolo",
  "genericesport", "csgo", "leagueoflegends", "overwatch"
* @apiParam {boolean} publicView should the league be viewable by people not playing in the league?
* @apiParam {boolean} publicJoin should the league be joinable by any team that has viewing rights?
* @apiParam {number} signupStart The unix timestamp of the start of the signup period
* @apiParam {number} signupEnd The unix timestamp of the end of the signup period
* @apiParam {number} leagueStart The unix timestamp of the start of the competition period
* @apiParam {number} leagueEnd The unix timestamp of the end of the competition period
*
* @apiSuccess {int} id the primary id of the created league
*
* @apiError notLoggedIn No user is logged in
* @apiError notAdmin Currently logged in user is not a league administrator
* @apiError nameTooLong The league name has exceeded 50 characters
* @apiError nameTooShort The league name is shorter than 2 characters
* @apiError descriptionTooLong The description has exceeded 500 characters
* @apiError gameStringNotValid The game string is not one of the allowed values
* @apiError nameInUse The league name is currently in use
*/
func updateLeagueInfo(ctx *gin.Context) {
	var lgRequest databaseAccess.LeagueDTO
	err := ctx.ShouldBindJSON(&lgRequest)
	if checkJsonErr(ctx, err) {
		return
	}
	lgRequest.Id = ctx.GetInt("leagueId")

	if failIfGameStringtNotValid(ctx, lgRequest.Game) {
		return
	}
	if failIfDescriptionTooLong(ctx, lgRequest.Description) {
		return
	}
	if failIfNameTooLong(ctx, lgRequest.Name) {
		return
	}
	if failIfNameTooShort(ctx, lgRequest.Name) {
		return
	}
	if failIfLeagueNameInUse(ctx, lgRequest.Id, lgRequest.Name) {
		return
	}

	err = LeaguesDAO.UpdateLeague(lgRequest)

	if checkErr(ctx, err) {
		return
	}

	ctx.Status(http.StatusOK)
}

/**
 * @api{POST} /api/leagues/setActiveLeague/:id Set Active League
 * @apiName setActiveLeague
 * @apiGroup Leagues
 * @apiDescription Attempt to set the active league to :id
 * @apiParam {int} id the primary id of the league
 *
 * @apiError leagueDoesNotExist The league with specified id does not exist
 * @apiError 403 Forbidden
 */
func setActiveLeague(ctx *gin.Context) {
	//TODO: check if league exists
	//get user Id (or -1 if not logged in)
	userId, err := ElmSessions.AuthenticateAndGetUserId(ctx)
	if checkErr(ctx, err) {
		return
	}

	viewable, err := LeaguesDAO.IsLeagueViewable(ctx.GetInt("urlId"), userId)
	if checkErr(ctx, err) {
		return
	}

	if !viewable {
		ctx.JSON(http.StatusForbidden, nil)
	} else {
		err := ElmSessions.SetActiveLeague(ctx, ctx.GetInt("urlId"))
		if checkErr(ctx, err) {
			return
		}
	}
}

/**
* @api{GET} /api/leagues Get Active League Information
* @apiGroup Leagues
* @apiDescription Get information about the currently selected league
*
* @apiSuccess {int} id The unique numerical identifier of the league
* @apiSuccess {string} name The name of the currently selected league
* @apiSuccess {string} description The description of the currently selected league
* @apiSuccess {bool} publicView True if league is publicly viewable
* @apiSuccess {bool} publicJoin True if league is publicly joinable
* @apiSuccess {number} signupStart The unix timestamp of the start of the signup period
* @apiSuccess {number} signupEnd The unix timestamp of the end of the signup period
* @apiSuccess {number} leagueStart The unix timestamp of the start of the competition period
* @apiSuccess {number} leagueStart The unix timestamp of the start of the competition period
* @apiSuccess {string} game The type of game. will be one of:
 "genericsport", "basketball", "curling", "football", "hockey", "rugby", "soccer", "volleyball", "waterpolo",
 "genericesport", "csgo", "leagueoflegends", "overwatch"
*
* @apiError noActiveLeague There is no active league selected
*/
func getActiveLeagueInformation(ctx *gin.Context) {
	leagueInfo, err := LeaguesDAO.GetLeagueInformation(ctx.GetInt("leagueId"))
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, leagueInfo)
}

/**
 * @api{GET} /api/leagues/teamSummary Get Team Summary
 * @apiGroup Leagues
 * @apiDescription Get the team summary of the current league, sorted by standings
 *
 * @apiSuccess {jsonArray} _ An array of JSON objects, each representing a team
 * @apiSuccess {int} _.id The unique numerical identifier of the team
 * @apiSuccess {int} _.name The name of the team
 * @apiSuccess {int} _.tag The tag of the team
 * @apiSuccess {int} _.wins The number of wins of the team
 * @apiSuccess {int} _.losses The number of losses of the team
 * @apiSuccess {int} _.iconSmall The small icon filename
 * @apiSuccess {int} _.iconLarge The large icon filename
 *
 * @apiError noActiveLeague There is no active league selected
 */
func getTeamSummary(ctx *gin.Context) {
	teamSummary, err := LeaguesDAO.GetTeamSummary(ctx.GetInt("leagueId"))
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, teamSummary)
}

/**
 * @api{POST} /api/leagues/join Join Active League
 * @apiGroup Leagues
 * @apiDescription Join the currently selected league as a manager
 *
 * @apiError notLoggedIn No user is logged in
 * @apiError noActiveLeague There is no active league selected
 * @apiError canNotJoin The active league is not accepting new members
 */
func joinActiveLeague(ctx *gin.Context) {
	err := LeaguesDAO.JoinLeague(ctx.GetInt("leagueId"), ctx.GetInt("userId"))
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

/**
 * @api{GET} /api/leagues/teamManagers Get Team Managers
 * @apiGroup Leagues
 * @apiDescription If logged in as a league administrator, see all users that have permissions to manage teams in this league
 *
 * @apiSuccess {jsonArray} _ An array of JSON objects, each representing a team
 * @apiSuccess {int} _.teamId The unique numerical identifier of the team
 * @apiSuccess {string} _.teamName The name of the team
 * @apiSuccess {string} _.teamTag The tag of the team
 * @apiSuccess {[]Object} _.managers The users on this team that have management permissions
 * @apiSuccess {int} _.managers.userId The unique numerical identifier of the user/manager
 * @apiSuccess {string} _.managers.userEmail The email of the user/manager
 * @apiSuccess {bool} _.managers.administrator True if this user can manage permissions of other users on the team
 * @apiSuccess {bool} _.managers.information True if this user can edit information about the team
 * @apiSuccess {bool} _.managers.players True if this user can edit players on this team
 * @apiSuccess {bool} _.managers.reportResults True if this user can report results for this team
 *
 * @apiError notLoggedIn No user is logged in
 * @apiError noActiveLeague There is no active league selected
 * @apiError notAdmin The currently logged in user is not a league administrator
 */
func getTeamManagers(ctx *gin.Context) {
	teamManagerInfo, err := LeaguesDAO.GetTeamManagerInformation(ctx.GetInt("leagueId"))
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, teamManagerInfo)
}

/**
 * @api{GET} /api/leagues/gameSummary Get Game Summary
 * @apiGroup Leagues
 * @apiDescription Get the game summary of the current league, in chronological order
 *
 * @apiSuccess {jsonArray} _ An array of JSON objects, each representing a game
 * @apiSuccess {int} _.id The unique numerical identifier of the game
 * @apiSuccess {int} _.team1Id The unique numerical identifier of the team in position 1
 * @apiSuccess {int} _.team2Id The unique numerical identifier of the team in position 2
 * @apiSuccess {int} _.gameTime The unix epoch time in seconds when the game is played
 * @apiSuccess {bool} _.complete A boolean indicating if the game is finished or not
 * @apiSuccess {int} _.winnerId The Id of the winning team, or -1 if the game is not complete
 * @apiSuccess {int} _.scoreTeam1 The score of the team in position 1
 * @apiSuccess {int} _.scoreTeam2 The score of the team in position 2
 *
 * @apiError noActiveLeague There is no active league selected
 */
func getGameSummary(ctx *gin.Context) {
	gameSummary, err := LeaguesDAO.GetGameSummary(ctx.GetInt("leagueId"))
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, gameSummary)
}

/**
 * @api{GET} /api/leagues/publicLeagues Get List of Publicly Viewable Leagues
 * @apiGroup Leagues
 * @apiDescription Get a list of all publicly viewable leagues
 *
 * @apiSuccess {jsonArray} _ An array of JSON objects, each representing a league
 * @apiSuccess {umber} _.id The unique numerical identifier of the league
 * @apiSuccess {string} _.name The name of the league
 * @apiSuccess {string} _.description The description of the league
 * @apiSuccess {bool} _.publicJoin A boolean that signifies if the league can be joined by the general public
 * @apiSuccess {number} _.signupStart The unix timestamp of the start of the signup period
 * @apiSuccess {number} _.signupEnd The unix timestamp of the end of the signup period
 * @apiSuccess {number} _.leagueStart The unix timestamp of the start of the competition period
 * @apiSuccess {number} _.leagueStart The unix timestamp of the start of the competition period
 * @apiSuccess {string} _.game The type of game. will be one of:
 *
 */
func getPublicLeagues(ctx *gin.Context) {
	leagueList, err := LeaguesDAO.GetPublicLeagueList()
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, leagueList)
}

/**
 * @api{POST} /api/leagues/setLeaguePermissions Set League Permissions
 * @apiGroup Leagues
 * @apiDescription Set the specified users league permissions in the currently active league
 *
 * @apiParam {number} id the unique numerical identifier of the user
 * @apiSuccess {bool} administrator if user is a league administrator
 * @apiSuccess {bool} createTeams if the user can create teams
 * @apiSuccess {bool} editTeams if the user can edit existing teams
 * @apiSuccess {bool} editGames if the user can edit games in this league
 *
 * @apiError notLoggedIn No user is logged in
 * @apiError noActiveLeague There is no active league selected
 * @apiError notAdmin The currently logged in user is not a league administrator
 */
func setLeaguePermissions(ctx *gin.Context) {
	//get parameters
	var permissionChange databaseAccess.UserPermissionsDTO
	err := ctx.ShouldBindJSON(&permissionChange)
	if checkJsonErr(ctx, err) {
		return
	}

	err = LeaguesDAO.SetLeaguePermissions(
		ctx.GetInt("leagueId"), permissionChange)

	if checkErr(ctx, err) {
		return
	}
	ctx.Status(http.StatusOK)
}

/**
* @api{POST} /api/leagues/markdown Provide leagues rules and information markdown
* @apiName createMarkdown
* @apiGroup Leagues
* @apiDescription Provide rules and information markdown for active league
*
* @apiParam {string} markdown The markdown to be stored
*
* @apiError notLoggedIn No user is logged in
* @apiError noActiveLeague There is no active league selected
* @apiError notAdmin Currently logged in user is not a league administrator
* @apiError markdownTooLong Markdown is larger than 50k characters
 */
func setLeagueMarkdown(ctx *gin.Context) {
	var md LeagueMarkdown
	err := ctx.ShouldBindJSON(&md)
	if checkJsonErr(ctx, err) {
		return
	}

	if failIfMdTooLong(ctx, md.Markdown) {
		return
	}

	oldFile, err := LeaguesDAO.GetMarkdownFile(ctx.GetInt("leagueId"))
	if checkErr(ctx, err) {
		return
	}

	fileName, err := MarkdownManager.StoreMarkdown(ctx.GetInt("leagueId"), md.Markdown, oldFile)
	if checkErr(ctx, err) {
		return
	}

	err = LeaguesDAO.SetMarkdownFile(ctx.GetInt("leagueId"), fileName)
	if checkErr(ctx, err) {
		return
	}

	ctx.Status(http.StatusOK)
}

/**
* @api{GET} /api/leagues/markdown Get leagues rules and information markdown
* @apiGroup Leagues
* @apiDescription Get rules and information markdown for active league
*
* @apiSuccess {string} markdown The unique numerical identifier of the league
*
* @apiError noActiveLeague There is no active league selected
 */
func getLeagueMarkdown(ctx *gin.Context) {
	fileName, err := LeaguesDAO.GetMarkdownFile(ctx.GetInt("leagueId"))
	if checkErr(ctx, err) {
		return
	}

	if fileName == "" {
		ctx.JSON(http.StatusOK, gin.H{"markdown": ""})
	} else {
		markdown, err := MarkdownManager.GetMarkdown(fileName)
		if checkErr(ctx, err) {
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"markdown": markdown})
	}
}

func RegisterLeagueHandlers(g *gin.RouterGroup) {
	g.POST("/", authenticate(), createNewLeague)
	g.PUT("/", authenticate(), getActiveLeague(), failIfNotLeagueAdmin(), updateLeagueInfo)
	g.POST("/setActiveLeague/:id", getUrlId(), failIfLeagueDoesNotExist(), setActiveLeague)
	g.POST("/join", authenticate(), getActiveLeague(), failIfCannotJoinLeague(), joinActiveLeague)
	g.GET("/", getActiveLeague(), getActiveLeagueInformation)
	g.GET("/publicLeagues", getPublicLeagues)
	g.GET("/teamSummary", getActiveLeague(), getTeamSummary)
	g.GET("/gameSummary", getActiveLeague(), getGameSummary)
	g.GET("/teamManagers", authenticate(), getActiveLeague(), failIfNotLeagueAdmin(), getTeamManagers)
	g.POST("/setLeaguePermissions",
		authenticate(), getActiveLeague(), failIfNotLeagueAdmin(), setLeaguePermissions)
	g.POST("/markdown", authenticate(), getActiveLeague(), failIfNotLeagueAdmin(), setLeagueMarkdown)
	g.GET("/markdown", getActiveLeague(), getLeagueMarkdown)
}
