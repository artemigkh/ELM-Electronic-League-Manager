import {Component, Inject} from "@angular/core";
import {ManagePlayersPopup, ManagePlayersTeamComponent, PlayerData} from "../manage-players-team";
import {LeagueOfLegendsPlayer, Player} from "../../../../interfaces/Player";
import {Action} from "../../../actions";
import {MAT_DIALOG_DATA, MatDialog, MatDialogRef} from "@angular/material";
import {PlayersService} from "../../../../httpServices/players.service";
import {Id} from "../../../../httpServices/api-return-schemas/id";
import {LeagueService} from "../../../../httpServices/leagues.service";
import {TeamsService} from "../../../../httpServices/teams.service";

const allPositions: string[] = ["top", "jungle", "middle", "bottom", "support"];

class PlayerDataLoL extends PlayerData {
    caller: ManagePlayersTeamLeagueOfLegendsComponent;
    availablePositions: string[];
    player: LeagueOfLegendsPlayer;
}

@Component({
    selector: 'app-manage-players-team-lol',
    templateUrl: './manage-players-team-lol.html',
    styleUrls: ['./manage-players-team-lol.scss'],
})
export class ManagePlayersTeamLeagueOfLegendsComponent extends ManagePlayersTeamComponent{
    constructor(public leagueService: LeagueService,
                public playersService: PlayersService,
                public teamsService: TeamsService,
                public dialog: MatDialog){
        super(leagueService, playersService, teamsService, dialog);
    }

    getAvailablePositions(currentPlayer: LeagueOfLegendsPlayer = null): string[] {
        let availablePositions: string[] = Object.assign([], allPositions);

        this.team.players.forEach((player: LeagueOfLegendsPlayer) => {
            if(player.position != "" && player != currentPlayer) {
                availablePositions.splice(availablePositions.indexOf(player.position.toLowerCase()), 1 );
            }
        });

        return availablePositions;
    }

    newPlayerPopup(teamId: number, mainRoster: boolean): void {
        const dialogRef = this.dialog.open(ManagePlayersPopupLeagueOfLegends, {
            width: '500px',
            data: {
                title: "Create New Player",
                player: {
                    name: "",
                    gameIdentifier: "",
                    position: ""
                },
                availablePositions: this.getAvailablePositions(),
                teamId: teamId,
                mainRoster: mainRoster,
                action: Action.Create,
                caller: this
            },
            autoFocus: false
        });
    }

    editPlayerPopup(player: LeagueOfLegendsPlayer, teamId: number, mainRoster: boolean): void {
        const dialogRef = this.dialog.open(ManagePlayersPopupLeagueOfLegends, {
            width: '500px',
            data: {
                title: "Edit Player Information",
                player: player,
                availablePositions: this.getAvailablePositions(player),
                teamId: teamId,
                mainRoster: mainRoster,
                action: Action.Edit,
                caller: this
            },
            autoFocus: false
        });
    }

    notifyCreateSuccessLoL(mainRoster: boolean, player: LeagueOfLegendsPlayer): void {
        if(mainRoster) {
            this.team.players.push(player);
        } else {
            this.team.substitutes.push(player);
        }
        console.log("component create success");
    }

    movePlayerRoleLoL(player: Player, teamId: number, mainRoster: boolean): void {
        if(!mainRoster) {
            player.position = "";
        }
        this.movePlayerRole(player, teamId, mainRoster);
    }

    getPositionIcon(player: LeagueOfLegendsPlayer): string {
        return "assets/leagueOfLegends/" + player.position.toLowerCase() + "_Icon.png";
    }
}

@Component({
    selector: 'manage-players-popup-lol',
    templateUrl: 'manage-players-popup-lol.html',
    styleUrls: ['./manage-players-popup-lol.scss'],
})
export class ManagePlayersPopupLeagueOfLegends {
    action: Action;
    player: LeagueOfLegendsPlayer;
    availablePositions: string[];
    mainRoster: boolean;

    constructor(
        public dialogRef: MatDialogRef<ManagePlayersPopupLeagueOfLegends>,
        @Inject(MAT_DIALOG_DATA) public data: PlayerDataLoL,
        private playersService: PlayersService) {
        this.mainRoster = data.mainRoster;
        this.action = data.action;
        this.player = data.player;
        this.availablePositions = data.availablePositions;
    }

    setPosition(pos: string): void {
        this.player.position = pos;
    }

    OnCancel(): void {
        this.dialogRef.close();
    }

    OnConfirm(): void {
        console.log("confirm called");
        console.log("action is", this.action);
        if(this.action == Action.Create) {
            this.playersService.addPlayer(
                this.data.teamId, this.data.mainRoster, this.player
            ).subscribe(
                (next: Id) => {
                    console.log("successfully added player");
                    this.player.id = next.id;
                    this.data.caller.notifyCreateSuccessLoL(
                        this.data.mainRoster, this.player
                    );
                    this.dialogRef.close();
                }, error => {
                    console.log("error during player creation");
                    console.log(error);
                    this.dialogRef.close();
                }
            );
        } else if(this.action == Action.Edit) {
            this.playersService.updatePlayer(
                this.data.teamId, this.data.mainRoster, this.player
            ).subscribe(
                next => {
                    console.log("successfully updated player");
                    this.data.caller.notifyUpdateSuccess(
                        this.data.player.id, this.data.teamId, this.player.name, this.player.gameIdentifier, this.data.mainRoster
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
