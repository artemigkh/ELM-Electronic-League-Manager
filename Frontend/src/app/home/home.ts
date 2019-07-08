import {Component, OnInit} from '@angular/core';
import {EmptySortedGames, SortedGames} from "../interfaces/Game";
import {TeamWithPlayers} from "../interfaces/Team";
import {GamesService} from "../httpServices/games.service";
import {TeamsService} from "../httpServices/teams.service";
import {sports} from "../shared/lookup.defs";
import {ElmState} from "../shared/state/state.service";
import {League} from "../interfaces/League";
import {NGXLogger} from "ngx-logger";

@Component({
    selector: 'app-home',
    templateUrl: './home.html',
    styleUrls: ['./home.scss']
})
export class HomeComponent implements OnInit {
    games: SortedGames;
    league: League;
    allTeams: TeamWithPlayers[];
    topTeams: TeamWithPlayers[];
    signupPeriod: boolean = false;

    constructor(private state: ElmState,
                private log: NGXLogger,
                private gamesService: GamesService,
                private teamsService: TeamsService) {
        this.games = EmptySortedGames();
        this.allTeams = [];
        this.topTeams = [];
    }

    ngOnInit(): void {
        this.state.subscribeLeague((league: League) => this.league = league);

        this.gamesService.getLeagueGames({limit: "5"}).subscribe(
            games => this.games = games,
            error => this.log.error(error)
        );

        this.teamsService.getLeagueTeams().subscribe(
            teams => {
                this.allTeams = teams;
                if (teams.length > 3) {
                    this.topTeams = teams.slice(0, 3);
                } else {
                    this.topTeams = teams;
                }
            },
            error => this.log.error(error)
        );
    }

    getGameLabel(): string {
        return sports[this.league.game];
    }
}
