package databaseAccess

import (
	"database/sql"
)

type LeagueInformation struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Game        string `json:"game"`
	PublicView  bool   `json:"publicView"`
	PublicJoin  bool   `json:"publicJoin"`
	SignupStart int    `json:"signupStart"`
	SignupEnd   int    `json:"signupEnd"`
	LeagueStart int    `json:"leagueStart"`
	LeagueEnd   int    `json:"leagueEnd"`
}

type LeaguePermissions struct {
	Administrator bool `json:"administrator"`
	CreateTeams   bool `json:"createTeams"`
	EditTeams     bool `json:"editTeams"`
	EditGames     bool `json:"editGames"`
}

type TeamSummaryInformation struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Tag       string `json:"tag"`
	Wins      int    `json:"wins"`
	Losses    int    `json:"losses"`
	IconSmall string `json:"iconSmall"`
	IconLarge string `json:"iconLarge"`
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
	err := psql.Insert("leagues").
		Columns("name", "description", "markdownLoc", "game", "publicView", "publicJoin",
			"signupStart", "signupEnd", "leagueStart", "leagueEnd").
		Values(name, description, "", game, publicView, publicJoin, signupStart, signupEnd, leagueStart, leagueEnd).
		Suffix("RETURNING \"id\"").
		RunWith(db).QueryRow().Scan(&leagueId)
	if err != nil {
		return -1, err
	}

	//create permissions entry linking current user Id as the league creator
	_, err = psql.Insert("leaguePermissions").
		Columns("userId", "leagueId",
			"administrator", "createTeams", "editTeams", "editGames").
		Values(userId, leagueId, true, true, true, true).
		RunWith(db).Exec()
	if err != nil {
		return -1, err
	}

	return leagueId, nil
}

func (d *PgLeaguesDAO) UpdateLeague(leagueId int, name, description, game string, publicView, publicJoin bool,
	signupStart, signupEnd, leagueStart, leagueEnd int) error {
	_, err := db.Exec(
		`
		UPDATE leagues SET name = $1, description = $2, game = $3, publicView = $4, publicJoin = $5,
			signupStart = $6, signupEnd = $7, leagueStart = $8, leagueEnd = $9
		WHERE id = $10
		`, name, description, game, publicView, publicJoin, signupStart, signupEnd, leagueStart, leagueEnd, leagueId)
	return err
}

func (d *PgLeaguesDAO) JoinLeague(leagueId, userId int) error {
	_, err := psql.Insert("leaguePermissions").
		Columns("userId", "leagueId",
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
		From("leagues").
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
	err := psql.Select("publicview").
		From("leagues").
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
		From("leaguePermissions").
		Where("userId = ? AND leagueId = ?", userId, leagueId).
		RunWith(db).QueryRow().Scan(&uid)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (d *PgLeaguesDAO) GetLeagueInformation(leagueId int) (*LeagueInformation, error) {
	var leagueInfo LeagueInformation
	err := psql.Select("id", "name", "description", "game", "publicView", "publicJoin",
		"signupStart", "signupEnd", "leagueStart", "leagueEnd").
		From("leagues").
		Where("id = ?", leagueId).
		RunWith(db).QueryRow().Scan(&leagueInfo.Id, &leagueInfo.Name, &leagueInfo.Description, &leagueInfo.Game,
		&leagueInfo.PublicView, &leagueInfo.PublicJoin, &leagueInfo.SignupStart, &leagueInfo.SignupEnd,
		&leagueInfo.LeagueStart, &leagueInfo.LeagueEnd)
	if err != nil {
		return nil, err
	}

	return &leagueInfo, nil
}

func (d *PgLeaguesDAO) GetTeamSummary(leagueId int) ([]TeamSummaryInformation, error) {
	rows, err := psql.Select("id", "name", "tag", "wins", "losses", "iconSmall", "iconLarge").From("teams").
		Where("leagueId = ?", leagueId).
		OrderBy("wins DESC, losses ASC").
		RunWith(db).Query()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []TeamSummaryInformation
	var team TeamSummaryInformation

	for rows.Next() {
		err := rows.Scan(&team.Id, &team.Name, &team.Tag, &team.Wins, &team.Losses, &team.IconSmall, &team.IconLarge)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return teams, nil
}

func (d *PgLeaguesDAO) GetGameSummary(leagueId int) ([]GameSummaryInformation, error) {
	rows, err := psql.Select("id", "team1Id", "team2Id", "gametime", "complete", "winnerId",
		"scoreteam1", "scoreteam2").From("games").
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
	err := psql.Select("publicJoin").
		From("leagues").
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
		From("teamPermissions").
		Join("users ON teamPermissions.userId = users.id").
		Join("teams ON teamPermissions.teamId = teams.id").
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
	rows, err := psql.Select("id", "name", "description", "publicJoin", "signupStart",
		"signupEnd", "leagueStart", "leagueEnd", "game").
		From("leagues").
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
		From("leaguePermissions").
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
		UPDATE leaguePermissions SET administrator = $1, createTeams = $2, editTeams = $3, editGames = $4
		WHERE leagueId = $5 AND userId = $6
		`, administrator, createTeams, editTeams, editGames, leagueId, userId)
	return err
}

func (d *PgLeaguesDAO) SetMarkdownFile(leagueId int, fileName string) error {
	_, err := db.Exec(
		`
		UPDATE leagues SET markdownLoc = $1
		WHERE id = $2
		`, fileName, leagueId)
	return err
}

func (d *PgLeaguesDAO) GetMarkdownFile(leagueId int) (string, error) {
	var markdownFile string
	err := psql.Select("markdownLoc").
		From("leagues").
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
	err := psql.Insert("leagueRecurringAvailabilities").
		Columns("leagueId", "weekday", "timezone", "hour", "minute", "duration",
			"constrained", "startUnixTime", "endUnixTime").
		Values(leagueId, weekday, timezone, hour, minute, duration, constrained, start, end).
		Suffix("RETURNING \"id\"").
		RunWith(db).QueryRow().Scan(&availabilityId)
	if err != nil {
		return -1, err
	}

	return availabilityId, nil
}

func (d *PgLeaguesDAO) RemoveRecurringAvailabilities(leagueId, availabilityId int) error {
	_, err := psql.Delete("leagueRecurringAvailabilities").
		Where("id = ? AND leagueId = ?", availabilityId, leagueId).
		RunWith(db).Exec()
	return err
}

func (d *PgLeaguesDAO) GetSchedulingAvailability(leagueId, availabilityId int) (*SchedulingAvailability, error) {
	var availability SchedulingAvailability
	err := psql.Select("id", "weekday", "timezone", "hour", "minute", "duration",
		"constrained", "startUnixTime", "endUnixTime").
		From("leagueRecurringAvailabilities").
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
		"constrained", "startUnixTime", "endUnixTime").
		From("leagueRecurringAvailabilities").
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
		UPDATE leagueRecurringAvailabilities SET weekday = $1, timezone = $2, hour = $3, minute = $4,
		duration = $5, constrained = $6, startUnixTime = $7, endUnixTime = $8
		WHERE leagueId = $9 AND id = $10
		`, weekday, timezone, hour, minute, duration, constrained, start, end, leagueId, availabilityId)
	return err
}
