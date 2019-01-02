import {Component} from "@angular/core";
import {ActivatedRoute} from "@angular/router";
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
    constructor(private route: ActivatedRoute,
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
}
