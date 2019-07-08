// import {Component, ViewEncapsulation} from '@angular/core';
// import {LeagueService} from "../httpServices/leagues.service";
import {CompetitionWeek, Game} from "../interfaces/Game";
import * as moment from "moment";
import {Moment} from "moment";
import {Component, OnInit, ViewEncapsulation} from "@angular/core";
import {GamesService} from "../httpServices/games.service";
import {NGXLogger} from "ngx-logger";
import {Week} from "../interfaces/League";
// import {GamesService} from "../httpServices/games.service";
//

@Component({
    selector: 'app-games',
    templateUrl: './games.html',
    styleUrls: ['./games.scss'],
    encapsulation: ViewEncapsulation.None
})
export class GamesComponent implements OnInit{
    weeks: Week[];
    currentWeek: Number;

    constructor(private log: NGXLogger, private gamesService: GamesService) {
        this.weeks = [];
    }

    ngOnInit(): void {
        this.gamesService.getGamesByWeek().subscribe(
            weeks => this.processCompetitionWeeks(weeks),
            error => this.log.error(error)
        );
    }

    private processCompetitionWeeks(competitionWeek: CompetitionWeek[]) {
        competitionWeek.forEach(week => {
            this.weeks.push(<Week>{
                start: moment.unix(week.weekStart).startOf('isoWeek'),
                end: moment.unix(week.weekStart).endOf('isoWeek'),
                games: week.games
            })
        });
        this.setCurrentWeek();
    }

    private setCurrentWeek() {
        this.currentWeek = 1;
        Array.from(Array(this.weeks.length).keys()).forEach(i => {
            if(moment().isBetween(this.weeks[i].start, this.weeks[i].end)) {
                this.currentWeek = i;
            }
        });
    }
}
