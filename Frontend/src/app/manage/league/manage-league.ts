import {Component} from "@angular/core";
import {MatSnackBar} from "@angular/material";
import {ConfirmationComponent} from "../../shared/confirmation/confirmation-component";
import {LeagueService} from "../../httpServices/leagues.service";
import {LeagueInformation} from "../../interfaces/LeagueInformation";

@Component({
    selector: 'app-manage-league',
    templateUrl: './manage-league.html',
    styleUrls: ['./manage-league.scss'],
})
export class ManageLeagueComponent {
    leagueInformation: LeagueInformation;
    constructor(public confirmation: MatSnackBar, private leagueService: LeagueService) {
        this.leagueInformation = {
            id: 0,
            name: "",
            description: "",
            publicView: false,
            publicJoin: false,
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
