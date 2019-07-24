import {Component, OnInit} from "@angular/core";
import {Router} from "@angular/router";
import {LeagueService} from "../httpServices/leagues.service";
import {sports} from "../shared/lookup.defs";
import {League} from "../interfaces/League";
import {NGXLogger} from "ngx-logger";

@Component({
    selector: 'app-leagues',
    templateUrl: './leagues.html',
    styleUrls: ['./leagues.scss']
})
export class LeaguesComponent implements OnInit{
    leagues: League[];
    constructor(private log: NGXLogger,
                private router: Router,
                private leagueService: LeagueService) {
    }

    ngOnInit(): void {
        this.leagueService.getPublicLeagues().subscribe(
            leagues => this.leagues = leagues,
            error => this.log.error(error)
        );
    }

    getGameLabel(sport: string): string {
        return sports[sport];
    }

    join(league: League) {
        // this.leagueService.setActiveLeague(league.id).subscribe(
        //     next=> {
        //         this.leagueService.joinActiveLeague().subscribe(
        //             next => {
        //                 this.router.navigate([""]);
        //                 this.confirmation.openFromComponent(ConfirmationComponent, {
        //                     duration: 1250,
        //                     panelClass: ['blue-snackbar'],
        //                     data: {
        //                         message: "Successfully joined league " + league.name
        //                     }
        //                 });
        //             }, error => {
        //                 console.log(error);
        //             }
        //         );
        //     }, error=> {
        //         console.log(error);
        //     }
        // );
    }

    view(league: League) {
        this.leagueService.setActiveLeague(league.leagueId).subscribe(
            next => this.router.navigate([""]),
            error => this.log.error(error)
        );
    }

    newLeague() {
        this.router.navigate(["leagueCreation"]);
    }


}
