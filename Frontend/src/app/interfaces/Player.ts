export class Player implements PlayerCore{
    playerId: number;
    name: string;
    gameIdentifier: string;
    mainRoster: boolean;
    constructor(mainRoster: boolean) {
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
