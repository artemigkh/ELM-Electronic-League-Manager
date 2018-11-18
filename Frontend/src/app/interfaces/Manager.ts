export interface Manager {
    userId: number;
    userEmail: string;
    administrator: boolean;
    information: boolean;
    players: boolean;
    reportResults: boolean;
}

export interface TeamManagers {
    teamId: number;
    teamName: string;
    teamTag: string;
    managers: Manager[];
}
