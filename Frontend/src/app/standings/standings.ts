import {Component, OnInit, ViewEncapsulation} from '@angular/core';
import {TeamWithPlayers} from "../interfaces/Team";
import {TeamsService} from "../httpServices/teams.service";
import {NGXLogger} from "ngx-logger";

@Component({
    selector: 'app-standings',
    templateUrl: './standings.html',
    styleUrls: ['./standings.scss'],
    encapsulation: ViewEncapsulation.None,
})
export class StandingsComponent implements OnInit {
    teams: TeamWithPlayers[];

    constructor(private log: NGXLogger, private teamsService: TeamsService) {
    }

    ngOnInit(): void {
        this.teamsService.getLeagueTeams().subscribe(
            team => this.teams = team,
            error => this.log.error(error)
        );
    }
}
