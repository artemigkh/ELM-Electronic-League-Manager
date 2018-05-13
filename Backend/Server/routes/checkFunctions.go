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

func failIfPasswordTooShort(password string, ctx iris.Context) bool {
	if len(password) < 8 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(errorResponse{Error: "passwordTooShort"})
		return true
	} else {
		return false
	}
}

func failIfEmailMalformed(email string, ctx iris.Context) bool {
	err := checkmail.ValidateFormat(email)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(errorResponse{Error: "emailMalformed"})
		return true
	} else {
		return false
	}
}

func failIfEmailInUse(emailToCheck string, ctx iris.Context, psql squirrel.StatementBuilderType, db *sql.DB ) bool {
	//check if email already exists
	var email string
	err := psql.Select("email").
		From("users").
		Where("email = ?", emailToCheck).
		RunWith(db).QueryRow().Scan(&email)
	if err != nil {
	 	return false
	} else {
		 ctx.JSON(errorResponse{Error: "emailInUse"})
		 ctx.StatusCode(iris.StatusBadRequest)
		 return true
	}
}

func failIfEmailNotInUse(emailToCheck string, ctx iris.Context, psql squirrel.StatementBuilderType, db *sql.DB ) bool {
	//check if email already exists
	var email string
	err := psql.Select("email").
		From("users").
		Where("email = ?", emailToCheck).
		RunWith(db).QueryRow().Scan(&email)

	println("got here")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(errorResponse{Error: "invalidLogin"})
		return true
	} else {
		return false
	}
}
