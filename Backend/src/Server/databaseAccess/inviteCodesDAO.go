package databaseAccess

import (
	"database/sql"
	"encoding/hex"
	"github.com/gorilla/securecookie"
	"time"
)

type PgInviteCodesDAO struct{}

type TeamManagerInviteCode struct {
	Code          string `json:"code"`
	CreationTime  int    `json:"creationTime"`
	LeagueId      int    `json:"leagueId"`
	TeamId        int    `json:"teamId"`
	Administrator bool   `json:"administrator"`
	Information   bool   `json:"information"`
	Players       bool   `json:"players"`
	ReportResults bool   `json:"reportResults"`
}

func (d *PgInviteCodesDAO) CreateTeamManagerInviteCode(leagueId, teamId int,
	administrator, information, players, reportResults bool) (string, error) {
	code := hex.EncodeToString(securecookie.GenerateRandomKey(32))
	_, err := psql.Insert("teamManagerInviteCodes").
		Columns("code", "creationTime", "leagueId", "teamId",
			"administrator", "information", "players", "reportResults").
		Values(code, int32(time.Now().Unix()), teamId, leagueId,
			administrator, information, players, reportResults).
		RunWith(db).Exec()

	return code, err
}

func (d *PgInviteCodesDAO) GetTeamManagerInviteCodeInformation(code string) (*TeamManagerInviteCode, error) {
	var codeInfo TeamManagerInviteCode
	err := psql.Select("code", "creationTime", "leagueId", "teamId",
		"administrator", "information", "players", "reportResults").
		From("teamManagerInviteCodes").
		Where("code = ?", code).
		RunWith(db).QueryRow().Scan(&codeInfo.Code, &codeInfo.CreationTime, &codeInfo.LeagueId, &codeInfo.TeamId,
		&codeInfo.Administrator, &codeInfo.Information, &codeInfo.Players, &codeInfo.ReportResults)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &codeInfo, nil
}

func (d *PgInviteCodesDAO) UseTeamManagerInviteCode(userId int, code string) error {
	// get code information and delete it
	var codeInfo TeamManagerInviteCode
	err := psql.Select("code", "creationTime", "leagueId", "teamId",
		"administrator", "information", "players", "reportResults").
		From("teamManagerInviteCodes").
		Where("code = ?", code).
		RunWith(db).QueryRow().Scan(&codeInfo.Code, &codeInfo.CreationTime, &codeInfo.LeagueId, &codeInfo.TeamId,
		&codeInfo.Administrator, &codeInfo.Information, &codeInfo.Players, &codeInfo.ReportResults)
	if err != nil {
		return err
	}

	_, err = psql.Delete("teamManagerInviteCodes").
		Where("code = ?", code).
		RunWith(db).Exec()
	if err != nil {
		return err
	}

	//create permissions entry
	_, err = psql.Insert("teamPermissions").
		Columns("userId", "teamId", "administrator", "information", "players", "reportResults").
		Values(userId, codeInfo.TeamId, codeInfo.Administrator,
			codeInfo.Information, codeInfo.Players, codeInfo.ReportResults).
		RunWith(db).Exec()
	if err != nil {
		return err
	}

	return nil
}

func (d *PgInviteCodesDAO) IsTeamManagerInviteCodeValid(code string) (bool, string, error) {
	var codeInfo TeamManagerInviteCode
	err := psql.Select("code", "creationTime", "leagueId", "teamId",
		"administrator", "information", "players", "reportResults").
		From("teamManagerInviteCodes").
		Where("code = ?", code).
		RunWith(db).QueryRow().Scan(&codeInfo.Code, &codeInfo.CreationTime, &codeInfo.LeagueId, &codeInfo.TeamId,
		&codeInfo.Administrator, &codeInfo.Information, &codeInfo.Players, &codeInfo.ReportResults)
	if err == sql.ErrNoRows {
		return false, "", nil
	} else if err != nil {
		return false, "", err
	}

	//check if code expired (more than 24 hours)
	if int(time.Now().Unix())-codeInfo.CreationTime > 24*3600 {
		return false, "expired", nil
	}

	//check if team still exists
	exists, err := doesTeamExist(codeInfo.LeagueId, codeInfo.TeamId)
	if err != nil {
		return false, "", err
	}
	if !exists {
		return false, "teamDeleted", nil
	}

	return true, "", nil
}
