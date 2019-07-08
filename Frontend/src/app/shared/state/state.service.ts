import {EmptyUser, User, UserWithPermissions} from "../../interfaces/User";
import {EmptyLeague, League, LeaguePermissionsCore} from "../../interfaces/League";
import {Injectable} from "@angular/core";
import {NGXLogger} from "ngx-logger";

type ChangeCallback = () => void;
type UserCallback = (user: UserWithPermissions) => void;
type LeagueCallback = (league: League) => void;

@Injectable({
    providedIn: 'root',
})
export class ElmState {
    private user: UserWithPermissions;
    private league: League;

    private leagueNotDefault;

    private changeListeners: ChangeCallback[] = [];
    private userListeners: UserCallback[] = [];
    private leagueListeners: LeagueCallback[] = [];

    constructor() {
        this.user = EmptyUser();
        this.league = EmptyLeague();
        this.leagueNotDefault = false;
    }

    public subscribeChanges(listener: ChangeCallback) {
        this.changeListeners.push(listener);
    }

    public subscribeUser(listener: UserCallback) {
        this.userListeners.push(listener);
        listener(this.user);
    }

    public subscribeLeague(listener: LeagueCallback, getStateOnlyIfSet: boolean = false) {
        this.leagueListeners.push(listener);
        if (!(getStateOnlyIfSet && this.leagueNotDefault)) {
            listener(this.league);
        }
    }

    public setUser(user: User) {
        this.user.userId = user.userId;
        this.user.email = user.email;
        this.changeListeners.forEach((listener: ChangeCallback) => listener());
        this.userListeners.forEach((listener: UserCallback) => listener(this.user));
    }

    public setUserWithPermissions(user: UserWithPermissions) {
        this.user = user;
        this.userListeners.forEach((listener: UserCallback) => listener(this.user));
    }

    public setLeague(league: League) {
        this.league = league;
        this.leagueNotDefault = true;
        this.changeListeners.forEach((listener: ChangeCallback) => listener());
        this.leagueListeners.forEach((listener: LeagueCallback) => listener(this.league));
    }
}
