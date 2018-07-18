package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TeamInformation struct {
	Name string `json:"name"`
	Tag string `json:"tag"`
}

func createNewTeam(ctx *gin.Context) {
	//get parameters
	var teamInfo TeamInformation
	err := ctx.ShouldBindJSON(&teamInfo)
	if checkJsonErr(ctx, err) {
		return
	}

	//must be logged in to create a team
	userID, err := ElmSessions.AuthenticateAndGetUserID(ctx)
	if checkErr(ctx, err) {
		return
	}
	if userID == -1 {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "notLoggedIn"})
		return
	}

	//must have an active league to create a team in it
	leagueId, err := ElmSessions.GetActiveLeague(ctx)
	if checkErr(ctx, err) {
		return
	}
	if leagueId == -1 {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "noActiveLeague"})
		return
	}

	//must have permissions to edit teams in the league to create one
	canEditTeams, err := LeaguesDAO.HasEditTeamsPermission(leagueId, userID)
	if checkErr(ctx, err) {
		return
	}
	if !canEditTeams {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "noEditLeaguePermissions"})
		return
	}

	if failIfTeamNameTooLong(ctx, teamInfo.Name) {
		return
	}
	if failIfTeamTagTooLong(ctx, teamInfo.Tag) {
		return
	}
	if failIfTeamInfoInUse(ctx, teamInfo.Name, teamInfo.Tag, leagueId) {
		return
	}

	teamID, err := TeamsDAO.CreateTeam(leagueId, userID, teamInfo.Name, teamInfo.Tag)
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": teamID})
}

func getTeamInformation(ctx *gin.Context) {
	teamID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "IdMustBeInteger"})
		return
	}

	if failIfTeamDoesNotExist(ctx, teamID, ctx.GetInt("")) {
		return
	}

	teamInfo, err := TeamsDAO.GetTeamInformation(teamID, ctx.GetInt("leagueID"))
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, teamInfo)
}

func RegisterTeamHandlers(g *gin.RouterGroup) {
	g.Use(getActiveLeague())

	g.POST("/", createNewTeam)
	g.GET("/:id", getTeamInformation)
}
