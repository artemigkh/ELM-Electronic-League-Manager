package databaseAccess

import (
	"github.com/gorilla/securecookie"
	"encoding/hex"
	"time"
	"database/sql"
)

type PgInviteCodesDAO struct{}

type TeamManagerInviteCode struct {
	Code string
	CreationTime int
	LeagueId int
	TeamId int
	editPermissions bool
	editTeamInfo bool
	editPlayers bool
	reportResults bool
}

func (d *PgInviteCodesDAO) CreateTeamManagerInviteCode(teamId int,
	editPermissions, editTeamInfo, editPlayers, reportResult bool) (string, error) {
	code := hex.EncodeToString(securecookie.GenerateRandomKey(32))
	_, err := psql.Insert("teamManagerInviteCodes").
		Columns("code", "creationTime", "teamId", "editPermissions", "editTeamInfo",
			"editPlayers", "reportResult").
		Values(code, int32(time.Now().Unix()), teamId, editPermissions, editTeamInfo,
			editPlayers, reportResult).
		RunWith(db).Exec()

	return code, err
}

func (d *PgInviteCodesDAO) GetTeamManagerInviteCodeInformation(code string) (*TeamManagerInviteCode, error) {
	return nil, nil
}

func (d *PgInviteCodesDAO) UseTeamManagerInviteCode(userId int, code string) error {
	return nil
}

func (d *PgInviteCodesDAO) IsTeamManagerInviteCodeValid(code string) (bool, string, error) {
	var codeInfo TeamManagerInviteCode
	err := psql.Select("code", "creationTime", "leagueId", "teamId",
		"editPermissions", "editTeamInfo", "editPlayers", "reportResult").
		From("teamManagerInviteCodes").
		Where("code = ?", code).
		RunWith(db).QueryRow().Scan(&codeInfo.Code, &codeInfo.CreationTime, &codeInfo.LeagueId, &codeInfo.TeamId,
			&codeInfo.editPermissions, &codeInfo.editTeamInfo, &codeInfo.editPlayers, &codeInfo.reportResults)
	if err == sql.ErrNoRows {
		return false, "", nil
	} else if err != nil {
		return false, "", err
	}

	//check if code expired (more than 24 hours)
	if int(time.Now().Unix()) - codeInfo.CreationTime > 24 * 3600 {
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