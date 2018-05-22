package routes

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
)

func authenticateAndGetCurrUserId(ctx iris.Context, session *sessions.Session) int {
	//check if user logged in
	if auth, _ := session.GetBoolean("authenticated"); !auth {
		ctx.StatusCode(iris.StatusForbidden)
		return -1
	}

	//get id
	userID, _ := session.GetInt("userID")
	return userID
}
