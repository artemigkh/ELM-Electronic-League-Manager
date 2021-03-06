package sessionManager

import (
	"Server/config"
	"github.com/gorilla/sessions"
)

func CreateCookieSessionManager(conf config.Config) SessionManager {
	authKey, encryptionKey := conf.GetKeys()
	cookieStore := sessions.NewCookieStore(authKey, encryptionKey)
	//println(hex.EncodeToString(securecookie.GenerateRandomKey(64)))
	//println(hex.EncodeToString(securecookie.GenerateRandomKey(32)))
	sessions.NewSession(cookieStore, "lm-session")

	return &CookieSessionManager{
		store: cookieStore,
	}

}
