package databaseAccess

import "database/sql"

type PgLeaguesDAO struct{}

// Modify League
func (d *PgLeaguesDAO) CreateLeague(userId int, leagueInfo LeagueCore) (int, error) {
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

func (d *PgLeaguesDAO) UpdateLeague(leagueId int, leagueInfo LeagueCore) error {
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

func (d *PgLeaguesDAO) JoinLeague(leagueId, userId int) error {
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

func (d *PgLeaguesDAO) SetLeaguePermissions(leagueId, userId int, permissions LeaguePermissionsCore) error {
	_, err := psql.Update("league_permissions").
		Set("administrator", permissions.Administrator).
		Set("create_teams", permissions.CreateTeams).
		Set("edit_teams", permissions.EditTeams).
		Set("edit_games", permissions.EditGames).
		Where("league_id = ? AND user_id = ?", leagueId, userId).
		RunWith(db).Exec()

	return err
}

func getLeaguePermissions(leagueId, userId int) (*LeaguePermissionsCore, error) {
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

func (d *PgLeaguesDAO) GetLeaguePermissions(leagueId, userId int) (*LeaguePermissionsCore, error) {
	return getLeaguePermissions(leagueId, userId)
}

func (d *PgLeaguesDAO) GetTeamManagerInformation(leagueId int) ([]*TeamWithManagers, error) {
	return nil, nil
	//var teamManagerInformation TeamManagerDTOArray
	//if err := ScanRows(psql.Select(
	//	"user_id",
	//	"email",
	//	"team_id",
	//	"name",
	//	"tag",
	//	"description",
	//	"icon_small",
	//	"administrator",
	//	"information",
	//	"players",
	//	"report_results",
	//).
	//	From("team_permissions").
	//	Join("user_ ON team_permissions.user_id = user.id").
	//	Join("team ON team_permissions.team_id = team.id").
	//	Where("league_id = ?", leagueId), &teamManagerInformation); err != nil {
	//	return nil, err
	//}
	//
	//return teamManagerInformation.rows, nil
}

func (d *PgLeaguesDAO) IsLeagueViewable(leagueId, userId int) (bool, error) {
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
func (d *PgLeaguesDAO) CanJoinLeague(leagueId, userId int) (bool, error) {
	var canJoin = false
	err := psql.Select("public_join").
		From("league").
		Where("league_id = ?", leagueId).
		RunWith(db).QueryRow().Scan(&canJoin)

	return canJoin, err
}

// Get Information About Leagues

func (d *PgLeaguesDAO) GetLeagueInformation(leagueId int) (*League, error) {
	row := getLeagueSelector().
		Where("league_id = ?", leagueId).
		RunWith(db).QueryRow()

	return GetScannedLeague(row)
}

func (d *PgLeaguesDAO) IsNameInUse(leagueId int, name string) (bool, error) {
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

func (d *PgLeaguesDAO) GetPublicLeagueList() ([]*League, error) {
	var leagueSummary LeagueArray
	if err := ScanRows(getLeagueSelector().
		Where("public_view = true"), &leagueSummary); err != nil {
		return nil, err
	}

	return leagueSummary.rows, nil
}

// Get Information About Entities in a League

func (d *PgLeaguesDAO) GetTeamSummary(leagueId int) ([]*TeamWithPlayers, error) {
	return nil, nil
	//var teamSummary TeamDTOArray
	//if err := ScanRows(psql.Select(
	//	"team_id",
	//	"name",
	//	"tag",
	//	"description",
	//	"wins",
	//	"losses",
	//	"icon_small",
	//	"icon_large",
	//).
	//	From("team").
	//	Where("league_id = ?", leagueId).
	//	OrderBy("wins DESC, losses ASC"), &teamSummary); err != nil {
	//	return nil, err
	//}
	//
	//return teamSummary.rows, nil
}

// Markdown
func (d *PgLeaguesDAO) GetMarkdownFile(leagueId int) (string, error) {
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

func (d *PgLeaguesDAO) SetMarkdownFile(leagueId int, fileName string) error {
	_, err := psql.Update("league").
		Set("markdown_path", fileName).
		Where("league_id = ?", leagueId).
		RunWith(db).Exec()

	return err
}

// Availabilities
func (d *PgLeaguesDAO) AddAvailability(leagueId int, availability AvailabilityCore) (int, error) {
	return 0, nil
}
func (d *PgLeaguesDAO) GetAvailabilities(leagueId int) ([]*Availability, error) {
	return nil, nil
}
func (d *PgLeaguesDAO) DeleteAvailability(leagueId, availabilityId int) error {
	return nil
}

func (d *PgLeaguesDAO) AddWeeklyAvailability(leagueId int, availability WeeklyAvailabilityCore) (int, error) {
	return 0, nil
}
func (d *PgLeaguesDAO) GetWeeklyAvailabilities(leagueId int) ([]*WeeklyAvailability, error) {
	return nil, nil
}
func (d *PgLeaguesDAO) EditWeeklyAvailability(leagueId, availabilityId int,
	availability WeeklyAvailabilityCore) (int, error) {
	return 0, nil
}
func (d *PgLeaguesDAO) DeleteWeeklyAvailability(leagueId, availabilityId int) error {
	return nil
}

//func (d *PgLeaguesDAO) AddRecurringAvailability(leagueId int, availability SchedulingAvailabilityDTO) (int, error) {
//	var availabilityId = -1
//	err := psql.Insert("league_recurring_availability").
//		Columns(
//			"league_id",
//			"weekday",
//			"timezone",
//			"hour",
//			"minute",
//			"duration",
//			"constrained",
//			"start_time",
//			"end_time",
//		).
//		Values(
//			leagueId,
//			availability.Weekday,
//			availability.Timezone,
//			availability.Hour,
//			availability.Minute,
//			availability.Duration,
//			availability.Constrained,
//			availability.Start,
//			availability.End,
//		).
//		Suffix("RETURNING \"id\"").
//		RunWith(db).QueryRow().Scan(&availabilityId)
//
//	return availabilityId, err
//}
//
//func (d *PgLeaguesDAO) EditRecurringAvailability(availability SchedulingAvailabilityDTO) error {
//	_, err := psql.Update("league_recurring_availability").
//		Set("weekday", availability.Weekday).
//		Set("timezone", availability.Timezone).
//		Set("hour", availability.Hour).
//		Set("minute", availability.Minute).
//		Set("duration", availability.Duration).
//		Set("constrained", availability.Constrained).
//		Set("start_time", availability.Start).
//		Set("end_time", availability.End).
//		Where("recurring_availability_id = ?", availability.Id).
//		RunWith(db).Exec()
//
//	return err
//}
//
//func (d *PgLeaguesDAO) RemoveRecurringAvailabilities(availabilityId int) error {
//	_, err := psql.Delete("league_recurring_availability").
//		Where("recurring_availability_id = ?", availabilityId).
//		RunWith(db).Exec()
//	return err
//}
//
//func (d *PgLeaguesDAO) GetSchedulingAvailability(availabilityId int) (*SchedulingAvailabilityDTO, error) {
//	return GetScannedSchedulingAvailabilityDTO(
//		psql.Select(
//			"recurring_availability_id",
//			"weekday",
//			"timezone",
//			"hour",
//			"minute",
//			"duration",
//			"constrained",
//			"start_time",
//			"end_time",
//		).
//			From("league_recurring_availability").
//			Where("recurring_availability_id = ?", availabilityId).
//			RunWith(db).QueryRow())
//}
//
//func (d *PgLeaguesDAO) GetSchedulingAvailabilities(leagueId int) ([]*SchedulingAvailabilityDTO, error) {
//	var availabilities SchedulingAvailabilityArray
//	if err := ScanRows(psql.Select(
//		"recurring_availability_id",
//		"weekday",
//		"timezone",
//		"hour",
//		"minute",
//		"duration",
//		"constrained",
//		"start_time",
//		"end_time",
//	).
//		From("league_recurring_availability").
//		Where("league_id = ?", leagueId), &availabilities); err != nil {
//		return nil, err
//	}
//
//	return availabilities.rows, nil
//}
