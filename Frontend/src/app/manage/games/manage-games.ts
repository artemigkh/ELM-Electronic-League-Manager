import {Component, Inject, ViewEncapsulation} from "@angular/core";
import {LeagueService} from "../../httpServices/leagues.service";
import {Team} from "../../interfaces/Team";
import {Game, GameCollection} from "../../interfaces/Game";
import {MAT_DIALOG_DATA, MatDialog, MatDialogRef} from "@angular/material";
import {WarningPopup} from "../warningPopup/warning-popup";
import {GamesService} from "../../httpServices/games.service";
import {ManageComponentInterface} from "../manage-component-interface";
import {gameSort, gameSortReverse} from "../../shared/elm-data-utils";
import * as moment from "moment";
import {Moment} from "moment";
// import * as _moment from 'moment';
// tslint:disable-next-line:no-duplicate-imports
// import {default as _rollupMoment} from 'moment';
import {MAT_MOMENT_DATE_FORMATS, MomentDateAdapter} from '@angular/material-moment-adapter';
import {DateAdapter, MAT_DATE_FORMATS, MAT_DATE_LOCALE} from '@angular/material/core';
import {FormControl} from "@angular/forms";

class GameData {
    caller: ManageGamesComponent;
    game: Game;
}

@Component({
    selector: 'app-manage-games',
    templateUrl: './manage-games.html',
    styleUrls: ['./manage-games.scss'],
    // encapsulation: ViewEncapsulation.None
})
export class ManageGamesComponent implements ManageComponentInterface{
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
            data: {
                caller: this,
                game: game
            },
            autoFocus: false
        });
    }

    newGamePopup(): void {
        const dialogRef = this.dialog.open(ManageGamePopup, {
            width: '500px',
            data: {
                title: "Schedule New Game",
                game: {
                    gameTime: null
                }
            },
            autoFocus: false
        });
    }

    editGamePopup(game: Game): void {
        const dialogRef = this.dialog.open(ManageGamePopup, {
            width: '500px',
            data: {
                title: "Edit Game",
                game: game,
                caller: this
            },
            autoFocus: false
        });
    }

    deletePopup(game: Game): void {
        const dialogRef = this.dialog.open(WarningPopup, {
            width: '500px',
            data: {
                entity: "game",
                name: game.team1.name + " vs " + game.team2.name
            },
            autoFocus: false
        });
    }

    notifyDelete: (id: number, id2?: number) => void;

    notifyComplete(game: Game) {
        this.upcomingGames = this.upcomingGames.filter((g: Game) => {
           return g.id != game.id;
        });
        this.completeGames.push(game);
        this.completeGames.sort(gameSortReverse);
        console.log("amend notified");
    }
}


@Component({
    selector: 'report-game-popup',
    templateUrl: 'report-game-popup.html',
    styleUrls: ['./report-game-popup.scss'],
})
export class ReportGamePopup {
    game: Game;

    constructor(
        public dialogRef: MatDialogRef<ReportGamePopup>,
        private gamesService: GamesService,
        @Inject(MAT_DIALOG_DATA) public data: GameData) {
        this.game = data.game;
    }

    OnCancel(): void {
        this.dialogRef.close();
    }

    OnConfirm(): void {
        console.log("confirm called");
        console.log(this.game);
        console.log("team1Score", this.game.scoreTeam1);
        console.log("team2Score", this.game.scoreTeam2);
        let numScoreTeam1: number = Number(this.game.scoreTeam1);
        let numScoreTeam2: number = Number(this.game.scoreTeam2);
        this.game.winnerId = numScoreTeam1 > numScoreTeam2 ? this.game.team1Id : this.game.team2Id;
        this.gamesService.reportResult(
            this.game.id,
            this.game.winnerId,
            numScoreTeam1,
            numScoreTeam2
        ).subscribe(
            next => {
                console.log("reported game");
                if (!this.game.complete) {
                    this.data.caller.notifyComplete(this.game);
                }
                this.dialogRef.close();
            }, error => {
                console.log(error);
                this.dialogRef.close();
            }
        )
    }
}


@Component({
    selector: 'manage-game-popup',
    templateUrl: 'manage-game-popup.html',
    styleUrls: ['./manage-game-popup.scss'],
    providers: [
        {provide: DateAdapter, useClass: MomentDateAdapter, deps: [MAT_DATE_LOCALE]},
        {provide: MAT_DATE_FORMATS, useValue: MAT_MOMENT_DATE_FORMATS},
    ],
})
export class ManageGamePopup {
    teams: Team[];
    time: Moment;
    date: FormControl;
    constructor(
        public dialogRef: MatDialogRef<ManageGamePopup>,
        @Inject(MAT_DIALOG_DATA) public data: GameData,
        private leagueService: LeagueService) {

        if (data.game.gameTime == null) {
            this.time = null;
            this.date = new FormControl()
        } else {
            this.time = moment.unix(data.game.gameTime);
            this.date = new FormControl(this.time);
        }

        this.leagueService.getTeamSummary().subscribe(
            teamSummary => {
                this.teams = teamSummary;
            }, error => {
                console.log(error);
            });
    }

    onCancel(): void {
        this.dialogRef.close();
    }

    onConfirm(): void {
        console.log(this.time.format("dddd, MMMM Do YYYY, h:mm:ss a"));
        this.dialogRef.close();
    }
}

