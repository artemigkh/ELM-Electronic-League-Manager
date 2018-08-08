package sessionManager

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
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
		println(err.Error())
		return err
	}

	session.Values["authenticated"] = true
	session.Values["userID"] = userID
	session.Save(ctx.Request, ctx.Writer)
	return nil
}

func (s *CookieSessionManager) LogOut(ctx *gin.Context) error {
	session, err := s.store.Get(ctx.Request, "lm-session")
	if err != nil {
		println(err.Error())
		return err
	}

	session.Values["authenticated"] = false
	session.Values["userID"] = -1
	session.Save(ctx.Request, ctx.Writer)
	return nil
}

func (s *CookieSessionManager) SetActiveLeague(ctx *gin.Context, leagueID int) error {
	session, err := s.store.Get(ctx.Request, "lm-session")
	if err != nil {
		println(err.Error())
		return err
	}

	session.Values["leagueID"] = leagueID
	session.Save(ctx.Request, ctx.Writer)
	return nil
}

func (s *CookieSessionManager) GetActiveLeague(ctx *gin.Context) (int, error) {
	session, err := s.store.Get(ctx.Request, "lm-session")
	if err != nil {
		println(err.Error())
		return -1, err
	}

	IDValue := session.Values["leagueID"]

	if IDValue == nil {
		return -1, nil
	}

	leagueID := IDValue.(int)

	return leagueID, nil
}
