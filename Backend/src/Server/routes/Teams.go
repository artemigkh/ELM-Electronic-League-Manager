package routes

import (
	"Server/databaseAccess"
	"github.com/gin-gonic/gin"
	"net/http"
)

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/createTeam
func createNewTeam() gin.HandlerFunc {
	var team databaseAccess.TeamCore
	return endpoint{
		Entity:     Team,
		AccessType: Create,
		BindData: func(ctx *gin.Context) bool {
			if ctx.ContentType() == "application/json" {
				return bindAndCheckErr(ctx, &team)
			} else if ctx.ContentType() == "multipart/form-data" {
				team.Name = ctx.PostForm("name")
				team.Tag = ctx.PostForm("tag")
				team.Description = ctx.PostForm("description")
				if team.Name == "" || team.Tag == "" {
					ctx.JSON(http.StatusBadRequest, gin.H{"error": "malformedInput"})
					return true
				} else {
					return false
				}
			} else {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "malformedInput"})
				return true
			}
		},
		IsDataInvalid: func(ctx *gin.Context) (bool, string, error) { return team.ValidateNew(getLeagueId(ctx)) },
		Core: func(ctx *gin.Context) (interface{}, error) {
			_, err := ctx.FormFile("icon")
			if err == nil {
				smallIcon, largeIcon, err := IconManager.StoreNewIcon(ctx)
				if checkErr(ctx, err) {
					return nil, err
				}

				teamId, err := TeamsDAO.CreateTeamWithIcon(getLeagueId(ctx), getUserId(ctx), team, smallIcon, largeIcon)
				return gin.H{"teamId": teamId}, err
			} else {
				teamId, err := TeamsDAO.CreateTeam(getLeagueId(ctx), getUserId(ctx), team)
				return gin.H{"teamId": teamId}, err
			}
		},
	}.createEndpointHandler()
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/getLeagueTeams
func getAllTeams() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		games, err := TeamsDAO.GetAllTeamsInLeague(getLeagueId(ctx))
		if checkErr(ctx, err) {
			return
		}

		ctx.JSON(http.StatusOK, games)
	}
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/getLeagueTeamsWithRosters
func getAllTeamsWithRosters() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		games, err := TeamsDAO.GetAllTeamsInLeagueWithRosters(getLeagueId(ctx))
		if checkErr(ctx, err) {
			return
		}

		ctx.JSON(http.StatusOK, games)
	}
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/getTeam
func getTeamInfo() gin.HandlerFunc {
	return endpoint{
		Entity:     Team,
		AccessType: View,
		Core: func(ctx *gin.Context) (interface{}, error) {
			return TeamsDAO.GetTeamInformation(getTeamId(ctx))
		},
	}.createEndpointHandler()
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/getTeamWithRosters
func getTeamWithRosters() gin.HandlerFunc {
	return endpoint{
		Entity:     Team,
		AccessType: View,
		Core: func(ctx *gin.Context) (interface{}, error) {
			return TeamsDAO.GetTeamWithRosters(getTeamId(ctx))
		},
	}.createEndpointHandler()
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/updateTeam
func editTeam() gin.HandlerFunc {
	var team databaseAccess.TeamCore
	return endpoint{
		Entity:     Team,
		AccessType: Edit,
		BindData: func(ctx *gin.Context) bool {
			if ctx.ContentType() == "application/json" {
				return bindAndCheckErr(ctx, &team)
			} else if ctx.ContentType() == "multipart/form-data" {
				team.Name = ctx.PostForm("name")
				team.Tag = ctx.PostForm("tag")
				team.Description = ctx.PostForm("description")
				if team.Name == "" || team.Tag == "" {
					ctx.JSON(http.StatusBadRequest, gin.H{"error": "malformedInput"})
					return true
				} else {
					return false
				}
			} else {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "malformedInput"})
				return true
			}
		},
		IsDataInvalid: func(ctx *gin.Context) (bool, string, error) {
			return team.ValidateEdit(getLeagueId(ctx), getTeamId(ctx))
		},
		Core: func(ctx *gin.Context) (interface{}, error) {
			err := TeamsDAO.UpdateTeam(getTeamId(ctx), team)
			if checkErr(ctx, err) {
				return nil, err
			}

			_, err = ctx.FormFile("icon")
			if err == nil {
				smallIcon, largeIcon, err := IconManager.StoreNewIcon(ctx)
				if checkErr(ctx, err) {
					return nil, err
				}

				return nil, TeamsDAO.UpdateTeamIcon(getTeamId(ctx), smallIcon, largeIcon)
			} else {
				return nil, nil
			}
		},
	}.createEndpointHandler()
}

func deleteTeam() gin.HandlerFunc {
	return endpoint{
		Entity:     Team,
		AccessType: Delete,
		Core: func(ctx *gin.Context) (interface{}, error) {
			return nil, TeamsDAO.DeleteTeam(getTeamId(ctx))
		},
	}.createEndpointHandler()
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/setTeamPermissions
func editTeamManagerPermissions() gin.HandlerFunc {
	var permissions databaseAccess.TeamPermissionsCore
	return endpoint{
		Entity:        Team,
		AccessType:    Edit, //TODO: check permissions on this one
		BindData:      func(ctx *gin.Context) bool { return bindAndCheckErr(ctx, &permissions) },
		IsDataInvalid: func(ctx *gin.Context) (bool, string, error) { return permissions.Validate() },
		Core: func(ctx *gin.Context) (interface{}, error) {
			return nil, TeamsDAO.ChangeManagerPermissions(getTeamId(ctx), getTargetUserId(ctx), permissions)
		},
	}.createEndpointHandler()
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/createPlayer
func createNewPlayer() gin.HandlerFunc {
	var player databaseAccess.PlayerCore
	return endpoint{
		Entity:     Player,
		AccessType: Create,
		BindData:   func(ctx *gin.Context) bool { return bindAndCheckErr(ctx, &player) },
		IsDataInvalid: func(ctx *gin.Context) (bool, string, error) {
			return player.ValidateNew(getLeagueId(ctx), getTeamId(ctx))
		},
		Core: func(ctx *gin.Context) (interface{}, error) {
			playerId, err := TeamsDAO.CreatePlayer(getLeagueId(ctx), getTeamId(ctx), player)
			return gin.H{"playerId": playerId}, err
		},
	}.createEndpointHandler()
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/updatePlayer
func updatePlayer() gin.HandlerFunc {
	var player databaseAccess.PlayerCore
	return endpoint{
		Entity:     Player,
		AccessType: Edit,
		BindData:   func(ctx *gin.Context) bool { return bindAndCheckErr(ctx, &player) },
		IsDataInvalid: func(ctx *gin.Context) (bool, string, error) {
			return player.ValidateEdit(getLeagueId(ctx), getTeamId(ctx), getPlayerId(ctx))
		},
		Core: func(ctx *gin.Context) (interface{}, error) {
			return nil, TeamsDAO.UpdatePlayer(getPlayerId(ctx), player)
		},
	}.createEndpointHandler()
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/deletePlayer
func deletePlayer() gin.HandlerFunc {
	return endpoint{
		Entity:     Player,
		AccessType: Delete,
		Core: func(ctx *gin.Context) (interface{}, error) {
			return nil, TeamsDAO.DeletePlayer(getPlayerId(ctx))
		},
	}.createEndpointHandler()
}

func RegisterTeamHandlers(g *gin.RouterGroup) {
	g.POST("/teams", createNewTeam())
	g.GET("/teams", getAllTeams())
	g.GET("/teamsWithRosters", getAllTeamsWithRosters())

	withTeamId := g.Group("/teams/:teamId", storeTeamId())
	withTeamId.GET("", getTeamInfo())
	withTeamId.GET("/withRosters", getTeamWithRosters())
	withTeamId.PUT("", editTeam())
	withTeamId.DELETE("", deleteTeam())

	withTeamId.PUT("/permissions/:userId", storeTargetUserId(), editTeamManagerPermissions()) //TODO: test this one

	withTeamId.POST("/players", createNewPlayer())
	withPlayerId := withTeamId.Group("/players/:playerId", storePlayerId())
	withPlayerId.PUT("", updatePlayer())
	withPlayerId.DELETE("", deletePlayer())
}
