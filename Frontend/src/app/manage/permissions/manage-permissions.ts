import {Component} from "@angular/core";
import {LeagueService} from "../../httpServices/leagues.service";
import {Manager, TeamManagers} from "../../interfaces/Manager";
import {forkJoin} from "rxjs/index";
import {TeamsService} from "../../httpServices/teams.service";

@Component({
    selector: 'app-manage-permissions',
    templateUrl: './manage-permissions.html',
    styleUrls: ['./manage-permissions.scss'],
})
export class ManagePermissionsComponent {
    teams: TeamManagers[];
    displayedColumns: string[] = ['userEmail', 'administrator', 'information', 'players', 'reportResults'];
    constructor(private leagueService: LeagueService, private teamsService: TeamsService){
        this.leagueService.getTeamManagers().subscribe(
            (next: TeamManagers[]) => {
                this.teams = next;
                console.log(this.teams);
            }, error => {
                console.log("error getting manager permissions: ", error);
            }
        )
    }

    updateTeamPermissions(team: TeamManagers): void {
        forkJoin(team.managers.map((manager: Manager) => {
            return this.teamsService.updateManagerPermissions(team.teamId, manager.userId, manager.administrator,
                manager.information, manager.players, manager.reportResults);
        })).subscribe(_ => {console.log("successfully updated permissions")}, error=>{console.log(error)});
    }
}
