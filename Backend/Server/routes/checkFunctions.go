package routes

import (
	"github.com/badoux/checkmail"
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
 * This file contains functions for checking user json input in endpoint handlers
 * If a check does not require information from json, it should instead be a middleware handler
 *
 * The functions in this file return true on fail so that the handlers can check outcome to exit early
 *
 * For consistency across all function signatures, each function should start with the gin context
 * Then, the numerical Ids which should be in order of magnitude of entity:
 * first should be league, then team, then game, then user, then player
 * Then, the remaining parameters should be included in a logical order
 */

const (
	MinPasswordLength = 8
	MaxNameLength     = 50
	MaxTagLength      = 5
)

// operator wrappers
func le(x, y int) bool {
	return x < y
}

func ge(x, y int) bool {
	return x > y
}

// General Cases
func checkJsonErr(ctx *gin.Context, err error) bool {
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "malformedInput"})
		return true
	} else {
		return false
	}
}

func checkErr(ctx *gin.Context, err error) bool {
	if err != nil {
		println(err.Error())
		ctx.Status(http.StatusInternalServerError)
		return true
	} else {
		return false
	}
}

func failIfBooleanConditionTrue(ctx *gin.Context, cond bool, err error, responseCode int, errorString string) bool {
	if checkErr(ctx, err) {
		return true
	}
	if cond {
		ctx.JSON(responseCode, gin.H{"error": errorString})
	}
	return cond
}

func failIfImproperLength(ctx *gin.Context, s string, length int, comparator func(int, int) bool, errorString string) bool {
	if comparator(len(s), length) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errorString})
		return true
	} else {
		return false
	}
}

// Input Length Checks
func failIfPasswordTooShort(ctx *gin.Context, password string) bool {
	return failIfImproperLength(ctx, password, MinPasswordLength, le, "passwordTooShort")
}

func failIfTeamTagTooLong(ctx *gin.Context, name string) bool {
	return failIfImproperLength(ctx, name, MaxTagLength, ge, "tagTooLong")
}

func failIfGameIdentifierTooLong(ctx *gin.Context, gameIdentifier string) bool {
	return failIfImproperLength(ctx, gameIdentifier, MaxNameLength, ge, "gameIdentifierTooLong")
}

func failIfNameTooLong(ctx *gin.Context, name string) bool {
	return failIfImproperLength(ctx, name, MaxNameLength, ge, "nameTooLong")
}

// Boolean Checks
func failIfEmailInUse(ctx *gin.Context, emailToCheck string) bool {
	inUse, err := UsersDAO.IsEmailInUse(emailToCheck)
	return failIfBooleanConditionTrue(ctx, inUse, err, http.StatusBadRequest, "emailInUse")
}

func failIfEmailNotInUse(ctx *gin.Context, emailToCheck string) bool {
	inUse, err := UsersDAO.IsEmailInUse(emailToCheck)
	return failIfBooleanConditionTrue(ctx, !inUse, err, http.StatusBadRequest, "invalidLogin")
}

func failIfLeagueNameInUse(ctx *gin.Context, name string) bool {
	inUse, err := LeaguesDAO.IsNameInUse(name)
	return failIfBooleanConditionTrue(ctx, inUse, err, http.StatusBadRequest, "nameInUse")
}

func failIfTeamInfoInUse(ctx *gin.Context, leagueId int, name, tag string) bool {
	inUse, errorMsg, err := TeamsDAO.IsInfoInUse(leagueId, name, tag)
	return failIfBooleanConditionTrue(ctx, inUse, err, http.StatusBadRequest, errorMsg)
}

func failIfTeamDoesNotExist(ctx *gin.Context, leagueId, teamId int) bool {
	exists, err := TeamsDAO.DoesTeamExist(leagueId, teamId)
	return failIfBooleanConditionTrue(ctx, !exists, err, http.StatusBadRequest, "teamDoesNotExist")
}

func failIfConflictExists(ctx *gin.Context, team1Id, team2Id, gameTime int) bool {
	conflictExists, err := GamesDAO.DoesExistConflict(team1Id, team2Id, gameTime)
	return failIfBooleanConditionTrue(ctx, conflictExists, err, http.StatusBadRequest, "conflictExists")
}

func failIfGameDoesNotExist(ctx *gin.Context, leagueId, gameId int) bool {
	gameInformation, err := GamesDAO.GetGameInformation(leagueId, gameId)
	return failIfBooleanConditionTrue(ctx, gameInformation == nil, err, http.StatusBadRequest, "gameDoesNotExist")
}

func failIfCannotEditPlayersOnTeam(ctx *gin.Context, leagueId, teamId, userId int) bool {
	canEditPlayers, err := TeamsDAO.HasPlayerEditPermissions(leagueId, teamId, userId)
	return failIfBooleanConditionTrue(ctx, !canEditPlayers, err, http.StatusForbidden, "canNotEditPlayers")
}

func failIfPlayerDoesNotExist(ctx *gin.Context, teamId, playerId int) bool {
	playerExists, err := TeamsDAO.DoesPlayerExist(teamId, playerId)
	return failIfBooleanConditionTrue(ctx, !playerExists, err, http.StatusBadRequest, "playerDoesNotExist")
}

// Misc. Checks
func failIfEmailMalformed(ctx *gin.Context, email string) bool {
	err := checkmail.ValidateFormat(email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "emailMalformed"})
		return true
	} else {
		return false
	}
}

func failIfGameIdentifierInUse(ctx *gin.Context, leagueId, teamId int, gameIdentifier string) bool {
	teamInfo, err := TeamsDAO.GetTeamInformation(leagueId, teamId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return true
	}

	for _, player := range teamInfo.Players {
		if player.GameIdentifier == gameIdentifier {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "gameIdentifierInUse"})
			return true
		}
	}

	return false
}
