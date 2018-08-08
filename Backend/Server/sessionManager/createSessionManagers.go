package sessionManager

import (
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

func CreateCookieSessionManager() SessionManager {
	cookieStore := sessions.NewCookieStore(
		securecookie.GenerateRandomKey(64),
		securecookie.GenerateRandomKey(32))
	sessions.NewSession(cookieStore, "lm-session")

	return &CookieSessionManager{
		store: cookieStore,
	}

}
