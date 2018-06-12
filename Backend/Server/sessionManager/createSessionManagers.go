package sessionManager

import (
	"github.com/gorilla/sessions"
	"github.com/gorilla/securecookie"
)

func CreateCookieSessionManager() SessionManager {
	cookieStore := sessions.NewCookieStore(
		securecookie.GenerateRandomKey(64),
		securecookie.GenerateRandomKey(32),)

	return &CookieSessionManager{
		store: cookieStore,
	}

}