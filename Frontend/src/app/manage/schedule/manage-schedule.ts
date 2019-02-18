import {Component, Inject} from "@angular/core";
import {LeagueInformation} from "../../interfaces/LeagueInformation";
import {MAT_DIALOG_DATA, MatDialog, MatDialogRef, MatSnackBar} from "@angular/material";
import {LeagueService} from "../../httpServices/leagues.service";
import {availability, schedule, scheduledGame, teamInfo} from "../../httpServices/api-return-schemas/schedule";
import * as moment from "moment";
import {Moment} from "moment";
import {Game} from "../../interfaces/Game";
import {Team} from "../../interfaces/Team";
import {Player} from "../../interfaces/Player";
import {Action} from "../actions";
import {PlayersService} from "../../httpServices/players.service";
import {Id} from "../../httpServices/api-return-schemas/id";
import {
    ManagePlayersPopup,
    ManagePlayersTeamComponent,
    PlayerData
} from "../players/manage-players-team/manage-players-team";
import {WarningPopup} from "../warningPopup/warning-popup";
import {FormControl} from "@angular/forms";
import {forkJoin} from "rxjs";
import {Manager} from "../../interfaces/Manager";
import {ConfirmationComponent} from "../../shared/confirmation/confirmation-component";
import {GamesService} from "../../httpServices/games.service";

class Availability {
    id: number;
    start: Moment;
    end: Moment;
    constrained: boolean;
    constraintStart: Moment;
    constraintEnd: Moment;
}

class Week {
    start: Moment;
    end: Moment;
    games: Game[];
}

export class AvailabilityData {
    title: string;
    action: Action;
    availability: Availability;
    caller: ManageScheduleComponent;
}

@Component({
    selector: 'app-manage-schedule',
    templateUrl: './manage-schedule.html',
    styleUrls: ['./manage-schedule.scss'],
})
export class ManageScheduleComponent {
    tournamentTypes: {display:string;value:string;}[];
    tournamentType: string;
    roundsPerWeek: number;
    concurrentGameNum: number;
    gameDuration: number;

    availabilities: Availability[];
    daysOfWeek: {display:string;value:number;}[];
    dayLookup: { [id: number] : string; };

    teamLookup: { [id: number] : teamInfo; };
    tentativeSchedule: scheduledGame[];

    weeks: Week[];
    constructor(public confirmation: MatSnackBar, private leagueService: LeagueService,
                private gamesService: GamesService,
                public dialog: MatDialog) {
        this.tentativeSchedule = [];
        this.teamLookup = {};
        this.tournamentTypes = [
            {
                display: "Round Robin",
                value: "roundrobin"
            }, {
                display: "Double Round Robin",
                value: "doubleroundrobin"
            }
        ];
        this.tournamentType = 'doubleroundrobin';
        this.roundsPerWeek = 2;
        this.concurrentGameNum = 1;
        this.gameDuration = 60;

        this.availabilities = [];
        this.leagueService.getSchedulingAvailabilities().subscribe(
            (next: availability[]) => {
                if(next == null) {
                    next = [];
                }
                next.forEach((avail: availability) => {
                    let startMoment = moment().
                        milliseconds(0).
                        seconds(0).
                        minute(avail.minute).
                        hour(avail.hour).
                        day(avail.weekday).
                        utcOffset(avail.timezone / 60);
                    this.availabilities.push({
                        id: avail.id,
                        start: startMoment,
                        end: startMoment.clone().add(avail.duration, 'm'),
                        constrained: avail.constrained,
                        constraintStart: moment.unix(avail.start),
                        constraintEnd: moment.unix(avail.end)
                    });
                });
            }, error => {
                console.log(error);
            }
        );
    }

    processGamesIntoWeeks(): void {
        this.weeks = [];
        let firstGame = moment.unix(this.tentativeSchedule[0].gameTime);
        let lastGame = moment.unix(this.tentativeSchedule[this.tentativeSchedule.length-1].gameTime);

        this.weeks.push({
            start: firstGame.clone().startOf('isoWeek'),
            end: firstGame.clone().endOf('isoWeek'),
            games: []
        });
        while(!lastGame.isBetween(this.weeks[this.weeks.length-1].start, this.weeks[this.weeks.length-1].end)) {
            this.weeks.push({
                start: this.weeks[this.weeks.length-1].start.clone().add(1, 'w'),
                end: this.weeks[this.weeks.length-1].end.clone().add(1, 'w'),
                games: []
            });
        }

        let weekIndex = 0;
        this.tentativeSchedule.forEach((game: scheduledGame) => {
            while(!moment.unix(game.gameTime).
            isBetween(this.weeks[weekIndex].start, this.weeks[weekIndex].end)) {
                weekIndex++;
            }
            this.weeks[weekIndex].games.push({
                id: 0,
                gameTime: game.gameTime,
                complete: false,
                winnerId: 0,
                scoreTeam1: 0,
                scoreTeam2: 0,
                team1Id: game.team1Id,
                team2Id: game.team2Id,
                team1: {
                    id: game.team1Id,
                    name: this.teamLookup[game.team1Id].name,
                    tag: this.teamLookup[game.team1Id].tag,
                    description: "",
                    wins: 0,
                    losses: 0,
                    iconSmall: this.teamLookup[game.team1Id].iconSmall,
                    iconLarge: "",
                    players: [],
                    substitutes: [],
                    visible: true
                },
                team2: {
                    id: game.team2Id,
                    name: this.teamLookup[game.team2Id].name,
                    tag: this.teamLookup[game.team2Id].tag,
                    description: "",
                    wins: 0,
                    losses: 0,
                    iconSmall: this.teamLookup[game.team2Id].iconSmall,
                    iconLarge: "",
                    players: [],
                    substitutes: [],
                    visible: true
                },
            });
        });

        console.log(this.weeks);
    }

    generateSchedule(): void {
        this.leagueService.generateSchedule(this.tournamentType, this.roundsPerWeek,
            this.concurrentGameNum, this.gameDuration).subscribe(
            (next: schedule) => {
                next.teams.forEach((team: teamInfo) => {
                   this.teamLookup[team.id] = team;
                });
                this.tentativeSchedule = next.games;
                this.processGamesIntoWeeks();
            }, error => {
                console.log(error);
            }
        );
    }

    acceptSchedule(): void {
        forkJoin(this.tentativeSchedule.map((game: scheduledGame) => {
            return this.gamesService.createNewGame(game.team1Id, game.team2Id, game.gameTime);
        })).subscribe(_ => {
            console.log("successfully updated permissions");
            this.confirmation.openFromComponent(ConfirmationComponent, {
                duration: 1250,
                panelClass: ['blue-snackbar'],
                data: {
                    message: "Games Successfully Scheduled"
                }
            });
        }, error=>{
            console.log(error);
            this.confirmation.openFromComponent(ConfirmationComponent, {
                duration: 2000,
                panelClass: ['red-snackbar'],
                data: {
                    message: "Scheduling Games Failed"
                }
            });
        });
    }

    newAvailabilityPopup(): void {
        const dialogRef = this.dialog.open(ManageAvailabilityPopup, {
            width: '500px',
            data: {
                title: "Create New Availability",
                availability: {
                    id: 0,
                    start: null,
                    end: null,
                    constrained: false,
                    constraintStart: null,
                    constraintEnd: null,
                },
                action: Action.Create,
                caller: this
            },
            autoFocus: false
        });
    }

    editAvailabilityPopup(availability: Availability): void {
        const dialogRef = this.dialog.open(ManageAvailabilityPopup, {
            width: '500px',
            data: {
                title: "Edit Availability",
                availability: availability,
                action: Action.Edit,
                caller: this
            },
            autoFocus: false
        });
    }

    warningPopup(availability: Availability): void {
        const dialogRef = this.dialog.open(WarningPopup, {
            width: '500px',
            data: {
                entity: "availability",
                name: availability.start.format("dddd") + " " +
                    availability.start.format("h:mm a") + " to " +
                    availability.end.format("h:mm a"),
                caller: this,
                Id: availability.id,
                Id2: -1
            },
            autoFocus: false
        });
    }

    notifyDelete(id: number) {
        console.log("notify deleted with id ", id);
        this.leagueService.deleteRecurringSchedulingAvailability(id).subscribe(
            next => {
                this.availabilities = this.availabilities.filter((g: Availability) => {
                    return g.id != id;
                });
            }, error => {
                console.log(error);
            }
        );
    }

    notifyCreateSuccess(availability: Availability): void {
        this.availabilities.push(availability);
    }

    notifyEditSuccess(availability: Availability): void {
        // this.availabilities.push(availability);
    }
}

@Component({
    selector: 'manage-availability-popup',
    templateUrl: 'manage-availability-popup.html',
    styleUrls: ['./manage-availability-popup.scss'],
})
export class ManageAvailabilityPopup {
    action: Action;
    availability: Availability;
    daysOfWeek: {display:string;value:number;}[];
    dayOfWeek: number;
    start: Moment;
    end: Moment;

    constructor(
        public dialogRef: MatDialogRef<ManageAvailabilityPopup>,
        @Inject(MAT_DIALOG_DATA) public data: AvailabilityData,
        private leagueService: LeagueService) {
        this.daysOfWeek = [{
            display: "Sunday", value: 0
        }, {
            display: "Monday", value: 1
        }, {
            display: "Tuesday", value: 2
        }, {
            display: "Wednesday", value: 3
        }, {
            display: "Thursday", value: 4
        },{
            display: "Friday", value: 5
        }, {
            display: "Saturday", value: 6
        }];
        this.action = data.action;
        this.availability = data.availability;
        this.end = moment(this.availability.end);

        if (this.availability.start != null) {
            this.start = moment(this.availability.start);
            this.end = moment(this.availability.end);
            this.dayOfWeek = this.availability.start.weekday();
        } else {
            this.start = null;
            this.end = null;
        }

    }

    OnCancel(): void {
        this.dialogRef.close();
    }

    OnConfirm(): void {
        console.log(this.start);
        console.log(this.end);
        if(typeof this.start == "string") {
            this.availability.start = moment(this.start, "hh-mm a").day(this.dayOfWeek).seconds(0).milliseconds(0);
        } else {
            this.availability.start = this.start;
        }
        if(typeof this.end == "string") {
            this.availability.end = moment(this.end, "hh-mm a").day(this.dayOfWeek).seconds(0).milliseconds(0);
        } else {
            this.availability.end = this.end;
        }
        console.log(this.availability.start);
        console.log(this.availability.end);

        if(this.availability.constraintStart == null) {
            this.availability.constraintStart = moment.unix(0);
            this.availability.constraintEnd = moment.unix(0);
        }
        console.log(this.availability.start.format());
        console.log(this.availability.end.format());
        console.log(this.availability.end.diff(this.availability.start));
        console.log(moment.duration(this.availability.end.diff(this.availability.start)).asMinutes());
        if(this.action == Action.Create) {
            this.leagueService.addRecurringSchedulingAvailability(
                this.availability.start.format("dddd").toLowerCase(), this.availability.start.utcOffset() * 60,
                this.availability.start.hour(), this.availability.start.minute(),
                moment.duration(this.availability.end.diff(this.availability.start)).asMinutes(),
                this.availability.constrained, this.availability.constraintStart.unix(),
                this.availability.constraintEnd.unix()
            ).subscribe(
                (next: Id) => {
                    console.log("successfully added constraint");
                    this.availability.id = next.id;
                    this.data.caller.notifyCreateSuccess(
                        this.availability
                    );
                    this.dialogRef.close();
                }, error => {
                    console.log("error during player creation");
                    console.log(error);
                    this.dialogRef.close();
                }
            );
        } else if(this.action == Action.Edit) {
            this.leagueService.editRecurringSchedulingAvailability(this.availability.id,
                this.availability.start.format("dddd").toLowerCase(), this.availability.start.utcOffset() * 60,
                this.availability.start.hour(), this.availability.start.minute(),
                moment.duration(this.availability.end.diff(this.availability.start)).asMinutes(),
                this.availability.constrained, this.availability.constraintStart.unix(),
                this.availability.constraintEnd.unix()
            ).subscribe(
                (next: Id) => {
                    console.log("successfully added constraint");
                    this.data.caller.notifyEditSuccess(
                        this.availability
                    );
                    this.dialogRef.close();
                }, error => {
                    console.log("error during player creation");
                    console.log(error);
                    this.dialogRef.close();
                }
            );
        }
    }
}
