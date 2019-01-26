export class Player {
    id: number;
    name: string;
    gameIdentifier: string;
}

export class LeagueOfLegendsPlayer extends Player {
    id: number;
    name: string;
    gameIdentifier: string;
    rank: string;
    tier: string;
}
