import {Team} from "../interfaces/Team";

export function getTeamName(id: number, teams: Team[]): string {
    let toReturn = "unknown name";
    teams.forEach(team => {
        if(team.id == id){
            toReturn = team.name;
        }
    });
    return toReturn;
}
