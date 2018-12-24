import {Component} from "@angular/core";
import {MAT_MOMENT_DATE_FORMATS, MomentDateAdapter} from '@angular/material-moment-adapter';
import {DateAdapter, MAT_DATE_FORMATS, MAT_DATE_LOCALE} from '@angular/material/core';
import {FormControl} from "@angular/forms";
import {LeagueService} from "../../httpServices/leagues.service";
import * as moment from "moment";
import {Moment} from "moment";
import {LeagueInformation} from "../../interfaces/LeagueInformation";
import {ConfirmationComponent} from "../../shared/confirmation/confirmation-component";
import {MatSnackBar} from "@angular/material";

@Component({
    selector: 'app-manage-dates',
    templateUrl: './manage-dates.html',
    styleUrls: ['./manage-dates.scss'],
})
export class ManageDatesComponent {
    constructor(public confirmation: MatSnackBar, private leagueService: LeagueService) {
        this.signupStart = new FormControl();
        this.signupEnd = new FormControl();
        this.leagueStart = new FormControl();
        this.leagueEnd = new FormControl();
        this.leagueInformation = null;
        this.leagueService.getLeagueInformation().subscribe(
            (next: LeagueInformation) => {
                this.leagueInformation = next;
                this.signupStart = new FormControl(moment.unix(next.signupStart));
                this.signupEnd = new FormControl(moment.unix(next.signupEnd));
                this.leagueStart = new FormControl(moment.unix(next.leagueStart));
                this.leagueEnd = new FormControl(moment.unix(next.leagueEnd));
            }, error => {
                console.log(error);
            }
        )
    }

    onUpdate(): void {
        this.leagueInformation.signupStart = this.signupStart.value.unix();
        this.leagueInformation.signupEnd = this.signupEnd.value.unix();
        this.leagueInformation.leagueStart = this.leagueStart.value.unix();
        this.leagueInformation.leagueEnd = this.leagueEnd.value.unix();

        this.leagueService.updateLeagueInformation(this.leagueInformation).subscribe(
            next => {
                this.confirmation.openFromComponent(ConfirmationComponent, {
                    duration: 1250,
                    panelClass: ['blue-snackbar'],
                    data: {
                        message: "League Dates Successfully Updated"
                    }
                });
            }, error => {
                console.log(error);
                this.confirmation.openFromComponent(ConfirmationComponent, {
                    duration: 2000,
                    panelClass: ['red-snackbar'],
                    data: {
                        message: "Update Failed"
                    }
                });
            }
        );
    }
}
