package databaseAccess

import (
	"Server/dataModel"
	"database/sql"
	"fmt"
)

type UserSqlDao struct{}

func (d *UserSqlDao) CreateUser(email, salt, hash string) error {
	_, err := psql.Insert("user_").
		Columns(
			"email",
			"salt",
			"hash",
		).
		Values(
			email,
			salt,
			hash,
		).RunWith(db).Exec()
	return err
}

func (d *UserSqlDao) IsEmailInUse(email string) (bool, error) {
	//TODO: check for equivalent emails
	var count int
	if err := psql.Select("count(email)").
		From("user_").
		Where("email = ?", email).
		RunWith(db).QueryRow().Scan(&count); err != nil {
		return false, err
	} else {
		return count > 0, nil
	}
}

func (d *UserSqlDao) GetAuthenticationInformation(email string) (*dataModel.UserAuthenticationDTO, error) {
	return GetScannedUserAuthenticationDTO(psql.Select(
		"user_id",
		"salt",
		"hash").
		From("user_").
		Where("email = ?", email).
		RunWith(db).QueryRow())
}

func getLeagueAndTeamPermissions(leagueId, teamId, userId int) (*dataModel.LeaguePermissionsCore, *dataModel.TeamPermissionsCore, error) {
	leaguePermissions, err := getLeaguePermissions(leagueId, userId)
	if err != nil {
		return nil, nil, err
	}

	teamPermissions, err := teamsDAO.GetTeamPermissions(teamId, userId)
	if err != nil {
		return nil, nil, err
	}

	return leaguePermissions, teamPermissions, nil
}

//
//func (d *UserSqlDao) GetPermissions(leagueId, userId int) (*UserPermissionsDTO, error) {
//	var userPermissions UserPermissionsDTO
//
//	leaguePermissions, err := getLeaguePermissions(leagueId, userId)
//	if err != nil {
//		return nil, err
//	}
//
//	var teamPermissions TeamPermissionsDTOArray
//	if err := ScanRows(psql.Select(
//		"administrator",
//		"information",
//		"players",
//		"report_results",
//	).
//		From("team_permissions").
//		Where("user_id = ?", userId), &teamPermissions); err != nil {
//		return nil, err
//	}
//
//	userPermissions.LeaguePermissions = leaguePermissions
//	userPermissions.TeamPermissions = teamPermissions.rows
//	return &userPermissions, nil
//}

func (d *UserSqlDao) GetUserProfile(userId int) (*dataModel.User, error) {
	return GetScannedUser(getUserSelector().Where("user_id = ?", userId).RunWith(db).QueryRow())
}

func (d *UserSqlDao) GetUserWithPermissions(leagueId, userId int) (*dataModel.UserWithPermissions, error) {
	fmt.Printf("league id: %v, userId: %v\n", leagueId, userId)
	var userBase *dataModel.User
	userBase, err := GetScannedUser(getUserSelector().
		Where("user_id = ?", userId).RunWith(db).QueryRow())
	if err == sql.ErrNoRows {
		userBase = &dataModel.User{}
	} else if err != nil {
		return nil, err
	}

	user := &dataModel.UserWithPermissions{
		UserId: userBase.UserId,
		Email:  userBase.Email,
	}

	leaguePermissions, err := getLeaguePermissions(leagueId, userId)
	if err != nil {
		return nil, err
	}
	user.LeaguePermissions = leaguePermissions

	var teamPermissions TeamPermissionsArray
	teamPermissions.rows = make([]*dataModel.TeamPermissions, 0)
	if err := ScanRows(getTeamPermissionsSelector().
		Where("team.league_id = ? AND team_permissions.user_id = ?", leagueId, userId), &teamPermissions); err != nil {
		return nil, err
	}
	user.TeamPermissions = teamPermissions.rows

	return user, nil
}
