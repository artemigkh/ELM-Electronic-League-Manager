import {Component, OnInit} from "@angular/core";
import {UserWithPermissions} from "../interfaces/User";
import {ElmState} from "../shared/state/state.service";
import {NGXLogger} from "ngx-logger";

@Component({
    selector: 'app-manage',
    templateUrl: './manage.html',
    styleUrls: ['./manage.scss'],
})
export class ManageComponent implements OnInit {
    user: UserWithPermissions;

    constructor(private state: ElmState,
                private log: NGXLogger) {
    }

    ngOnInit(): void {
        this.state.subscribeUser(user => this.user = user);
    }

    private displayLeagueControls(): boolean {
        return this.user.leaguePermissions.administrator;
    }

    private displayTeamControls(): boolean {
        return this.user.leaguePermissions.administrator ||
            this.user.leaguePermissions.editTeams ||
            this.user.teamPermissions.map(t => t.administrator || t.information).reduce((p, c) => p || c, false);
    }

    private displayGameControls(): boolean {
        return this.user.leaguePermissions.administrator ||
            this.user.leaguePermissions.editGames ||
            this.user.teamPermissions.map(t => t.administrator || t.games).reduce((p, c) => p || c, false);
    }
}
