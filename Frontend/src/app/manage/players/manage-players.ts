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
                private playersService: PlayersService,
                public dialog: MatDialog) {
        this.leagueService.getTeamSummary().subscribe(
            teamSummary => {
                this.teams = teamSummary;
                console.log(this.teams);
            }, error => {
                console.log(error);
            });
    }
}
