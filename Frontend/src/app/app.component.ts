import {Component, ViewEncapsulation} from '@angular/core';
import {LeagueService} from './httpServices/leagues.service';
import {TestingConfig} from "../../testingConfig";

@Component({
    selector: 'elm-app',
    templateUrl: './app.component.html',
    styleUrls: ['./app.component.scss'],
    encapsulation: ViewEncapsulation.None,
})
// and password:
export class AppComponent {
    constructor(private leagueService: LeagueService) {
        this.leagueService.setActiveLeague(TestingConfig.leagueId).subscribe(
            success => {
                console.log('successful set league');
                console.log(success);
                this.leagueService.login(TestingConfig.email,
                    TestingConfig.password).subscribe(
                    success => {
                        console.log('successful login');
                        console.log(success);
                    }, error => {
                    console.log('error');
                    console.log(error);
                });
            },
            error => {
                console.log('error');
                console.log(error);
            });
    }
}
