package routes

import (
	"github.com/badoux/checkmail"
	"github.com/gin-gonic/gin"
	"net/http"
)

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

func failIfTeamInfoInUse(ctx *gin.Context, name, tag string, leagueId int) bool {
	inUse, errorMsg, err := TeamsDAO.IsInfoInUse(name, tag, leagueId)
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

func failIfTeamDoesNotExist(ctx *gin.Context, teamId, leagueId int) bool {
	exists, err := TeamsDAO.DoesTeamExist(teamId, leagueId)
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

func failIfGameDoesNotExist(ctx *gin.Context) bool {
	gameInformation, err := GamesDAO.GetGameInformation(ctx.GetInt("urlId"), ctx.GetInt("leagueId"))
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

func failIfCannotEditPlayersOnTeam(ctx *gin.Context, userId, teamId, leagueId int) bool {
	canEditPlayers, err := TeamsDAO.HasPlayerEditPermissions(teamId, userId, leagueId)
	if err != nil {
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

func failIfGameIdentifierInUse(ctx *gin.Context, gameIdentifier string, teamId, leagueId int) bool {
	teamInfo, err := TeamsDAO.GetTeamInformation(teamId, leagueId)
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

func failIfCannotJoinLeague(ctx *gin.Context, userId, leagueId int) bool {
	canJoin, err := LeaguesDAO.CanJoinLeague(userId, leagueId)
	if err != nil {
		println(err.Error())
		ctx.JSON(http.StatusInternalServerError, nil)
		return true
	} else if !canJoin {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "canNotJoin"})
		return true
	} else {
		return false
	}
}