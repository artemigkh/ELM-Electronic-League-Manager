package databaseAccess

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
)

var db *sql.DB
var psql squirrel.StatementBuilderType

//var usersDAO databaseAccess.UsersDAO
var Leagues LeaguesDAO = &PgLeaguesDAO{}
var teamsDAO TeamsDAO = &PgTeamsDAO{}
var usersDAO UsersDAO = &PgUsersDAO{}

//var gamesDAO databaseAccess.GamesDAO
//var leagueOfLegendsDAO databaseAccess.LeagueOfLegendsDAO
