package sessionManager

import "github.com/gin-gonic/gin"

type SessionManager interface {
	LogIn(ctx *gin.Context, userId int) error
	LogOut(ctx *gin.Context) error
	AuthenticateAndGetUserId(ctx *gin.Context) (int, error)
	SetActiveLeague(ctx *gin.Context, leagueId int) error
	GetActiveLeague(ctx *gin.Context) (int, error)
}
