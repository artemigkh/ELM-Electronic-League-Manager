import {HttpClient} from "@angular/common/http";
import {Injectable} from "@angular/core";
import {Observable} from "rxjs/Rx";
import {httpOptions} from "./http-options";
import {Game, GameCollection} from "../interfaces/Game";
import {of} from "rxjs/index";
import {TeamsService} from "./teams.service";

@Injectable()
export class GamesService {
    constructor(private http: HttpClient,
                private teamsService: TeamsService) {}

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

    public getAllGames(): Observable<GameCollection> {
        return new Observable(observer => {
            //get the game summary from the server
            this.http.get('http://localhost:8080/api/leagues/gameSummary', httpOptions).subscribe(
                (games: Game[]) => {
                    //get the team summary from the server
                    this.teamsService.getTeamSummary().subscribe(
                        teams => {
                            this.teamsService.addTeamInformation(games, teams);

                            let completeGames = [];
                            let upcomingGames = [];

                            games.forEach(game => {
                                if(game.complete) {
                                    completeGames.push(game);
                                } else {
                                    upcomingGames.push(game);
                                }
                            });
                            upcomingGames.sort((a,b)=>
                                (a.gameTime > b.gameTime) ? 1 :
                                    ((a.gameTime < b.gameTime) ? -1 : 0));
                            observer.next({
                                upcomingGames: upcomingGames,
                                completeGames: completeGames
                            });
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
