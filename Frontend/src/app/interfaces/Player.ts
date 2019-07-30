export class Player implements PlayerCore{
    playerId: number;
    name: string;
    gameIdentifier: string;
    mainRoster: boolean;
    constructor(mainRoster: boolean) {
        this.mainRoster = mainRoster;
    }
}

export function createUniquePlayer(mainRoster: boolean, uniqueIdFunc: () => number): Player {
    let newPlayer = new Player(mainRoster);
    newPlayer.playerId = uniqueIdFunc();
    return newPlayer;
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

export function createUniqueLoLPlayer(mainRoster: boolean, uniqueIdFunc: () => number): LoLPlayer {
    let newPlayer = new LoLPlayer(mainRoster);
    newPlayer.playerId = uniqueIdFunc();
    return newPlayer;
}

export interface PlayerCore {
    name: string;
    gameIdentifier: string;
    mainRoster: boolean;
}
export interface PlayerId {
    playerId: number;
}
