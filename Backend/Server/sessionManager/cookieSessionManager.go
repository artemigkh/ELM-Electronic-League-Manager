package sessionManager

import "github.com/gorilla/sessions"

type SessionManager interface {
	AuthenticateAndGetUserID() (int, error)
}

type CookieSessionManager struct {
	store *sessions.CookieStore
}

func (s *CookieSessionManager) AuthenticateAndGetUserID() (int, error) {
	return 0, nil
}
