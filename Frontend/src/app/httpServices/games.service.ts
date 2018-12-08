import {HttpClient} from "@angular/common/http";
import {Injectable} from "@angular/core";
import {Observable} from "rxjs/Rx";
import {httpOptions} from "./http-options";

@Injectable()
export class GamesService {
    constructor(private http: HttpClient) {
    }

    public reportResult(gameId: number, winnerId: number,
                        scoreTeam1: number, scoreTeam2: number): Observable<Object> {
        return this.http.post('http://localhost:8080/api/games/report/' + gameId, {
            winnerId: winnerId,
            scoreTeam1: scoreTeam1,
            scoreTeam2: scoreTeam2
        }, httpOptions)
    }

    public createNewGame(team1Id: number, team2Id: number, gameTime: number): Observable<Object> {
        return this.http.post('http://localhost:8080/api/games/', {
            team1Id: team1Id,
            team2Id: team2Id,
            gameTime: gameTime
        }, httpOptions)
    }

    public rescheduleGame(gameId: number, gameTime: number): Observable<Object> {
        return this.http.put('http://localhost:8080/api/games/', {
            id: gameId,
            gameTime: gameTime
        }, httpOptions)
    }

    public deleteGame(teamId: number) {
        return this.http.request('delete', 'http://localhost:8080/api/games/' + teamId,
            httpOptions);
    }
}
