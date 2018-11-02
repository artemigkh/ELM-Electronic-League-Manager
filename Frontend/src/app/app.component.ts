import {Component, ViewEncapsulation} from '@angular/core';
import {LeagueService} from './httpServices/leagues.service';

@Component({
    selector: 'elm-app',
    templateUrl: './app.component.html',
    styleUrls: ['./app.component.scss'],
    encapsulation: ViewEncapsulation.None,
})
export class AppComponent {
    constructor(private leagueService: LeagueService) {
        this.leagueService.setActiveLeague(11).subscribe(
            success => {
                console.log('success');
                console.log(success);
            },
            error => {
                console.log('error');
                console.log(error);
            });
    }
}