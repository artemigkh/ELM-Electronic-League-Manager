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
	profile, err := UsersDAO.GetUserProfile(ctx.GetInt("getUserId"))
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, profile)
}

/**
 * @api{GET} /api/users/permissions Get Permissions
 * @apiGroup Users
 * @apiDescription Get permissions in currently active league and any teams in league if applicable
 *
 * @apiSuccess {Object} leaguePermissions The permissions the user has in this league
 * @apiSuccess {bool} leaguePermissions.administrator if user is a league administrator
 * @apiSuccess {bool} leaguePermissions.createTeams if the user can create teams
 * @apiSuccess {bool} leaguePermissions.editTeams if the user can edit existing teams
 * @apiSuccess {bool} leaguePermissions.editGames if the user can edit games in this league
 * @apiSuccess {[]Object} teamPermissions The permissions this user has for teams in the league
 * @apiSuccess {number} teamPermissions.id The unique numerical identifer of the team
 * @apiSuccess {bool} teamPermissions.administrator If the user is an admin of this team
 * @apiSuccess {bool} teamPermissions.information If the user can change team information
 * @apiSuccess {bool} teamPermissions.players If the user can edit players on the team
 * @apiSuccess {bool} teamPermissions.reportResults If the user can report game results of this team
 *
 * @apiError notLoggedIn No user is logged in
 * @apiError noActiveLeague There is no active league selected
 */
//func getUserPermissions(ctx *gin.Context) {
//	userPermissions, err := UsersDAO.GetPermissions(
//		ctx.GetInt("getLeagueId"), ctx.GetInt("getUserId"))
//
//	if checkErr(ctx, err) {
//		return
//	}
//	ctx.JSON(http.StatusOK, userPermissions)
//}

func RegisterUserHandlers(g *gin.RouterGroup) {
	g.POST("/", createNewUser)
	g.GET("/profile", getProfile)
	//g.GET("/permissions", authenticate(), getActiveLeague(), getUserPermissions)
}
