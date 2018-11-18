package databaseAccess

import (
	"database/sql"
)

type LeagueInformation struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type LeaguePermissions struct {
	Administrator bool
	CreateTeams   bool
	EditTeams     bool
	EditGames     bool
}

type TeamSummaryInformation struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Tag    string `json:"tag"`
	Wins   int    `json:"wins"`
	Losses int    `json:"losses"`
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
	UserId          int    `json:"userId"`
	UserEmail       string `json:"userEmail"`
	EditPermissions bool   `json:"editPermissions"`
	EditTeamInfo    bool   `json:"editTeamInfo"`
	EditPlayers     bool   `json:"editPlayers"`
	ReportResult    bool   `json:"reportResult"`
}

type PublicLeagueInformation struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PublicJoin  bool   `json:"publicJoin"`
}

type PgLeaguesDAO struct{}

func (d *PgLeaguesDAO) CreateLeague(userId int, name, description string, publicView, publicJoin bool) (int, error) {
	var leagueId int
	err := psql.Insert("leagues").Columns("name", "description", "publicView", "publicJoin").
		Values(name, description, publicView, publicJoin).Suffix("RETURNING \"id\"").
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

func (d *PgLeaguesDAO) IsNameInUse(name string) (bool, error) {
	err := psql.Select("name").
		From("leagues").
		Where("name = ?", name).
		RunWith(db).QueryRow().Scan(&name)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
	return false, nil
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
	err := psql.Select("id", "name", "description").
		From("leagues").
		Where("id = ?", leagueId).
		RunWith(db).QueryRow().Scan(&leagueInfo.Id, &leagueInfo.Name, &leagueInfo.Description)
	if err != nil {
		return nil, err
	}

	return &leagueInfo, nil
}

func (d *PgLeaguesDAO) GetTeamSummary(leagueId int) ([]TeamSummaryInformation, error) {
	rows, err := psql.Select("id", "name", "tag", "wins", "losses").From("teams").
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
		err := rows.Scan(&team.Id, &team.Name, &team.Tag, &team.Wins, &team.Losses)
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
		userId          int
		teamId          int
		email           string
		name            string
		tag             string
		editPermissions bool
		editTeamInfo    bool
		editPlayers     bool
		reportResult    bool
	)

	//iterate through the rows returned from database
	for rows.Next() {
		//scan the variables from the sql row into local variables
		err := rows.Scan(&userId, &teamId, &email, &name,
			&tag, &editPermissions, &editTeamInfo, &editPlayers, &reportResult)
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
			UserId:          userId,
			UserEmail:       email,
			EditPermissions: editPermissions,
			EditTeamInfo:    editTeamInfo,
			EditPlayers:     editPlayers,
			ReportResult:    reportResult,
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
	rows, err := psql.Select("id", "name", "description", "publicJoin").
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
		err := rows.Scan(&league.Id, &league.Name, &league.Description, &league.PublicJoin)
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

func (d *PgLeaguesDAO) GetLeaguePermissions(leagueId, userId int) (*LeaguePermissions, error) {
	var lp LeaguePermissions
	err := psql.Select("administrator", "createTeams", "editTeams", "editGames").
		From("leaguePermissions").
		Where("userId = ? AND leagueId = ?", userId, leagueId).
		RunWith(db).QueryRow().Scan(&lp.Administrator, &lp.CreateTeams, &lp.EditTeams, &lp.EditGames)
	if err != nil {
		return nil, err
	}
	return &lp, nil
}