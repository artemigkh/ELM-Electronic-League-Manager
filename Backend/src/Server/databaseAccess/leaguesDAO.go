package databaseAccess

import (
	"database/sql"
)

type LeaguePermissions struct {
	Administrator bool `json:"administrator"`
	CreateTeams   bool `json:"createTeams"`
	EditTeams     bool `json:"editTeams"`
	EditGames     bool `json:"editGames"`
}

type GameSummaryInformation struct {
	Id         int  `json:"id"`
	Team1Id    int  `json:"team1Id"`
	Team2Id    int  `json:"team2Id"`
	GameTime   int  `json:"gameTime"`
	Complete   bool `json:"complete"`
	WinnerId   int  `json:"winnerId"`
	ScoreTeam1 int  `json:"scoreTeam1"`
	ScoreTeam2 int  `json:"scoreTeam2"`
}

type TeamManagerInformation struct {
	TeamId   int                  `json:"teamId"`
	TeamName string               `json:"teamName"`
	TeamTag  string               `json:"teamTag"`
	Managers []ManagerInformation `json:"managers"`
}

type ManagerInformation struct {
	UserId        int    `json:"userId"`
	UserEmail     string `json:"userEmail"`
	Administrator bool   `json:"administrator"`
	Information   bool   `json:"information"`
	Players       bool   `json:"players"`
	ReportResults bool   `json:"reportResults"`
}

type PublicLeagueInformation struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PublicJoin  bool   `json:"publicJoin"`
	SignupStart int    `json:"signupStart"`
	SignupEnd   int    `json:"signupEnd"`
	LeagueStart int    `json:"leagueStart"`
	LeagueEnd   int    `json:"leagueEnd"`
	Game        string `json:"game"`
}

type SchedulingAvailability struct {
	Id          int  `json:"id"`
	Weekday     int  `json:"weekday"`
	Timezone    int  `json:"timezone"`
	Hour        int  `json:"hour"`
	Minute      int  `json:"minute"`
	Duration    int  `json:"duration"`
	Constrained bool `json:"constrained"`
	Start       int  `json:"start"`
	End         int  `json:"end"`
}

type PgLeaguesDAO struct{}

func (d *PgLeaguesDAO) CreateLeague(userId int, name, description, game string, publicView, publicJoin bool,
	signupStart, signupEnd, leagueStart, leagueEnd int) (int, error) {

	var leagueId int
	err := psql.Insert("league").
		Columns("name", "description", "markdown_path", "game", "public_view", "public_join",
			"signup_start", "signup_end", "league_start", "league_end").
		Values(name, description, "", game, publicView, publicJoin, signupStart, signupEnd, leagueStart, leagueEnd).
		Suffix("RETURNING \"id\"").
		RunWith(db).QueryRow().Scan(&leagueId)
	if err != nil {
		return -1, err
	}

	//create permissions entry linking current user Id as the league creator
	_, err = psql.Insert("league_permissions").
		Columns("userId", "league_id",
			"administrator", "createTeams", "editTeams", "editGames").
		Values(userId, leagueId, true, true, true, true).
		RunWith(db).Exec()
	if err != nil {
		return -1, err
	}

	return leagueId, nil
}

func (d *PgLeaguesDAO) UpdateLeague(leagueInfo LeagueInformationDTO) error {
	//_, err := psql.Update("league").
	//	Set("complete", true).
	//	Set("winnerId", winnerId).
	//	Set("scoreteam1", scoreTeam1).
	//	Set("scoreteam2", scoreTeam2).
	//	Where("id = ? AND leagueId = ?", gameId, leagueId).RunWith(db).Exec()
	//if err != nil {
	//	return err
	//}
	//_, err := db.Exec(
	//	`
	//	UPDATE leagues SET name = $1, description = $2, game = $3, publicView = $4, publicJoin = $5,
	//		signupStart = $6, signupEnd = $7, leagueStart = $8, leagueEnd = $9
	//	WHERE id = $10
	//	`, name, description, game, publicView, publicJoin, signupStart, signupEnd, leagueStart, leagueEnd, leagueId)
	//return err
	return nil
}

func (d *PgLeaguesDAO) JoinLeague(leagueId, userId int) error {
	_, err := psql.Insert("league_permissions").
		Columns("userId", "league_id",
			"administrator", "createTeams", "editTeams", "editGames").
		Values(userId, leagueId, false, true, false, false).
		RunWith(db).Exec()
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (d *PgLeaguesDAO) IsNameInUse(leagueId int, name string) (bool, error) {
	var leagueIdOfMatch int
	err := psql.Select("name, id").
		From("league").
		Where("name = ?", name).
		RunWith(db).QueryRow().Scan(&name, &leagueIdOfMatch)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	} else if leagueIdOfMatch != leagueId {
		return true, nil
	} else {
		return false, nil
	}
}

func (d *PgLeaguesDAO) IsLeagueViewable(leagueId, userId int) (bool, error) {
	//check if publicly viewable
	var publicView bool
	err := psql.Select("public_view").
		From("league").
		Where("id = ?", leagueId).
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
	err = psql.Select("userId").
		From("league_permissions").
		Where("userId = ? AND leagueId = ?", userId, leagueId).
		RunWith(db).QueryRow().Scan(&uid)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (d *PgLeaguesDAO) GetLeagueInformation(leagueId int) (*LeagueInformationDTO, error) {
	var leagueInfo LeagueInformationDTO
	err := psql.Select("id", "name", "description", "game", "public_view", "public_join",
		"signup_start", "signup_end", "league_start", "league_end").
		From("league").
		Where("id = ?", leagueId).
		RunWith(db).QueryRow().Scan(&leagueInfo.Id, &leagueInfo.Name, &leagueInfo.Description, &leagueInfo.Game,
		&leagueInfo.PublicView, &leagueInfo.PublicJoin, &leagueInfo.SignupStart, &leagueInfo.SignupEnd,
		&leagueInfo.LeagueStart, &leagueInfo.LeagueEnd)
	if err != nil {
		return nil, err
	}

	return &leagueInfo, err
}

//ScanRows(statement *squirrel.SelectBuilder, out RowArr)
func (d *PgLeaguesDAO) GetTeamSummary(leagueId int) ([]*TeamSummaryInformationDTO, error) {
	var teamSummary TeamSummaryInformationArray
	if err := ScanRows(psql.Select("id", "name", "tag", "wins", "losses", "iconSmall", "iconLarge").From("teams").
		Where("leagueId = ?", leagueId).
		OrderBy("wins DESC, losses ASC"), &teamSummary); err != nil {
		return nil, err
	}

	return teamSummary.rows, nil
}

func (d *PgLeaguesDAO) GetGameSummary(leagueId int) ([]GameSummaryInformation, error) {
	rows, err := psql.Select("id", "team1Id", "team2Id", "gametime", "complete", "winnerId",
		"scoreteam1", "scoreteam2").From("game").
		Where("leagueId = ?", leagueId).
		OrderBy("gametime DESC").
		RunWith(db).Query()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []GameSummaryInformation
	var game GameSummaryInformation

	for rows.Next() {
		err := rows.Scan(&game.Id, &game.Team1Id, &game.Team2Id, &game.GameTime, &game.Complete, &game.WinnerId,
			&game.ScoreTeam1, &game.ScoreTeam2)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return games, nil
}

//TODO: make invite system for private leagues, check if user invited in this function
//TODO: make ordering consistent
func (d *PgLeaguesDAO) CanJoinLeague(leagueId, userId int) (bool, error) {
	var canJoin bool
	err := psql.Select("public_join").
		From("league").
		Where("id = ?", leagueId).
		RunWith(db).QueryRow().Scan(&canJoin)
	if err != nil {
		return false, err
	}

	return canJoin, nil
}

func (d *PgLeaguesDAO) GetTeamManagerInformation(leagueId int) ([]TeamManagerInformation, error) {
	rows, err := psql.Select("userId", "teamId", "email", "name", "tag",
		"administrator", "information", "players", "reportResults").
		From("team_permissions").
		Join("users ON team_permissions.userId = users.id").
		Join("team ON team_permissions.teamId = team.id").
		Where("leagueId = ?", leagueId).
		RunWith(db).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//make a map of team IDs to team information objects
	var teams = make(map[int]*TeamManagerInformation)
	var (
		userId        int
		teamId        int
		email         string
		name          string
		tag           string
		administrator bool
		information   bool
		players       bool
		reportResults bool
	)

	//iterate through the rows returned from database
	for rows.Next() {
		//scan the variables from the sql row into local variables
		err := rows.Scan(&userId, &teamId, &email, &name,
			&tag, &administrator, &information, &players, &reportResults)
		if err != nil {
			return nil, err
		}

		//if the map does not have an entry for this team Id, create it
		if _, hasEntry := teams[teamId]; !hasEntry {
			teams[teamId] = &TeamManagerInformation{
				TeamId:   teamId,
				TeamName: name,
				TeamTag:  tag,
				Managers: make([]ManagerInformation, 0),
			}
		}

		//add the manager to this team representation
		teams[teamId].Managers = append(teams[teamId].Managers, ManagerInformation{
			UserId:        userId,
			UserEmail:     email,
			Administrator: administrator,
			Information:   information,
			Players:       players,
			ReportResults: reportResults,
		})
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	//create an array of the values of the teams map and return it
	teamsReps := make([]TeamManagerInformation, 0, len(teams))
	for _, team := range teams {
		teamsReps = append(teamsReps, *team)
	}

	return teamsReps, nil
}

func (d *PgLeaguesDAO) GetPublicLeagueList() ([]PublicLeagueInformation, error) {
	rows, err := psql.Select("id", "name", "description", "public_join", "signup_start",
		"signup_end", "league_start", "league_end", "game").
		From("league").
		Where("publicView = true").
		RunWith(db).Query()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leagues []PublicLeagueInformation
	var league PublicLeagueInformation

	for rows.Next() {
		err := rows.Scan(&league.Id, &league.Name, &league.Description, &league.PublicJoin, &league.SignupStart,
			&league.SignupEnd, &league.LeagueStart, &league.LeagueEnd, &league.Game)
		if err != nil {
			return nil, err
		}
		leagues = append(leagues, league)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return leagues, nil
}

func getLeaguePermissions(leagueId, userId int) (*LeaguePermissions, error) {
	var lp LeaguePermissions
	err := psql.Select("administrator", "createTeams", "editTeams", "editGames").
		From("league_permissions").
		Where("userId = ? AND leagueId = ?", userId, leagueId).
		RunWith(db).QueryRow().Scan(&lp.Administrator, &lp.CreateTeams, &lp.EditTeams, &lp.EditGames)
	if err == sql.ErrNoRows {
		return &LeaguePermissions{
			Administrator: false,
			CreateTeams:   false,
			EditTeams:     false,
			EditGames:     false,
		}, nil
	} else if err != nil {
		return nil, err
	}
	return &lp, nil
}

func (d *PgLeaguesDAO) GetLeaguePermissions(leagueId, userId int) (*LeaguePermissions, error) {
	return getLeaguePermissions(leagueId, userId)
}

func (d *PgLeaguesDAO) SetLeaguePermissions(leagueId, userId int,
	administrator, createTeams, editTeams, editGames bool) error {

	_, err := db.Exec(
		`
		UPDATE league_permissions SET administrator = $1, createTeams = $2, editTeams = $3, editGames = $4
		WHERE leagueId = $5 AND userId = $6
		`, administrator, createTeams, editTeams, editGames, leagueId, userId)
	return err
}

func (d *PgLeaguesDAO) SetMarkdownFile(leagueId int, fileName string) error {
	_, err := db.Exec(
		`
		UPDATE leagues SET markdown_path = $1
		WHERE id = $2
		`, fileName, leagueId)
	return err
}

func (d *PgLeaguesDAO) GetMarkdownFile(leagueId int) (string, error) {
	var markdownFile string
	err := psql.Select("markdown_path").
		From("league").
		Where("id = ?", leagueId).
		RunWith(db).QueryRow().Scan(&markdownFile)
	if err != nil {
		return "", err
	}
	return markdownFile, nil
}

func (d *PgLeaguesDAO) AddRecurringAvailability(leagueId int, weekday int, timezone int,
	hour, minute, duration int, constrained bool, start, end int) (int, error) {

	var availabilityId int
	err := psql.Insert("league_recurring_availability").
		Columns("league_id", "weekday", "timezone", "hour", "minute", "duration",
			"constrained", "start_time", "end_time").
		Values(leagueId, weekday, timezone, hour, minute, duration, constrained, start, end).
		Suffix("RETURNING \"id\"").
		RunWith(db).QueryRow().Scan(&availabilityId)
	if err != nil {
		return -1, err
	}

	return availabilityId, nil
}

func (d *PgLeaguesDAO) RemoveRecurringAvailabilities(leagueId, availabilityId int) error {
	_, err := psql.Delete("league_recurring_availability").
		Where("id = ? AND leagueId = ?", availabilityId, leagueId).
		RunWith(db).Exec()
	return err
}

func (d *PgLeaguesDAO) GetSchedulingAvailability(leagueId, availabilityId int) (*SchedulingAvailability, error) {
	var availability SchedulingAvailability
	err := psql.Select("id", "weekday", "timezone", "hour", "minute", "duration",
		"constrained", "start_time", "end_time").
		From("league_recurring_availability").
		Where("leagueId = ? AND id = ?", leagueId, availabilityId).
		RunWith(db).QueryRow().Scan(&availability.Id, &availability.Weekday, &availability.Timezone,
		&availability.Hour, &availability.Minute, &availability.Duration, &availability.Constrained,
		&availability.Start, &availability.End)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &availability, nil
}

func (d *PgLeaguesDAO) GetSchedulingAvailabilities(leagueId int) ([]SchedulingAvailability, error) {
	rows, err := psql.Select("id", "weekday", "timezone", "hour", "minute", "duration",
		"constrained", "start_time", "end_time").
		From("league_recurring_availability").
		Where("leagueId = ?", leagueId).
		RunWith(db).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var availabilities []SchedulingAvailability
	var availability SchedulingAvailability

	for rows.Next() {
		err := rows.Scan(&availability.Id, &availability.Weekday, &availability.Timezone,
			&availability.Hour, &availability.Minute, &availability.Duration, &availability.Constrained,
			&availability.Start, &availability.End)
		if err != nil {
			return nil, err
		}
		availabilities = append(availabilities, availability)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return availabilities, nil
}

func (d *PgLeaguesDAO) EditRecurringAvailability(leagueId, availabilityId int, weekday int, timezone int,
	hour, minute, duration int, constrained bool, start, end int) error {
	_, err := db.Exec(
		`
		UPDATE league_recurring_availability SET weekday = $1, timezone = $2, hour = $3, minute = $4,
		duration = $5, constrained = $6, start_time = $7, end_time = $8
		WHERE leagueId = $9 AND id = $10
		`, weekday, timezone, hour, minute, duration, constrained, start, end, leagueId, availabilityId)
	return err
}
