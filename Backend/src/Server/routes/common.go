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
var DataValidator databaseAccess.DataValidator
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

// Functions
func getLeagueAndTeamPermissions(leagueId, teamId, userId int) (
	*databaseAccess.LeaguePermissionsDTO, *databaseAccess.TeamPermissionsDTO, error) {
	leaguePermissions, err := LeaguesDAO.GetLeaguePermissions(leagueId, userId)
	if err != nil {
		return nil, nil, err
	}

	teamPermissions, err := TeamsDAO.GetTeamPermissions(teamId, userId)
	if err != nil {
		return nil, nil, err
	}

	return leaguePermissions, teamPermissions, nil
}

// context helpers

func leagueId(ctx *gin.Context) int {
	return ctx.GetInt("leagueId")
}

func userId(ctx *gin.Context) int {
	return ctx.GetInt("userId")
}

//type Context gin.Context
//
//func (ctx *Context) LeagueId() int {
//	return 0
//}
