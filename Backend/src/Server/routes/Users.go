package routes

import (
	"Server/databaseAccess"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"golang.org/x/crypto/scrypt"
)

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/createUser
func createNewUser() gin.HandlerFunc {
	var user databaseAccess.UserCreationInformation
	return endpoint{
		Entity:        User,
		AccessType:    Create,
		BindData:      func(ctx *gin.Context) bool { return bindAndCheckErr(ctx, &user) },
		IsDataInvalid: func(ctx *gin.Context) (bool, string, error) { return user.Validate() },
		Core: func(ctx *gin.Context) (interface{}, error) {
			//create users password information
			salt := securecookie.GenerateRandomKey(32)
			hash, err := scrypt.Key([]byte(user.Password), salt, 32768, 8, 1, 64)
			if checkErr(ctx, err) {
				return nil, err
			}

			//create user in database
			return nil, UsersDAO.CreateUser(user.Email, hex.EncodeToString(salt), hex.EncodeToString(hash))
		},
	}.createEndpointHandler()
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/getUser
func getProfile() gin.HandlerFunc {
	return endpoint{
		Entity:     User,
		AccessType: View,
		Core: func(ctx *gin.Context) (interface{}, error) {
			return UsersDAO.GetUserProfile(getUserId(ctx))
		},
	}.createEndpointHandler()
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/getUserLeaguePermissions
func getUserLeaguePermissions() gin.HandlerFunc {
	return endpoint{
		Entity:     User,
		AccessType: View,
		Core: func(ctx *gin.Context) (interface{}, error) {
			return UsersDAO.GetUserWithPermissions(getLeagueId(ctx), getUserId(ctx))
		},
	}.createEndpointHandler()
}

func RegisterUserHandlers(g *gin.RouterGroup) {
	g.POST("", createNewUser())
	g.GET("", getProfile())
	g.GET("leaguePermissions", getUserLeaguePermissions())
}
