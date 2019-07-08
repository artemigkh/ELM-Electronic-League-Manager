import {Component, Inject, OnInit} from "@angular/core";
import {
    getStartMoment,
    updateFromMoments,
    SchedulingParameters,
    WeeklyAvailability,
} from "../../interfaces/Availability";
import {NGXLogger} from "ngx-logger";
import {LeagueService} from "../../httpServices/leagues.service";
import {daysOfWeekDef, physicalSportsDef, tournamentTypesDef} from "../../shared/lookup.defs";
import {Option} from "../../interfaces/UI";
import {Action} from "../actions";
import {MAT_DIALOG_DATA, MatDialog, MatDialogRef} from "@angular/material";
import {Moment} from "moment";
import * as moment from "moment";
import {EventDisplayerService} from "../../shared/eventDisplayer/event-displayer.service";
import {League, Week} from "../../interfaces/League";
import {ElmState} from "../../shared/state/state.service";
import {WarningPopup, WarningPopupData} from "../warningPopup/warning-popup";
import {Game, GameCore, GameCreationInformation} from "../../interfaces/Game";
import {forkJoin} from "rxjs";
import {GamesService} from "../../httpServices/games.service";

class AvailabilityData {
    action: Action;
    availability: WeeklyAvailability;
    onSuccess: (availability?: WeeklyAvailability) => void;
}

@Component({
    selector: 'app-manage-schedule',
    templateUrl: './manage-schedule.html',
    styleUrls: ['./manage-schedule.scss'],
})
export class ManageScheduleComponent implements OnInit {
    availabilities: WeeklyAvailability[];
    schedulingParameters: SchedulingParameters;
    tournamentTypes: Option[];
    league: League;
    tentativeSchedule: GameCore[];
    weeks: Week[];

    constructor(private log: NGXLogger,
                private state: ElmState,
                private leagueService: LeagueService,
                private gamesService: GamesService,
                private eventDisplayer: EventDisplayerService,
                private dialog: MatDialog) {
        this.schedulingParameters = new SchedulingParameters();
        this.tentativeSchedule = [];
        this.weeks = [];
    }

    ngOnInit(): void {
        this.tournamentTypes = Object.entries(tournamentTypesDef).map(o => <Option>{value: o[0], display: o[1]});
        this.state.subscribeLeague(league => this.league = league);
        this.leagueService.getWeeklyAvailabilities().subscribe(
            availabilities => this.availabilities = availabilities,
            error => this.log.error(error)
        );
    }

    weekdayDisplay(weekday: string): string {
        return weekday[0].toLocaleUpperCase() + weekday.substr(1);
    }

    timeDisplay(a: WeeklyAvailability): string {
        let start = getStartMoment(a);
        return start.format("h:mm a") + " to " + start.clone().add(a.duration, 'm').format("h:mm a");
    }

    createAvailability() {
        this.dialog.open(ManageAvailabilityPopup, {
            data: <AvailabilityData>{
                action: Action.Create,
                availability: new WeeklyAvailability(this.league),
                onSuccess: availability => {
                    this.eventDisplayer.displaySuccess("Availability Successfully Created");
                    this.availabilities.push(availability);
                    console.log(this.availabilities);
                }
            },
            autoFocus: false, width: '500px'
        });
    }

    editAvailability(availability: WeeklyAvailability) {
        this.dialog.open(ManageAvailabilityPopup, {
            data: <AvailabilityData>{
                action: Action.Edit,
                availability: availability,
                onSuccess: () => this.eventDisplayer.displaySuccess("Availability Successfully Updated")
            },
            autoFocus: false, width: '500px'
        });
    }

    _deleteAvailability(availabilityId: number) {
        this.leagueService.deleteWeeklyAvailability(availabilityId).subscribe(
            () => {
                this.eventDisplayer.displaySuccess("Availability Successfully Deleted");
                this.availabilities = this.availabilities.filter(a => a.availabilityId != availabilityId);
            }, error => this.eventDisplayer.displayError(error)
        );
    }

    deleteAvailability(availability: WeeklyAvailability) {
        this.dialog.open(WarningPopup, {
            data: <WarningPopupData>{
                entity: "availability",
                name: this.timeDisplay(availability),
                onAccept: () => this._deleteAvailability(availability.availabilityId)
            },
            autoFocus: false, width: '500px'
        });
    }

    _generateSchedule(tentativeSchedule: GameCore[]) {
        this.tentativeSchedule = tentativeSchedule;
        if (tentativeSchedule.length > 0) {
            this.weeks.push(new Week(tentativeSchedule[0].gameTime));
            tentativeSchedule.forEach(game => {
                if (moment.unix(game.gameTime).isAfter(this.weeks[this.weeks.length - 1].end)) {
                    this.weeks.push(new Week(game.gameTime));
                }
                this.weeks[this.weeks.length - 1].games.push(new Game(game));
            })
        }
    }

    generateSchedule() {
        this.leagueService.generateSchedule(this.schedulingParameters).subscribe(
            tentativeSchedule => this._generateSchedule(tentativeSchedule),
            error => this.eventDisplayer.displayError(error)
        );
    }

    acceptSchedule() {
        forkJoin(this.tentativeSchedule.map(game => {
            return this.gamesService.createGame(<GameCreationInformation>{
                team1Id: game.team1.teamId,
                team2Id: game.team2.teamId,
                gameTime: game.gameTime
            });
        })).subscribe(
            () => this.eventDisplayer.displaySuccess("All games successfully scheduled"),
            error => this.log.error(error)
        );
    }
}

@Component({
    selector: 'manage-availability-popup',
    templateUrl: 'manage-availability-popup.html',
    styleUrls: ['./manage-availability-popup.scss'],
})
export class ManageAvailabilityPopup {
    title: string;
    daysOfWeek: Option[];
    start: Moment;
    end: Moment;

    constructor(
        @Inject(MAT_DIALOG_DATA) public data: AvailabilityData,
        public dialogRef: MatDialogRef<ManageAvailabilityPopup>,
        private log: NGXLogger,
        private leagueService: LeagueService) {
        this.title = this.data.action == Action.Create ? "Create New Availability" : "Edit Availability";
        this.daysOfWeek = Object.values(daysOfWeekDef).map(o => <Option>{value: o.toLowerCase(), display: o});
        this.start = getStartMoment(this.data.availability);
        this.end = this.start.clone().add(this.data.availability.duration, 'm');
    }

    onCancel() {
        this.dialogRef.close();
    }

    saveAvailability() {
        ['start', 'end'].forEach(k => {
            if (typeof this[k] == "string") {
                this[k] = moment(this[k], "hh-mm a").seconds(0).milliseconds(0);
            }
        });
        updateFromMoments(this.data.availability, this.start, this.end);

        if (this.data.action == Action.Create) {
            this.leagueService.createWeeklyAvailability(this.data.availability).subscribe(
                res => {
                    this.data.availability.availabilityId = res.availabilityId;
                    this.data.onSuccess(this.data.availability);
                    this.dialogRef.close();
                }, error => {
                    this.log.error(error);
                    this.dialogRef.close();
                }
            );
        } else {
            this.leagueService.updateWeeklyAvailability(this.data.availability.availabilityId, this.data.availability).subscribe(
                () => {
                    this.data.onSuccess();
                    this.dialogRef.close();
                }, error => {
                    this.log.error(error);
                    this.dialogRef.close();
                }
            );
        }
    }
}
