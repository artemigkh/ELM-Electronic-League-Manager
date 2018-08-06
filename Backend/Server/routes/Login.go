package routes

import (
	"github.com/gin-gonic/gin"
	"encoding/hex"
	"golang.org/x/crypto/scrypt"
	"bytes"
	"net/http"
)

type loginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

/**
  * @api{POST} /login/ Get authentication cookie
  * @apiName createNewUser
  * @apiGroup login
  * @apiDescription Provide user email and password to get login authorization
  *
  * @apiParam {string} email
  * @apiParam {string} password
  *
  * @apiError passwordTooShort The password was too short
  * @apiError emailMalformed The email was not formed correctly
  * @apiError invalidLogin The user does not exist or password was incorrect
  */

func login(ctx *gin.Context) {
	//get parameters
	var request loginRequest
	err := ctx.ShouldBindJSON(&request)
	if checkJsonErr(ctx, err) {
		return
	}

	//check parameters
	if failIfPasswordTooShort(ctx, request.Password) {
		return
	}
	if failIfEmailMalformed(ctx, request.Email) {
		return
	}
	if failIfEmailNotInUse(ctx, request.Email) {
		return
	}

	id, salt, storedHash, err := UsersDAO.GetAuthenticationInformation(request.Email)
	if checkErr(ctx, err) {
		return
	}

	//check if password matches
	saltBin, err := hex.DecodeString(salt)
	if checkErr(ctx, err) {
		return
	}

	storedHashBin, err := hex.DecodeString(storedHash)
	if checkErr(ctx, err) {
		return
	}

	hash, err := scrypt.Key([]byte(request.Password), saltBin, 32768, 8, 1, 64)
	if checkErr(ctx, err) {
		return
	}

	if !bytes.Equal(hash[:], storedHashBin[:]) {
		//passwords do not match
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalidLogin"})
		return
	}

	err = ElmSessions.LogIn(ctx, id)
	if checkErr(ctx, err) {
		return
	}
}

func logout(ctx *gin.Context) {
	err := ElmSessions.LogOut(ctx)
	if checkErr(ctx, err) {
		return
	}
}

func RegisterLoginHandlers(g *gin.RouterGroup) {
	g.POST("/login", login)
	g.POST("/logout", logout)
}
