import {Injectable} from "@angular/core";
import {HttpClient} from "@angular/common/http";
import {httpOptions} from "./http-options";
import {Observable} from "rxjs/Rx";
import {LeagueService} from "./leagues.service";

@Injectable()
export class PlayersService {
    constructor(private http: HttpClient, private leagueSevice: LeagueService) {}

    public addPlayer(teamId: number, name: string, gameIdentifier: string, mainRoster: boolean): Observable<Object> {
        let url = "";
        switch(this.leagueSevice.getGame()) {
            case 'leagueoflegends': {
                url = 'http://localhost:8080/api/league-of-legends/teams/addPlayer';
                break;
            }
            default: {
                url = 'http://localhost:8080/api/teams/addPlayer';
            }
        }
        return this.http.post(url, {
            teamId: teamId,
            name: name,
            gameIdentifier: gameIdentifier,
            mainRoster: mainRoster
        }, httpOptions)
    }

    public removePlayer(teamId: number, playerId: number) {
        return this.http.request('delete', 'http://localhost:8080/api/teams/removePlayer',
            {
                withCredentials: httpOptions.withCredentials,
                headers: httpOptions.headers,
                body: {
                    teamId: teamId,
                    playerId: playerId
                }
            });
    }

    public updatePlayer(teamId: number, playerId: number, name: string, gameIdentifier: string, mainRoster: boolean): Observable<Object> {
        return this.http.put('http://localhost:8080/api/teams/updatePlayer', {
            teamId: teamId,
            playerId: playerId,
            name: name,
            gameIdentifier: gameIdentifier,
            mainRoster: mainRoster
        }, httpOptions)
    }
}
