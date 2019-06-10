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

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/createLeague
func createNewLeague(ctx *gin.Context) {
	var league databaseAccess.LeagueCore
	if bindAndCheckErr(ctx, &league) {
		return
	}

	valid, problem, err := league.ValidateNew()
	if dataInvalid(ctx, valid, problem, err) {
		return
	}

	leagueId, err := LeaguesDAO.CreateLeague(getUserId(ctx), league)
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": leagueId})
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/updateLeague
func updateLeagueInfo(ctx *gin.Context) {
	hasPermissions, err := Access.League(databaseAccess.Edit, getLeagueId(ctx), getUserId(ctx))
	if accessForbidden(ctx, hasPermissions, err) {
		return
	}

	var league databaseAccess.LeagueCore
	if bindAndCheckErr(ctx, &league) {
		return
	}

	valid, problem, err := league.ValidateEdit(getLeagueId(ctx))
	if dataInvalid(ctx, valid, problem, err) {
		return
	}

	err = LeaguesDAO.UpdateLeague(getLeagueId(ctx), league)
	if checkErr(ctx, err) {
		return
	}

	ctx.Status(http.StatusOK)
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/setLeagueMd
func setLeagueMarkdown(ctx *gin.Context) {
	hasPermissions, err := Access.League(databaseAccess.Edit, getLeagueId(ctx), getUserId(ctx))
	if accessForbidden(ctx, hasPermissions, err) {
		return
	}

	var md databaseAccess.Markdown
	err = ctx.ShouldBindJSON(&md)
	if checkJsonErr(ctx, err) {
		return
	}

	valid, problem, err := md.Validate()
	if dataInvalid(ctx, valid, problem, err) {
		return
	}

	oldFile, err := LeaguesDAO.GetMarkdownFile(ctx.GetInt("getLeagueId"))
	if checkErr(ctx, err) {
		return
	}

	fileName, err := MarkdownManager.StoreMarkdown(ctx.GetInt("getLeagueId"), md.Markdown, oldFile)
	if checkErr(ctx, err) {
		return
	}

	err = LeaguesDAO.SetMarkdownFile(ctx.GetInt("getLeagueId"), fileName)
	if checkErr(ctx, err) {
		return
	}

	ctx.Status(http.StatusOK)
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/setActiveLeague
func setActiveLeague(ctx *gin.Context) {
	hasPermissions, err := Access.League(
		databaseAccess.View, getTargetLeagueId(ctx), getUserId(ctx))
	if accessForbidden(ctx, hasPermissions, err) {
		return
	}

	err = ElmSessions.SetActiveLeague(ctx, getTargetLeagueId(ctx))
	if checkErr(ctx, err) {
		return
	}

	ctx.Status(http.StatusOK)
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/getLeagueInfo
func getActiveLeagueInformation(ctx *gin.Context) {
	hasPermissions, err := Access.League(databaseAccess.View, getLeagueId(ctx), getUserId(ctx))
	if accessForbidden(ctx, hasPermissions, err) {
		return
	}

	leagueInfo, err := LeaguesDAO.GetLeagueInformation(getLeagueId(ctx))
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, leagueInfo)
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/joinLeague
func joinActiveLeague(ctx *gin.Context) {
	allowedJoin, err := LeaguesDAO.CanJoinLeague(getLeagueId(ctx), getUserId(ctx))
	if accessForbidden(ctx, allowedJoin, err) {
		return
	}

	err = LeaguesDAO.JoinLeague(getLeagueId(ctx), getUserId(ctx))
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/getLeagueTeamManagers
func getTeamManagers(ctx *gin.Context) {
	hasPermissions, err := Access.League(databaseAccess.Edit, getLeagueId(ctx), getUserId(ctx))
	if accessForbidden(ctx, hasPermissions, err) {
		return
	}

	teamManagerInfo, err := LeaguesDAO.GetTeamManagerInformation(ctx.GetInt("getLeagueId"))
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, teamManagerInfo)
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/getPublicLeagues
func getPublicLeagues(ctx *gin.Context) {
	leagueList, err := LeaguesDAO.GetPublicLeagueList()
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, leagueList)
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/setLeaguePermissions
func setLeaguePermissions(ctx *gin.Context) {
	hasPermissions, err := Access.League(databaseAccess.Edit, getLeagueId(ctx), getUserId(ctx))
	if accessForbidden(ctx, hasPermissions, err) {
		return
	}

	var permissions databaseAccess.LeaguePermissionsCore
	if bindAndCheckErr(ctx, &permissions) {
		return
		return
	}

	valid, problem, err := permissions.Validate()
	if dataInvalid(ctx, valid, problem, err) {
		return
	}

	err = LeaguesDAO.SetLeaguePermissions(
		getLeagueId(ctx), getTargetUserId(ctx), permissions)
	if checkErr(ctx, err) {
		return
	}

	ctx.Status(http.StatusOK)
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/getLeagueMd
func getLeagueMarkdown(ctx *gin.Context) {
	fileName, err := LeaguesDAO.GetMarkdownFile(ctx.GetInt("getLeagueId"))
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
	// League Manage
	g.POST("/", createNewLeague)

	g.PUT("/", updateLeagueInfo)
	g.POST("/markdown", setLeagueMarkdown)
	g.GET("/teamManagers", getTeamManagers)
	g.PUT("/setLeaguePermissions/:userId", storeTargetUserId(), setLeaguePermissions)

	// League Interact
	g.POST("/setActiveLeague/:leagueId", storeTargetLeagueId(), setActiveLeague)
	g.POST("/join", joinActiveLeague)

	// League Information
	g.GET("/", getActiveLeagueInformation)
	g.GET("/markdown", getLeagueMarkdown)
	g.GET("/publicLeagues", getPublicLeagues)
}
