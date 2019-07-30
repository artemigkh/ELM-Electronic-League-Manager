// import {Component, Inject} from "@angular/core";
// import {ManagePlayersPopup, ManagePlayersTeamComponent, PlayerData} from "../manage-players-team";
// import {LeagueOfLegendsPlayer, Player} from "../../../../interfaces/Player";
// import {Action} from "../../../actions";
// import {MAT_DIALOG_DATA, MatDialog, MatDialogRef} from "@angular/material";
// import {PlayersService} from "../../../../httpServices/players.service";
// import {Id} from "../../../../httpServices/api-return-schemas/id";
// import {LeagueService} from "../../../../httpServices/leagues.service";
// import {TeamsService} from "../../../../httpServices/teams.service";
//


import {Component, Inject, OnInit} from "@angular/core";
import {ManagePlayersTeamComponent} from "../manage-players-team";
import {LoLTeamWithRosters, TeamWithRosters} from "../../../../interfaces/Team";
import {ElmState} from "../../../../shared/state/state.service";
import {NGXLogger} from "ngx-logger";
import {TeamsService} from "../../../../httpServices/teams.service";
import {EventDisplayerService} from "../../../../shared/eventDisplayer/event-displayer.service";
import {MAT_DIALOG_DATA, MatDialog, MatDialogRef} from "@angular/material";
import {createUniqueLoLPlayer, createUniquePlayer, LoLPlayer, Player} from "../../../../interfaces/Player";
import {FormBuilder, FormGroup, Validators} from "@angular/forms";
import {Action} from "../../../actions";

const allPositions: string[] = ["top", "jungle", "middle", "bottom", "support"];
class Position {
    value: string;
    display: string;
    available: boolean;
    constructor(value: string) {
        this.value = value;
        this.display = value.charAt(0).toUpperCase() + value.slice(1);
        this.available = true;
    }
}

class LoLPlayerData {
    action: Action;
    player: LoLPlayer;
    teamId: number;
    positions: Position[];
    onSuccess: (player?: LoLPlayer) => void;
    sendToServer: boolean;
    localPlayers: Player[];
    mainRoster: boolean;
}

@Component({
    selector: 'app-manage-players-team-lol',
    templateUrl: './manage-players-team-lol.html',
    styleUrls: ['./manage-players-team-lol.scss'],
})
export class ManagePlayersTeamLeagueOfLegendsComponent extends ManagePlayersTeamComponent{
    team: LoLTeamWithRosters;
    positions: Position[];
    constructor(public state: ElmState,
                public log: NGXLogger,
                public teamsService: TeamsService,
                public eventDisplayer: EventDisplayerService,
                public dialog: MatDialog) {
        super(state, log, teamsService, eventDisplayer, dialog);
        this.positions = allPositions.map(pos => new Position(pos));
    }

    selectPosition(newPos: string, oldPos?: string) {
        this.positions.forEach(pos => {
           if (pos.value == newPos) {
               pos.available = false;
           } else if (pos.value == oldPos) {
               pos.available = true;
           }
        });
    }

    setTeam(team: TeamWithRosters) {
        if (this.sendToServer) {
            this.teamsService.getTeamWithRosters(team.teamId.toString()).subscribe(
                (team: LoLTeamWithRosters) => {
                    this.team = team;
                    this.team.mainRoster.forEach(player => {
                        this.selectPosition(player.position);
                    });
                },
                error => this.log.error(error)
            );
        } else {
            this.team = <LoLTeamWithRosters>team;
        }
    }

    getPositionIcon(player: LoLPlayer): string {
        return "assets/leagueOfLegends/" + player.position + "_Icon.png";
    }

    createLoLPlayer(mainRoster: boolean): void {
        this.dialog.open(ManageLoLPlayerPopup, {
            data: <LoLPlayerData>{
                sendToServer: this.sendToServer,
                localPlayers: this.team.mainRoster.concat(this.team.substituteRoster),
                mainRoster: mainRoster,
                action: Action.Create,
                player: createUniqueLoLPlayer(mainRoster, () => this.internalPlayerIndex--),
                teamId: this.team.teamId,
                positions: this.positions,
                onSuccess: (player => {
                    this.eventDisplayer.displaySuccess("Player Successfully Created");
                    mainRoster? this.team.mainRoster.push(player) : this.team.substituteRoster.push(player);
                    this.selectPosition(player.position);
                })
            },
            autoFocus: false, width: '500px'
        });
    }

    editLoLPlayer(player: LoLPlayer) {
        let oldPosition = player.position;
        this.dialog.open(ManageLoLPlayerPopup, {
            data: <LoLPlayerData>{
                sendToServer: this.sendToServer,
                localPlayers: this.team.mainRoster.concat(this.team.substituteRoster),
                mainRoster: player.mainRoster,
                action: Action.Edit,
                player: player,
                teamId: this.team.teamId,
                positions: this.positions,
                onSuccess: (() => {
                    this.eventDisplayer.displaySuccess("Player Successfully Updated");
                    this.selectPosition(player.position, oldPosition);
                })
            },
            autoFocus: false, width: '500px'
        });
    }

    deletePlayer(player: LoLPlayer) {
        super.deletePlayer(player);
        this.selectPosition(null, player.position);
    }

    movePlayerRoster(player: LoLPlayer, isDestinationMainRoster: boolean) {
        super.movePlayerRoster(player, isDestinationMainRoster);
        isDestinationMainRoster ? this.selectPosition(player.position) : this.selectPosition(null, player.position);
    }
}

@Component({
    selector: 'manage-players-popup-lol',
    templateUrl: 'manage-players-popup-lol.html',
    styleUrls: ['./manage-players-popup-lol.scss'],
})
export class ManageLoLPlayerPopup {
    title: string;
    playerForm: FormGroup;

    constructor(
        @Inject(MAT_DIALOG_DATA) public data: LoLPlayerData,
        public dialogRef: MatDialogRef<ManageLoLPlayerPopup>,
        private log: NGXLogger,
        private teamsService: TeamsService,
        private formBuilder: FormBuilder) {
        this.title = this.data.action == Action.Create ? "Create New Player" : "Edit Player";
        this.playerForm = this.formBuilder.group({
            'gameIdentifier': [this.data.player.gameIdentifier, [Validators.required, Validators.minLength(2), Validators.maxLength(25)],
                this.teamsService.validateGameIdentifierUniqueness(this.data.player.playerId, this.data.localPlayers)],
            'position': [this.data.player.position, [p => {
                console.log(p.value);
                return !this.data.mainRoster || (p.value != null && p.value.length > 0) ? null : {'mainRosterMustHavePosition': true}
            }]],
        })
    }

    onCancel(): void {
        this.dialogRef.close();
    }

    savePlayer(): void {
        ['gameIdentifier', 'position'].forEach(k => this.data.player[k] = this.playerForm.value[k]);
        if (this.data.sendToServer) {
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
        } else {
            if(this.data.action == Action.Create) {
                this.data.onSuccess(this.data.player);
            } else if(this.data.action == Action.Edit) {
                this.data.onSuccess();
            }
            this.dialogRef.close();
        }
    }
}
