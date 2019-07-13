import {Component, OnInit} from "@angular/core";
import {LeagueService} from "../../httpServices/leagues.service";
import {forkJoin} from "rxjs/index";
import {TeamsService} from "../../httpServices/teams.service";
import {TeamManager, TeamWithManagers} from "../../interfaces/Team";
import {ElmState} from "../../shared/state/state.service";
import {NGXLogger} from "ngx-logger";
import {EventDisplayerService} from "../../shared/eventDisplayer/event-displayer.service";

@Component({
    selector: 'app-manage-permissions',
    templateUrl: './manage-permissions.html',
    styleUrls: ['./manage-permissions.scss'],
})
export class ManagePermissionsComponent implements OnInit {
    teams: TeamWithManagers[];
    displayedColumns: string[] = ['userEmail', 'administrator', 'information', 'games'];

    constructor(private state: ElmState,
                private log: NGXLogger,
                private leagueService: LeagueService,
                private teamsService: TeamsService,
                private eventDisplayer: EventDisplayerService) {
    }

    ngOnInit(): void {
        this.leagueService.getTeamManagers().subscribe(
            teams => {this.teams = teams; console.log(this.teams)},
            error => this.log.error(error)
        );
    }

    updateTeamPermissions(team: TeamWithManagers): void {
        forkJoin(team.managers.map((manager: TeamManager) => {
            return this.teamsService.updateTeamManagerPermissions(team.teamId, manager.userId, manager);
        })).subscribe(
            () => this.eventDisplayer.displaySuccess("Manager permissions successfully updated"),
            error => this.eventDisplayer.displayError(error))
    }


}
