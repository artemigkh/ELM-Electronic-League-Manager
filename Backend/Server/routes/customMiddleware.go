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
		} else if userId == -1 {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "notLoggedIn"})
		} else {
			ctx.Set("userId", userId)
			ctx.Next()
		}
	}
}

func getActiveLeague() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		leagueId, err := ElmSessions.GetActiveLeague(ctx)
		if checkErr(ctx, err) {
			ctx.Abort()
		} else if leagueId == -1 {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "noActiveLeague"})
		} else {
			ctx.Set("leagueId", leagueId)
			ctx.Next()
		}
	}
}

func getUrlId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		urlId, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "IdMustBeInteger"})
		} else {
			ctx.Set("urlId", urlId)
			ctx.Next()
		}
	}
}

func getTeamCreatePermissions() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		canEditTeams, err := LeaguesDAO.HasCreateTeamsPermission(ctx.GetInt("leagueId"), ctx.GetInt("userId"))
		if checkErr(ctx, err) {
			ctx.Abort()
		} else if !canEditTeams {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "noEditTeamPermissions"})
		} else {
			ctx.Next()
		}
	}
}

func getTeamEditPermissions() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		canEditTeams, err := LeaguesDAO.HasEditTeamsPermission(ctx.GetInt("leagueId"), ctx.GetInt("userId"))
		if checkErr(ctx, err) {
			ctx.Abort()
		} else if !canEditTeams {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "noEditTeamPermissions"})
		} else {
			ctx.Next()
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
		} else if !canReportResult {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "noReportResultPermissions"})
		} else {
			ctx.Next()
		}
	}
}

func failIfCannotJoinLeague() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		canJoin, err := LeaguesDAO.CanJoinLeague(ctx.GetInt("leagueId"), ctx.GetInt("userId"))
		if checkErr(ctx, err) {
			ctx.Abort()
		} else if !canJoin {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "canNotJoin"})
		} else {
			ctx.Next()
		}
	}
}

func failIfNotLeagueAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		isLeagueAdmin, err := LeaguesDAO.IsLeagueAdmin(ctx.GetInt("leagueId"), ctx.GetInt("userId"))
		if checkErr(ctx, err) {
			ctx.Abort()
		} else if !isLeagueAdmin {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "notAdmin"})
		} else {
			ctx.Next()
		}
	}
}
