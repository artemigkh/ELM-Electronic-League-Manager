import {HttpClient, HttpParams} from "@angular/common/http";
import {Injectable, Query} from "@angular/core";
import {Observable} from "rxjs/Rx";
import {
    CompetitionWeek,
    Game,
    GameCreationInformation,
    GameId,
    GameResult,
    GameTime,
    SortedGames
} from "../interfaces/Game";
import * as moment from "moment";
import {Moment} from "moment";
import {httpOptions} from "./http-options";

@Injectable()
export class GamesService {
    constructor(private http: HttpClient) {
    }

    public getLeagueGames(args: {limit?: string; teamId?: string;}): Observable<SortedGames> {
        let queryParams = new HttpParams();
        if (args.limit) {
            queryParams = queryParams.set("limit", args.limit);
        }
        if (args.teamId) {
            queryParams = queryParams.set("limit", args.teamId);
        }
        return this.http.get<SortedGames>('http://localhost:8080/api/v1/sortedGames', {
            withCredentials: true,
            params: queryParams
        })
    }

    public getGamesByWeek(): Observable<CompetitionWeek[]> {
        return this.http.get<CompetitionWeek[]>('http://localhost:8080/api/v1/gamesByWeek', {
            withCredentials: true,
            params: <any>{
                timeZoneOffset: moment().utcOffset() * 60
            }
        })
    }

    public createGame(game: GameCreationInformation): Observable<GameId> {
        return this.http.post<GameId>('http://localhost:8080/api/v1/games', game, httpOptions);
    }

    public rescheduleGame(gameId: number, gameTime: GameTime): Observable<GameId> {
        return this.http.post<GameId>('http://localhost:8080/api/v1/games/' + gameId + '/reschedule', gameTime, httpOptions);
    }

    public reportResult(gameId: number, result: GameResult): Observable<null> {
        return this.http.post<null>('http://localhost:8080/api/v1/games/' + gameId + '/report', result, httpOptions)
    }

    public deleteGame(gameId: number): Observable<null> {
        return this.http.delete<null>('http://localhost:8080/api/v1/games/' + gameId, httpOptions);
    }

}
