package databaseAccess

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
)

var db *sql.DB
var psql squirrel.StatementBuilderType

//var usersDAO databaseAccess.UsersDAO
var leaguesDAO LeaguesDAO = &PgLeaguesDAO{}
var teamsDAO TeamsDAO = &PgTeamsDAO{}

//var gamesDAO databaseAccess.GamesDAO
//var leagueOfLegendsDAO databaseAccess.LeagueOfLegendsDAO
