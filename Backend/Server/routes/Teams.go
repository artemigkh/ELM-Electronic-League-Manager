package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type TeamInformation struct {
	Name string `json:"name"`
	Tag  string `json:"tag"`
}

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

func addPlayerToTeam(ctx *gin.Context) {

}

func RegisterTeamHandlers(g *gin.RouterGroup) {
	g.Use(getActiveLeague())

	g.POST("/", authenticate(), getTeamEditPermissions(), createNewTeam)
	g.GET("/:id", getUrlId(), getTeamInformation)
}
