package sessionManager

import "github.com/gin-gonic/gin"

type SessionManager interface {
	LogIn(ctx *gin.Context, userID int) error
	LogOut(ctx *gin.Context) error
	AuthenticateAndGetUserID(ctx *gin.Context) (int, error)
	SetActiveLeague(ctx *gin.Context, leagueID int) error
	GetActiveLeague(ctx *gin.Context) (int, error)
}