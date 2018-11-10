import {Component} from "@angular/core";
import {LeagueService} from "../../httpServices/leagues.service";
import {TeamManagers} from "../../interfaces/Manager";

@Component({
    selector: 'app-manage-permissions',
    templateUrl: './manage-permissions.html',
    styleUrls: ['./manage-permissions.scss'],
})
export class ManagePermissionsComponent {
    teams: TeamManagers[];
    displayedColumns: string[] = ['userEmail', 'editPermissions', 'editTeamInfo', 'editPlayers', 'reportResult'];
    constructor(private leagueService: LeagueService){
        this.leagueService.getTeamManagers().subscribe(
            (next: TeamManagers[]) => {
                this.teams = next;
                console.log(this.teams);
            }, error => {
                console.log("error getting manager permissions: ", error);
            }
        )
    }
}
