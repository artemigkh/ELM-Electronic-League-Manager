package routes

import (
	"github.com/badoux/checkmail"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const (
	MIN_PASSWORD_LENGTH = 8
)

func checkErr(ctx *gin.Context, err error) bool {
	if err != nil {
		log.Fatal(err)
		println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "malformed input"})
		return true
	} else {
		return false
	}
}

func failIfPasswordTooShort(ctx *gin.Context, password string) bool {
	if len(password) < 8 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "passwordTooShort"})
		return true
	} else {
		return false
	}
}

func failIfEmailMalformed(ctx *gin.Context, email string) bool {
	err := checkmail.ValidateFormat(email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "emailMalformed"})
		return true
	} else {
		return false
	}
}

func failIfEmailInUse(ctx *gin.Context, emailToCheck string) bool {
	//check if email already exists
	inUse, err := UsersDAO.IsEmailInUse(emailToCheck)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return true
	} else if inUse {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "emailInUse"})
		return true
	} else {
		return false
	}
}

//func failIfEmailNotInUse(emailToCheck string, ctx gin.Context, psql squirrel.StatementBuilderType, db *sql.DB ) bool {
//	//check if email already exists
//	var email string
//	err := psql.Select("email").
//		From("users").
//		Where("email = ?", emailToCheck).
//		RunWith(db).QueryRow().Scan(&email)
//
//	if err != nil {
//		ctx.StatusCode(iris.StatusBadRequest)
//		ctx.JSON(errorResponse{Error: "invalidLogin"})
//		return true
//	} else {
//		return false
//	}
//}
//
//func failIfLeagueNameInUse(leagueName string, ctx gin.Context, psql squirrel.StatementBuilderType, db *sql.DB ) bool {
//	var name string
//	err := psql.Select("name").
//		From("leagues").
//		Where("name = ?", leagueName).
//		RunWith(db).QueryRow().Scan(&name)
//	if err != nil {
//		return false
//	} else {
//		ctx.JSON(errorResponse{Error: "nameInUse"})
//		ctx.StatusCode(iris.StatusBadRequest)
//		return true
//	}
//}
//
//func failIfLeagueNameTooLong(leagueName string, ctx gin.Context) bool {
//	if len(leagueName) > 50 {
//		ctx.StatusCode(iris.StatusBadRequest)
//		ctx.JSON(errorResponse{Error: "nameTooLong"})
//		return true
//	} else {
//		return false
//	}
//}
