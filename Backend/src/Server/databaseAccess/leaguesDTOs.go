package databaseAccess

import "database/sql"

// In and out of DAO

type LeagueInformationDTO struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Game        string `json:"game"`
	PublicView  bool   `json:"publicView"`
	PublicJoin  bool   `json:"publicJoin"`
	SignupStart int    `json:"signup_start"`
	SignupEnd   int    `json:"signupEnd"`
	LeagueStart int    `json:"leagueStart"`
	LeagueEnd   int    `json:"leagueEnd"`
}

type TeamSummaryInformationDTO struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Tag       string `json:"tag"`
	Wins      int    `json:"wins"`
	Losses    int    `json:"losses"`
	IconSmall string `json:"iconSmall"`
	IconLarge string `json:"iconLarge"`
}

type TeamSummaryInformationArray struct {
	rows []*TeamSummaryInformationDTO
}

func (r *TeamSummaryInformationArray) Scan(rows *sql.Rows) error {
	var row TeamSummaryInformationDTO
	err := rows.Scan(
		&row.Id,
		&row.Name,
		&row.Tag,
		&row.Wins, &row.Losses,
		&row.IconSmall,
		&row.IconLarge)
	r.rows = append(r.rows, &row)
	return err
}
