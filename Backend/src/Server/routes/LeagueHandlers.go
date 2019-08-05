package routes

import (
	"Server/dataModel"
	"github.com/gin-gonic/gin"
	"net/http"
)

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/createLeague
func createNewLeague() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var league dataModel.LeagueCore
		endpoint{
			Entity:        League,
			AccessType:    Create,
			BindData:      func(ctx *gin.Context) bool { return bindAndCheckErr(ctx, &league) },
			IsDataInvalid: func(ctx *gin.Context) (bool, string, error) { return league.ValidateNew(LeagueDAO) },
			Core: func(ctx *gin.Context) (interface{}, error) {
				leagueId, err := LeagueDAO.CreateLeague(getUserId(ctx), league)
				return gin.H{"leagueId": leagueId}, err
			},
		}.createEndpointHandler()(ctx)
	}
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/updateLeague
func updateLeagueInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var league dataModel.LeagueCore
		endpoint{
			Entity:     League,
			AccessType: Edit,
			BindData:   func(ctx *gin.Context) bool { return bindAndCheckErr(ctx, &league) },
			IsDataInvalid: func(ctx *gin.Context) (bool, string, error) {
				return league.ValidateEdit(getLeagueId(ctx), LeagueDAO, GameDAO)
			},
			Core: func(ctx *gin.Context) (interface{}, error) {
				return nil, LeagueDAO.UpdateLeague(getLeagueId(ctx), league)
			},
		}.createEndpointHandler()(ctx)
	}
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/updateLeagueMd
func setLeagueMarkdown() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var markdown dataModel.Markdown
		endpoint{
			Entity:        League,
			AccessType:    Edit,
			BindData:      func(ctx *gin.Context) bool { return bindAndCheckErr(ctx, &markdown) },
			IsDataInvalid: func(ctx *gin.Context) (bool, string, error) { return markdown.Validate() },
			Core: func(ctx *gin.Context) (interface{}, error) {
				oldFile, err := LeagueDAO.GetMarkdownFile(getLeagueId(ctx))
				if err != nil {
					return nil, err
				}

				fileName, err := MarkdownManager.StoreMarkdown(markdown.Markdown, oldFile)
				if err != nil {
					return nil, err
				}

				return nil, LeagueDAO.SetMarkdownFile(getLeagueId(ctx), fileName)
			},
		}.createEndpointHandler()(ctx)
	}
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/setActiveLeague
func setActiveLeague() gin.HandlerFunc {
	return endpoint{
		Entity:     League,
		AccessType: View,
		Core: func(ctx *gin.Context) (interface{}, error) {
			err := ElmSessions.SetActiveLeague(ctx, getLeagueId(ctx))
			if err != nil {
				return nil, err
			}
			return LeagueDAO.GetLeagueInformation(getLeagueId(ctx))
		},
	}.createEndpointHandler()
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/getLeague
func getActiveLeagueInformation() gin.HandlerFunc {
	return endpoint{
		Entity:     League,
		AccessType: View,
		Core: func(ctx *gin.Context) (interface{}, error) {
			return LeagueDAO.GetLeagueInformation(getLeagueId(ctx))
		},
	}.createEndpointHandler()
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/joinLeague
func joinActiveLeague(ctx *gin.Context) {
	allowedJoin, err := LeagueDAO.CanJoinLeague(getLeagueId(ctx), getUserId(ctx))
	if accessForbidden(ctx, allowedJoin, err) {
		return
	}

	err = LeagueDAO.JoinLeague(getLeagueId(ctx), getUserId(ctx))
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/getLeagueTeamManagers
func getTeamManagers() gin.HandlerFunc {
	return endpoint{
		Entity:     League,
		AccessType: Edit,
		Core: func(ctx *gin.Context) (interface{}, error) {
			return LeagueDAO.GetTeamManagerInformation(getLeagueId(ctx))
		},
	}.createEndpointHandler()
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/getPublicLeagues
func getPublicLeagues(ctx *gin.Context) {
	leagueList, err := LeagueDAO.GetPublicLeagueList()
	if checkErr(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, leagueList)
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/updateLeaguePermissions
func setLeaguePermissions() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var permissions dataModel.LeaguePermissionsCore
		endpoint{
			Entity:        League,
			AccessType:    Edit,
			BindData:      func(ctx *gin.Context) bool { return bindAndCheckErr(ctx, &permissions) },
			IsDataInvalid: func(ctx *gin.Context) (bool, string, error) { return permissions.Validate() },
			Core: func(ctx *gin.Context) (interface{}, error) {
				return nil, LeagueDAO.SetLeaguePermissions(getLeagueId(ctx), getTargetUserId(ctx), permissions)
			},
		}.createEndpointHandler()(ctx)
	}
}

// https://artemigkh.github.io/ELM-Electronic-League-Manager/#operation/getLeagueMd
func getLeagueMarkdown() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var fileName string
		var err error
		endpoint{
			Entity:     League,
			AccessType: View,
			BindData: func(ctx *gin.Context) bool {
				fileName, err = LeagueDAO.GetMarkdownFile(getLeagueId(ctx))
				if checkErr(ctx, err) {
					return true
				} else {
					return false
				}
			},
			Core: func(ctx *gin.Context) (interface{}, error) {
				markdown, err := MarkdownManager.GetMarkdown(fileName)
				return gin.H{"markdown": markdown}, err
			},
		}.createEndpointHandler()(ctx)
	}
}

func RegisterLeagueHandlers(g *gin.RouterGroup) {
	// League Manage
	g.POST("", createNewLeague())

	g.PUT("", updateLeagueInfo())
	g.PUT("/markdown", setLeagueMarkdown())
	g.GET("/teamManagers", getTeamManagers())
	g.PUT("/permissions/:userId", storeTargetUserId(), setLeaguePermissions()) //TODO: test this one in integrat

	// League Interact
	g.POST("/setActiveLeague/:leagueId", storeTargetLeagueId(), setActiveLeague())
	g.POST("/join", joinActiveLeague)

	// League Information
	g.GET("", getActiveLeagueInformation())
	g.GET("/markdown", getLeagueMarkdown())
	g.GET("/publicLeagues", getPublicLeagues)
}
