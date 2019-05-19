package databaseAccess

import "github.com/Masterminds/squirrel"

type UserDTO struct {
	UserId int    `json:"userId"`
	Email  string `json:"email"`
}

func GetScannedUserDTO(rows squirrel.RowScanner) (*UserDTO, error) {
	var user UserDTO
	if err := rows.Scan(
		&user.UserId,
		&user.Email,
	); err != nil {
		return nil, err
	} else {
		return &user, nil
	}
}

type UserPermissionsDTO struct {
	UserId            int                   `json:"userId"`
	LeaguePermissions *LeaguePermissionsDTO `json:"leaguePermissions"`
	TeamPermissions   []*TeamPermissionsDTO `json:"teamPermissions"`
}

type UserAuthenticationDTO struct {
	UserId int    `json:"userId"`
	Salt   string `json:"salt"`
	Hash   string `json:"hash"`
}

func GetScannedUserAuthenticationDTO(rows squirrel.RowScanner) (*UserAuthenticationDTO, error) {
	var authenticationInformation UserAuthenticationDTO
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
