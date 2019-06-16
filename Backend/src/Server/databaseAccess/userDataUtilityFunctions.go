package databaseAccess

import (
	"github.com/Masterminds/squirrel"
	"github.com/badoux/checkmail"
)

// UserCreationInformation
func (user *UserCreationInformation) Validate() (bool, string, error) {
	return validate(
		user.email(),
		user.uniqueness(),
		user.password())
}

func (user *UserCreationInformation) email() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		if err := checkmail.ValidateFormat(user.Email); err != nil {
			return false
		} else {
			return true
		}
	}
}

func (user *UserCreationInformation) uniqueness() ValidateFunc {
	return func(problemDest *string, errorDest *error) bool {
		valid := false
		inUse, err := usersDAO.IsEmailInUse(user.Email)
		if err != nil {
			*errorDest = err
		} else if inUse {
			*problemDest = EmailInUse
		} else {
			valid = true
		}
		return valid
	}
}

func (user *UserCreationInformation) password() ValidateFunc {
	return func(problemDest *string, _ *error) bool {
		valid := false
		if len(user.Password) > MaxPasswordLength {
			*problemDest = PasswordTooLong
		} else if len(user.Password) < MinPasswordLength {
			*problemDest = PasswordTooShort
		} else {
			valid = true
		}
		return valid
	}
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

// User
func getUserSelector() squirrel.SelectBuilder {
	return psql.Select(
		"user_.user_id",
		"user_.email",
		"league_permissions.administrator",
		"league_permissions.create_teams",
		"league_permissions.edit_teams",
		"league_permissions.edit_games",
	).
		From("user_").
		Join("league_permissions ON user_.user_id = league_permissions.league_id")
}

func GetScannedUser(rows squirrel.RowScanner) (*User, error) {
	var user User
	var leaguePermissions LeaguePermissionsCore
	if err := rows.Scan(
		&user.UserId,
		&user.Email,
		&leaguePermissions.Administrator,
		&leaguePermissions.CreateTeams,
		&leaguePermissions.EditTeams,
		&leaguePermissions.EditGames,
	); err != nil {
		return nil, err
	} else {
		user.LeaguePermissions = leaguePermissions
		return &user, nil
	}
}
