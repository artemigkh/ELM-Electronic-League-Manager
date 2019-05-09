package databaseAccess

import (
	"database/sql"
)

type PgUsersDAO struct{}

type UserInformation struct {
	Email string `json:"email"`
}

type TeamPermissionsInformation struct {
	Id            int  `json:"id"`
	Administrator bool `json:"administrator"`
	Information   bool `json:"information"`
	Players       bool `json:"players"`
	ReportResults bool `json:"reportResults"`
}

type UserPermissions struct {
	LeaguePermissions LeaguePermissions            `json:"leaguePermissions"`
	TeamPermissions   []TeamPermissionsInformation `json:"teamPermissions"`
}

func (d *PgUsersDAO) CreateUser(email, salt, hash string) error {
	_, err := psql.Insert("user_").Columns("email", "salt", "hash").
		Values(email, salt, hash).RunWith(db).Exec()
	return err
}

func (d *PgUsersDAO) IsEmailInUse(email string) (bool, error) {
	err := psql.Select("email").
		From("user_").
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

	err := psql.Select("id", "salt", "hash").From("user_").Where("email = ?", email).
		RunWith(db).QueryRow().Scan(&id, &salt, &storedHash)
	if err != nil {
		return 0, "", "", err
	}

	return id, salt, storedHash, nil
}

func (d *PgUsersDAO) GetUserProfile(userId int) (*UserInformation, error) {
	var profile UserInformation

	err := psql.Select("email").From("user_").Where("id = ?", userId).
		RunWith(db).QueryRow().Scan(&profile.Email)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (d *PgUsersDAO) GetPermissions(leagueId, userId int) (*UserPermissions, error) {
	var userPermissions UserPermissions
	userPermissions.TeamPermissions = make([]TeamPermissionsInformation, 0)
	// get users league permissions
	lp, err := getLeaguePermissions(leagueId, userId)
	if err != nil {
		return nil, err
	}
	userPermissions.LeaguePermissions = *lp

	// get permissions from teams in league this user has entries for
	rows, err := db.Query(`
SELECT teamId, administrator, information, players, reportResults FROM team_permissions
WHERE userId = $1
AND teamId IN (SELECT id FROM team WHERE leagueId = $2)
	`, userId, leagueId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tpi TeamPermissionsInformation

	for rows.Next() {
		err := rows.Scan(&tpi.Id, &tpi.Administrator, &tpi.Information, &tpi.Players, &tpi.ReportResults)
		if err != nil {
			return nil, err
		}
		userPermissions.TeamPermissions = append(userPermissions.TeamPermissions, tpi)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &userPermissions, nil
}
