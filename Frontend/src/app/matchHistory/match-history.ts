import {Component} from '@angular/core';
import {LeagueService} from "../httpServices/leagues.service";
import {Game} from "../interfaces/Game";

@Component({
    selector: 'app-match-history',
    templateUrl: './match-history.html',
    styleUrls: ['./match-history.scss']
})
export class MatchHistoryComponent {
    games: Game[];

    constructor(private leagueService: LeagueService) {
        this.leagueService.getCompleteGames().subscribe(
            gameSummary => {
                console.log('success');
                console.log(gameSummary);
                this.games = gameSummary;
            },
            error => {
                console.log('error');
                console.log(error);
            });
    }
}
