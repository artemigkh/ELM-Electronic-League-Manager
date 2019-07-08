import {Component} from "@angular/core";
import {ElmState} from "../../shared/state/state.service";
import {NGXLogger} from "ngx-logger";
import {TeamsService} from "../../httpServices/teams.service";
import {EventDisplayerService} from "../../shared/eventDisplayer/event-displayer.service";
import {MatDialog} from "@angular/material";
import {TeamWithRosters} from "../../interfaces/Team";
import {UserWithPermissions} from "../../interfaces/User";

@Component({
    selector: 'app-manage-players',
    templateUrl: './manage-players.html',
    styleUrls: ['./manage-players.scss'],
})
export class ManagePlayersComponent {
    teams: TeamWithRosters[];
    user: UserWithPermissions;

    constructor(private state: ElmState,
                private log: NGXLogger,
                private teamsService: TeamsService,
                private eventDisplayer: EventDisplayerService,
                private dialog: MatDialog) {
    }

    ngOnInit(): void {
        this.state.subscribeUser(user => {
            this.user = user;
            this.getEditableTeamList();
        });
    }

    private getEditableTeamList() {
        this.teamsService.getLeagueTeamsWithRosters().subscribe(
            teams => {
                if (this.user.leaguePermissions.administrator || this.user.leaguePermissions.editTeams) {
                    this.teams = teams;
                } else {
                    this.teams = teams.filter(team => {
                        this.user.teamPermissions
                            .filter(t => t.administrator || t.information)
                            .map(t => t.teamId)
                            .includes(team.teamId);
                    });
                }
            }, error => this.eventDisplayer.displayError(error)
        );
    }
}
