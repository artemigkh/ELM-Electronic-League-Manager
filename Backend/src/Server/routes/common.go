package routes

import (
	"Server/databaseAccess"
	"Server/icons"
	"Server/lolApi"
	"Server/markdown"
	"Server/sessionManager"
	"github.com/gin-gonic/gin"
)

// Objects
var UsersDAO databaseAccess.UsersDAO
var LeaguesDAO databaseAccess.LeaguesDAO
var TeamsDAO databaseAccess.TeamsDAO
var GamesDAO databaseAccess.GamesDAO
var LeagueOfLegendsDAO databaseAccess.LeagueOfLegendsDAO
var Access databaseAccess.Access

var ElmSessions sessionManager.SessionManager

var IconManager icons.IconManager
var MarkdownManager markdown.MdManager

var LoLApi lolApi.LoLApi

// Structs
type userInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// context helpers
func getLeagueId(ctx *gin.Context) int {
	return ctx.GetInt("leagueId")
}

func getTargetLeagueId(ctx *gin.Context) int {
	return ctx.GetInt("targetLeagueId")
}

func getUserId(ctx *gin.Context) int {
	return ctx.GetInt("userId")
}

func getGameId(ctx *gin.Context) int {
	return ctx.GetInt("gameId")
}
func getTargetUserId(ctx *gin.Context) int {
	return ctx.GetInt("targetUserId")
}

func getTeamId(ctx *gin.Context) int {
	return ctx.GetInt("teamId")
}

func getPlayerId(ctx *gin.Context) int {
	return ctx.GetInt("playerId")
}
