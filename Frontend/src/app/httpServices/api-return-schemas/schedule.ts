export interface teamInfo {
    id: number;
    name: string;
    tag: string;
    iconSmall: string;
}

export interface scheduledGame {
    team1Id: number;
    team2Id: number;
    gameTime: number;
}

export interface schedule {
    teams: teamInfo[];
    games: scheduledGame[];
}

export interface availability {
    id: number;
    weekday: number;
    timezone: number;
    hour: number;
    minute: number;
    duration: number;
    constrained: boolean;
    start: number;
    end: number;
}
