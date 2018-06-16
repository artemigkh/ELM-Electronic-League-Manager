package sessionManager

import (
	"github.com/gorilla/sessions"
	"github.com/gin-gonic/gin"
)

type CookieSessionManager struct {
	store *sessions.CookieStore
}

func (s *CookieSessionManager) AuthenticateAndGetUserID(ctx *gin.Context) (int, error) {
	session, err := s.store.Get(ctx.Request, "lm-session")
	if err != nil {
		return -1, err
	}

	authValue := session.Values["authenticated"]
	if authValue == nil {
		return -1, nil
	}

	authenticated := authValue.(bool)
	if !authenticated {
		return -1, nil
	}

	IDValue := session.Values["userID"]
	userID := IDValue.(int)

	return userID, nil
}

func (s *CookieSessionManager) LogIn(ctx *gin.Context, userID int) error {
	session, err := s.store.Get(ctx.Request, "lm-session")
	if err != nil {
		return err
	}

	session.Values["authenticated"] = true
	session.Values["userID"] = userID
	session.Save(ctx.Request, ctx.Writer)
	return nil
}