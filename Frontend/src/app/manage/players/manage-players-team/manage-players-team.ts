import {Component, Inject, Input} from "@angular/core";
import {LeagueService} from "../../../httpServices/leagues.service";
import {Team} from "../../../interfaces/Team";
import {forkJoin} from "rxjs";
import {MAT_DIALOG_DATA, MatDialog, MatDialogRef} from "@angular/material";
import {Player} from "../../../interfaces/Player";
import {WarningPopup} from "../../warningPopup/warning-popup";
import {ManageComponentInterface} from "../../manage-component-interface";
import {Action} from "../../actions";
import {PlayersService} from "../../../httpServices/players.service";
import {Id} from "../../../httpServices/api-return-schemas/id";
import {isUndefined} from "util";
class PlayerData {
    title: string;
    action: Action;
    player: Player;
    caller: ManagePlayersTeamComponent;
    teamId: number;
    mainRoster: boolean;
}

@Component({
    selector: 'app-manage-players-team',
    templateUrl: './manage-players-team.html',
    styleUrls: ['./manage-players-team.scss'],
})
export class ManagePlayersTeamComponent implements ManageComponentInterface {
    @Input() team: Team;

    constructor(private leagueService: LeagueService,
                private playersService: PlayersService,
                public dialog: MatDialog) {
    }

    ngOnChanges() {
        if(isUndefined(this.team)){return;}
        console.log(this.team);
        this.leagueService.addPlayerInformationToTeam(this.team).subscribe(
            (next: Team)=>{
                this.team = next;
            },error=>{
                console.log(error);
            }
        );
    }

    newPlayerPopup(teamId: number, mainRoster: boolean): void {
        const dialogRef = this.dialog.open(ManagePlayersPopup, {
            width: '500px',
            data: {
                title: "Create New Player",
                player: {
                    name: "",
                    gameIdentifier: ""
                },
                teamId: teamId,
                mainRoster: mainRoster,
                action: Action.Create,
                caller: this
            },
            autoFocus: false
        });
    }

    editPlayerPopup(player: Player, teamId: number, mainRoster: boolean): void {
        const dialogRef = this.dialog.open(ManagePlayersPopup, {
            width: '500px',
            data: {
                title: "Edit Player Information",
                player: player,
                teamId: teamId,
                mainRoster: mainRoster,
                action: Action.Edit,
                caller: this
            },
            autoFocus: false
        });
    }

    movePlayerRole(player: Player, teamId: number, mainRoster: boolean): void {
        this.playersService.updatePlayer(
            teamId, player.id, player.name, player.gameIdentifier, mainRoster
        ).subscribe(
            next => {
                console.log("successfully updated player");
                let oldList: Player[];
                let newList: Player[];
                if(mainRoster) {
                    oldList = this.team.substitutes;
                    newList = this.team.players;
                } else {
                    oldList = this.team.players;
                    newList = this.team.substitutes;
                }

                if(oldList && newList) {
                    let index = 0;
                    let movedPlayer: Player;
                    oldList.forEach((p: Player) => {
                        if(p.id == player.id) {
                            movedPlayer = p;
                            oldList.splice(index, 1);
                        }
                        index++;
                    });
                    if(movedPlayer) {
                        newList.push(movedPlayer);
                    } else {
                        console.log("something went wrong");
                        return;
                    }
                } else {
                    console.log("something went wrong");
                    return;
                }
            }, error => {
                console.log("error during player updated");
                console.log(error);
            }
        );
    }

    warningPopup(player: Player, teamId: number): void {
        const dialogRef = this.dialog.open(WarningPopup, {
            width: '500px',
            data: {
                entity: "player",
                name: player.name + " (" + player.gameIdentifier + ")",
                caller: this,
                Id: teamId,
                Id2: player.id
            },
            autoFocus: false
        });
    }

    notifyCreateSuccess(id: number, teamId: number, name: string, gameIdentifier: string, mainRoster: boolean): void {
        let newPlayer = {
            id: id,
            name: name,
            gameIdentifier: gameIdentifier
        };
        if(mainRoster) {
            this.team.players.push(newPlayer);
        } else {
            this.team.substitutes.push(newPlayer);
        }
        console.log("component create success");
    }

    notifyUpdateSuccess(id: number, teamId: number, name: string, gameIdentifier: string, mainRoster: boolean): void {
        let playerList: Player[];
        if(mainRoster) {
            playerList = this.team.players;
        } else {
            playerList = this.team.substitutes;
        }
        if(playerList) {
            playerList.forEach((player: Player) => {
                if(player.id == id) {
                    player.name = name;
                    player.gameIdentifier = gameIdentifier;
                }
            });
            console.log("component update success");
        } else {
            console.log("something went wrong");
            return;
        }
    }

    notifyDelete(id: number, id2: number): void {
        console.log("component delete start", id, id2);
        this.playersService.removePlayer(id, id2).subscribe(
            next => {
                console.log("removed player with id ", id2);
                let index = 0;
                this.team.players.forEach((player: Player) => {
                    if(player.id == id2) {
                        this.team.players.splice(index, 1);
                        return;
                    }
                    index++;
                });
                index = 0;
                this.team.substitutes.forEach((player: Player) => {
                    if(player.id == id2) {
                        this.team.substitutes.splice(index, 1);
                        return;
                    }
                    index++;
                });
            }, error => {
                console.log('failed to delete player, reason:', error);
            }
        )
    }
}


@Component({
    selector: 'manage-players-popup',
    templateUrl: 'manage-players-popup.html',
    styleUrls: ['./manage-players-popup.scss'],
})
export class ManagePlayersPopup {
    action: Action;
    name: string;
    gameIdentifier: string;

    constructor(
        public dialogRef: MatDialogRef<ManagePlayersPopup>,
        @Inject(MAT_DIALOG_DATA) public data: PlayerData,
        private playersService: PlayersService) {
        this.action = data.action;
        this.name = data.player.name;
        this.gameIdentifier = data.player.gameIdentifier;
    }

    OnCancel(): void {
        this.dialogRef.close();
    }

    OnConfirm(): void {
        console.log("confirm called");
        console.log("action is", this.action);
        if(this.action == Action.Create) {
            this.playersService.addPlayer(
                this.data.teamId, this.name, this.gameIdentifier, this.data.mainRoster
            ).subscribe(
                (next: Id) => {
                    console.log("successfully added player");
                    this.data.caller.notifyCreateSuccess(
                        next.id, this.data.teamId, this.name, this.gameIdentifier, this.data.mainRoster
                    );
                    this.dialogRef.close();
                }, error => {
                    console.log("error during player creation");
                    console.log(error);
                    this.dialogRef.close();
                }
            );
        } else if(this.action = Action.Edit) {
            this.playersService.updatePlayer(
                this.data.teamId, this.data.player.id, this.name, this.gameIdentifier, this.data.mainRoster
            ).subscribe(
                next => {
                    console.log("successfully updated player");
                    this.data.caller.notifyUpdateSuccess(
                        this.data.player.id, this.data.teamId, this.name, this.gameIdentifier, this.data.mainRoster
                    );
                    this.dialogRef.close();
                }, error => {
                    console.log("error during player updated");
                    console.log(error);
                    this.dialogRef.close();
                }
            );
        }
    }
}

