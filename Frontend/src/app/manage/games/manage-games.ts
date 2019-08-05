// import {Component, Inject, ViewEncapsulation} from "@angular/core";
// import {LeagueService} from "../../httpServices/leagues.service";
// import {Team} from "../../interfaces/Team";
// import {Game, GameCollection} from "../../interfaces/Game";
// import {MAT_DIALOG_DATA, MatDialog, MatDialogRef} from "@angular/material";
// import {WarningPopup} from "../warningPopup/warning-popup";
// import {GamesService} from "../../httpServices/games.service";
// import {ManageComponentInterface} from "../manage-component-interface";
// import {gameSort, gameSortReverse} from "../../shared/elm-data-utils";
// import * as moment from "moment";
// import {Moment} from "moment";
// // import * as _moment from 'moment';
// // tslint:disable-next-line:no-duplicate-imports
// // import {default as _rollupMoment} from 'moment';
// import {MAT_MOMENT_DATE_FORMATS, MomentDateAdapter} from '@angular/material-moment-adapter';
// import {DateAdapter, MAT_DATE_FORMATS, MAT_DATE_LOCALE} from '@angular/material/core';
// import {FormControl} from "@angular/forms";
// import {Action} from "../actions";
// import {Id} from "../../httpServices/api-return-schemas/id";
// import {TeamsService} from "../../httpServices/teams.service";
// import {TeamPermissions, UserPermissions} from "../../httpServices/api-return-schemas/permissions";
// import {UserService} from "../../httpServices/user.service";

import {Action} from "../actions";
import {Component, Inject, OnInit} from "@angular/core";
import {MAT_MOMENT_DATE_FORMATS, MomentDateAdapter} from "@angular/material-moment-adapter";
import {
    DateAdapter,
    MAT_DATE_FORMATS,
    MAT_DATE_LOCALE,
    MAT_DIALOG_DATA,
    MatDialog,
    MatDialogRef
} from "@angular/material";
import {UserWithPermissions} from "../../interfaces/User";
import {EmptySortedGames, Game, GameCreationInformation, SortedGames} from "../../interfaces/Game";
import {NGXLogger} from "ngx-logger";
import {ElmState} from "../../shared/state/state.service";
import {GamesService} from "../../httpServices/games.service";
import {EventDisplayerService} from "../../shared/eventDisplayer/event-displayer.service";
import {Team} from "../../interfaces/Team";
import {TeamsService} from "../../httpServices/teams.service";
import {Option} from "../../interfaces/UI";
import {Moment} from "moment";
import * as moment from "moment";
import {FormBuilder, FormControl, FormGroup, Validators} from "@angular/forms";
import {WarningPopup, WarningPopupData} from "../warningPopup/warning-popup";
import {League} from "../../interfaces/League";
import {TournamentCodePopup} from "./league-of-legends/tournament-code-popup";

class GameData {
    action?: Action;
    game?: Game;
    onSuccess: (game?: Game) => void;
    teams?: { [teamId: number] : Team; };
}

@Component({
    selector: 'app-manage-games',
    templateUrl: './manage-games.html',
    styleUrls: ['./manage-games.scss'],
    // encapsulation: ViewEncapsulation.None
})
export class ManageGamesComponent implements OnInit{
    user: UserWithPermissions;
    league: League;
    games: SortedGames;
    teamsMap: { [teamId: number] : Team; };
    selectedTeamId: number;

    constructor(private log: NGXLogger,
                private state: ElmState,
                private eventDisplayer: EventDisplayerService,
                private gamesService: GamesService,
                private teamsService: TeamsService,
                private dialog: MatDialog) {
        this.games = EmptySortedGames();
        this.teamsMap = {};
        this.league = null;
    }

    ngOnInit(): void {
        this.teamsService.getLeagueTeams().subscribe(
            teams => this.createTeamsMap(teams),
            error => this.log.error(error)
        );
        this.gamesService.getLeagueGames({}).subscribe(
            games => this.games = games,
            error => this.log.error(error)
        );
        this.state.subscribeUser(user => this.user = user);
        this.state.subscribeLeague(league => this.league = league);
    }

    createTeamsMap(teams: Team[]) {
        teams.forEach(team => {
            this.teamsMap[team.teamId] = team;
        });
    }

    hasPermissions(game: Game): boolean {
        return this.user.leaguePermissions.administrator || this.user.leaguePermissions.editGames || // league-wide permissions check
            this.user.teamPermissions.filter(t => t.administrator || t.games).map(t => t.teamId) // get list of editable team IDs;
                .map(teamId => game.team1.teamId == teamId || game.team2.teamId == teamId) // map to true if at least one of the teams in game correspond to this team Id
                .reduce((p, c) => p || c, false) // reduce to return true if there is at least on true in this list
    }

    isGameViewable(game: Game): boolean {
        return this.hasPermissions(game) && this.selectedTeamId == null ? true :
            game.team1.teamId == this.selectedTeamId || game.team2.teamId == this.selectedTeamId;
    }

    moveGameToCompleted(game: Game): void {
        game.complete = true;
        this.games.upcomingGames = this.games.upcomingGames.filter(g => g.gameId != game.gameId);
        this.games.completedGames.push(game);
    }

    reportGame(game: Game): void {
        this.dialog.open(ReportGamePopup, {
            data: <GameData>{
                game: game,
                onSuccess: game => {
                    this.eventDisplayer.displaySuccess("Game Successfully Reported");
                    if (!game.complete) {
                        this.moveGameToCompleted(game);
                    }
                }
            },
            autoFocus: false, width: '500px'
        });
    }

    newGame() {
        this.dialog.open(ManageGamePopup, {
            data: <GameData>{
                action: Action.Create,
                teams: this.teamsMap,
                onSuccess: game => {
                    this.eventDisplayer.displaySuccess("Game Successfully Scheduled");
                    this.games.upcomingGames.push(game);
                }
            },
            width: '500px', autoFocus: false
        });
    }

    rescheduleGame(game: Game) {
        this.dialog.open(ManageGamePopup, {
            data: <GameData>{
                action: Action.Edit,
                teams: this.teamsMap,
                game: game,
                onSuccess: () => this.eventDisplayer.displaySuccess("Game Successfully Rescheduled")
            },
            width: '500px', autoFocus: false
        });
    }

    leagueGame(): string {
        return this.league == null ? '' : this.league.game;
    }

    openTournamentCodePopup(game: Game) {
        this.dialog.open(TournamentCodePopup, {
            data: <GameData>{
                game: game
            },
            width: '750px', autoFocus: false
        });
    }

    private _deleteGame(gameId: number) {
        this.gamesService.deleteGame(gameId).subscribe(
            () => {
                this.eventDisplayer.displaySuccess("Game Successfully Deleted");
                this.games.upcomingGames = this.games.upcomingGames.filter(game => game.gameId != gameId);
            }, error => this.eventDisplayer.displayError(error)
        );
    }

    private deleteGame(game: Game) {
        this.dialog.open(WarningPopup, {
            data: <WarningPopupData>{
                entity: "game",
                name: game.team1.name + " vs " + game.team2.name,
                onAccept: () => this._deleteGame(game.gameId)
            },
            autoFocus: false, width: '500px'
        });
    }
}

//     newGamePopup(): void {
//         const dialogRef = this.dialog.open(ManageGamePopup, {
//             width: '500px',
//             data: {
//                 action: Action.Create,
//                 title: "Schedule New Game",
//                 game: {
//                     gameTime: null
//                 },
//                 caller: this
//             },
//             autoFocus: false
//         });
//     }
//
//     editGamePopup(game: Game): void {
//         const dialogRef = this.dialog.open(ManageGamePopup, {
//             width: '500px',
//             data: {
//                 action: Action.Edit,
//                 title: "Edit Game",
//                 game: game,
//                 caller: this
//             },
//             autoFocus: false
//         });
//     }
//
//     deletePopup(game: Game): void {
//         const dialogRef = this.dialog.open(WarningPopup, {
//             width: '500px',
//             data: {
//                 entity: "game",
//                 name: game.team1.name + " vs " + game.team2.name,
//                 caller: this,
//                 Id: game.id
//             },
//             autoFocus: false
//         });
//     }
//
// }
//
//
@Component({
    selector: 'report-game-popup',
    templateUrl: 'report-game-popup.html',
    styleUrls: ['./report-game-popup.scss'],
})
export class ReportGamePopup {
    constructor(
        public dialogRef: MatDialogRef<ReportGamePopup>,
        private log: NGXLogger,
        private gamesService: GamesService,
        @Inject(MAT_DIALOG_DATA) public data: GameData) {
    }

    OnCancel(): void {
        this.dialogRef.close();
    }

    OnConfirm(): void {
        if (this.data.game.scoreTeam1 > this.data.game.scoreTeam2) {
            this.data.game.winnerId = this.data.game.team1.teamId;
            this.data.game.loserId = this.data.game.team2.teamId;
        } else {
            this.data.game.winnerId = this.data.game.team2.teamId;
            this.data.game.loserId = this.data.game.team1.teamId;
        }
        this.gamesService.reportResult(this.data.game.gameId, this.data.game).subscribe(
            () => {
                this.data.onSuccess(this.data.game);
                this.dialogRef.close();
            },
            error => {
                this.log.error(error);
                this.dialogRef.close();
            }
        );
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
    title: string;
    teams: Option[];
    time: Moment;
    date: FormControl;
    creationInformation: GameCreationInformation;
    constructor(
        private eventDisplayer: EventDisplayerService,
        public dialogRef: MatDialogRef<ReportGamePopup>,
        private log: NGXLogger,
        private gamesService: GamesService,
        @Inject(MAT_DIALOG_DATA) public data: GameData) {
        this.title = this.data.action == Action.Create ? "Schedule New Game" : "Reschedule Game";
        this.teams = Object.entries(this.data.teams).map(o => <Option>{value: parseInt(o[0]), display: o[1].name});
        this.creationInformation = new GameCreationInformation();
        if (this.data.game) {
            this.time = moment.unix(this.data.game.gameTime);
            this.date = new FormControl(this.time.clone());
            this.creationInformation.team1Id = this.data.game.team1.teamId;
            this.creationInformation.team2Id = this.data.game.team2.teamId;
        } else {
            this.time = moment();
            this.date = new FormControl(moment());
        }
    }

    onCancel(): void {
        this.dialogRef.close();
    }

    requiredFieldsPresent(): boolean {
        return this.creationInformation.team1Id != null && this.creationInformation.team1Id > 0 &&
            this.creationInformation.team2Id != null && this.creationInformation.team2Id > 0 &&
            this.creationInformation.team1Id != this.creationInformation.team2Id &&
            this.time != null && this.date.value != null;
    }

    onConfirm(): void {
        if (typeof this.time == "string") {
            this.time = moment(this.time, "hh-mm a").seconds(0).milliseconds(0);
        }
        this.creationInformation.updateFromMoments(this.date.value, this.time);

        if (this.data.action == Action.Create) {
            this.gamesService.createGame(this.creationInformation).subscribe(
                res => {
                    this.data.onSuccess(<Game>{
                        gameId: res.gameId,
                        complete: false,
                        gameTime: this.creationInformation.gameTime,
                        team1: this.data.teams[this.creationInformation.team1Id],
                        team2: this.data.teams[this.creationInformation.team2Id],
                        winnerId: -1, loserId: -1, scoreTeam1: 0, scoreTeam2: 0
                    })
                }, error => {
                    this.eventDisplayer.displayError(error);
                    this.dialogRef.close();
                }
            );
            this.dialogRef.close();
        } else {
            this.data.game.gameTime = this.creationInformation.gameTime;
            this.gamesService.rescheduleGame(this.data.game.gameId,
                {gameTime: this.creationInformation.gameTime}).subscribe(
                () => {
                    this.data.onSuccess(this.data.game);
                    this.dialogRef.close();
                }, error => {
                    this.eventDisplayer.displayError(error);
                    this.dialogRef.close();
                }
            );

        }
    }
}

