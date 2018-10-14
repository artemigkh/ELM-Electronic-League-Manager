import {Component} from "@angular/core";
import {LeagueService} from "../../httpServices/leagues.service";
import {Team} from "../../interfaces/Team";
import {forkJoin} from "rxjs";

@Component({
    selector: 'app-manage-players',
    templateUrl: './manage-players.html',
    styleUrls: ['./manage-players.scss'],
})
export class ManagePlayersComponent {
    teams: Team[];

    constructor(private leagueService: LeagueService) {
        this.leagueService.getTeamSummary().subscribe(
            teamSummary => {
                this.teams = teamSummary;
                console.log(this.teams);
                forkJoin(this.teams.map(team=> {
                    return leagueService.addPlayerInformationToTeam(team);
                })).subscribe(results=> {
                    console.log(results);
                    this.teams = results;
                });


            }, error => {
                console.log(error);
            });
    }
}
