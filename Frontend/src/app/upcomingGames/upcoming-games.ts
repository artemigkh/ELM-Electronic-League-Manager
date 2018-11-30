import {Component} from '@angular/core';
import {LeagueService} from "../httpServices/leagues.service";
import {Game} from "../interfaces/Game";

@Component({
    selector: 'app-upcoming-games',
    templateUrl: './upcoming-games.html',
    styleUrls: ['./upcoming-games.scss']
})
export class UpcomingGamesComponent {
    games: Game[];

    constructor(private leagueService: LeagueService) {
        this.leagueService.getUpcomingGames().subscribe(
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
