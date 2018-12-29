import {Component} from '@angular/core';
import {Game} from "../interfaces/Game";
import {GamesService} from "../httpServices/games.service";

@Component({
    selector: 'app-upcoming-games',
    templateUrl: './upcoming-games.html',
    styleUrls: ['./upcoming-games.scss']
})
export class UpcomingGamesComponent {
    games: Game[];

    constructor(private gamesService: GamesService) {
        this.gamesService.getUpcomingGames().subscribe(
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
