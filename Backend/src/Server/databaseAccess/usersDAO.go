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

func (d *PgUsersDAO) GetPermissions(leagueId, userId int) (*UserPermissionsDTO, error) {
	var userPermissions UserPermissionsDTO

	leaguePermissions, err := getLeaguePermissions(leagueId, userId)
	if err != nil {
		return nil, err
	}

	var teamPermissions TeamPermissionsDTOArray
	if err := ScanRows(psql.Select(
		"administrator",
		"information",
		"players",
		"report_results",
	).
		From("team_permissions").
		Where("user_id = ?", userId), &teamPermissions); err != nil {
		return nil, err
	}

	userPermissions.LeaguePermissions = leaguePermissions
	userPermissions.TeamPermissions = teamPermissions.rows
	return &userPermissions, nil
}

func (d *PgUsersDAO) GetUserProfile(userId int) (*UserDTO, error) {
	return GetScannedUserDTO(psql.Select(
		"id, email").
		From("user_").
		Where("user_id = ?", userId).
		RunWith(db).QueryRow())
}
