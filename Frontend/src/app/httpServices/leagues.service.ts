import {Injectable} from '@angular/core';
import {HttpClient, HttpHeaders} from '@angular/common/http';
import {Team} from "../interfaces/Team";
import {Observable} from "rxjs/index";
import { of } from 'rxjs';
import {Game, GameCollection} from "../interfaces/Game";
import {getTeamName} from "../shared/elm-data-utils";
import {GtiTeam} from "./api-return-schemas/get-team-information";
import {Player} from "../interfaces/Player";
import {User} from "../interfaces/User";
import {httpOptions} from "./http-options";
import {Id} from "./api-return-schemas/id";
import {NavBar} from "../shared/navbar/navbar";



@Injectable()
export class LeagueService {
    teams: Team[];
    gameSummaryLoaded: boolean;
    completeGames: Game[];
    upcomingGames: Game[];
    user: User;
    navBar: NavBar;

    constructor(private http: HttpClient) {
        this.teams = null;
        this.gameSummaryLoaded = false;
        this.completeGames = null;
        this.upcomingGames = null;
    }

    public login(email: string, password: string): Observable<User> {
        return new Observable(observer => {
            this.http.post('http://localhost:8080/login', {
                email: email,
                password: password
            }, httpOptions).subscribe(
                (next: Id) => {
                    console.log(this.navBar);
                    this.navBar.notifyLogin();
                    this.user = {
                        id: next.id,
                        email: email
                    };
                    observer.next(this.user)
                }, error => {
                    observer.error(error);
                }
            )
        })
    }

    public signup(email: string, password: string): Observable<boolean> {
        return new Observable(observer => {
            this.http.post('http://localhost:8080/api/users/', {
                email: email,
                password: password
            }, httpOptions).subscribe(
                next => {observer.next(true);},
                error => {observer.next(false);}
            )
        })
    }

    public logout(): Observable<Object> {
        return this.http.post('http://localhost:8080/logout', httpOptions);
    }

    public checkIfLoggedIn(): Observable<boolean> {
        return new Observable(observer => {
            this.http.get('http://localhost:8080/api/users/profile', httpOptions).subscribe(
                next => {observer.next(true);},
                error => {observer.next(false);}
            )
        });
    }

    public getCurrentUser() {
        return this.user;
    }

    public registerNavBar(navBar: NavBar) {
        this.navBar = navBar;
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
