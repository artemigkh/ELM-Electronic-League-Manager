import {Component, OnInit} from "@angular/core";
import {FormControl} from "@angular/forms";
import {ElmState} from "../../shared/state/state.service";
import {LeagueService} from "../../httpServices/leagues.service";
import {EventDisplayerService} from "../../shared/eventDisplayer/event-displayer.service";
import * as moment from "moment";
import {League} from "../../interfaces/League";

@Component({
    selector: 'app-manage-dates',
    templateUrl: './manage-dates.html',
    styleUrls: ['./manage-dates.scss'],
})
export class ManageDatesComponent implements OnInit {
    league: League;
    signupStart: FormControl;
    signupEnd: FormControl;
    leagueStart: FormControl;
    leagueEnd: FormControl;
    timeKeys: string[] = ['signupStart', 'signupEnd', 'leagueStart', 'leagueEnd'];

    constructor(private state: ElmState,
                private leagueService: LeagueService,
                private eventDisplayer: EventDisplayerService) {
    }

    ngOnInit(): void {
        this.state.subscribeLeague(league => {
            this.league = league;
            this.timeKeys.forEach(k => this[k] = new FormControl(moment.unix(league[k])));
        });
    }

    updateAtServer() {
        this.timeKeys.forEach(k => this.league[k] = this[k].value.unix());
        this.leagueService.updateLeagueInformation(this.league).subscribe(
            success => this.eventDisplayer.displaySuccess("League Dates Successfully Updated"),
            error => this.eventDisplayer.displayError(error)
        );
    }
}
