package routes

import (
	"Server/dataModel"
	"github.com/gin-gonic/gin"
	"net/http"
)

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/createTeam
func createNewTeam() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var team dataModel.TeamCore
		endpoint{
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
			IsDataInvalid: func(ctx *gin.Context) (bool, string, error) { return team.ValidateNew(getLeagueId(ctx), TeamDAO) },
			Core: func(ctx *gin.Context) (interface{}, error) {
				_, err := ctx.FormFile("icon")
				if err == nil {
					smallIcon, largeIcon, err := IconManager.StoreNewIconFromForm(ctx)
					if checkErr(ctx, err) {
						return nil, err
					}

					teamId, err := TeamDAO.CreateTeamWithIcon(getLeagueId(ctx), getUserId(ctx), team, smallIcon, largeIcon)
					return gin.H{"teamId": teamId}, err
				} else {
					teamId, err := TeamDAO.CreateTeam(getLeagueId(ctx), getUserId(ctx), team)
					return gin.H{"teamId": teamId}, err
				}
			},
		}.createEndpointHandler()(ctx)
	}

}

func createNewTeamWithPlayers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var team dataModel.TeamWithPlayersCore
		endpoint{
			Entity:     Team,
			AccessType: Create,
			BindData: func(ctx *gin.Context) bool {
				return bindAndCheckErr(ctx, &team)
			},
			IsDataInvalid: func(ctx *gin.Context) (bool, string, error) {
				return team.Validate(getLeagueId(ctx), TeamDAO)
			},
			Core: func(ctx *gin.Context) (interface{}, error) {
				var err error
				smallIcon := ""
				largeIcon := ""
				if len(team.Icon) > 0 {
					smallIcon, largeIcon, err = IconManager.StoreNewIconFromBase64String(team.Icon)
					if checkErr(ctx, err) {
						return nil, err
					}
				}
				return TeamDAO.CreateTeamWithPlayers(
					getLeagueId(ctx), getUserId(ctx), team.Team, team.Players, smallIcon, largeIcon)
			},
		}.createEndpointHandler()(ctx)
	}
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/getLeagueTeams
func getAllTeams() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		games, err := TeamDAO.GetAllTeamsInLeague(getLeagueId(ctx))
		if checkErr(ctx, err) {
			return
		}

		ctx.JSON(http.StatusOK, games)
	}
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/getLeagueTeamsWithRosters
func getAllTeamsWithRosters() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		games, err := TeamDAO.GetAllTeamsInLeagueWithRosters(getLeagueId(ctx))
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
			return TeamDAO.GetTeamInformation(getTeamId(ctx))
		},
	}.createEndpointHandler()
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/getTeamWithRosters
func getTeamWithRosters() gin.HandlerFunc {
	return endpoint{
		Entity:     Team,
		AccessType: View,
		Core: func(ctx *gin.Context) (interface{}, error) {
			return TeamDAO.GetTeamWithRosters(getTeamId(ctx))
		},
	}.createEndpointHandler()
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/updateTeam
func editTeam() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var team dataModel.TeamCore
		endpoint{
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
				return team.ValidateEdit(getLeagueId(ctx), getTeamId(ctx), TeamDAO)
			},
			Core: func(ctx *gin.Context) (interface{}, error) {
				err := TeamDAO.UpdateTeam(getTeamId(ctx), team)
				if checkErr(ctx, err) {
					return nil, err
				}

				_, err = ctx.FormFile("icon")
				if err == nil {
					smallIcon, largeIcon, err := IconManager.StoreNewIconFromForm(ctx)
					if checkErr(ctx, err) {
						return nil, err
					}

					return nil, TeamDAO.UpdateTeamIcon(getTeamId(ctx), smallIcon, largeIcon)
				} else {
					return nil, nil
				}
			},
		}.createEndpointHandler()(ctx)
	}
}

func deleteTeam() gin.HandlerFunc {
	return endpoint{
		Entity:     Team,
		AccessType: Delete,
		Core: func(ctx *gin.Context) (interface{}, error) {
			return nil, TeamDAO.DeleteTeam(getTeamId(ctx))
		},
	}.createEndpointHandler()
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/setTeamPermissions
func editTeamManagerPermissions() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var permissions dataModel.TeamPermissionsCore
		endpoint{
			Entity:        Team,
			AccessType:    Edit, //TODO: check permissions on this one
			BindData:      func(ctx *gin.Context) bool { return bindAndCheckErr(ctx, &permissions) },
			IsDataInvalid: func(ctx *gin.Context) (bool, string, error) { return permissions.Validate() },
			Core: func(ctx *gin.Context) (interface{}, error) {
				return nil, TeamDAO.ChangeManagerPermissions(getTeamId(ctx), getTargetUserId(ctx), permissions)
			},
		}.createEndpointHandler()(ctx)
	}
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/createPlayer
func createNewPlayer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var player dataModel.PlayerCore
		endpoint{
			Entity:     Player,
			AccessType: Create,
			BindData:   func(ctx *gin.Context) bool { return bindAndCheckErr(ctx, &player) },
			IsDataInvalid: func(ctx *gin.Context) (bool, string, error) {
				return player.ValidateNew(getLeagueId(ctx), getTeamId(ctx), TeamDAO)
			},
			Core: func(ctx *gin.Context) (interface{}, error) {
				playerId, err := TeamDAO.CreatePlayer(getLeagueId(ctx), getTeamId(ctx), player)
				return gin.H{"playerId": playerId}, err
			},
		}.createEndpointHandler()(ctx)
	}
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/updatePlayer
func updatePlayer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var player dataModel.PlayerCore
		endpoint{
			Entity:     Player,
			AccessType: Edit,
			BindData:   func(ctx *gin.Context) bool { return bindAndCheckErr(ctx, &player) },
			IsDataInvalid: func(ctx *gin.Context) (bool, string, error) {
				return player.ValidateEdit(getLeagueId(ctx), getTeamId(ctx), getPlayerId(ctx), TeamDAO)
			},
			Core: func(ctx *gin.Context) (interface{}, error) {
				return nil, TeamDAO.UpdatePlayer(getPlayerId(ctx), player)
			},
		}.createEndpointHandler()(ctx)
	}
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/deletePlayer
func deletePlayer() gin.HandlerFunc {
	return endpoint{
		Entity:     Player,
		AccessType: Delete,
		Core: func(ctx *gin.Context) (interface{}, error) {
			return nil, TeamDAO.DeletePlayer(getPlayerId(ctx))
		},
	}.createEndpointHandler()
}

func RegisterTeamHandlers(g *gin.RouterGroup) {
	g.POST("/teams", createNewTeam())
	g.POST("/teamsWithPlayers", createNewTeamWithPlayers())
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
