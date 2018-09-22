import {Team} from "./Team";

export interface Game {
    id: number;
    gameTime: number;
    complete: boolean;
    winnerId: number;
    scoreTeam1: number;
    scoreTeam2: number;
    team1: Team;
    team2: Team;
}
