import {Component} from "@angular/core";
import {LeagueCore} from "../interfaces/League";
import {Option} from "../interfaces/UI";
import {FormControl} from "@angular/forms";
import {LeagueService} from "../httpServices/leagues.service";
import {EventDisplayerService} from "../shared/eventDisplayer/event-displayer.service";
import {eSportsDef, physicalSportsDef} from "../shared/lookup.defs";
import * as moment from "moment";
import {Router} from "@angular/router";

@Component({
    selector: 'app-league-creation',
    templateUrl: './league-creation.html',
    styleUrls: ['./league-creation.scss']
})
export class LeagueCreationComponent {
    league: LeagueCore;
    physicalSports: Option[];
    eSports: Option[];

    signupStart: FormControl;
    signupEnd: FormControl;
    leagueStart: FormControl;
    leagueEnd: FormControl;
    timeKeys: string[] = ['signupStart', 'signupEnd', 'leagueStart', 'leagueEnd'];

    constructor(private leagueService: LeagueService,
                private router: Router,
                private eventDisplayer: EventDisplayerService) {
        this.league = new LeagueCore();
        this.physicalSports = Object.entries(physicalSportsDef).map(o => <Option>{value: o[0], display: o[1]});
        this.eSports = Object.entries(eSportsDef).map(o => <Option>{value: o[0], display: o[1]});
        this.timeKeys.forEach(k => this[k] = new FormControl(moment.unix(this.league[k])));

    }

    create() {
        this.timeKeys.forEach(k => this.league[k] = this[k].value.unix());
        this.leagueService.createLeague(this.league).subscribe(
            res => this.navigateToCreatedLeague(res.leagueId),
            error => this.eventDisplayer.displayError(error)
        );
    }

    private navigateToCreatedLeague(leagueId: number) {
        this.eventDisplayer.displaySuccess("League successfully created");
        this.leagueService.setActiveLeague(leagueId).subscribe(
            () => this.router.navigate([""]),
            error => this.eventDisplayer.displayError(error)
        );
    }
}
