package routes

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/scrypt"
	"encoding/hex"
	"github.com/gorilla/securecookie"
	"net/http"
)

type userProfile struct {
	Id int `json:"id"`
}

/**
 * @api{POST} /api/users/ Create a new user
 * @apiName createNewUser
 * @apiGroup users
 * @apiDescription Register a new user in the database
 *
 * @apiParam {string} email
 * @apiParam {string} password
 *
 * @apiError passwordTooShort 400 The password was too short
 * @apiError emailMalformed 400 The email was not formed correctly
 * @apiError emailInUse 400 This email is already in use
 */
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
 * @api{POST} /api/users/profile Get current users profile
 * @apiName getUserProfile
 * @apiGroup users
 * @apiDescription If a user is logged in, get their profile information
 *
 * @apiSuccess {int} id the userID
 *
 * @apiError notLoggedIn 403 No user is currently logged in
 */
func getProfile(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, userProfile{
		Id: ctx.GetInt("userID"),
	})
}

func RegisterUserHandlers(g *gin.RouterGroup) {
	g.POST("/", createNewUser)
	g.GET("/profile", authenticate(), getProfile)
}
