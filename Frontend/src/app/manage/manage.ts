import {Component} from "@angular/core";
import {LeagueService} from "../httpServices/leagues.service";
import {UserService} from "../httpServices/user.service";
import {TeamPermissions, UserPermissions} from "../httpServices/api-return-schemas/permissions";
import {isUndefined} from "util";

@Component({
    selector: 'app-manage',
    templateUrl: './manage.html',
    styleUrls: ['./manage.scss'],
})
export class ManageComponent {
    permissions: UserPermissions;
    constructor(private userService: UserService) {
        this.userService.getUserPermissions().subscribe(
            (next: UserPermissions) => {
                console.log(next);
                this.permissions = next;
            }, error => {
                console.log(error);
            }
        );
    }

    displayLeagueControls(): boolean {
        if(isUndefined(this.permissions)) {
            return false;
        } else {
            return this.permissions.leaguePermissions.administrator;
        }
    }

    displayTeamControls(): boolean {
        if(isUndefined(this.permissions)) {
            return false;
        } else {
            let display = this.permissions.leaguePermissions.administrator ||
                this.permissions.leaguePermissions.editTeams;
            this.permissions.teamPermissions.forEach((team: TeamPermissions) => {
                display = display || team.administrator || team.information;
            });
            return display;
        }
    }

    displayPlayerControls(): boolean {
        if(isUndefined(this.permissions)) {
            return false;
        } else {
            let display = this.permissions.leaguePermissions.administrator ||
                this.permissions.leaguePermissions.editTeams;
            this.permissions.teamPermissions.forEach((team: TeamPermissions) => {
                display = display || team.administrator || team.players;
            });
            return display;
        }
    }

    displayGameControls(): boolean {
        if(isUndefined(this.permissions)) {
            return false;
        } else {
            let display = this.permissions.leaguePermissions.administrator ||
                this.permissions.leaguePermissions.editGames;
            this.permissions.teamPermissions.forEach((team: TeamPermissions) => {
                display = display || team.administrator || team.reportResults;
            });
            return display;
        }
    }
}
