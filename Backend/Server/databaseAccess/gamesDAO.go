package databaseAccess

import "github.com/Masterminds/squirrel"

type GameInformation struct {
	Team1ID int `json:"team1Id"`
	Team2ID int `json:"team2Id"`
	GameTime int `json:"gameTime"`
	Complete bool `json:"complete"`
	WinnerID int `json:"winnerId"`
	ScoreTeam1 int `json:"scoreTeam1"`
	ScoreTeam2 int `json:"scoreTeam2"`
}

type PgGamesDAO struct {
	psql squirrel.StatementBuilderType
}

func (d *PgGamesDAO) CreateGame(team1ID, team2ID, gameTime int) (int, error) {
	return 0, nil
}

func (d *PgGamesDAO) DoesExistConflict(team1ID, team2ID, gameTime int) (bool, error) {
	return false, nil
}

func (d *PgGamesDAO) GetGameInformation(gameID int) (*GameInformation, error) {
	return nil, nil
}
