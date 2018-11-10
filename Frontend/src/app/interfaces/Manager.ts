export interface Manager {
   userId: number;
   userEmail: string;
   editPermissions: boolean;
   editTeamInfo: boolean;
   editPlayers: boolean;
   reportResult: boolean;
}

export interface TeamManagers {
    teamId: number;
    teamName: string;
    teamTag: string;
    managers: Manager[];
}
