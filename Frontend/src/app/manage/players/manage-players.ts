import {Component, Inject} from "@angular/core";
import {LeagueService} from "../../httpServices/leagues.service";
import {Team} from "../../interfaces/Team";
import {forkJoin} from "rxjs";
import {MAT_DIALOG_DATA, MatDialog, MatDialogRef} from "@angular/material";
import {Player} from "../../interfaces/Player";
import {WarningPopup} from "../warningPopup/warning-popup";

@Component({
    selector: 'app-manage-players',
    templateUrl: './manage-players.html',
    styleUrls: ['./manage-players.scss'],
})
export class ManagePlayersComponent {
    teams: Team[];

    constructor(private leagueService: LeagueService, public dialog: MatDialog) {
        this.leagueService.getTeamSummary().subscribe(
            teamSummary => {
                this.teams = teamSummary;
                console.log(this.teams);
                forkJoin(this.teams.map(team=> {
                    return leagueService.addPlayerInformationToTeam(team);
                })).subscribe(results=> {
                    console.log(results);
                    this.teams = results;
                });


            }, error => {
                console.log(error);
            });
    }

    newPlayerPopup(): void {
        const dialogRef = this.dialog.open(ManagePlayersPopup, {
            width: '500px',
            data: {
                title: "Create New Player",
                player: {
                    name: "",
                    gameIdentifier: ""
                }
            },
            autoFocus: false
        });
    }

    editPlayerPopup(player: Player): void {
        const dialogRef = this.dialog.open(ManagePlayersPopup, {
            width: '500px',
            data: {
                title: "Edit Player Information",
                player: player
            },
            autoFocus: false
        });
    }

    deletePopup(player: Player): void {
        const dialogRef = this.dialog.open(WarningPopup, {
            width: '500px',
            data: {
                entity: "player",
                name: player.name + " (" + player.gameIdentifier + ")"
            },
            autoFocus: false
        });
    }
}


@Component({
    selector: 'manage-players-popup',
    templateUrl: 'manage-players-popup.html',
    styleUrls: ['./manage-players-popup.scss'],
})
export class ManagePlayersPopup {
    constructor(
        public dialogRef: MatDialogRef<ManagePlayersPopup>,
        @Inject(MAT_DIALOG_DATA) public data: Object,
        private leagueService: LeagueService) {}

    onNoClick(): void {
        this.dialogRef.close();
    }
}

