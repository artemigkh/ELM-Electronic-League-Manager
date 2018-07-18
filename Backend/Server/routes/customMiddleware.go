package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
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