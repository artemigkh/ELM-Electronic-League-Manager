package databaseAccess

import (
	"database/sql"
	_ "github.com/jackc/pgx/stdlib"
	"log"
	"esports-league-manager/Backend/Server/config"
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
}