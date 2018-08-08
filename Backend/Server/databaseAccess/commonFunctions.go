package databaseAccess

import (
	"database/sql"
	"esports-league-manager/Backend/Server/config"
	"github.com/Masterminds/squirrel"
	_ "github.com/jackc/pgx/stdlib"
	"log"
)

func openDB(conf config.Config) {
	var err error
	db, err = sql.Open("pgx", conf.GetDbConnString())
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func Init(conf config.Config) {
	openDB(conf)
	psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}
