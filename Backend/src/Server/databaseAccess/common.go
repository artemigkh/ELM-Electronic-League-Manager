package databaseAccess

import (
	"Server/config"
	"database/sql"
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

type RowArr interface {
	Scan(*sql.Rows) error
}

func ScanRows(statement squirrel.SelectBuilder, out RowArr) error {
	rows, err := statement.RunWith(db).Query()
	if err != nil {
		return nil
	}

	defer rows.Close()

	for rows.Next() {
		if err := out.Scan(rows); err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}
