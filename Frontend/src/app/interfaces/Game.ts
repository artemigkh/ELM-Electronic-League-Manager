import {TeamDisplay} from "./Team";

import {Moment} from "moment";
import * as moment from "moment";

export class Game implements GameResult{
    gameId: number;
    complete: boolean;
    /**
     * The start time of the game in seconds since unix epoch
     */
    gameTime: number;
    team1: TeamDisplay;
    team2: TeamDisplay;
    winnerId: number;
    loserId: number;
    scoreTeam1: number;
    scoreTeam2: number;

    constructor(g: GameCore) {
        Object.entries(g).forEach(o => this[o[0]] = o[1]);
    }

}

export interface SortedGames {
    completedGames: Game[];
    upcomingGames: Game[];
}

export interface CompetitionWeek {
    weekStart: number;
    games: Game[];
}

export function EmptySortedGames(): SortedGames {
    return {
        completedGames: [],
        upcomingGames: []
    }
}

export interface GameCore {
    /**
     * The start time of the game in seconds since unix epoch
     */
    gameTime: number;
    team1: TeamDisplay;
    team2: TeamDisplay;
}
export class GameCreationInformation {
    team1Id: number;
    team2Id: number;
    /**
     * The start time of the game in seconds since unix epoch
     */
    gameTime: number;

    constructor() {
        this.gameTime = moment().unix();
    }

    updateFromMoments(date: Moment, time: Moment) {
        this.gameTime = date.clone().minute(time.minute()).hour(time.hour()).unix();
    }

}
export interface GameId {
    gameId: number;
}
export interface GameResult {
    winnerId: number;
    loserId: number;
    scoreTeam1: number;
    scoreTeam2: number;
}
export interface GameTime {
    /**
     * The start time of the game in seconds since unix epoch
     */
    gameTime: number;
}
