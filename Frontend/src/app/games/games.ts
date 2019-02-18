import {Component, ViewEncapsulation} from '@angular/core';
import {LeagueService} from "../httpServices/leagues.service";
import {Game} from "../interfaces/Game";
import * as moment from "moment";
import {Moment} from "moment";
import {GamesService} from "../httpServices/games.service";

class Week {
    start: Moment;
    end: Moment;
    games: Game[];
}

@Component({
    selector: 'app-games',
    templateUrl: './games.html',
    styleUrls: ['./games.scss'],
    encapsulation: ViewEncapsulation.None
})
export class GamesComponent {
    games: Game[];
    weeks: Week[];
    currentWeek: Number;
    processGamesIntoWeeks(): void {
        this.weeks = [];
        this.games.sort((a,b)=>
            (a.gameTime > b.gameTime) ? 1 :
            ((a.gameTime < b.gameTime) ? -1 : 0));
        console.log(this.games);
        let firstGame = moment.unix(this.games[0].gameTime);
        let lastGame = moment.unix(this.games[this.games.length-1].gameTime);

        this.weeks.push({
            start: firstGame.clone().startOf('isoWeek'),
            end: firstGame.clone().endOf('isoWeek'),
            games: []
        });
        while(!lastGame.isBetween(this.weeks[this.weeks.length-1].start, this.weeks[this.weeks.length-1].end)) {
            this.weeks.push({
                start: this.weeks[this.weeks.length-1].start.clone().add(1, 'w'),
                end: this.weeks[this.weeks.length-1].end.clone().add(1, 'w'),
                games: []
            });
        }

        let weekIndex = 0;
        this.games.forEach((game: Game) => {
            while(!moment.unix(game.gameTime).
                          isBetween(this.weeks[weekIndex].start, this.weeks[weekIndex].end)) {
                weekIndex++;
            }
            this.weeks[weekIndex].games.push(game);
        });

        this.currentWeek = 1;
        weekIndex = 0;
        this.weeks.forEach((week: Week) => {
           if(moment().isBetween(week.start, week.end)) {
               this.currentWeek = weekIndex + 1;
           }
           weekIndex++;
        });
        console.log(this.weeks);
    }

    constructor(private gamesService: GamesService) {
        this.gamesService.getAllGames().subscribe(
            gameSummary => {
                console.log('success');
                console.log(gameSummary);
                this.games = gameSummary.upcomingGames.concat(gameSummary.completeGames);
                if(this.games.length > 0) {
                    this.processGamesIntoWeeks();
                }
            },
            error => {
                console.log('error');
                console.log(error);
            });


    }
}
