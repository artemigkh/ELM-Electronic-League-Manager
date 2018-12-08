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
}
