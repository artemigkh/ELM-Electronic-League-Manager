package sessionManager

import (
	"github.com/gorilla/sessions"
	"github.com/gin-gonic/gin"
)



type CookieSessionManager struct {
	store *sessions.CookieStore
}

func (s *CookieSessionManager) AuthenticateAndGetUserID(ctx gin.Context) (int, error) {
	session, err := s.store.Get(ctx.Request, "elm-session")
	if err != nil {
		return -1, err
	}
	authValue := session.Values["authenticated"]
	authenticated := authValue.(bool)

	if !authenticated {
		return -1, nil
	}

	IDValue := session.Values["userID"]
	userID := IDValue.(int)

	return userID, nil
}

func (s *CookieSessionManager) LogIn(ctx gin.Context, userID int) error {
	session, err := s.store.Get(ctx.Request, "elm-session")
	if err != nil {
		return err
	}

	session.Values["authenticated"] = true
	session.Values["userID"] = userID
	return nil
}