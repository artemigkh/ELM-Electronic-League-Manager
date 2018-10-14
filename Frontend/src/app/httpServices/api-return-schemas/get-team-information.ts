export interface GtiPlayer {
    id: number;
    name: string;
    gameIdentifier: string;
    mainRoster: boolean;
}

export interface GtiTeam {
    name: string;
    tag: string;
    wins: number;
    losses: number;
    players: GtiPlayer[];
}
