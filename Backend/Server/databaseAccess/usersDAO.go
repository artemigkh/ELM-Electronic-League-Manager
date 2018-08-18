package databaseAccess

import (
	"database/sql"
)

type PgUsersDAO struct{}

type UserInformation struct {
	Email string `json:"email"`
}

func (d *PgUsersDAO) CreateUser(email, salt, hash string) error {
	_, err := psql.Insert("users").Columns("email", "salt", "hash").
		Values(email, salt, hash).RunWith(db).Exec()
	return err
}

func (d *PgUsersDAO) IsEmailInUse(email string) (bool, error) {
	err := psql.Select("email").
		From("users").
		Where("email = ?", email).
		RunWith(db).QueryRow().Scan(&email)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (d *PgUsersDAO) GetAuthenticationInformation(email string) (int, string, string, error) {
	var id int
	var salt string
	var storedHash string

	err := psql.Select("id", "salt", "hash").From("users").Where("email = ?", email).
		RunWith(db).QueryRow().Scan(&id, &salt, &storedHash)
	if err != nil {
		return 0, "", "", err
	}

	return id, salt, storedHash, nil
}

func (d *PgUsersDAO) GetUserProfile(userId int) (*UserInformation, error) {
	var profile UserInformation

	err := psql.Select("email").From("users").Where("id = ?", userId).
		RunWith(db).QueryRow().Scan(&profile.Email)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}