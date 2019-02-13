import {Component, Inject} from "@angular/core";
import {LeagueService} from "../../httpServices/leagues.service";
import {Team} from "../../interfaces/Team";
import {forkJoin} from "rxjs";
import {MAT_DIALOG_DATA, MatDialog, MatDialogRef} from "@angular/material";
import {Player} from "../../interfaces/Player";
import {WarningPopup} from "../warningPopup/warning-popup";
import {ManageComponentInterface} from "../manage-component-interface";
import {Action} from "../actions";
import {PlayersService} from "../../httpServices/players.service";
import {Id} from "../../httpServices/api-return-schemas/id";
import {TeamsService} from "../../httpServices/teams.service";
import {UserService} from "../../httpServices/user.service";
import {TeamPermissions, UserPermissions} from "../../httpServices/api-return-schemas/permissions";

class PlayerData {
    title: string;
    action: Action;
    player: Player;
    caller: ManagePlayersComponent;
    teamId: number;
    mainRoster: boolean;
}

@Component({
    selector: 'app-manage-players',
    templateUrl: './manage-players.html',
    styleUrls: ['./manage-players.scss'],
})
export class ManagePlayersComponent {
    teams: Team[];

    constructor(private leagueService: LeagueService,
                private teamsService: TeamsService,
                private userService: UserService,
                public dialog: MatDialog) {
        this.teams = [];
        this.teamsService.getTeamSummary().subscribe(
            teamSummary => {
                let teams = teamSummary;
                this.userService.getUserPermissions().subscribe(
                    (next: UserPermissions) => {
                        this.teams = [];
                        teams.forEach((team: Team) => {
                            if(next.leaguePermissions.administrator || next.leaguePermissions.editTeams) {
                                this.teamsService.addPlayerInformationToTeam(team).subscribe(
                                    (next: Team)=>{
                                        this.teams.push(next);
                                    },error=>{
                                        console.log(error);
                                    }
                                );
                            } else {
                                next.teamPermissions.forEach((teamPermission: TeamPermissions) => {
                                    if(team.id == teamPermission.id &&
                                        (teamPermission.administrator || teamPermission.players)) {
                                        this.teamsService.addPlayerInformationToTeam(team).subscribe(
                                            (next: Team)=>{
                                                this.teams.push(next);
                                            },error=>{
                                                console.log(error);
                                            }
                                        );
                                    }
                                });
                            }
                        });
                    }, error => {
                        console.log(error);
                    }
                );
            }, error => {
                console.log(error);
            });
    }
}
