package routes

import (
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"golang.org/x/crypto/scrypt"
	"net/http"
)

type userProfile struct {
	Id int `json:"id"`
}

/**
 * @api{POST} /api/users/ Create a new user
 * @apiGroup Users
 * @apiDescription Register a new user in the database
 *
 * @apiParam {string} email The email of the user to be created
 * @apiParam {string} password The password of the user to be created
 *
 * @apiError passwordTooShort 400 The password was too short
 * @apiError emailMalformed 400 The email was not formed correctly
 * @apiError emailInUse 400 This email is already in use
 */
//TODO: Add Maximum (64 char) password length
func createNewUser(ctx *gin.Context) {
	//get parameters
	var usrInfo userInfo
	err := ctx.ShouldBindJSON(&usrInfo)
	if checkJsonErr(ctx, err) {
		return
	}

	//check parameters
	if failIfPasswordTooShort(ctx, usrInfo.Password) {
		return
	}
	if failIfEmailMalformed(ctx, usrInfo.Email) {
		return
	}
	if failIfEmailInUse(ctx, usrInfo.Email) {
		return
	}

	//create users password information
	salt := securecookie.GenerateRandomKey(32)
	hash, err := scrypt.Key([]byte(usrInfo.Password), salt, 32768, 8, 1, 64)
	if checkErr(ctx, err) {
		return
	}

	//create user in database
	err = UsersDAO.CreateUser(usrInfo.Email, hex.EncodeToString(salt), hex.EncodeToString(hash))
	if checkErr(ctx, err) {
		return
	}
}

/**
 * @api{GET} /api/users/profile Get Profile
 * @apiGroup Users
 * @apiDescription If a user is logged in, get their profile information
 *
 * @apiSuccess {string} email The email of the currently logged in user
 *
 * @apiError notLoggedIn 403 No user is currently logged in
 */
func getProfile(ctx *gin.Context) {
	profile, err := UsersDAO.GetUserProfile(ctx.GetInt("userId"))
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, profile)
}

func RegisterUserHandlers(g *gin.RouterGroup) {
	g.POST("/", createNewUser)
	g.GET("/profile", authenticate(), getProfile)
}
