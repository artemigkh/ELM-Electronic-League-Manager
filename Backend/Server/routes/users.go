package routes

import (
	"github.com/kataras/iris"
	sq "github.com/Masterminds/squirrel"
	"database/sql"
	"github.com/gorilla/securecookie"
	"golang.org/x/crypto/scrypt"
	"encoding/hex"
	"github.com/kataras/iris/sessions"
)

type userProfile struct {
	Id int `json:"id"`
}

func RegisterUserHandlers(app iris.Party, db *sql.DB, sessions *sessions.Sessions) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	/**
	  * @api{POST} /api/users/ Create a new user
	  * @apiName createNewUser
	  * @apiGroup users
	  * @apiDescription Register a new user in the database
	  *
	  * @apiParam {string} email
	  * @apiParam {string} password
	  *
	  * @apiError passwordTooShort The password was too short
	  * @apiError emailMalformed The email was not formed correctly
	  * @apiError emailInUse This email is already in use
	  */
	app.Post("/", func(ctx iris.Context) {
		//get params
		var usrInfo userInfo
		err := ctx.ReadJSON(&usrInfo)
		if checkErr(ctx, err) {return}
		//check parameters
		if failIfPasswordTooShort(usrInfo.Password, ctx) {return}
		if failIfEmailMalformed(usrInfo.Email, ctx) {return}
		if failIfEmailInUse(usrInfo.Email, ctx, psql, db) {return}

		//create users password information
		salt := securecookie.GenerateRandomKey(32)
		hash, err := scrypt.Key([]byte(usrInfo.Password), salt, 32768, 8, 1, 64)
		if checkErr(ctx, err) {return}

		//create user in database
		_, err = psql.Insert("users").Columns("email", "salt", "hash").
			Values(usrInfo.Email, hex.EncodeToString(salt), hex.EncodeToString(hash)).RunWith(db).Exec()
		if checkErr(ctx, err) {return}
	})//post

	app.Get("/profile", func(ctx iris.Context) {
		session := sessions.Start(ctx)
		//check if user logged in
		if auth, _ := session.GetBoolean("authenticated"); !auth {
			ctx.StatusCode(iris.StatusForbidden)
			return
		}

		//get id
		userID, _ := session.GetInt("userID")

		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(userProfile{Id: userID})
	})
}//register function
