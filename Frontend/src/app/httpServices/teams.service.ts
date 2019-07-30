import {Injectable} from "@angular/core";
import {HttpClient} from "@angular/common/http";
import {httpOptions, httpOptionsForm} from "./http-options";
import {Observable} from "rxjs/Rx";
import {
    LoLTeamWithRosters, Team, TeamCore,
    TeamCoreWithIcon,
    TeamId,
    TeamPermissionsCore,
    TeamWithPlayers,
    TeamWithRosters
} from "../interfaces/Team";
import {Player, PlayerCore, PlayerId} from "../interfaces/Player";
import {ElmState} from "../shared/state/state.service";
import {AbstractControl, AsyncValidatorFn, ValidationErrors} from "@angular/forms";

// function sortMainRosterByPosition(team: Team) {
//     let sortedRoster = [];
//     ['top', 'jungle', 'middle', 'support', 'bottom'].forEach((role: string) => {
//         team.players.forEach((player: Player) => {
//             if(player.position.toLowerCase() == role) {
//                 sortedRoster.push(player);
//             }
//         });
//     });
//     team.players = sortedRoster;
// }

@Injectable()
export class TeamsService {
    constructor(private state: ElmState, private http: HttpClient) {
    }

    private gameSpecificCall<T>(call: (game: string) => Observable<T>): Observable<T> {
        return new Observable(observer => {
            this.state.subscribeLeague(league => {
                console.log("request with game = " + league.game);
                call(league.game).subscribe(
                    (next: T) => observer.next(next),
                    error => observer.error(error)
                );
            });
        });
    }

    public getTeamWithRosters(teamId: string): Observable<TeamWithRosters> {
        return this.gameSpecificCall<TeamWithRosters>(game => {
            let apiPrefix = '';
            if (game == 'leagueoflegends') {
                apiPrefix = '/lol'
            }
            return this.http.get<LoLTeamWithRosters>(
                'http://localhost:8080/api/v1' + apiPrefix + '/teams/' + teamId + "/withRosters", httpOptions);
        });
    }

    public createPlayer(teamId: number, player: PlayerCore): Observable<PlayerId> {
        return this.gameSpecificCall<PlayerId>(game => {
            let apiPrefix = '';
            if (game == 'leagueoflegends') {
                apiPrefix = '/lol'
            }
            return this.http.post<PlayerId>('http://localhost:8080/api/v1' + apiPrefix + '/teams/' + teamId + '/players',
                player, httpOptions)
        });
    }

    public updatePlayer(teamId: number, playerId: number, player: PlayerCore): Observable<null> {
        return this.gameSpecificCall<null>(game => {
            let apiPrefix = '';
            if (game == 'leagueoflegends') {
                apiPrefix = '/lol'
            }
            return this.http.put<null>('http://localhost:8080/api/v1' + apiPrefix + '/teams/' + teamId + '/players/' + playerId,
                player, httpOptions)
        });
    }

    public createTeamWithPlayers(team: TeamCore, b64icon: string, players: PlayerCore[]): Observable<TeamId> {
        return this.gameSpecificCall<TeamId>(game => {
            let apiPrefix = '';
            if (game == 'leagueoflegends') {
                apiPrefix = '/lol'
            }
            return this.http.post<TeamId>('http://localhost:8080/api/v1' + apiPrefix + '/teamsWithPlayers', {
                'team': team,
                'icon': b64icon,
                'players': players
            }, httpOptions)
        });
    }

    public getLeagueTeams(): Observable<TeamWithPlayers[]> {
        return this.http.get<TeamWithPlayers[]>('http://localhost:8080/api/v1/teams', httpOptions)
    }

    public getLeagueTeamsWithRosters(): Observable<TeamWithRosters[]> {
        return this.http.get<TeamWithRosters[]>('http://localhost:8080/api/v1/teamsWithRosters', httpOptions)
    }

    public createTeam(form: FormData): Observable<TeamId> {
        return this.http.post<TeamId>('http://localhost:8080/api/v1/teams', form, httpOptionsForm)
    }

    public updateTeam(teamId: number, form: FormData): Observable<null> {
        return this.http.put<null>('http://localhost:8080/api/v1/teams/' + teamId, form, httpOptionsForm)
    }

    public deleteTeam(teamId: number): Observable<null> {
        return this.http.delete<null>('http://localhost:8080/api/v1/teams/' + teamId, httpOptions)
    }

    public deletePlayer(teamId: number, playerId: number,): Observable<null> {
        return this.http.delete<null>('http://localhost:8080/api/v1/teams/' + teamId + '/players/' + playerId,
            httpOptions)
    }

    public updateTeamManagerPermissions(teamId: number, userId: number, permissions: TeamPermissionsCore) {
        return this.http.put<null>('http://localhost:8080/api/v1/teams/' + teamId + '/permissions/' + userId,
            permissions, httpOptions)
    }

    public validateTeamNameUniqueness(teamToValidateId: number): AsyncValidatorFn {
        return (c: AbstractControl): Observable<ValidationErrors> | null => {
            return new Observable(observer => {
                this.getLeagueTeams().subscribe(
                    teams => {
                        teams.forEach(team => {
                            if (team.teamId != teamToValidateId) {
                                if (team.name.toLowerCase().replace(/\s/g,'') == c.value.toLowerCase().replace(/\s/g,'')) {
                                    observer.next({'nameInUse': true});
                                    observer.complete();
                                    return;
                                }
                            }
                            observer.next(null);
                            observer.complete();
                        });
                    }, error => observer.error(error)
                );
            });
        };
    }

    public validateTagUniqueness(teamToValidateId: number): AsyncValidatorFn {
        return (c: AbstractControl): Observable<ValidationErrors> | null => {
            return new Observable(observer => {
                this.getLeagueTeams().subscribe(
                    teams => {
                        teams.forEach(team => {
                            if (team.teamId != teamToValidateId) {
                                if (team.tag.toLowerCase().replace(/\s/g,'')  == c.value.toLowerCase().replace(/\s/g,'')) {
                                    observer.next({'tagInUse': true});
                                    observer.complete();
                                    return;
                                }
                            }
                            observer.next(null);
                            observer.complete();
                        });
                    }, error => observer.error(error)
                );
            });
        };
    }

    public validateGameIdentifierUniqueness(playerToValidateId: number, localPlayers: Player[]): AsyncValidatorFn {
        return (c: AbstractControl): Observable<ValidationErrors> | null => {
            return new Observable(observer => {
                localPlayers.forEach(player => {
                    if (player.playerId != playerToValidateId) {
                        if (player.gameIdentifier.toLowerCase().replace(/\s/g,'')  == c.value.toLowerCase().replace(/\s/g,'')) {
                            observer.next({'gameIdentifierInUse': true});
                            observer.complete();
                            return;
                        }
                    }
                });

                this.getLeagueTeams().subscribe(
                    teams => {
                        teams.forEach(team => {
                            team.players.forEach(player => {
                                if (player.playerId != playerToValidateId) {
                                    if (player.gameIdentifier.toLowerCase().replace(/\s/g,'')  == c.value.toLowerCase().replace(/\s/g,'')) {
                                        observer.next({'gameIdentifierInUse': true});
                                        observer.complete();
                                        return;
                                    }
                                }
                            });
                        });
                        observer.next(null);
                        observer.complete();
                    }, error => observer.error(error)
                );
            });
        };
    }

    //
    // public getTeamManagers(): Observable<any> {
    //     return this.http.get('http://localhost:8080/api/leagues/teamManagers', httpOptions);
    // }
    //
    // public updateManagerPermissions(teamId: number, userId: number, administrator: boolean, information: boolean,
    //                                 players: boolean, reportResults: boolean) {
    //     return this.http.put('http://localhost:8080/api/teams/updatePermissions', {
    //         teamId: teamId,
    //         userId : userId,
    //         administrator: administrator,
    //         information: information,
    //         players: players,
    //         reportResults: reportResults
    //     }, httpOptions)
    // }
    //
    // public getTeamInformation(teamId: number): Observable<Object> {
    //     let url = "";
    //     switch(this.leagueService.getGame()) {
    //         case 'leagueoflegends': {
    //             url = 'http://localhost:8080/api/league-of-legends/teams/';
    //             break;
    //         }
    //         default: {
    //             url = 'http://localhost:8080/api/teams/';
    //         }
    //     }
    //     return new Observable(observer => {
    //         this.http.get(url + teamId, httpOptions).subscribe(
    //         (next: Team) => {
    //                 console.log(next);
    //                 let players = next.players;
    //                 console.log(players);
    //                 let team = next;
    //                 team.substitutes = [];
    //                 team.players = [];
    //                 if(players) {
    //                     players.forEach((player: any)=> {
    //                         if(player.mainRoster) {
    //                             team.players.push(player);
    //                         } else {
    //                             team.substitutes.push(player);
    //                         }
    //                     });
    //                 }
    //                 if(this.leagueService.getGame() == 'leagueoflegends') {
    //                     sortMainRosterByPosition(team);
    //                 }
    //                 team.id = teamId;
    //                 observer.next(team);
    //                 observer.complete();
    //             }, error => {
    //                 observer.error(error);
    //                 console.log(error);
    //             }
    //         );
    //     });
    // }
    //
    // public addPlayerInformationToTeam(team: Team): Observable<Team> {
    //     return new Observable(observer => {
    //         this.http.get('http://localhost:8080/api/teams/' + team.id, httpOptions).subscribe(
    //             (next: GtiTeam) => {
    //                 if(next.players) {
    //                     next.players.forEach(player=> {
    //                         if(player.mainRoster) {
    //                             team.players.push(player);
    //                         } else {
    //                             team.substitutes.push(player);
    //                         }
    //                     });
    //                 } else {
    //                     team.players = [];
    //                     team.substitutes = [];
    //                 }
    //
    //
    //                 observer.next(team)
    //             }, error => {
    //                 observer.error(error);
    //                 console.log(error);
    //             }
    //         );
    //     });
    // }
    //
    // public getTeamSummary(): Observable<Team[]> {
    //     return new Observable(observer => {
    //         this.http.get('http://localhost:8080/api/leagues/teamSummary', httpOptions).subscribe(
    //             (next: Team[]) => {
    //                 if(next == null) {
    //                     observer.next([]);
    //                 } else {
    //                     let teams = next;
    //                     teams.forEach(team => {
    //                         team.players = [];
    //                         team.substitutes = [];
    //                     });
    //                     observer.next(teams)
    //                 }
    //             }, error => {
    //                 console.log(error);
    //                 observer.error(error);
    //             }
    //         );
    //     });
    // }
    //
    // public addTeamInformation(games: Game[], teams: Team[]) {
    //     games.forEach(game => {
    //         teams.forEach(team => {
    //             if(game.team1Id == team.id) {
    //                 game.team1 = team;
    //             } else if (game.team2Id == team.id) {
    //                 game.team2 = team;
    //             }
    //         })
    //     })
    // }
}
