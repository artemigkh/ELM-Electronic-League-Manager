import {Component} from '@angular/core';
import {ActivatedRoute} from "@angular/router";
import {TeamsService} from "../httpServices/teams.service";
import {Team} from "../interfaces/Team";
import { ViewEncapsulation } from '@angular/core';
import {Game} from "../interfaces/Game";
import {doesGameHaveTeam, gameSort} from "../shared/elm-data-utils";
import {LeagueService} from "../httpServices/leagues.service";
import {GamesService} from "../httpServices/games.service";
@Component({
    selector: 'app-teams',
    templateUrl: './teams.html',
    styleUrls: ['./teams.scss'],
    encapsulation: ViewEncapsulation.None
})
export class TeamsComponent {
    team: Team;
    pastGames: Game[];
    upcomingGames: Game[];

    filterGamesByTeam(): void {
        console.log("team id: ", this.team.id);
        this.pastGames = this.pastGames.filter(doesGameHaveTeam(this.team.id));
        this.upcomingGames = this.upcomingGames.filter(doesGameHaveTeam(this.team.id));
        this.pastGames = this.pastGames.sort(gameSort);
        this.upcomingGames = this.upcomingGames.sort(gameSort);
        console.log(this.pastGames);
        console.log(this.upcomingGames);
    }

    constructor(private route: ActivatedRoute,
                private leagueService: LeagueService,
                private teamsService: TeamsService,
                private gamesService: GamesService) {
        this.route.params.subscribe(params => {
            this.teamsService.getTeamInformation(+params['id']).subscribe(
                (next: Team) => {
                    this.team = next;
                    this.team.id = +params['id'];
                    console.log('test new');
                    console.log(next);
                    this.gamesService.getAllGames().subscribe(
                        gameSummary => {
                            console.log('success');
                            console.log(gameSummary);
                            this.pastGames = gameSummary.completeGames;
                            this.upcomingGames = gameSummary.upcomingGames;
                            this.filterGamesByTeam();
                        },
                        error => {
                            console.log('error');
                            console.log(error);
                        });
                }, error => {
                    console.log(error);
                }
            )
        });
    }
}
