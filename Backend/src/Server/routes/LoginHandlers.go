package routes

import (
	"Server/dataModel"
	"bytes"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/scrypt"
	"net/http"
)

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/logIn
func login(ctx *gin.Context) {
	var request dataModel.LoginRequest
	err := ctx.ShouldBindJSON(&request)
	if checkJsonErr(ctx, err) {
		return
	}

	valid, problem, err := request.Validate()
	if DataInvalid(ctx, valid, problem, err) {
		return
	}

	authInfo, err := UserDAO.GetAuthenticationInformation(request.Email)
	if checkErr(ctx, err) {
		return
	}

	//check if password matches
	saltBin, err := hex.DecodeString(authInfo.Salt)
	if checkErr(ctx, err) {
		return
	}

	storedHashBin, err := hex.DecodeString(authInfo.Hash)
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

	err = ElmSessions.LogIn(ctx, authInfo.UserId)
	if checkErr(ctx, err) {
		return
	}
	user, err := UserDAO.GetUserProfile(authInfo.UserId)
	ctx.JSON(http.StatusOK, user)
}

/**
 * @api{POST} /logout/ Logout
 * @apiGroup Login
 * @apiDescription Set client state to logged out
 */
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
