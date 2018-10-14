import {Component} from "@angular/core";
import {LeagueService} from "../../httpServices/leagues.service";
import {Team} from "../../interfaces/Team";

@Component({
    selector: 'app-manage-teams',
    templateUrl: './manage-teams.html',
    styleUrls: ['./manage-teams.scss'],
})

export class ManageTeamsComponent {
    displayedColumns: string[] = ['team'];
    teams: Team[];

    constructor(private leagueService: LeagueService) {
        this.leagueService.getTeamSummary().subscribe(
            teamSummary => {
                this.teams = teamSummary;
            }, error => {
                console.log(error);
            });
    }
}
