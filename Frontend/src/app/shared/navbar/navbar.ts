import {Component, OnInit} from "@angular/core";
import {UserService} from "../../httpServices/user.service";
import {UserWithPermissions} from "../../interfaces/User";
import {League} from "../../interfaces/League";
import {ElmState} from "../state/state.service";
import {NGXLogger} from "ngx-logger";
import {gamesWithStatsPage} from "../lookup.defs";
import * as moment from "moment";

@Component({
    selector: 'app-navbar',
    templateUrl: './navbar.html',
    styleUrls: ['./navbar.scss']
})
export class NavBar implements OnInit {
    user: UserWithPermissions;
    league: League;

    constructor(private state: ElmState,
                private log: NGXLogger,
                private userService: UserService) {
        this.user = null;
    }

    ngOnInit(): void {
        this.state.subscribeUser(user => this.user = user);
        this.state.subscribeLeague(league => this.league = league);
    }

    private hasStatsPage(): boolean {
        return gamesWithStatsPage.includes(this.league.game);
    }

    private isRegistrationPeriod(): boolean {
        return moment().isBetween(moment.unix(this.league.signupStart), moment.unix(this.league.signupEnd));
    }

    private isManager(): boolean {
        if (this.user == null) {
            return false;
        } else {
            return ['administrator', 'createTeams', 'editTeams', 'editGames']
                    .map(k => this.user.leaguePermissions[k]) // create an array of boolean permission values
                    .reduce((p, c) => p || c, false) || // return true if at least one is true
                this.user.teamPermissions.map(teamPermissions => ['administrator', 'information', 'games']
                    .map(k => teamPermissions[k]) // create an array of boolean permission values)
                    .reduce((p, c) => p || c, false)) // return true if at least one is true
                    .reduce((p, c) => p || c, false) // return true if at least one team has a true
        }
    }

    private loggedIn(): boolean {
        return this.user != null && this.user.userId > 0;
    }

    private logout() {
        this.userService.logout();
    }
}
