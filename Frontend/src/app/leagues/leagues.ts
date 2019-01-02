import {Component} from "@angular/core";
import {Router} from "@angular/router";
import {LeagueService} from "../httpServices/leagues.service";
import {LeagueInformation} from "../interfaces/LeagueInformation";
import {esportsDef, physicalSportsDef, sports} from "../shared/sports.defs";
import {MatDialog, MatSnackBar} from "@angular/material";
import {ConfirmationComponent} from "../shared/confirmation/confirmation-component";

@Component({
    selector: 'app-leagues',
    templateUrl: './leagues.html',
    styleUrls: ['./leagues.scss']
})
export class LeaguesComponent {
    leagues: LeagueInformation[];
    constructor(public confirmation: MatSnackBar,
                public dialog: MatDialog,
                private router: Router,
                private leagueService: LeagueService) {
        this.leagueService.getListOfLeagues().subscribe(
            (next: LeagueInformation[]) => {
                this.leagues = next;
            }, error => {
                console.log(error);
            }
        );
    }

    getGameLabel(sport: string): string {
        return sports[sport];
    }

    join(league: LeagueInformation) {
        this.leagueService.setActiveLeague(league.id).subscribe(
            next=> {
                this.leagueService.joinActiveLeague().subscribe(
                    next => {
                        this.router.navigate([""]);
                        this.confirmation.openFromComponent(ConfirmationComponent, {
                            duration: 1250,
                            panelClass: ['blue-snackbar'],
                            data: {
                                message: "Successfully joined league " + league.name
                            }
                        });
                    }, error => {
                        console.log(error);
                    }
                );
            }, error=> {
                console.log(error);
            }
        );
    }

    view(league: LeagueInformation) {
        this.leagueService.setActiveLeague(league.id).subscribe(
            next=> {
                this.router.navigate([""]);
            }, error=> {
                console.log(error);
            }
        );
    }

    newLeaguePopup() {
        this.router.navigate(["leagueCreation"]);
    }
}
