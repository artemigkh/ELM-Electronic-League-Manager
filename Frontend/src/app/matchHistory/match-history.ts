import {Component} from '@angular/core';
import {Game} from "../interfaces/Game";
import {GamesService} from "../httpServices/games.service";

@Component({
    selector: 'app-match-history',
    templateUrl: './match-history.html',
    styleUrls: ['./match-history.scss']
})
export class MatchHistoryComponent {
    games: Game[];

    constructor(private gamesService: GamesService) {
        this.gamesService.getCompleteGames().subscribe(
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
