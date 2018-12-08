import {Team} from "../interfaces/Team";
import {Game} from "../interfaces/Game";

export function getTeamName(id: number, teams: Team[]): string {
    let toReturn = "unknown name";
    teams.forEach(team => {
        if(team.id == id){
            toReturn = team.name;
        }
    });
    return toReturn;
}

export function doesGameHaveTeam(teamId: number) {
    return function (game: Game) {
        return game.team1Id == teamId || game.team2Id == teamId;
    };
}

export function gameSort(a: Game, b: Game): number {
    return (a.gameTime > b.gameTime) ? 1 :
           ((a.gameTime < b.gameTime) ? -1 : 0);
}
