package routes

import (
	"github.com/gin-gonic/gin"
	"esports-league-manager/Backend/Server/databaseAccess"
	"golang.org/x/crypto/scrypt"
	"encoding/hex"
	"github.com/gorilla/securecookie"
)

type userProfile struct {
	Id int `json:"id"`
}

var UsersDAO databaseAccess.UsersDAO

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
	//get params
	var usrInfo userInfo
	err := ctx.ShouldBindJSON(&usrInfo)
	if checkErr(ctx, err) {
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
	err = UsersDAO.InsertUser(usrInfo.Email, hex.EncodeToString(salt), hex.EncodeToString(hash))
	if checkErr(ctx, err) {
		return
	}
}

func RegisterUserHandlers(g *gin.RouterGroup) {
	g.POST("/", createNewUser)

	//app.Get("/profile", func(ctx iris.Context) {
	//	session := sessions.Start(ctx)
	//	userID := authenticateAndGetCurrUserId(ctx, session)
	//	if userID == -1 {
	//		return
	//	}
	//
	//	ctx.StatusCode(iris.StatusOK)
	//	ctx.JSON(userProfile{Id: userID})
	//})
} //register function
