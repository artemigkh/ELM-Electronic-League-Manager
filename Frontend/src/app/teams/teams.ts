import {Component, OnInit, ViewEncapsulation} from "@angular/core";
import {EmptySortedGames, Game, SortedGames} from "../interfaces/Game";
import {Team, TeamWithPlayers, TeamWithRosters} from "../interfaces/Team";
import {ElmState} from "../shared/state/state.service";
import {NGXLogger} from "ngx-logger";
import {GamesService} from "../httpServices/games.service";
import {TeamsService} from "../httpServices/teams.service";
import {ActivatedRoute} from "@angular/router";
import {FormControl} from "@angular/forms";

@Component({
    selector: 'app-teams',
    templateUrl: './teams.html',
    styleUrls: ['./teams.scss'],
    encapsulation: ViewEncapsulation.None
})
export class TeamsComponent implements OnInit{
    selected = new FormControl(0);

    team: TeamWithRosters;
    games: SortedGames;

    constructor(private state: ElmState,
                private log: NGXLogger,
                private route: ActivatedRoute,
                private gamesService: GamesService,
                private teamsService: TeamsService) {

        this.games = EmptySortedGames();
    }

    ngOnInit(): void {
        this.route.params.subscribe(
            params => {
                this.teamsService.getTeamWithRosters(params['teamId']).subscribe(
                    team => {this.team = team; this.selected.setValue(0);},
                    error => this.log.error(error)
                );
                this.gamesService.getLeagueGames({teamId: params['teamId']}).subscribe(
                    games => this.games = games,
                    error => this.log.error(error)
                );
            },
            error => this.log.error(error)
        );
    }
}
