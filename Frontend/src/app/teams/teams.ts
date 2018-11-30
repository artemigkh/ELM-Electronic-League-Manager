import {Component} from '@angular/core';
import {ActivatedRoute} from "@angular/router";
import {TeamsService} from "../httpServices/teams.service";
import {Team} from "../interfaces/Team";
import { ViewEncapsulation } from '@angular/core';
@Component({
    selector: 'app-teams',
    templateUrl: './teams.html',
    styleUrls: ['./teams.scss'],
    encapsulation: ViewEncapsulation.None
})
export class TeamsComponent {
    team: Team;

    constructor(private route: ActivatedRoute, private teamsService: TeamsService) {
        this.route.params.subscribe(params => {
            this.teamsService.getTeamInformation(+params['id']).subscribe(
                (next: Team) => {
                    this.team = next;
                    console.log('test new');
                    console.log(next);
                }, error => {
                    console.log(error);
                }
            )
        });
    }
}
