import {Component} from "@angular/core";
import {MatSnackBar} from "@angular/material";
import {ConfirmationComponent} from "../../shared/confirmation/confirmation-component";
import {LeagueService} from "../../httpServices/leagues.service";
import {LeagueInformation} from "../../interfaces/LeagueInformation";
import {esportsDef, physicalSportsDef} from "../../shared/sports.defs";

@Component({
    selector: 'app-manage-league',
    templateUrl: './manage-league.html',
    styleUrls: ['./manage-league.scss'],
})
export class ManageLeagueComponent {
    leagueInformation: LeagueInformation;
    physicalSportsArray: {value: string; display: string}[];
    eSportsArray: {value: string; display: string}[];
    constructor(public confirmation: MatSnackBar, private leagueService: LeagueService) {
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
            game: "genericsport",
            publicView: false,
            publicJoin: false,
            signupStart: 0,
            signupEnd: 0,
            leagueStart: 0,
            leagueEnd: 0
        };
        this.leagueService.getLeagueInformation().subscribe(
            (next: LeagueInformation) => {
                console.log(next);
                this.leagueInformation = next;
            }, error => {
                console.log(error);
            }
        );
    }
    updateAtServer() {
        this.leagueService.updateLeagueInformation(this.leagueInformation).subscribe(
            next => {
                this.confirmation.openFromComponent(ConfirmationComponent, {
                    duration: 1250,
                    panelClass: ['blue-snackbar'],
                    data: {
                        message: "League Information Successfully Updated"
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
