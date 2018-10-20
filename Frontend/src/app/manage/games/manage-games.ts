import {Component, Inject} from "@angular/core";
import {LeagueService} from "../../httpServices/leagues.service";
import {Team} from "../../interfaces/Team";
import {Game, GameCollection} from "../../interfaces/Game";
import {MAT_DIALOG_DATA, MatDialog, MatDialogRef} from "@angular/material";

@Component({
    selector: 'app-manage-games',
    templateUrl: './manage-games.html',
    styleUrls: ['./manage-games.scss'],
})
export class ManageGamesComponent {
    teams: Team[];
    teamVisibility: {[id: number] : boolean;} = {};
    upcomingGames: Game[];
    completeGames: Game[];

    constructor(private leagueService: LeagueService, public dialog: MatDialog) {
        this.leagueService.getTeamSummary().subscribe(
            teamSummary => {
                teamSummary.forEach(team => {
                   this.teamVisibility[team.id] = true;
                });
                this.teams = teamSummary;
                console.log(this.teams);

                this.leagueService.getAllGames().subscribe(
                    (games: GameCollection) => {
                        this.upcomingGames = games.upcomingGames;
                        this.completeGames = games.completeGames;
                        console.log(games);
                    }, error => {
                        console.log(error);
                    }
                )

            }, error => {
                console.log(error);
            });
    }

    swapVisibility(id: number): void {
        this.teamVisibility[id] = !this.teamVisibility[id];
    }

    deselectAll(): void {
        this.teams.forEach(team => {
            this.teamVisibility[team.id] = false;
        });
    }

    selectAll(): void {
        this.teams.forEach(team => {
            this.teamVisibility[team.id] = true;
        });
    }

    reportGamePopup(game: Game): void {
        const dialogRef = this.dialog.open(ReportGamePopup, {
            width: '500px',
            data: game
        });
    }

}


@Component({
    selector: 'report-game-popup',
    templateUrl: 'report-game-popup.html',
    styleUrls: ['./report-game-popup.scss'],
})
export class ReportGamePopup {

    constructor(
        public dialogRef: MatDialogRef<ReportGamePopup>,
        @Inject(MAT_DIALOG_DATA) public data: Game) {}

    onNoClick(): void {
        this.dialogRef.close();
    }
}
