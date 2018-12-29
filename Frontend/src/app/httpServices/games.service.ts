import {HttpClient} from "@angular/common/http";
import {Injectable} from "@angular/core";
import {Observable} from "rxjs/Rx";
import {httpOptions} from "./http-options";
import {Game, GameCollection} from "../interfaces/Game";
import {of} from "rxjs/index";
import {LeagueService} from "./leagues.service";
import {TeamsService} from "./teams.service";

@Injectable()
export class GamesService {
    gameSummaryLoaded: boolean;
    completeGames: Game[];
    upcomingGames: Game[];

    constructor(private http: HttpClient,
                private teamsService: TeamsService) {
        this.gameSummaryLoaded = false;
        this.completeGames = null;
        this.upcomingGames = null;
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

    private loadGameSummary(): Observable<boolean> {
        //if game summary already loaded, no need to do anything
        if(this.gameSummaryLoaded) {
            return of(true);
        } else {
            return new Observable(observer => {
                //get the game summary from the server
                this.http.get('http://localhost:8080/api/leagues/gameSummary', httpOptions).subscribe(
                    (games: Game[]) => {
                        //get the team summary from the server (or cached)
                        this.teamsService.getTeamSummary().subscribe(
                            teams => {
                                this.teamsService.addTeamInformation(games, teams);

                                this.completeGames = [];
                                this.upcomingGames = [];

                                games.forEach(game => {
                                    if(game.complete) {
                                        this.completeGames.push(game);
                                    } else {
                                        this.upcomingGames.push(game);
                                    }
                                });
                                this.upcomingGames.sort((a,b)=>
                                    (a.gameTime > b.gameTime) ? 1 :
                                        ((a.gameTime < b.gameTime) ? -1 : 0));
                                observer.next(true);
                            }, error => {
                                observer.error(error);
                                console.log(error);
                            }
                        );
                    }, error => {
                        observer.error(error);
                        console.log(error);
                    }
                )
            })
        }
    }

    public getGames(): Observable<any> {
        return this.http.get('http://localhost:8080/api/leagues/gameSummary', httpOptions)
    }

    public getCompleteGames(): Observable<Game[]> {
        if(this.gameSummaryLoaded) {
            return of(this.completeGames);
        } else {
            return new Observable(observer => {
                this.loadGameSummary().subscribe(
                    next => {
                        observer.next(this.completeGames);
                    }, error => {
                        observer.error(error);
                    }
                );
            });
        }
    }

    public getUpcomingGames(): Observable<Game[]> {
        if(this.gameSummaryLoaded) {
            return of(this.upcomingGames);
        } else {
            return new Observable(observer => {
                this.loadGameSummary().subscribe(
                    next => {
                        observer.next(this.upcomingGames);
                    }, error => {
                        observer.error(error);
                    }
                );
            });
        }
    }

    public getAllGames(): Observable<GameCollection> {
        if(this.gameSummaryLoaded) {
            return of({
                upcomingGames: this.upcomingGames,
                completeGames: this.completeGames
            });
        } else {
            return new Observable(observer => {
                this.loadGameSummary().subscribe(
                    next => {
                        observer.next({
                            upcomingGames: this.upcomingGames,
                            completeGames: this.completeGames
                        });
                    }, error => {
                        observer.error(error);
                    }
                );
            });
        }
    }
}
