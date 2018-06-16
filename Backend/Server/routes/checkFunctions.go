package routes

import (
	"github.com/badoux/checkmail"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	MIN_PASSWORD_LENGTH = 8
)

func checkJsonErr(ctx *gin.Context, err error) bool {
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "malformed input"})
		return true
	} else {
		return false
	}
}

func checkErr(ctx *gin.Context, err error) bool {
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return true
	} else {
		return false
	}
}

func failIfPasswordTooShort(ctx *gin.Context, password string) bool {
	if len(password) < 8 {
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
	//check if email already exists
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
	//check if email already exists
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
