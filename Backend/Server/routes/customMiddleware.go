package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, err := ElmSessions.AuthenticateAndGetUserID(ctx)
		if checkErr(ctx, err) {
			ctx.Abort()
			return
		}

		if userID == -1 {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "notLoggedIn"})
			ctx.Abort()
			return
		}

		ctx.Set("userID", userID)
		ctx.Next()
	}
}

func getActiveLeague() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		leagueID, err := ElmSessions.GetActiveLeague(ctx)
		if checkErr(ctx, err) {
			ctx.Abort()
			return
		}

		if leagueID == -1 {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "noActiveLeague"})
			ctx.Abort()
			return
		}

		ctx.Set("leagueID", leagueID)
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
		canEditTeams, err := LeaguesDAO.HasEditTeamsPermission(ctx.GetInt("leagueID"), ctx.GetInt("userID"))
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