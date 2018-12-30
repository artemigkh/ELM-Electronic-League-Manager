import {Component, ViewEncapsulation} from '@angular/core';
import {LeagueService} from '../httpServices/leagues.service';
import {Team} from "../interfaces/Team";
import {TeamsService} from "../httpServices/teams.service";

@Component({
    selector: 'app-standings',
    templateUrl: './standings.html',
    styleUrls: ['./standings.scss'],
    encapsulation: ViewEncapsulation.None,
})
export class StandingsComponent {
    displayedColumns: string[] = ['Position', 'Icon', 'Name', 'Score'];
    teams: Team[];

    constructor(private teamsService: TeamsService) {
        this.teamsService.getTeamSummary().subscribe(
            teamSummary => {
                this.teams = teamSummary;
            }, error => {
                console.log(error);
            });
    }
}
