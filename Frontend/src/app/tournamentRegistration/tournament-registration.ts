import {LeagueService} from "../httpServices/leagues.service";
import {Component} from "@angular/core";

@Component({
    selector: 'app-tournament-registration',
    templateUrl: './tournament-registration.html',
    styleUrls: ['./tournament-registration.scss']
})
export class TournamentRegistrationComponent {
    constructor(private leagueService: LeagueService) {

    }
}
