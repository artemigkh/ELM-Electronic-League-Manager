package routes

import (
	"Server/databaseAccess"
	"Server/icons"
	"Server/lolApi"
	"Server/markdown"
	"Server/sessionManager"
	"github.com/gin-gonic/gin"
	"net/http"
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

// Each endpoint does a subset of the following:
// 1. Check Permissions of logged in user against action
// 2. Binds and validates input data
// 3. Performs operations with input data
// 4. Returns output data
//type endpoint interface {
//	HasPermissions() (bool, error)
//	BindData() bool
//	IsDataInvalid() (bool, string, error)
//	Core() error
//	Return()
//}

type Entity int

const (
	User   Entity = iota
	League Entity = iota
	Team   Entity = iota
	Player Entity = iota
	Game   Entity = iota
	Report Entity = iota
)

const (
	View   = databaseAccess.View
	Edit   = databaseAccess.Edit
	Create = databaseAccess.Create
	Delete = databaseAccess.Delete
)

func HasPermissions(ctx *gin.Context, entity Entity, accessType databaseAccess.AccessType) (bool, error) {
	var hasPermissions bool
	var err error
	var entityId = 0
	var leagueId = getLeagueId(ctx)
	var userId = getUserId(ctx)

	switch entity {
	case User:
		//TODO: validate user access things
		hasPermissions = true
		err = nil
	case League:
		if accessType != Create {
			entityId = getLeagueId(ctx)
		}
		hasPermissions, err = Access.League(accessType, entityId, userId)
	case Team:
		if accessType != Create {
			entityId = getTeamId(ctx)
		}
		hasPermissions, err = Access.Team(accessType, leagueId, entityId, userId)
	case Player:
		if accessType != Create {
			entityId = getPlayerId(ctx)
		}
		hasPermissions, err = Access.Player(accessType, leagueId, getTeamId(ctx), entityId, userId)
	case Game:
		if accessType != Create {
			entityId = getGameId(ctx)
		}
		hasPermissions, err = Access.Game(accessType, leagueId, entityId, userId)
	case Report:
		hasPermissions, err = Access.Report(leagueId, getGameId(ctx), userId)
	}

	if err != nil {
		return false, err
	} else {
		return hasPermissions, nil
	}
}

type endpoint struct {
	Entity        Entity
	AccessType    databaseAccess.AccessType
	BindData      func(ctx *gin.Context) bool
	IsDataInvalid func(ctx *gin.Context) (bool, string, error)
	Core          func(ctx *gin.Context) (interface{}, error)
}

func (e endpoint) createEndpointHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Check Permissions of logged in user against action
		hasPermissions, err := HasPermissions(ctx, e.Entity, e.AccessType)
		if accessForbidden(ctx, hasPermissions, err) {
			return
		}
		// Binds input data
		if e.BindData != nil && e.BindData(ctx) {
			return
		}

		// Validate input data
		if e.IsDataInvalid != nil {
			valid, problem, err := e.IsDataInvalid(ctx)
			if DataInvalid(ctx, valid, problem, err) {
				return
			}
		}
		// Perform the core action of the endpoint
		returnData, err := e.Core(ctx)
		if checkErr(ctx, err) {
			return
		}

		// Return status and data if exists to router
		if returnData == nil {
			ctx.Status(http.StatusOK)
		} else {
			if e.AccessType == Create {
				ctx.JSON(http.StatusCreated, returnData)
			} else {
				ctx.JSON(http.StatusOK, returnData)
			}
		}
	}
}
