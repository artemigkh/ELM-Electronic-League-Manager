package databaseAccess

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
)

var db *sql.DB
var psql squirrel.StatementBuilderType