import {Component} from "@angular/core";
import {LeagueInformation} from "../interfaces/LeagueInformation";
import {esportsDef, physicalSportsDef} from "../shared/sports.defs";
import {MatSnackBar} from "@angular/material";
import {MAT_MOMENT_DATE_FORMATS, MomentDateAdapter} from '@angular/material-moment-adapter';
import {DateAdapter, MAT_DATE_FORMATS, MAT_DATE_LOCALE} from '@angular/material/core';
import {FormControl} from "@angular/forms";
import * as moment from "moment";
import {Moment} from "moment";
import {LeagueService} from "../httpServices/leagues.service";
import {ConfirmationComponent} from "../shared/confirmation/confirmation-component";
import {Router} from "@angular/router";
import {Id} from "../httpServices/api-return-schemas/id";

@Component({
    selector: 'app-league-creation',
    templateUrl: './league-creation.html',
    styleUrls: ['./league-creation.scss']
})
export class LeagueCreationComponent {
    leagueInformation: LeagueInformation;
    signupStart: FormControl;
    signupEnd: FormControl;
    leagueStart: FormControl;
    leagueEnd: FormControl;
    physicalSportsArray: {value: string; display: string}[];
    eSportsArray: {value: string; display: string}[];
    constructor(public confirmation: MatSnackBar, private leagueService: LeagueService, private router: Router,) {
        this.signupStart = new FormControl();
        this.signupEnd = new FormControl();
        this.leagueStart = new FormControl();
        this.leagueEnd = new FormControl();
        this.physicalSportsArray = [];
        Object.keys(physicalSportsDef).forEach((key: string) => {
            this.physicalSportsArray.push({
                value: key,
                display: physicalSportsDef[key]
            });
        });
        this.eSportsArray = [];
        Object.keys(esportsDef).forEach((key: string) => {
            this.eSportsArray.push({
                value: key,
                display: esportsDef[key]
            });
        });
        this.leagueInformation = {
            id: 0,
            name: "",
            description: "",
            game: null,
            publicView: false,
            publicJoin: false,
            signupStart: 0,
            signupEnd: 0,
            leagueStart: 0,
            leagueEnd: 0
        };
    }

    create() {
        this.leagueInformation.signupStart = this.signupStart.value.unix();
        this.leagueInformation.signupEnd = this.signupEnd.value.unix();
        this.leagueInformation.leagueStart = this.leagueStart.value.unix();
        this.leagueInformation.leagueEnd = this.leagueEnd.value.unix();

        this.leagueService.createLeague(this.leagueInformation).subscribe(
            (next: Id) => {
                this.leagueService.setActiveLeague(next.id).subscribe(
                    next=> {
                        this.router.navigate([""]);
                        this.confirmation.openFromComponent(ConfirmationComponent, {
                            duration: 1250,
                            panelClass: ['blue-snackbar'],
                            data: {
                                message: "League " + this.leagueInformation.name + " Successfully Created"
                            }
                        });
                    }, error=> {
                        console.log(error);
                    }
                );
            }, error => {
                console.log(error);
                this.confirmation.openFromComponent(ConfirmationComponent, {
                    duration: 2000,
                    panelClass: ['red-snackbar'],
                    data: {
                        message: "League Creation Failed"
                    }
                });
            }
        );
    }
}
