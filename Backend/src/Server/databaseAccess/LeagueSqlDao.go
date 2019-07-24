package databaseAccess

import (
	"Server/dataModel"
	"database/sql"
)

type LeagueSqlDao struct{}

// Modify League
func (d *LeagueSqlDao) CreateLeague(userId int, leagueInfo dataModel.LeagueCore) (int, error) {
	var leagueId = -1
	err := db.QueryRow("SELECT create_league($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)",
		leagueInfo.Name,
		leagueInfo.Description,
		leagueInfo.PublicView,
		leagueInfo.PublicJoin,
		leagueInfo.SignupStart,
		leagueInfo.SignupEnd,
		leagueInfo.LeagueStart,
		leagueInfo.LeagueEnd,
		leagueInfo.Game,
		userId,
	).Scan(&leagueId)

	return leagueId, err
}

func (d *LeagueSqlDao) UpdateLeague(leagueId int, leagueInfo dataModel.LeagueCore) error {
	_, err := psql.Update("league").
		Set("name", leagueInfo.Name).
		Set("description", leagueInfo.Description).
		Set("game", leagueInfo.Game).
		Set("public_view", leagueInfo.PublicView).
		Set("public_join", leagueInfo.PublicJoin).
		Set("signup_start", leagueInfo.SignupStart).
		Set("signup_end", leagueInfo.SignupEnd).
		Set("league_start", leagueInfo.LeagueStart).
		Set("league_end", leagueInfo.LeagueEnd).
		Where("league_id = ?", leagueId).
		RunWith(db).Exec()

	return err
}

func (d *LeagueSqlDao) JoinLeague(leagueId, userId int) error {
	_, err := psql.Insert("league_permissions").
		Columns(
			"user_id",
			"league_id",
			"administrator",
			"create_teams",
			"edit_teams",
			"edit_games",
		).
		Values(
			userId,
			leagueId,
			false,
			true,
			false,
			false,
		).
		RunWith(db).Exec()

	return err
}

// Permissions

func (d *LeagueSqlDao) SetLeaguePermissions(leagueId, userId int, permissions dataModel.LeaguePermissionsCore) error {
	_, err := psql.Update("league_permissions").
		Set("administrator", permissions.Administrator).
		Set("create_teams", permissions.CreateTeams).
		Set("edit_teams", permissions.EditTeams).
		Set("edit_games", permissions.EditGames).
		Where("league_id = ? AND user_id = ?", leagueId, userId).
		RunWith(db).Exec()

	return err
}

func getLeaguePermissions(leagueId, userId int) (*dataModel.LeaguePermissionsCore, error) {
	return GetScannedLeaguePermissionsCore(psql.Select(
		"administrator",
		"create_teams",
		"edit_teams",
		"edit_games",
	).
		From("league_permissions").
		Where("league_id = ? AND user_id = ?", leagueId, userId).
		RunWith(db).QueryRow())
}

func (d *LeagueSqlDao) GetLeaguePermissions(leagueId, userId int) (*dataModel.LeaguePermissionsCore, error) {
	return getLeaguePermissions(leagueId, userId)
}

func (d *LeagueSqlDao) GetTeamManagerInformation(leagueId int) ([]*dataModel.TeamWithManagers, error) {
	rows, err := getTeamWithManagersSelector().
		Where("team.league_id = ?", leagueId).
		RunWith(db).Query()
	if err != nil {
		return nil, err
	}
	return GetScannedAllTeamWithManagers(rows)
}

func (d *LeagueSqlDao) IsLeagueViewable(leagueId, userId int) (bool, error) {
	//check if publicly viewable
	var publicView bool
	err := psql.Select("public_view").
		From("league").
		Where("league_id = ?", leagueId).
		RunWith(db).QueryRow().Scan(&publicView)
	if err != nil {
		return false, err
	}

	if publicView {
		return true, nil
	}

	//if not publicly viewable, see if user has permission to view it. This is checked by seeing if there is a
	//leaguePermissions row with that userId and leagueId, if there is they have at least the base (viewing) privileges
	var uid int
	err = psql.Select("user_id").
		From("league_permissions").
		Where("league_id = ? AND user_id = ?", leagueId, userId).
		RunWith(db).QueryRow().Scan(&uid)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

//TODO: make invite system for private leagues, check if user invited in this function
//TODO: make ordering consistent
func (d *LeagueSqlDao) CanJoinLeague(leagueId, userId int) (bool, error) {
	var canJoin = false
	err := psql.Select("public_join").
		From("league").
		Where("league_id = ?", leagueId).
		RunWith(db).QueryRow().Scan(&canJoin)

	return canJoin, err
}

// Get Information About Leagues

func (d *LeagueSqlDao) DoesLeagueExist(leagueId int) (bool, error) {
	var count int
	if err := psql.Select("count(*)").
		From("league").
		Where("league_id = ?", leagueId).
		RunWith(db).QueryRow().Scan(&count); err != nil {
		return false, err
	} else {
		return count > 0, nil
	}
}

func (d *LeagueSqlDao) GetLeagueInformation(leagueId int) (*dataModel.League, error) {
	return GetScannedLeague(getLeagueSelector().
		Where("league_id = ?", leagueId).
		RunWith(db).QueryRow())
}

func (d *LeagueSqlDao) IsNameInUse(leagueId int, name string) (bool, error) {
	var count int
	if err := psql.Select("count(name)").
		From("league").
		Where("league_id != ? AND name = ?", leagueId, name).
		RunWith(db).QueryRow().Scan(&count); err != nil {
		return false, err
	} else {
		return count > 0, nil
	}
}

func (d *LeagueSqlDao) GetPublicLeagueList() ([]*dataModel.League, error) {
	var leagueSummary LeagueArray
	if err := ScanRows(getLeagueSelector().
		Where("public_view = true"), &leagueSummary); err != nil {
		return nil, err
	}

	return leagueSummary.rows, nil
}

// Markdown
func (d *LeagueSqlDao) GetMarkdownFile(leagueId int) (string, error) {
	var markdownFile = ""
	err := psql.Select("markdown_path").
		From("league").
		Where("league_id = ?", leagueId).
		RunWith(db).QueryRow().Scan(&markdownFile)
	if err == sql.ErrNoRows {
		return "", nil
	} else {
		return markdownFile, err
	}
}

func (d *LeagueSqlDao) SetMarkdownFile(leagueId int, fileName string) error {
	_, err := psql.Update("league").
		Set("markdown_path", fileName).
		Where("league_id = ?", leagueId).
		RunWith(db).Exec()

	return err
}

// Availabilities
func (d *LeagueSqlDao) AddAvailability(leagueId int, availability dataModel.AvailabilityCore) (int, error) {
	var availabilityId = -1
	err := psql.Insert("availability").
		Columns(
			"league_id",
			"start_time",
			"end_time",
			"is_recurring_weekly",
		).
		Values(
			leagueId,
			availability.StartTime,
			availability.EndTime,
			false,
		).
		Suffix("RETURNING \"availability_id\"").
		RunWith(db).QueryRow().Scan(&availabilityId)

	return availabilityId, err
}

func (d *LeagueSqlDao) GetAvailabilities(leagueId int) ([]*dataModel.Availability, error) {
	availabilities := AvailabilityArray{rows: make([]*dataModel.Availability, 0)}
	if err := ScanRows(getAvailabilitySelector(leagueId), &availabilities); err != nil {
		return nil, err
	}

	return availabilities.rows, nil
}
func (d *LeagueSqlDao) DeleteAvailability(availabilityId int) error {
	_, err := psql.Delete("availability").
		Where("availability_id = ?", availabilityId).
		RunWith(db).Exec()
	return err
}

func (d *LeagueSqlDao) AddWeeklyAvailability(leagueId int, availability dataModel.WeeklyAvailabilityCore) (int, error) {
	var availabilityId int
	err := psql.Insert("availability").
		Columns(
			"league_id",
			"start_time",
			"end_time",
			"is_recurring_weekly",
		).
		Values(
			leagueId,
			availability.StartTime,
			availability.EndTime,
			true,
		).
		Suffix("RETURNING \"availability_id\"").
		RunWith(db).QueryRow().Scan(&availabilityId)
	if err != nil {
		return -1, err
	}

	_, err = psql.Insert("weekly_recurrence").
		Columns(
			"availability_id",
			"weekday",
			"timezone",
			"hour",
			"minute",
			"duration",
		).
		Values(
			availabilityId,
			availability.Weekday,
			availability.Timezone,
			availability.Hour,
			availability.Minute,
			availability.Duration,
		).
		RunWith(db).Exec()
	if err != nil {
		return -1, err
	}

	return availabilityId, nil
}

func (d *LeagueSqlDao) GetWeeklyAvailabilities(leagueId int) ([]*dataModel.WeeklyAvailability, error) {
	availabilities := WeeklyAvailabilityArray{rows: make([]*dataModel.WeeklyAvailability, 0)}
	if err := ScanRows(getWeeklyAvailabilitySelector(leagueId), &availabilities); err != nil {
		return nil, err
	}

	return availabilities.rows, nil
}

func (d *LeagueSqlDao) EditWeeklyAvailability(availabilityId int,
	availability dataModel.WeeklyAvailabilityCore) error {
	_, err := psql.Update("availability").
		Set("start_time", availability.StartTime).
		Set("end_time", availability.EndTime).
		Where("availability_id = ?", availabilityId).
		RunWith(db).Exec()
	if err != nil {
		return err
	}

	_, err = psql.Update("weekly_recurrence").
		Set("weekday", availability.Weekday).
		Set("timezone", availability.Timezone).
		Set("hour", availability.Hour).
		Set("minute", availability.Minute).
		Set("duration", availability.Duration).
		Where("availability_id = ?", availabilityId).
		RunWith(db).Exec()
	return err
}

func (d *LeagueSqlDao) DeleteWeeklyAvailability(availabilityId int) error {
	_, err := psql.Delete("weekly_recurrence").
		Where("availability_id = ?", availabilityId).
		RunWith(db).Exec()
	if err != nil {
		return err
	}

	_, err = psql.Delete("availability").
		Where("availability_id = ?", availabilityId).
		RunWith(db).Exec()
	return err
}

func (d *LeagueSqlDao) DoesAvailabilityExistInLeague(leagueId, availabilityId int) (bool, error) {
	var count int
	if err := psql.Select("count(*)").
		From("availability").
		Where("league_id = ? AND availability_id = ?", leagueId, availabilityId).
		RunWith(db).QueryRow().Scan(&count); err != nil {
		return false, err
	} else {
		return count > 0, nil
	}
}
