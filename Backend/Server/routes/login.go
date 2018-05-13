package routes

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/kataras/iris"
	"database/sql"
	"github.com/kataras/iris/sessions"
	"log"
	"encoding/hex"
	"golang.org/x/crypto/scrypt"
	"bytes"
)

func RegisterLoginHandlers(app iris.Party, db *sql.DB, sessions *sessions.Sessions) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
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
	app.Post("/", func(ctx iris.Context) {
		//get params
		var usrInfo userInfo
		err := ctx.ReadJSON(&usrInfo)
		if checkErr(ctx, err) {return}

		//check parameters
		if failIfPasswordTooShort(usrInfo.Password, ctx) {return}
		if failIfEmailMalformed(usrInfo.Email, ctx) {return}
		if failIfEmailNotInUse(usrInfo.Email, ctx, psql, db) {return}
		//var email string

		//read variables from db
		var id int
		var salt string
		var storedHash string

		err = psql.Select("id", "salt", "hash").From("users").Where("email = ?", usrInfo.Email).
		RunWith(db).QueryRow().Scan(&id, &salt, &storedHash)

		if err != nil {
			log.Fatal(err)
			println(err.Error())
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(errorResponse{Error: "invalidLogin"})
			return
		}

		//check if password matches
		saltBin, err := hex.DecodeString(salt)
		if checkErr(ctx, err) {return}
		storedHashBin, err := hex.DecodeString(storedHash)
		if checkErr(ctx, err) {return}

		hash, err := scrypt.Key([]byte(usrInfo.Password), saltBin, 32768, 8, 1, 64)
		if checkErr(ctx, err) {return}

		if !bytes.Equal(hash[:], storedHashBin[:]) {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(errorResponse{Error: "invalidLogin"})
			return
		}

		//set authenticated to true and stored ID of logged in user in session tracker
		session := sessions.Start(ctx)
		session.Set("authenticated", true)
		session.Set("userID", id)

		ctx.StatusCode(iris.StatusOK)
	})//post
}