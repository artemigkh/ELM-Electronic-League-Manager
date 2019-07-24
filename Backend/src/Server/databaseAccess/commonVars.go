package databaseAccess

import (
	"Server/dataModel"
	"database/sql"
	"github.com/Masterminds/squirrel"
)

var db *sql.DB
var psql squirrel.StatementBuilderType

//var usersDAO databaseAccess.UsersDAO
var Leagues dataModel.LeagueDAO = &LeagueSqlDao{}
var teamsDAO dataModel.TeamDAO = &TeamSqlDao{}
var usersDAO dataModel.UserDAO = &UserSqlDao{}

//var gamesDAO databaseAccess.GamesDAO
//var leagueOfLegendsDAO databaseAccess.LeagueOfLegendsDAO
