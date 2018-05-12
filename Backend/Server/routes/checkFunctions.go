package routes

import (
	"github.com/kataras/iris"
	"github.com/badoux/checkmail"
	"database/sql"
	"log"
	"github.com/Masterminds/squirrel"
)

const(
	MIN_PASSWORD_LENGTH = 8
)

func checkErr(ctx iris.Context, err error) bool {
	if err != nil {
		log.Fatal(err)
		println(err.Error())
		ctx.StatusCode(iris.StatusBadRequest)
		return true
	} else {
		return false
	}
}

func checkPasswordLength(password string, ctx iris.Context) bool {
	if len(password) < 8 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(errorResponse{Error: "passwordTooShort"})
		return true
	} else {
		return false
	}
}

func checkEmailWellFormed(email string, ctx iris.Context) bool {
	err := checkmail.ValidateFormat(email)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(errorResponse{Error: "emailMalformed"})
		return true
	} else {
		return false
	}
}

func checkEmailDoesntExist(emailToCheck string, ctx iris.Context, psql squirrel.StatementBuilderType, db sql.DB ) bool {
	//check if email already exists
	var email sql.NullString
	err := psql.Select("email").
		From("users").
		Where("email = ?", emailToCheck).
		RunWith(&db).QueryRow().Scan(&email)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
		ctx.StatusCode(iris.StatusBadRequest)
		return true
	} else if email.Valid {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(errorResponse{Error: "emailInUse"})
		return true
	} else {
		return false
	}
}

func checkEmailExists(emailToCheck string, ctx iris.Context, psql squirrel.StatementBuilderType, db sql.DB ) bool {
	//check if email already exists
	var email sql.NullString
	err := psql.Select("email").
		From("users").
		Where("email = ?", emailToCheck).
		RunWith(&db).QueryRow().Scan(&email)

	if err == sql.ErrNoRows{
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(errorResponse{Error: "invalidLogin"})
		return true
	} else if err != nil {
		log.Fatal(err)
		ctx.StatusCode(iris.StatusBadRequest)
		return true
	} else {
		return false
	}
}
