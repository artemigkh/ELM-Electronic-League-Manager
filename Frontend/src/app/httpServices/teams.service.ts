import {Injectable} from "@angular/core";
import {HttpClient} from "@angular/common/http";
import {httpOptions} from "./http-options";
import {Observable} from "rxjs/Rx";
import {GtiTeam} from "./api-return-schemas/get-team-information";
import {Player} from "../interfaces/Player";
import {Team} from "../interfaces/Team";

@Injectable()
export class TeamsService {
    constructor(private http: HttpClient) {}

    public createNewTeam(name: string, tag: string): Observable<Object> {
        return this.http.post('http://localhost:8080/api/teams/', {
            name: name,
            tag: tag
        }, httpOptions)
    }

    public updateTeam(id: number, name: string, tag: string): Observable<Object> {
        return this.http.put('http://localhost:8080/api/teams/updateTeam/' + id, {
            name: name,
            tag: tag
        }, httpOptions)
    }

    public deleteTeam(id: number): Observable<Object> {
        return this.http.delete('http://localhost:8080/api/teams/removeTeam/' + id, httpOptions)
    }

    public updateManagerPermissions(teamId: number, userId: number, administrator: boolean, information: boolean,
                                    players: boolean, reportResults: boolean) {
        return this.http.put('http://localhost:8080/api/teams/updatePermissions', {
            teamId: teamId,
            userId : userId,
            administrator: administrator,
            information: information,
            players: players,
            reportResults: reportResults
        }, httpOptions)
    }

    public getTeamInformation(teamId: number): Observable<Object> {
        return new Observable(observer => {
            this.http.get('http://localhost:8080/api/teams/' + teamId, httpOptions).subscribe(
                (next: Team) => {
                    let players = next.players;
                    let team = next;
                    team.substitutes = [];
                    team.players = [];
                    if(players) {
                        players.forEach((player: any)=> {
                            let tempPlayer: Player = {
                                id: player.id,
                                name: player.name,
                                gameIdentifier: player.gameIdentifier
                            };

                            if(player.mainRoster) {
                                team.players.push(tempPlayer);
                            } else {
                                team.substitutes.push(tempPlayer);
                            }
                        });
                    }
                    observer.next(team);
                }, error => {
                    observer.error(error);
                    console.log(error);
                }
            );
        });
    }
}
