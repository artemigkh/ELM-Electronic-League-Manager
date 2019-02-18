import {Injectable} from '@angular/core';
import {HttpClient, HttpHeaders} from '@angular/common/http';
import {Observable, of} from "rxjs/index";
import {httpOptions} from "./http-options";
import {LeagueInformation} from "../interfaces/LeagueInformation";
import {TestingConfig} from "../../../testingConfig";
import * as moment from "moment";
import {Moment} from "moment";

@Injectable()
export class LeagueService {
    game: string;
    constructor(private http: HttpClient) {
        this.game = 'leagueoflegends';
    }

    public getGame(): string {
        console.log(this.game);
        return this.game;
    }

    public setActiveLeague(leagueId: number): Observable<any> {
        return new Observable(observer => {
            this.http.post('http://localhost:8080/api/leagues/setActiveLeague/' + leagueId, null, httpOptions).subscribe(
                next => {
                    this.getLeagueInformation().subscribe(
                        (next: LeagueInformation) => {
                            console.log("setting active game to " + next.game);
                            this.game = next.game;
                            observer.next();
                        }, error => {
                            observer.error(error);
                        }
                    );
                }, error => {
                    observer.error(error);
                }
            )
        });
    }

    public joinActiveLeague(): Observable<any> {
        return this.http.post('http://localhost:8080/api/leagues/join', null, httpOptions);
    }

    public getLeagueInformation(): Observable<Object> {
        return this.http.get('http://localhost:8080/api/leagues/', httpOptions);
    }

    public isLeagueRegistrationPeriod(): Observable<boolean> {
        if(TestingConfig.testing && TestingConfig.forceRegistrationPeriod) {
            return of(true);
        }
        return new Observable(observer => {
            this.http.get('http://localhost:8080/api/leagues/', httpOptions).subscribe(
                (next: LeagueInformation) => {
                    let now = moment();
                    let start = moment.unix(next.signupStart);
                    let end = moment.unix(next.signupEnd);
                    observer.next(now.isBetween(start, end));
                }, error => {
                    observer.error(error);
                    console.log(error);
                }
            );
        });
    }

    public updateLeagueInformation(leagueInfo: LeagueInformation) {
        return new Observable(observer => {
            this.http.put('http://localhost:8080/api/leagues/', {
                name: leagueInfo.name,
                description: leagueInfo.description,
                game: leagueInfo.game,
                publicView: leagueInfo.publicView,
                publicJoin: leagueInfo.publicJoin,
                signupStart: leagueInfo.signupStart,
                signupEnd: leagueInfo.signupEnd,
                leagueStart: leagueInfo.leagueStart,
                leagueEnd: leagueInfo.leagueEnd
            }, httpOptions).subscribe(
                next => {
                    this.game = leagueInfo.game;
                    observer.next();
                }, error => {
                    observer.error(error);
                }
            );
        });
    }

    public createLeague(leagueInfo: LeagueInformation) {
        return this.http.post('http://localhost:8080/api/leagues/', {
            name: leagueInfo.name,
            description: leagueInfo.description,
            game: leagueInfo.game,
            publicView: leagueInfo.publicView,
            publicJoin: leagueInfo.publicJoin,
            signupStart: leagueInfo.signupStart,
            signupEnd: leagueInfo.signupEnd,
            leagueStart: leagueInfo.leagueStart,
            leagueEnd: leagueInfo.leagueEnd
        }, httpOptions)
    }

    public getListOfLeagues(): Observable<any>  {
        return this.http.get('http://localhost:8080/api/leagues/publicLeagues', httpOptions);
    }

    public getMarkdown(): Observable<any>  {
        return this.http.get('http://localhost:8080/api/leagues/markdown', httpOptions);
    }

    public setMarkdown(markdown: string): Observable<any>  {
        return this.http.post('http://localhost:8080/api/leagues/markdown', {
            markdown: markdown
        }, httpOptions);
    }

    public generateSchedule(tournamentType: string, roundsPerWeek: number, concurrentGameNum: number,
                            gameDuration: number): Observable<any> {
        return this.http.request('post', 'http://localhost:8080/api/scheduling/schedule', {
            body: {
                tournamentType: tournamentType,
                roundsPerWeek: roundsPerWeek,
                concurrentGameNum: concurrentGameNum,
                gameDuration: gameDuration
            }, withCredentials: true,
            headers: new HttpHeaders({
                'Content-Type': 'application/json'
            })
        });
    }

    public getSchedulingAvailabilities(): Observable<any> {
        return this.http.get('http://localhost:8080/api/scheduling/availabilities', httpOptions);
    }

    public addRecurringSchedulingAvailability(weekday: string, timezone: number, hour: number,
                                               minute: number, duration: number, constrained: boolean,
                                               start: number, end: number): Observable<any> {
        return this.http.post('http://localhost:8080/api/scheduling/recurringAvailability', {
            weekday: weekday,
            timezone: timezone,
            hour: hour,
            minute: minute,
            duration: duration,
            constrained: constrained,
            start: start,
            end: end
        }, httpOptions);
    }

    public editRecurringSchedulingAvailability(id: number, weekday: string, timezone: number, hour: number,
                                              minute: number, duration: number, constrained: boolean,
                                              start: number, end: number): Observable<any> {
        return this.http.put('http://localhost:8080/api/scheduling/recurringAvailability', {
            id: id,
            weekday: weekday,
            timezone: timezone,
            hour: hour,
            minute: minute,
            duration: duration,
            constrained: constrained,
            start: start,
            end: end
        }, httpOptions);
    }

    public deleteRecurringSchedulingAvailability(id: number) {
        return this.http.request('delete', 'http://localhost:8080/api/scheduling/recurringAvailability/' + id,
            httpOptions);
    }
}
