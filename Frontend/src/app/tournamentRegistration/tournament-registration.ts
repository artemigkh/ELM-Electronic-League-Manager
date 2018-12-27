import {LeagueService} from "../httpServices/leagues.service";
import {Component} from "@angular/core";
import {FormBuilder, FormGroup, Validators} from "@angular/forms";

@Component({
    selector: 'app-tournament-registration',
    templateUrl: './tournament-registration.html',
    styleUrls: ['./tournament-registration.scss']
})
export class TournamentRegistrationComponent {
    name: string;
    tag: string;
    description: string;
    firstFormGroup: FormGroup;
    secondFormGroup: FormGroup;

    constructor(private leagueService: LeagueService, private _formBuilder: FormBuilder) {

    }

    ngOnInit() {
        this.firstFormGroup = this._formBuilder.group({
            firstCtrl: ['', Validators.required]
        });
        this.secondFormGroup = this._formBuilder.group({
            secondCtrl: ['', Validators.required]
        });
    }
}
