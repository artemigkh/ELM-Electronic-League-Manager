export interface UserPermissions {
    leaguePermissions: LeaguePermissions;
    teamPermissions: TeamPermissions[];
}

export interface LeaguePermissions {
    administrator: boolean;
    createTeams: boolean;
    editTeams: boolean;
    editGames: boolean;
}

export interface TeamPermissions {
    id: number;
    administrator: boolean;
    information: boolean;
    players: boolean;
    reportResults: boolean;
}
