import {Component} from '@angular/core';
import {LeagueService} from "../httpServices/leagues.service";
import {Game} from "../interfaces/Game";
import {LeagueInformation} from "../interfaces/LeagueInformation";
import {Team} from "../interfaces/Team";
import {GamesService} from "../httpServices/games.service";
import {TeamsService} from "../httpServices/teams.service";
import {sports} from "../shared/sports.defs";
import {isUndefined} from "util";

@Component({
    selector: 'app-home',
    templateUrl: './home.html',
    styleUrls: ['./home.scss']
})
export class HomeComponent {
    completeGames: Game[];
    upcomingGames: Game[];
    leagueInformation: LeagueInformation;
    teams: Team[];
    constructor(private leagueService: LeagueService,
                private gamesService: GamesService,
                private teamsService: TeamsService) {
        this.leagueInformation = {
            id: 0,
            name: "",
            description: "",
            game: "genericsport",
            publicView: false,
            publicJoin: false,
            signupStart: 0,
            signupEnd: 0,
            leagueStart: 0,
            leagueEnd: 0
        };

        this.gamesService.getAllGames().subscribe(
            gameSummary => {
                this.completeGames = gameSummary.completeGames.slice(0,5);
                this.upcomingGames = gameSummary.upcomingGames.slice(0,5);
            },
            error => {
                console.log(error);
            });

        this.leagueService.getLeagueInformation().subscribe(
            (next: LeagueInformation) => {
                this.leagueInformation = next;
            }, error => {
                console.log(error);
            }
        );

        this.teamsService.getTeamSummary().subscribe(
            teamSummary => {
                if(teamSummary.length > 3) {
                    this.teams = teamSummary.slice(0,3);
                } else {
                    this.teams = teamSummary;
                }

            }, error => {
                console.log(error);
            });
    }

    getGameLabel(): string {
        return sports[this.leagueInformation.game];
    }
}
