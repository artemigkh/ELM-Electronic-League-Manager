package routes

import (
	"esports-league-manager/Backend/Server/databaseAccess"
	"esports-league-manager/Backend/Server/sessionManager"
)

var UsersDAO databaseAccess.UsersDAO
var LeaguesDAO databaseAccess.LeaguesDAO
var TeamsDAO databaseAccess.TeamsDAO
var GamesDAO databaseAccess.GamesDAO

var ElmSessions sessionManager.SessionManager
