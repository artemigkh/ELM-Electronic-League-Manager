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
import {Action} from "../actions";
import {Id} from "../../httpServices/api-return-schemas/id";

class GameData {
    title: string;
    action: Action;
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

    constructor(private leagueService: LeagueService, public dialog: MatDialog,
                private gamesService: GamesService) {

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
                action: Action.Create,
                title: "Schedule New Game",
                game: {
                    gameTime: null
                },
                caller: this
            },
            autoFocus: false
        });
    }

    editGamePopup(game: Game): void {
        const dialogRef = this.dialog.open(ManageGamePopup, {
            width: '500px',
            data: {
                action: Action.Edit,
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
                name: game.team1.name + " vs " + game.team2.name,
                caller: this,
                Id: game.id
            },
            autoFocus: false
        });
    }

    notifyDelete(id: number) {
        console.log("notify deleted with id ", id);
        this.gamesService.deleteGame(id).subscribe(
            next => {
                this.upcomingGames = this.upcomingGames.filter((g: Game) => {
                    return g.id != id;
                });
                this.completeGames = this.completeGames.filter((g: Game) => {
                    return g.id != id;
                });
            }, error => {
                console.log(error);
            }
        );
    }

    notifyComplete(game: Game) {
        this.upcomingGames = this.upcomingGames.filter((g: Game) => {
           return g.id != game.id;
        });
        this.completeGames.push(game);
        this.completeGames.sort(gameSortReverse);
        console.log("amend notified");
    }

    notifyCreated(team1: Team, team2: Team, gameId: number, gameTime: number) {
        console.log("notify created");
        this.upcomingGames.push({
            id: gameId,
            gameTime: gameTime,
            complete: false,
            winnerId: -1,
            scoreTeam1: 0,
            scoreTeam2: 0,
            team1Id: team1.id,
            team2Id: team2.id,
            team1: team1,
            team2: team2
        });
        this.upcomingGames.sort(gameSort);
    }

    notifyRescheduled() {
        console.log("notify reschedule");
        this.upcomingGames.sort(gameSort);
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
    team1: Team;
    team2: Team;
    constructor(
        public dialogRef: MatDialogRef<ManageGamePopup>,
        @Inject(MAT_DIALOG_DATA) public data: GameData,
        private leagueService: LeagueService,
        private gamesService: GamesService) {

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

    setTeam1(team: Team) {
        this.team1 = team;
    }

    setTeam2(team: Team) {
        this.team2 = team;
    }

    onCancel(): void {
        this.dialogRef.close();
    }

    onConfirm(): void {
        console.log(this.time);
        let mhTime = moment(this.time, "hh-mm a");
        console.log(mhTime.format("dddd, MMMM Do YYYY, h:mm:ss a"));
        let newTime = this.date.value.clone();
        newTime.minute(mhTime.minute());
        newTime.hour(mhTime.hour());
        console.log(newTime.format("dddd, MMMM Do YYYY, h:mm:ss a"));
        this.dialogRef.close();
        if(this.data.action == Action.Create) {
            this.gamesService.createNewGame(this.team1.id, this.team2.id, newTime.unix()).subscribe(
                (next: Id) => {
                    this.data.caller.notifyCreated(this.team1, this.team2, next.id, newTime.unix());
                }, error => {
                    console.log(error);
                }
            )
        } else if (this.data.action == Action.Edit) {
            this.gamesService.rescheduleGame(this.data.game.id, newTime.unix()).subscribe(
                next => {
                    this.data.game.gameTime = newTime.unix();
                    this.data.caller.notifyRescheduled();
                }, error => {
                    console.log(error);
                }
            )
        }
    }
}

