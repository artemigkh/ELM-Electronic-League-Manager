package databaseAccess

import (
	"Server/dataModel"
	"github.com/Masterminds/squirrel"
)

func GetScannedUserAuthenticationDTO(rows squirrel.RowScanner) (*dataModel.UserAuthenticationDTO, error) {
	var authenticationInformation dataModel.UserAuthenticationDTO
	if err := rows.Scan(
		&authenticationInformation.UserId,
		&authenticationInformation.Salt,
		&authenticationInformation.Hash,
	); err != nil {
		return nil, err
	} else {
		return &authenticationInformation, nil
	}
}

// User
func getUserSelector() squirrel.SelectBuilder {
	return psql.Select(
		"user_.user_id",
		"user_.email",
	).
		From("user_")
}

func GetScannedUser(rows squirrel.RowScanner) (*dataModel.User, error) {
	var user dataModel.User
	if err := rows.Scan(
		&user.UserId,
		&user.Email,
	); err != nil {
		println("error here")
		println(err.Error())
		return nil, err
	} else {
		return &user, nil
	}
}
