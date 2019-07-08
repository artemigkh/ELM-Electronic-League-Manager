import {Injectable, OnInit} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from "rxjs/index";
import {httpOptions} from "./http-options";
import {ElmState} from "../shared/state/state.service";
import {NGXLogger} from "ngx-logger";
import {League, LeagueCore, LeagueId, Markdown} from "../interfaces/League";
import {GameCore, SortedGames} from "../interfaces/Game";
import {
    AvailabilityId,
    SchedulingParameters,
    WeeklyAvailability,
    WeeklyAvailabilityCore
} from "../interfaces/Availability";
import {TeamWithManagers} from "../interfaces/Team";

@Injectable()
export class LeagueService implements OnInit {
    league: League;

    constructor(private state: ElmState,
                private log: NGXLogger,
                private http: HttpClient) {
    }

    ngOnInit(): void {
        this.state.subscribeLeague((league: League) => this.league = league);
    }

    public setActiveLeague(leagueId: number): Observable<League> {
        return new Observable(observer => {
            this.http.post<League>('http://localhost:8080/api/v1/leagues/setActiveLeague/' + leagueId, null, httpOptions).subscribe(
                league => {
                    this.log.debug("League successfully set", league);
                    this.state.setLeague(league);
                    observer.next(league);
                }, error => observer.error(error)
            )
        });
    }

    public joinActiveLeague(): Observable<any> {
        return this.http.post('http://localhost:8080/api/v1/leagues/join', null, httpOptions);
    }

    public getLeagueInformation(): Observable<League> {
        return new Observable(observer => {
            this.http.get<League>('http://localhost:8080/api/v1/leagues', httpOptions).subscribe(
                league => {
                    this.state.setLeague(league);
                    observer.next(league);
                }, error => observer.error(error)
            )
        });
    }

    public getTeamManagers(): Observable<TeamWithManagers[]> {
        return this.http.get<TeamWithManagers[]>('http://localhost:8080/api/v1/leagues/teamManagers', httpOptions);
    }

    public getPublicLeagues(): Observable<League[]> {
        return this.http.get<League[]>('http://localhost:8080/api/v1/leagues/publicLeagues', httpOptions);
    }

    public createLeague(league: LeagueCore): Observable<LeagueId> {
        return this.http.post<LeagueId>('http://localhost:8080/api/v1/leagues', league, httpOptions);
    }

    public updateLeagueInformation(league: LeagueCore): Observable<Object> {
        return this.http.put('http://localhost:8080/api/v1/leagues', league, httpOptions);
    }

    public getLeagueMarkdown(): Observable<Markdown> {
        return this.http.get<Markdown>('http://localhost:8080/api/v1/leagues/markdown', httpOptions);
    }

    public setLeagueMarkdown(markdown: Markdown): Observable<Object> {
        return this.http.put('http://localhost:8080/api/v1/leagues/markdown', markdown, httpOptions);
    }

    public getWeeklyAvailabilities(): Observable<WeeklyAvailability[]> {
        return this.http.get<WeeklyAvailability[]>('http://localhost:8080/api/v1/weeklyAvailabilities', httpOptions);
    }

    public createWeeklyAvailability(availability: WeeklyAvailabilityCore): Observable<AvailabilityId> {
        return this.http.post<AvailabilityId>('http://localhost:8080/api/v1/weeklyAvailabilities', availability, httpOptions);
    }

    public updateWeeklyAvailability(availabilityId: number, availability: WeeklyAvailabilityCore): Observable<null> {
        return this.http.put<null>('http://localhost:8080/api/v1/weeklyAvailabilities/' + availabilityId, availability, httpOptions);
    }

    public deleteWeeklyAvailability(availabilityId: number): Observable<null> {
        return this.http.delete<null>('http://localhost:8080/api/v1/weeklyAvailabilities/' + availabilityId, httpOptions);
    }

    public generateSchedule(schedulingParameters: SchedulingParameters): Observable<GameCore[]> {
        return this.http.post<GameCore[]>('http://localhost:8080/api/v1/schedule', schedulingParameters, httpOptions);
    }
}
