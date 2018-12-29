import {Injectable} from '@angular/core';
import {HttpClient, HttpHeaders} from '@angular/common/http';
import {Team} from "../interfaces/Team";
import {Observable} from "rxjs/index";
import { of } from 'rxjs';
import {Game, GameCollection} from "../interfaces/Game";
import {GtiTeam} from "./api-return-schemas/get-team-information";
import {Player} from "../interfaces/Player";
import {httpOptions} from "./http-options";
import {Id} from "./api-return-schemas/id";
import {TeamManagers} from "../interfaces/Manager";
import {LeagueInformation} from "../interfaces/LeagueInformation";



@Injectable()
export class LeagueService {
    teams: Team[];
    gameSummaryLoaded: boolean;
    completeGames: Game[];
    upcomingGames: Game[];

    constructor(private http: HttpClient) {
        this.teams = null;
        this.gameSummaryLoaded = false;
        this.completeGames = null;
        this.upcomingGames = null;
    }

    public getTeamManagers(): Observable<any> {
        return this.http.get('http://localhost:8080/api/leagues/teamManagers', httpOptions);
    }

    public setActiveLeague(leagueId: number): Observable<any> {
        return this.http.post('http://localhost:8080/api/leagues/setActiveLeague/' + leagueId, null, httpOptions);
    }

    public getTeamSummary(useCache = true): Observable<Team[]> {
        if(this.teams != null && useCache) {
            return of(this.teams);
        } else {
            return new Observable(observer => {
                this.http.get('http://localhost:8080/api/leagues/teamSummary', httpOptions).subscribe(
                    (next: Team[]) => {
                        this.teams = next;
                        this.teams.forEach(team => {
                            team.players = [];
                            team.substitutes = [];
                        });
                        observer.next(this.teams)
                    }, error => {
                        console.log(error);
                        observer.error(error);
                    }
                );
            });
        }
    }

    public addPlayerInformationToTeam(team: Team): Observable<Team> {
        return new Observable(observer => {
            this.http.get('http://localhost:8080/api/teams/' + team.id, httpOptions).subscribe(
                (next: GtiTeam) => {
                    if(next.players) {
                        next.players.forEach(player=> {
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
                    } else {
                        team.players = [];
                        team.substitutes = [];
                    }


                    observer.next(team)
                }, error => {
                    observer.error(error);
                    console.log(error);
                }
            );
        });
    }

    private addTeamInformation(games: Game[], teams: Team[]) {
        games.forEach(game => {
            teams.forEach(team => {
                if(game.team1Id == team.id) {
                    game.team1 = team;
                } else if (game.team2Id == team.id) {
                    game.team2 = team;
                }
            })
        })
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
                        this.getTeamSummary().subscribe(
                            teams => {
                                this.addTeamInformation(games, teams);

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

    public getLeagueInformation(): Observable<Object> {
        return this.http.get('http://localhost:8080/api/leagues/', httpOptions);
    }

    public updateLeagueInformation(leagueInfo: LeagueInformation) {
        return this.http.put('http://localhost:8080/api/leagues/', {
            name: leagueInfo.name,
            description: leagueInfo.description,
            publicView: leagueInfo.publicView,
            publicJoin: leagueInfo.publicJoin,
            signupStart: leagueInfo.signupStart,
            signupEnd: leagueInfo.signupEnd,
            leagueStart: leagueInfo.leagueStart,
            leagueEnd: leagueInfo.leagueEnd
        }, httpOptions)
    }
}
