package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, err := ElmSessions.AuthenticateAndGetUserId(ctx)
		if checkErr(ctx, err) {
			ctx.Abort()
			return
		}

		if userId == -1 {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "notLoggedIn"})
			ctx.Abort()
			return
		}

		ctx.Set("userId", userId)
		ctx.Next()
	}
}

func getActiveLeague() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		leagueId, err := ElmSessions.GetActiveLeague(ctx)
		if checkErr(ctx, err) {
			ctx.Abort()
			return
		}

		if leagueId == -1 {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "noActiveLeague"})
			ctx.Abort()
			return
		}

		ctx.Set("leagueId", leagueId)
		ctx.Next()
	}
}

func getUrlId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		println("id is: ", ctx.Param("id"))
		urlId, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "IdMustBeInteger"})
			ctx.Abort()
			return
		}

		ctx.Set("urlId", urlId)
		ctx.Next()
	}
}

func getTeamEditPermissions() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		canEditTeams, err := LeaguesDAO.HasEditTeamsPermission(ctx.GetInt("leagueId"), ctx.GetInt("userId"))
		if checkErr(ctx, err) {
			ctx.Abort()
			return
		}
		if !canEditTeams {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "noEditLeaguePermissions"})
			ctx.Abort()
			return
		}
	}
}

func getReportResultPermissions() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		canReportResult, err := GamesDAO.HasReportResultPermissions(
			ctx.GetInt("leagueId"),
			ctx.GetInt("urlId"),
			ctx.GetInt("userId"),
		)
		if checkErr(ctx, err) {
			ctx.Abort()
			return
		}
		if !canReportResult {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "noReportResultPermissions"})
			ctx.Abort()
			return
		}
	}
}
