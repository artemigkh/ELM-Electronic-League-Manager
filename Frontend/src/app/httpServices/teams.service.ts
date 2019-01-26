import {Injectable} from "@angular/core";
import {HttpClient} from "@angular/common/http";
import {httpOptions, httpOptionsForm} from "./http-options";
import {Observable} from "rxjs/Rx";
import {GtiTeam} from "./api-return-schemas/get-team-information";
import {Player} from "../interfaces/Player";
import {Team} from "../interfaces/Team";
import {Game} from "../interfaces/Game";
import {of} from "rxjs/index";
import {LeagueService} from "./leagues.service";

@Injectable()
export class TeamsService {
    constructor(private http: HttpClient, private leagueService: LeagueService) {}

    public createNewTeam(name: string, tag: string, description = ""): Observable<Object> {
        return this.http.post('http://localhost:8080/api/teams/', {
            name: name,
            tag: tag,
            description: description
        }, httpOptions)
    }

    public createNewTeamWithIcon(form: FormData): Observable<Object> {
        return this.http.post('http://localhost:8080/api/teams/withIcon', form, httpOptionsForm)
    }

    public updateTeam(id: number, name: string, tag: string, description = ""): Observable<Object> {
        return this.http.put('http://localhost:8080/api/teams/updateTeam/' + id, {
            name: name,
            tag: tag,
            description: description
        }, httpOptions)
    }

    public updateTeamWithIcon(teamId: number, form: FormData): Observable<Object> {
        return this.http.put('http://localhost:8080/api/teams/updateTeamWithIcon/' +teamId, form, httpOptionsForm)
    }

    public deleteTeam(id: number): Observable<Object> {
        return this.http.delete('http://localhost:8080/api/teams/removeTeam/' + id, httpOptions)
    }

    public getTeamManagers(): Observable<any> {
        return this.http.get('http://localhost:8080/api/leagues/teamManagers', httpOptions);
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
        let url = "";
        switch(this.leagueService.getGame()) {
            case 'leagueoflegends': {
                url = 'http://localhost:8080/api/league-of-legends/teams/';
                break;
            }
            default: {
                url = 'http://localhost:8080/api/teams/';
            }
        }
        return new Observable(observer => {
            this.http.get(url + teamId, httpOptions).subscribe(
            (next: Team) => {
                    console.log(next);
                    let players = next.players;
                    console.log(players);
                    let team = next;
                    team.substitutes = [];
                    team.players = [];
                    if(players) {
                        players.forEach((player: any)=> {
                            if(player.mainRoster) {
                                team.players.push(player);
                            } else {
                                team.substitutes.push(player);
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

    public addPlayerInformationToTeam(team: Team): Observable<Team> {
        return new Observable(observer => {
            this.http.get('http://localhost:8080/api/teams/' + team.id, httpOptions).subscribe(
                (next: GtiTeam) => {
                    if(next.players) {
                        next.players.forEach(player=> {
                            if(player.mainRoster) {
                                team.players.push(player);
                            } else {
                                team.substitutes.push(player);
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

    public getTeamSummary(): Observable<Team[]> {
        return new Observable(observer => {
            this.http.get('http://localhost:8080/api/leagues/teamSummary', httpOptions).subscribe(
                (next: Team[]) => {
                    if(next == null) {
                        observer.next([]);
                    } else {
                        let teams = next;
                        teams.forEach(team => {
                            team.players = [];
                            team.substitutes = [];
                        });
                        observer.next(teams)
                    }
                }, error => {
                    console.log(error);
                    observer.error(error);
                }
            );
        });
    }

    public addTeamInformation(games: Game[], teams: Team[]) {
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
}
