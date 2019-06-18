package routes

import (
	"bytes"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/scrypt"
	"net/http"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

/**
 * @api{POST} /login/ Login
 * @apiGroup Login
 * @apiDescription Provide user email and password to get login authorization
 *
 * @apiParam {string} email
 * @apiParam {string} password
 *
 * @apiSuccess {int} id The unique numerical identifier of the user that successfully logged in
 *
 * @apiError passwordTooShort The password was too short
 * @apiError emailMalformed The email was not formed correctly
 * @apiError invalidLogin The user does not exist or password was incorrect
 */
//TODO: Add Maximum (64 char) password length
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

	authInfo, err := UsersDAO.GetAuthenticationInformation(request.Email)
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
	user, err := UsersDAO.GetUserProfile(authInfo.UserId)
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
