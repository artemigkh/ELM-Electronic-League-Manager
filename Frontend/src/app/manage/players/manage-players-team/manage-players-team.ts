import {Component, Inject} from "@angular/core";
import {ManageComponentInterface} from "../../manage-component-interface";
import {LoLTeamWithRosters, Team, TeamWithRosters} from "../../../interfaces/Team";
import {ElmState} from "../../../shared/state/state.service";
import {NGXLogger} from "ngx-logger";
import {TeamsService} from "../../../httpServices/teams.service";
import {EventDisplayerService} from "../../../shared/eventDisplayer/event-displayer.service";
import {MAT_DIALOG_DATA, MatDialog, MatDialogRef} from "@angular/material";
import {FormBuilder, FormGroup, Validators} from "@angular/forms";
import {Action} from "../../actions";
import {Player} from "../../../interfaces/Player";
import {ManageTeamPopup} from "../../teams/manage-teams";
import {WarningPopup, WarningPopupData} from "../../warningPopup/warning-popup";

class PlayerData {
    action: Action;
    player: Player;
    teamId: number;
    onSuccess: (player?: Player) => void;
}

@Component({
    selector: 'app-manage-players-team',
    templateUrl: './manage-players-team.html',
    styleUrls: ['./manage-players-team.scss'],
})
export class ManagePlayersTeamComponent implements ManageComponentInterface{
    team: TeamWithRosters;
    constructor(public state: ElmState,
                public log: NGXLogger,
                public teamsService: TeamsService,
                public eventDisplayer: EventDisplayerService,
                public dialog: MatDialog) {
    }

    setTeam(team: TeamWithRosters) {
        this.team = team;
    }

    private createPlayer(mainRoster: boolean) {
        this.dialog.open(ManagePlayerPopup, {
            data: <PlayerData>{
                action: Action.Create,
                player: new Player(mainRoster),
                teamId: this.team.teamId,
                onSuccess: (player => {
                    this.eventDisplayer.displaySuccess("Player Successfully Created");
                    mainRoster? this.team.mainRoster.push(player) : this.team.substituteRoster.push(player);
                })
            },
            autoFocus: false, width: '500px'
        });
    }

    private editPlayer(player: Player) {
        this.dialog.open(ManagePlayerPopup, {
            data: <PlayerData>{
                action: Action.Edit,
                player: player,
                teamId: this.team.teamId,
                onSuccess: () => this.eventDisplayer.displaySuccess("Player Successfully Updated")
            },
            autoFocus: false, width: '500px'
        });
    }

    removePlayerFromRoster(playerId: number, fromMainRoster: boolean) {
        let roster = fromMainRoster ? this.team.mainRoster : this.team.substituteRoster;
        roster = roster.filter(rosterPlayer => rosterPlayer.playerId != playerId);
        if (fromMainRoster) {
            this.team.mainRoster = roster;
        } else {
            this.team.substituteRoster = roster;
        }
    }

    _deletePlayer(player: Player) {
        this.teamsService.deletePlayer(this.team.teamId, player.playerId).subscribe(
            () => {
                this.eventDisplayer.displaySuccess("Player Successfully Deleted");
                this.removePlayerFromRoster(player.playerId, player.mainRoster);
            }, error => this.eventDisplayer.displayError(error)
        );
    }

    deletePlayer(player: Player) {
        this.dialog.open(WarningPopup, {
            data: <WarningPopupData>{
                entity: "player",
                name: player.name,
                onAccept: () => this._deletePlayer(player)
            },
            autoFocus: false, width: '500px'
        });
    }

    movePlayerRoster(player: Player, isDestinationMainRoster: boolean) {
        player.mainRoster = isDestinationMainRoster;
        this.teamsService.updatePlayer(this.team.teamId, player.playerId, player).subscribe(
            () => {
                this.eventDisplayer.displaySuccess("Player Successfully Moved");
                this.removePlayerFromRoster(player.playerId, !isDestinationMainRoster);
                isDestinationMainRoster? this.team.mainRoster.push(player) : this.team.substituteRoster.push(player);
            }, error => this.log.error(error)
        );
    }
}

@Component({
    selector: 'manage-players-popup',
    templateUrl: 'manage-players-popup.html',
    styleUrls: ['./manage-players-popup.scss'],
})
export class ManagePlayerPopup {
    title: string;
    playerForm: FormGroup;

    constructor(
        @Inject(MAT_DIALOG_DATA) public data: PlayerData,
        public dialogRef: MatDialogRef<ManagePlayerPopup>,
        private log: NGXLogger,
        private teamsService: TeamsService,
        private formBuilder: FormBuilder) {
        this.title = this.data.action == Action.Create ? "Create New Player" : "Edit Player";
        this.playerForm = this.formBuilder.group({
            'name': [this.data.player.name, [Validators.required, Validators.minLength(3), Validators.maxLength(25)]],
            'gameIdentifier': [this.data.player.gameIdentifier, [Validators.required, Validators.minLength(1), Validators.maxLength(25)]],
        });
    }

    onCancel(): void {
        this.dialogRef.close();
    }

    savePlayer(): void {
        ['name', 'gameIdentifier'].forEach(k => this.data.player[k] = this.playerForm.value[k]);
        if(this.data.action == Action.Create) {
            this.teamsService.createPlayer(this.data.teamId, this.data.player).subscribe(
                res => {
                    this.data.player.playerId = res.playerId;
                    this.data.onSuccess(this.data.player);
                    this.dialogRef.close();
                }, error => {
                    this.log.error(error);
                    this.dialogRef.close();
                }
            );
        } else if(this.data.action == Action.Edit) {
            this.teamsService.updatePlayer(this.data.teamId, this.data.player.playerId, this.data.player).subscribe(
                () => {
                    this.data.onSuccess();
                    this.dialogRef.close();
                }, error => {
                    this.log.error(error);
                    this.dialogRef.close();
                }
            );
        }
    }
}
