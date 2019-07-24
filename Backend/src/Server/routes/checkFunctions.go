package routes

import (
	"github.com/badoux/checkmail"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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

func bindRepeatedAndCheckErr(ctx *gin.Context, obj interface{}) bool {
	if err := ctx.ShouldBindBodyWith(&obj, binding.JSON); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "malformedInput"})
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

func failIfEmailNotInUse(ctx *gin.Context, emailToCheck string) bool {
	inUse, err := UserDAO.IsEmailInUse(emailToCheck)
	return failIfBooleanConditionTrue(ctx, !inUse, err, http.StatusBadRequest, "invalidLogin")
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
