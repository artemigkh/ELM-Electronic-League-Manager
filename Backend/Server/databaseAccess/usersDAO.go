package databaseAccess

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	"log"
)

type PgUsersDAO struct {
	psql squirrel.StatementBuilderType
}

func (u *PgUsersDAO) InsertUser(email, salt, hash string) error {
	_, err := u.psql.Insert("users").Columns("email", "salt", "hash").
		Values(email, salt, hash).RunWith(db).Exec()
	return err
}

func (u *PgUsersDAO) IsEmailInUse(email string) (bool, error) {
	err := u.psql.Select("email").
		From("users").
		Where("email = ?", email).
		RunWith(db).QueryRow().Scan(&email)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		log.Fatal("PgUsersDAO.isEmailInUse: ", err)
		return false, err
	} else {
		return true, nil
	}
}

func (u *PgUsersDAO) GetAuthenticationInformation(email string) (int, string, string, error) {
	var id int
	var salt string
	var storedHash string

	err := u.psql.Select("id", "salt", "hash").From("users").Where("email = ?", email).
					RunWith(db).QueryRow().Scan(&id, &salt, &storedHash)
	if err != nil {
		return 0, "", "", err
	}

	return id, salt, storedHash, nil
}
