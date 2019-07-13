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
    private user: User;
    private userWithPermissions: UserWithPermissions;
    private league: League;

    private changeListeners: ChangeCallback[] = [];
    private userListeners: UserCallback[] = [];
    private leagueListeners: LeagueCallback[] = [];

    constructor() {
        this.user = null;
        this.userWithPermissions = null;
        this.league = null;
    }

    public subscribeChanges(listener: ChangeCallback) {
        this.changeListeners.push(listener);
    }

    public subscribeUser(listener: UserCallback) {
        this.userListeners.push(listener);
        if (this.userWithPermissions != null) {
            listener(this.userWithPermissions);
        }
    }

    public subscribeLeague(listener: LeagueCallback) {
        this.leagueListeners.push(listener);
        if (this.league != null) {
            listener(this.league);
        }
    }

    public setUser(user: User) {
        this.user = user;
        this.changeListeners.forEach((listener: ChangeCallback) => listener());
    }

    public setUserWithPermissions(userWithPermissions: UserWithPermissions) {
        this.userWithPermissions = userWithPermissions;
        this.userListeners.forEach((listener: UserCallback) => listener(this.userWithPermissions));
    }

    public setLeague(league: League) {
        console.log("setting league", league);
        this.league = league;
        this.changeListeners.forEach((listener: ChangeCallback) => listener());
        this.leagueListeners.forEach((listener: LeagueCallback) => listener(this.league));
    }
}
