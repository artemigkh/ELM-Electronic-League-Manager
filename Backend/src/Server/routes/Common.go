package routes

import (
	"Server/dataModel"
	"Server/icons"
	"Server/lolApi"
	"Server/markdown"
	"Server/sessionManager"
	"Server/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Objects
var UserDAO dataModel.UserDAO
var LeagueDAO dataModel.LeagueDAO
var TeamDAO dataModel.TeamDAO
var GameDAO dataModel.GameDAO
var LeagueOfLegendsDAO dataModel.LeagueOfLegendsDAO
var Access validation.Access

var ElmSessions sessionManager.SessionManager

var IconManager icons.IconManager
var MarkdownManager markdown.MdManager

var LoLApi lolApi.LoLApi

// context helpers
func getLeagueId(ctx *gin.Context) int {
	return ctx.GetInt("leagueId")
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

func getAvailabilityId(ctx *gin.Context) int {
	return ctx.GetInt("availabilityId")
}

func getExternalId(ctx *gin.Context) string {
	return ctx.GetString("externalId")
}

func getExternalGameId(ctx *gin.Context) *string {
	if id := ctx.GetString("externalGameId"); id != "" {
		return &id
	} else {
		return nil
	}
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
	User         Entity = iota
	League       Entity = iota
	Team         Entity = iota
	Player       Entity = iota
	Game         Entity = iota
	Report       Entity = iota
	Availability Entity = iota
)

// Re-declare so that don't have to use the package prefix to make it look nicer
const (
	View   = validation.View
	Edit   = validation.Edit
	Create = validation.Create
	Delete = validation.Delete
)

func HasPermissions(ctx *gin.Context, entity Entity, accessType validation.AccessType) (bool, error) {
	var hasPermissions bool
	var err error
	var entityId = 0
	var leagueId = getLeagueId(ctx)

	permissions, err := UserDAO.GetUserWithPermissions(getLeagueId(ctx), getUserId(ctx))
	if err != nil {
		return false, err
	}

	switch entity {
	case User:
		hasPermissions = true
		err = nil
		if accessType == View {
			entityId, err = ElmSessions.AuthenticateAndGetUserId(ctx)
			if entityId == 0 {
				hasPermissions = false
			}
		}

	case League:
		if accessType != Create {
			entityId = getLeagueId(ctx)
		}
		hasPermissions, err = Access.League(accessType, permissions, LeagueDAO, entityId)
	case Team:
		if accessType != Create {
			entityId = getTeamId(ctx)
		}
		hasPermissions, err = Access.Team(accessType, permissions, TeamDAO, LeagueDAO, leagueId, entityId)
	case Player:
		if accessType != Create {
			entityId = getPlayerId(ctx)
		}
		hasPermissions, err = Access.Player(accessType, permissions, TeamDAO, LeagueDAO, leagueId, getTeamId(ctx), entityId)
	case Game:
		if accessType != Create {
			entityId = getGameId(ctx)
		}
		hasPermissions, err = Access.Game(accessType, permissions, GameDAO, leagueId, entityId)
	case Availability:
		if accessType != Create {
			entityId = getAvailabilityId(ctx)
		}
		hasPermissions, err = Access.Availability(accessType, permissions, LeagueDAO, leagueId, entityId)
	}

	if err != nil {
		return false, err
	} else {
		return hasPermissions, nil
	}
}

type endpoint struct {
	Entity        Entity
	AccessType    validation.AccessType
	BindData      func(ctx *gin.Context) bool
	IsDataInvalid func(ctx *gin.Context) (bool, string, error)
	Core          func(ctx *gin.Context) (interface{}, error)
}

func (e endpoint) createEndpointHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//buf := make([]byte, 1024)
		//num, _ := ctx.Request.Body.Read(buf)
		//fmt.Printf("%+v\n", string(buf[0:num]))
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
		if returnData == nil && checkErr(ctx, err) {
			return
		}
		// Return status and data if exists to router
		if returnData == nil {
			if e.AccessType == Create {
				ctx.Status(http.StatusCreated)
			} else {
				ctx.Status(http.StatusOK)
			}
		} else if returnData != nil && err != nil {
			ctx.JSON(http.StatusBadRequest, returnData)
		} else {
			if e.AccessType == Create {
				ctx.JSON(http.StatusCreated, returnData)
			} else {
				ctx.JSON(http.StatusOK, returnData)
			}
		}
	}
}
