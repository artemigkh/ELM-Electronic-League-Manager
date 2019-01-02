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
    constructor(private http: HttpClient) {}

    public setActiveLeague(leagueId: number): Observable<any> {
        return this.http.post('http://localhost:8080/api/leagues/setActiveLeague/' + leagueId, null, httpOptions);
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
        return this.http.put('http://localhost:8080/api/leagues/', {
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

}
