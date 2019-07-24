import {Component, OnInit, ViewEncapsulation} from '@angular/core';
import {LeagueService} from './httpServices/leagues.service';
import {TestingConfig} from "../../testingConfig";
import {UserService} from "./httpServices/user.service";
import {ElmState} from "./shared/state/state.service";
import {NGXLogger} from "ngx-logger";

@Component({
    selector: 'elm-app',
    templateUrl: './app.component.html',
    styleUrls: ['./app.component.scss'],
    encapsulation: ViewEncapsulation.None,
})

export class AppComponent implements OnInit {
    constructor(private state: ElmState,
                private log: NGXLogger,
                private leagueService: LeagueService, private userService: UserService) {
    }

    ngOnInit(): void {
        if (TestingConfig.testing) {
            this.userService.login(
                TestingConfig.email,
                TestingConfig.password
            ).subscribe(success => {
                    this.leagueService.setActiveLeague(TestingConfig.leagueId).subscribe(
                        success => {
                        },
                        error => {
                            this.log.error(error)
                        })
                },
                error => {
                    this.log.warn(error)
                });
        } else {
            this.userService.getCurrentUser().subscribe(success => {
                    this.leagueService.getLeagueInformation().subscribe(
                        success => {
                        },
                        error => {
                            this.log.error(error)
                        })
                },
                error => {
                    this.log.warn(error)
                });
        }
    }
}
