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
 * For consistency across all function signatures, each function should start with the gin context
 * Then, the numerical Ids which should be in order of magnitude of entity:
 * first should be league, then team, then game, then user, then player
 * Then, the remaining parameters should be included in a logical order
 */

const (
	MIN_PASSWORD_LENGTH = 8
	MAX_LEAGUE_LENGTH   = 50
	MAX_NAME_LENGTH     = 50
	MAX_TAG_LENGTH      = 5
)

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
		ctx.JSON(http.StatusInternalServerError, nil)
		return true
	} else {
		return false
	}
}

func failIfPasswordTooShort(ctx *gin.Context, password string) bool {
	if len(password) < MIN_PASSWORD_LENGTH {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "passwordTooShort"})
		return true
	} else {
		return false
	}
}

func failIfEmailMalformed(ctx *gin.Context, email string) bool {
	err := checkmail.ValidateFormat(email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "emailMalformed"})
		return true
	} else {
		return false
	}
}

func failIfEmailInUse(ctx *gin.Context, emailToCheck string) bool {
	inUse, err := UsersDAO.IsEmailInUse(emailToCheck)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return true
	} else if inUse {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "emailInUse"})
		return true
	} else {
		return false
	}
}

func failIfEmailNotInUse(ctx *gin.Context, emailToCheck string) bool {
	inUse, err := UsersDAO.IsEmailInUse(emailToCheck)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return true
	} else if !inUse {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalidLogin"})
		return true
	} else {
		return false
	}
}

func failIfLeagueNameTooLong(ctx *gin.Context, name string) bool {
	if len(name) > MAX_LEAGUE_LENGTH {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "nameTooLong"})
		return true
	} else {
		return false
	}
}

func failIfLeagueNameInUse(ctx *gin.Context, name string) bool {
	inUse, err := LeaguesDAO.IsNameInUse(name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return true
	} else if inUse {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "nameInUse"})
		return true
	} else {
		return false
	}
}

func failIfTeamNameTooLong(ctx *gin.Context, name string) bool {
	if len(name) > MAX_NAME_LENGTH {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "nameTooLong"})
		return true
	} else {
		return false
	}
}

func failIfTeamTagTooLong(ctx *gin.Context, name string) bool {
	if len(name) > MAX_TAG_LENGTH {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "tagTooLong"})
		return true
	} else {
		return false
	}
}

func failIfTeamInfoInUse(ctx *gin.Context, leagueId int, name, tag string) bool {
	inUse, errorMsg, err := TeamsDAO.IsInfoInUse(leagueId, name, tag)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return true
	} else if inUse {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return true
	} else {
		return false
	}
}

func failIfTeamDoesNotExist(ctx *gin.Context, leagueId, teamId int) bool {
	exists, err := TeamsDAO.DoesTeamExist(leagueId, teamId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return true
	} else if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "teamDoesNotExist"})
		return true
	} else {
		return false
	}
}

func failIfConflictExists(ctx *gin.Context, team1Id, team2Id, gameTime int) bool {
	conflictExists, err := GamesDAO.DoesExistConflict(team1Id, team2Id, gameTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return true
	} else if conflictExists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "conflictExists"})
		return true
	} else {
		return false
	}
}

func failIfGameDoesNotExist(ctx *gin.Context, leagueId, teamId int) bool {
	gameInformation, err := GamesDAO.GetGameInformation(leagueId, teamId)
	if checkErr(ctx, err) {
		ctx.JSON(http.StatusInternalServerError, nil)
		return true
	}

	if gameInformation == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "gameDoesNotExist"})
		return true
	}

	return false
}

func failIfCannotEditPlayersOnTeam(ctx *gin.Context, leagueId, teamId, userId int) bool {
	canEditPlayers, err := TeamsDAO.HasPlayerEditPermissions(leagueId, teamId, userId)
	if err != nil {
		println(err.Error())
		ctx.JSON(http.StatusInternalServerError, nil)
		return true
	} else if !canEditPlayers {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "canNotEditPlayers"})
		return true
	} else {
		return false
	}
}

func failIfGameIdentifierTooLong(ctx *gin.Context, gameIdentifier string) bool {
	if len(gameIdentifier) > MAX_NAME_LENGTH {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "gameIdentifierTooLong"})
		return true
	} else {
		return false
	}
}

func failIfNameTooLong(ctx *gin.Context, name string) bool {
	if len(name) > MAX_NAME_LENGTH {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "nameTooLong"})
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

func failIfPlayerDoesNotExist(ctx *gin.Context, teamId, playerId int) bool {
	playerExists, err := TeamsDAO.DoesPlayerExist(teamId, playerId)
	if err != nil {
		println(err.Error())
		ctx.JSON(http.StatusInternalServerError, nil)
		return true
	} else if !playerExists {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "playerDoesNotExist"})
		return true
	} else {
		return false
	}
}
