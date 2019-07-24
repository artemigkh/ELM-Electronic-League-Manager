import {Component, OnInit} from "@angular/core";
import {LeagueCore} from "../interfaces/League";
import {Option} from "../interfaces/UI";
import {FormBuilder, FormControl, FormGroup, Validators} from "@angular/forms";
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
export class LeagueCreationComponent implements OnInit{
    league: LeagueCore;
    physicalSports: Option[];
    eSports: Option[];

    leagueForm: FormGroup;
    timeKeys: string[] = ['signupStart', 'signupEnd', 'leagueStart', 'leagueEnd'];

    constructor(private leagueService: LeagueService,
                private router: Router,
                private eventDisplayer: EventDisplayerService,
                private _formBuilder: FormBuilder) {
        this.league = new LeagueCore();
        this.physicalSports = Object.entries(physicalSportsDef).map(o => <Option>{value: o[0], display: o[1]});
        this.eSports = Object.entries(eSportsDef).map(o => <Option>{value: o[0], display: o[1]});
        this.timeKeys.forEach(k => this[k] = new FormControl(moment.unix(this.league[k])));
    }

    onSubmit() {
        this.timeKeys.forEach(k => this.league[k] = this.leagueForm.get('dates').get(k).value.unix());
        this.league.name = this.leagueForm.get('leagueInformation').get('name').value;
        this.league.description = this.leagueForm.get('leagueInformation').get('description').value;
        console.log(this.league);
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

    ngOnInit(): void {

        this.leagueForm = this._formBuilder.group({
            leagueInformation: new FormGroup({
                name: new FormControl(this.league.name, Validators.compose([
                    Validators.required,
                    Validators.minLength(3),
                    Validators.maxLength(25)])),
                description: new FormControl(this.league.description, Validators.maxLength(500)),
                publicView: new FormControl(true, Validators.required),
                publicJoin: new FormControl(true, Validators.required)
            }),
            dates: new FormGroup({
                signupStart: new FormControl(moment.unix(this.league.signupStart)),
                signupEnd: new FormControl(moment.unix(this.league.signupEnd)),
                leagueStart: new FormControl(moment.unix(this.league.leagueStart)),
                leagueEnd: new FormControl(moment.unix(this.league.leagueEnd))
            }, Validators.compose([
                (g: FormGroup) => {
                    return g.get('signupEnd').value.isAfter(g.get('signupStart').value) ? null : {'outOfOrder': true}
                },
                (g: FormGroup) => {
                    return g.get('leagueStart').value.isSameOrAfter(g.get('signupEnd').value) ? null : {'outOfOrder': true}
                },
                (g: FormGroup) => {
                    return g.get('leagueEnd').value.isAfter(g.get('leagueStart').value) ? null : {'outOfOrder': true}
                },
            ]))
        });
    }
}
