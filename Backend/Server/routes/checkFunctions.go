package routes

import (
	"github.com/badoux/checkmail"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	MIN_PASSWORD_LENGTH = 8
	MAX_LEAGUE_LENGTH = 50
	MAX_TEAM_LENGTH = 50
	MAX_TAG_LENGTH = 5
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
	if len(name) > MAX_TEAM_LENGTH {
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