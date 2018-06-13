package routes

import (
	"esports-league-manager/Backend/Server/sessionManager"
	"esports-league-manager/Backend/Server/databaseAccess"
)

var UsersDAO databaseAccess.UsersDAO
var ElmSessions sessionManager.SessionManager
