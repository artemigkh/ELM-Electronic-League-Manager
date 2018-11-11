package routes

/**
 * @api{POST} /api/inviteCodes/team Create Team Manager Invite Code
 * @apiGroup InviteCodes
 * @apiDescription Create Team Manager Invite Code. Logged in user must be a team admin or league admin
 * (has editPermissions in either). Code expires after 24h
 *
 * @apiParam {number} teamId The unique numerical identifier of the team
 * @apiParam {boolean} editPermissions Whether the invited user can edit permissions of other managers on the team (admin)
 * @apiParam {boolean} editTeamInfo Whether the invited user can edit team information
 * @apiParam {boolean} editPlayers Whether the invited user can edit players on the team
 * @apiParam {boolean} reportResult Whether the invited user can report game results of this team
 *
 * @apiSuccess {string} A 16 character string that is the invite code
 *
 * @apiError notLoggedIn No user is logged in
 * @apiError noActiveLeague There is no active league selected
 * @apiError teamDoesNotExist The specified team does not exist
 */

/**
 * @api{GET} /api/inviteCodes/team Create Team Manager Invic
 * @apiName createNewLeague
 * @apiGroup Leagues
 * @apiDescription Register a new league
 *
 * @apiParam {string} name the name of the league
 * @apiParam {string} description A brief (<500) char description of the league
 * @apiParam {boolean} publicView should the league be viewable by people not playing in the league?
 * @apiParam {boolean} publicJoin should the league be joinable by any team that has viewing rights?
 *
 * @apiSuccess {int} id the primary id of the created league
 *
 * @apiError notLoggedIn No user is logged in
 * @apiError nameTooLong The league name has exceeded 50 characters
 * @apiError descriptionTooLong The description has exceeded 500 characters
 * @apiError nameInUse The league name is currently in use
 */
