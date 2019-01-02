import {Component} from "@angular/core";
import {Router} from "@angular/router";
import {LeagueService} from "../httpServices/leagues.service";
import {LeagueInformation} from "../interfaces/LeagueInformation";
import {sports} from "../shared/sports.defs";

@Component({
    selector: 'app-leagues',
    templateUrl: './leagues.html',
    styleUrls: ['./leagues.scss']
})
export class LeaguesComponent {
    leagues: LeagueInformation[];
    constructor(private router: Router,
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
}
