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
	MinPasswordLength    = 8
	MaxNameLength        = 50
	MaxTagLength         = 5
	MaxDescriptionLength = 500
	MinInformationLength = 2
	MaxMdLength          = 50000
	MaxOffsetSize        = 24*3600 - 1
	MaxHour              = 24
	MaxMinute            = 60
	MaxDuration          = 60 * 24
)

var ValidGameStrings = [...]string{
	"genericsport",
	"basketball",
	"curling",
	"football",
	"hockey",
	"rugby",
	"soccer",
	"volleyball",
	"waterpolo",
	"genericesport",
	"csgo",
	"leagueoflegends",
	"overwatch",
}

// operator wrappers
func le(x, y int) bool {
	return x < y
}

func ge(x, y int) bool {
	return x > y
}

// General Cases
func bindAndCheckErr(ctx *gin.Context, obj interface{}) bool {
	if err := ctx.ShouldBindJSON(&obj); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "malformedInput"})
		return true
	} else {
		return false
	}
}

func dataInvalid(
	ctx *gin.Context, valid bool, problem string, err error) bool {
	if checkErr(ctx, err) {
		return true
	} else if !valid {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": problem})
		return true
	} else {
		return false
	}
}

func accessForbidden(ctx *gin.Context, allowed bool, err error) bool {
	if checkErr(ctx, err) {
		return true
	} else if !allowed {
		ctx.Status(http.StatusForbidden)
		return true
	} else {
		return false
	}
}

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

func failIfNameTooShort(ctx *gin.Context, name string) bool {
	return failIfImproperLength(ctx, name, MinInformationLength, le, "nameTooShort")
}

func failIfTagTooShort(ctx *gin.Context, tag string) bool {
	return failIfImproperLength(ctx, tag, MinInformationLength, le, "tagTooShort")
}

func failIfGameIdentifierTooShort(ctx *gin.Context, tag string) bool {
	return failIfImproperLength(ctx, tag, MinInformationLength, le, "gameIdentifierTooShort")
}

func failIfGameIdentifierTooLong(ctx *gin.Context, gameIdentifier string) bool {
	return failIfImproperLength(ctx, gameIdentifier, MaxNameLength, ge, "gameIdentifierTooLong")
}

func failIfNameTooLong(ctx *gin.Context, name string) bool {
	return failIfImproperLength(ctx, name, MaxNameLength, ge, "nameTooLong")
}

func failIfMdTooLong(ctx *gin.Context, name string) bool {
	return failIfImproperLength(ctx, name, MaxMdLength, ge, "markdownTooLong")
}

func failIfDescriptionTooLong(ctx *gin.Context, description string) bool {
	return failIfImproperLength(ctx, description, MaxDescriptionLength, ge, "descriptionTooLong")
}
func failIfTimezoneOffsetTooLarge(ctx *gin.Context, offset float64) bool {
	return failIfBooleanConditionTrue(ctx, offset > MaxOffsetSize, nil, http.StatusBadRequest, "offsetTooLarge")
}
func failIfMinuteTooLarge(ctx *gin.Context, minute int) bool {
	return failIfBooleanConditionTrue(ctx, minute > MaxMinute, nil, http.StatusBadRequest, "invalidMinute")
}

func failIfHourTooLarge(ctx *gin.Context, hour int) bool {
	return failIfBooleanConditionTrue(ctx, hour > MaxHour, nil, http.StatusBadRequest, "invalidHour")
}

func failIfDurationTooLarge(ctx *gin.Context, duration int) bool {
	return failIfBooleanConditionTrue(ctx, duration > MaxDuration, nil, http.StatusBadRequest, "invalidDuration")
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

func failIfLeagueNameInUse(ctx *gin.Context, leagueId int, name string) bool {
	inUse, err := LeaguesDAO.IsNameInUse(leagueId, name)
	return failIfBooleanConditionTrue(ctx, inUse, err, http.StatusBadRequest, "nameInUse")
}

func failIfTeamInfoInUse(ctx *gin.Context, leagueId, teamId int, name, tag string) bool {
	inUse, errorMsg, err := TeamsDAO.IsInfoInUse(leagueId, teamId, name, tag)
	return failIfBooleanConditionTrue(ctx, inUse, err, http.StatusBadRequest, errorMsg)
}

func failIfTeamDoesNotExist(ctx *gin.Context, leagueId, teamId int) bool {
	exists, err := TeamsDAO.DoesTeamExistInLeague(leagueId, teamId)
	return failIfBooleanConditionTrue(ctx, !exists, err, http.StatusBadRequest, "teamDoesNotExist")
}

func failIfManagerDoesNotExist(ctx *gin.Context, teamId, userId int) bool {
	tp, err := TeamsDAO.GetTeamPermissions(teamId, userId)
	return failIfBooleanConditionTrue(ctx, !(tp.Administrator || tp.ReportResults || tp.Players || tp.Information),
		err, http.StatusBadRequest, "managerDoesNotExist")
}

func failIfConflictExists(ctx *gin.Context, team1Id, team2Id, gameTime int) bool {
	conflictExists, err := GamesDAO.DoesExistConflict(team1Id, team2Id, gameTime)
	return failIfBooleanConditionTrue(ctx, conflictExists, err, http.StatusBadRequest, "conflictExists")
}

func failIfGameDoesNotExist(ctx *gin.Context, gameId int) bool {
	gameInformation, err := GamesDAO.GetGameInformation(gameId)
	return failIfBooleanConditionTrue(ctx, gameInformation == nil, err, http.StatusBadRequest, "gameDoesNotExist")
}

//func failIfAvailabilityDoesNotExist(ctx *gin.Context, availabilityId int) bool {
//	availability, err := LeaguesDAO.GetSchedulingAvailability(availabilityId)
//	return failIfBooleanConditionTrue(ctx, availability == nil, err, http.StatusBadRequest, "AvailabilityDoesNotExist")
//}
//
//func failIfGameDoesNotContainWinner(ctx *gin.Context, gameId, winnerId int) bool {
//	gi, err := GamesDAO.GetGameInformation(gameId)
//	return failIfBooleanConditionTrue(ctx,
//		!(gi.Team1Id == winnerId || gi.Team2Id == winnerId), err, http.StatusBadRequest, "gameDoesNotContainWinner")
//}
func failIfPlayerDoesNotExist(ctx *gin.Context, teamId, playerId int) bool {
	playerExists, err := TeamsDAO.DoesPlayerExistInTeam(teamId, playerId)
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

func failIfGameStringtNotValid(ctx *gin.Context, game string) bool {
	for _, g := range ValidGameStrings {
		if g == game {
			return false
		}
	}
	ctx.JSON(http.StatusBadRequest, gin.H{"error": "gameStringNotValid"})
	return true
}

func failWithError(ctx *gin.Context, error string) bool {
	ctx.JSON(http.StatusBadRequest, gin.H{"error": "gameStringNotValid"})
	return true
}

func booleanConditionFailed(ctx *gin.Context, cond bool, err error, errorString string) bool {
	if checkErr(ctx, err) {
		return true
	}
	if cond {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errorString})
	}
	return cond
}

func gameStringFailed(ctx *gin.Context, game string) bool {
	for _, g := range ValidGameStrings {
		if g == game {
			return false
		}
	}
	return failWithError(ctx, "invalid game string")
}

func nameStringFailed(ctx *gin.Context, name string) bool {
	if len(name) < MinInformationLength {
		return failWithError(ctx, "name too short")
	} else if len(name) > MaxNameLength {
		return failWithError(ctx, "name too long")
	} else {
		return false
	}
}

func descriptionStringFailed(ctx *gin.Context, description string) bool {
	if len(description) > MaxDescriptionLength {
		return failWithError(ctx, "description too long")
	} else {
		return false
	}
}
