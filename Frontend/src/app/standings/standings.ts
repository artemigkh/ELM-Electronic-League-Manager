import {Component, ViewEncapsulation} from '@angular/core';
import {LeagueService} from '../httpServices/leagues.service';
import {Team} from "../interfaces/Team";

@Component({
    selector: 'app-standings',
    templateUrl: './standings.html',
    styleUrls: ['./standings.scss'],
    encapsulation: ViewEncapsulation.None,
})
export class StandingsComponent {
    displayedColumns: string[] = ['Position', 'Icon', 'Name', 'Score'];
    dataSource: Team[];

    constructor(private leagueService: LeagueService) {
        this.leagueService.getTeamSummary().subscribe(
            teamSummary => {
                this.dataSource = teamSummary;
            }, error => {
                console.log(error);
            });
    }
}
