import {Injectable} from '@angular/core';
import {HttpClient, HttpHeaders} from '@angular/common/http';
import {Observable} from "rxjs/index";
import {httpOptions} from "./http-options";
import {LeagueInformation} from "../interfaces/LeagueInformation";

@Injectable()
export class LeagueService {
    constructor(private http: HttpClient) {}

    public setActiveLeague(leagueId: number): Observable<any> {
        return this.http.post('http://localhost:8080/api/leagues/setActiveLeague/' + leagueId, null, httpOptions);
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
