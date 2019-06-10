package sessionManager

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

type CookieSessionManager struct {
	store *sessions.CookieStore
}

func (s *CookieSessionManager) AuthenticateAndGetUserId(ctx *gin.Context) (int, error) {
	session, err := s.store.Get(ctx.Request, "lm-session")
	if err != nil {
		return -1, err
	}

	authValue := session.Values["authenticated"]
	if authValue == nil {
		return 0, nil
	}

	authenticated := authValue.(bool)
	if !authenticated {
		return 0, nil
	}

	IdValue := session.Values["userId"]
	userId := IdValue.(int)

	return userId, nil
}

func (s *CookieSessionManager) LogIn(ctx *gin.Context, userId int) error {
	session, err := s.store.Get(ctx.Request, "lm-session")
	if err != nil {
		println(err.Error())
		return err
	}

	session.Values["authenticated"] = true
	session.Values["userId"] = userId
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
	session.Values["userId"] = -1
	session.Save(ctx.Request, ctx.Writer)
	return nil
}

func (s *CookieSessionManager) SetActiveLeague(ctx *gin.Context, leagueId int) error {
	session, err := s.store.Get(ctx.Request, "lm-session")
	if err != nil {
		println(err.Error())
		return err
	}

	session.Values["leagueId"] = leagueId
	session.Save(ctx.Request, ctx.Writer)
	return nil
}

func (s *CookieSessionManager) GetActiveLeague(ctx *gin.Context) (int, error) {
	session, err := s.store.Get(ctx.Request, "lm-session")
	if err != nil {
		println(err.Error())
		return -1, err
	}

	IdValue := session.Values["leagueId"]

	if IdValue == nil {
		return 0, nil
	}

	return IdValue.(int), nil
}
