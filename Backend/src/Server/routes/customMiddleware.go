package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, err := ElmSessions.AuthenticateAndGetUserId(ctx)
		if checkErr(ctx, err) {
			ctx.Abort()
			return
		} else {
			ctx.Set("userId", userId)
		}

		leagueId, err := ElmSessions.GetActiveLeague(ctx)
		if checkErr(ctx, err) {
			ctx.Abort()
		} else {
			ctx.Set("leagueId", leagueId)
			ctx.Next()
		}
	}
}

func storeUrlId(urlParam, storedName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		urlId, err := strconv.Atoi(ctx.Param(urlParam))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "IdMustBeInteger"})
		} else {
			ctx.Set(storedName, urlId)
			ctx.Next()
		}
	}
}

func storeTargetLeagueId() gin.HandlerFunc {
	return storeUrlId("leagueId", "targetLeagueId")
}

func storeTargetUserId() gin.HandlerFunc {
	return storeUrlId("userId", "targetUserId")
}

func storeGameId() gin.HandlerFunc {
	return storeUrlId("gameId", "gameId")
}

func storeTeamId() gin.HandlerFunc {
	return storeUrlId("teamId", "teamId")
}

func storePlayerId() gin.HandlerFunc {
	return storeUrlId("playerId", "playerId")
}
