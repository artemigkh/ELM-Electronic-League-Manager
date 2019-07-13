export class Player implements PlayerCore{
    playerId: number;
    name: string;
    gameIdentifier: string;
    mainRoster: boolean;
    constructor(mainRoster: boolean) {
        this.mainRoster = mainRoster;
    }
}

export class LoLPlayer extends Player {
    playerId: number;
    name: string;
    gameIdentifier: string;
    mainRoster: boolean;
    position: string;
    rank: string;
    tier: string;
    constructor(mainRoster: boolean) {
        super(mainRoster);
        this.mainRoster = mainRoster;
    }
}
export interface PlayerCore {
    name: string;
    gameIdentifier: string;
    mainRoster: boolean;
}
export interface PlayerId {
    playerId: number;
}
