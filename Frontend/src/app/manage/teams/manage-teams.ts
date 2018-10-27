import {Component, Inject} from "@angular/core";
import {LeagueService} from "../../httpServices/leagues.service";
import {Team} from "../../interfaces/Team";
import {MAT_DIALOG_DATA, MatDialog, MatDialogRef} from "@angular/material";
import {Game} from "../../interfaces/Game";
import {WarningPopup} from "../warningPopup/warning-popup";

@Component({
    selector: 'app-manage-teams',
    templateUrl: './manage-teams.html',
    styleUrls: ['./manage-teams.scss'],
})

export class ManageTeamsComponent {
    displayedColumns: string[] = ['team'];
    teams: Team[];

    constructor(private leagueService: LeagueService, public dialog: MatDialog) {
        this.leagueService.getTeamSummary().subscribe(
            teamSummary => {
                this.teams = teamSummary;
            }, error => {
                console.log(error);
        });
    }

    newTeamPopup(): void {
        const dialogRef = this.dialog.open(ManageTeamPopup, {
            width: '500px',
            data: {
                title: "Create New Team",
                team: {
                    name: "",
                    tag: ""
                }
            },
            autoFocus: false
        });
    }

    editTeamPopup(team: Team): void {
        const dialogRef = this.dialog.open(ManageTeamPopup, {
            width: '500px',
            data: {
                title: "Edit Team Information",
                team: team
            },
            autoFocus: false
        });
    }

    deletePopup(team: Team): void {
        const dialogRef = this.dialog.open(WarningPopup, {
            width: '500px',
            data: {
                entity: "team",
                name: team.name
            },
            autoFocus: false
        });
    }
}

@Component({
    selector: 'manage-teams-popup',
    templateUrl: 'manage-teams-popup.html',
    styleUrls: ['./manage-teams-popup.scss'],
})
export class ManageTeamPopup {
    teams: Team[];
    constructor(
        public dialogRef: MatDialogRef<ManageTeamPopup>,
        @Inject(MAT_DIALOG_DATA) public data: Team,
        private leagueService: LeagueService) {
        this.leagueService.getTeamSummary().subscribe(
            teamSummary => {
                this.teams = teamSummary;
            }, error => {
                console.log(error);
            });
    }

    onNoClick(): void {
        this.dialogRef.close();
    }
}

