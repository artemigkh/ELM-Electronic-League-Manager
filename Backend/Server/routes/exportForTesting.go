package routes

// Middleware
var Testing_Export_authenticate = authenticate
var Testing_Export_getActiveLeague = getActiveLeague
var Testing_Export_getUrlId = getUrlId
var Testing_Export_getTeamEditPermissions = getTeamEditPermissions
var Testing_Export_failIfNoTeamCreatePermissions = failIfNoTeamCreatePermissions
var Testing_Export_getReportResultPermissions = getReportResultPermissions
var Testing_Export_failIfCannotJoinLeague = failIfCannotJoinLeague
var Testing_Export_failIfNotLeagueAdmin = failIfNotLeagueAdmin

// Users
var Testing_Export_createNewUser = createNewUser
var Testing_Export_getProfile = getProfile
var Testing_Export_login = login

// Leagues
var Testing_Export_createNewLeague = createNewLeague
var Testing_Export_joinActiveLeague = joinActiveLeague
var Testing_Export_setActiveLeague = setActiveLeague
var Testing_Export_getActiveLeagueInformation = getActiveLeagueInformation
var Testing_Export_getTeamSummary = getTeamSummary
var Testing_Export_getTeamManagers = getTeamManagers

// Teams
var Testing_Export_createNewTeam = createNewTeam
var Testing_Export_getTeamInformation = getTeamInformation
var Testing_Export_addPlayerToTeam = addPlayerToTeam
var Testing_Export_removePlayerFromTeam = removePlayerFromTeam

// Games
var Testing_Export_createNewGame = createNewGame
var Testing_Export_getGameInformation = getGameInformation
var Testing_Export_reportGameResult = reportGameResult
