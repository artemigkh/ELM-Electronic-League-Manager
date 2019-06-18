package databaseAccess

type PgUsersDAO struct{}

func (d *PgUsersDAO) CreateUser(email, salt, hash string) error {
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

func (d *PgUsersDAO) IsEmailInUse(email string) (bool, error) {
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

func (d *PgUsersDAO) GetAuthenticationInformation(email string) (*UserAuthenticationDTO, error) {
	return GetScannedUserAuthenticationDTO(psql.Select(
		"user_id",
		"salt",
		"hash").
		From("user_").
		Where("email = ?", email).
		RunWith(db).QueryRow())
}

func getLeagueAndTeamPermissions(leagueId, teamId, userId int) (*LeaguePermissionsCore, *TeamPermissionsCore, error) {
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
//func (d *PgUsersDAO) GetPermissions(leagueId, userId int) (*UserPermissionsDTO, error) {
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

func (d *PgUsersDAO) GetUserProfile(userId int) (*User, error) {
	return GetScannedUser(getUserSelector().Where("user_id = ?", userId).RunWith(db).QueryRow())
}

func (d *PgUsersDAO) GetUserWithPermissions(leagueId, userId int) (*UserWithPermissions, error) {
	userBase, err := GetScannedUser(getUserSelector().
		Where("user_id = ?", userId).RunWith(db).QueryRow())
	if err != nil {
		return nil, err
	}

	user := &UserWithPermissions{
		UserId: userBase.UserId,
		Email:  userBase.Email,
	}

	leaguePermissions, err := getLeaguePermissions(leagueId, userId)
	if err != nil {
		return nil, err
	}
	user.LeaguePermissions = leaguePermissions

	var teamPermissions TeamPermissionsArray
	if err := ScanRows(getTeamPermissionsSelector().
		Where("team.league_id = ? AND user_.user_id = ", leagueId, userId), &teamPermissions); err != nil {
		return nil, err
	}
	user.TeamPermissions = teamPermissions.rows

	return user, nil
}
