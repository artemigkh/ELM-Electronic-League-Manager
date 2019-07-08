import * as moment from "moment";
import {duration} from "moment";
import {Moment} from "moment";
import {Game} from "./Game";

export class Week {
    start: Moment;
    end: Moment;
    games: Game[];
    constructor(gameStart: number) {
        this.start = moment.unix(gameStart).startOf('isoWeek');
        this.end = this.start.clone().endOf('isoWeek');
        this.games = [];
    }
}

export class LeagueCore {
    name: string;
    description: string;
    game: "genericsport" | "basketball" | "curling" | "football" | "hockey" | "rugby" | "soccer" | "volleyball" | "waterpolo" | "genericesport" | "csgo" | "leagueoflegends" | "overwatch";
    /**
     * true if this league can be viewed by the public
     */
    publicView: boolean;
    /**
     * true if this league can be joined by anyone from the public
     */
    publicJoin: boolean;
    /**
     * Beginning of the signup period in seconds since unix epoch
     */
    signupStart: number;
    /**
     * End of the signup period in seconds since unix epoch
     */
    signupEnd: number;
    /**
     * Start of the competition period in seconds since unix epoch
     */
    leagueStart: number;
    /**
     * End of the competition period in seconds since unix epoch
     */
    leagueEnd: number;
    constructor() {
        this.name = "";
        this.description = "";
        this.game = "genericsport";
        this.publicView = true;
        this.publicJoin = true;
        this.signupStart = moment().unix();
        this.signupEnd = moment().add(1, 'w').endOf('isoWeek').unix();
        this.leagueStart = this.signupEnd;
        this.leagueEnd = moment().add(3, 'w').endOf('isoWeek').unix();
    }
}

export class League extends LeagueCore{
    constructor(leagueId: number, core: LeagueCore) {
        super();
        this.leagueId = leagueId;
        this.name = core.name;
        this.description = core.description;
        this.game = core.game;
        this.publicView = core.publicView;
        this.publicJoin = core.publicJoin;
        this.signupStart = core.signupStart;
        this.signupEnd = core.signupEnd;
        this.leagueStart = core.leagueStart;
        this.leagueEnd = core.leagueEnd;
    }

    leagueId: number;
    name: string;
    description: string;
    game: "genericsport" | "basketball" | "curling" | "football" | "hockey" | "rugby" | "soccer" | "volleyball" | "waterpolo" | "genericesport" | "csgo" | "leagueoflegends" | "overwatch";
    /**
     * true if this league can be viewed by the public
     */
    publicView: boolean;
    /**
     * true if this league can be joined by anyone from the public
     */
    publicJoin: boolean;
    /**
     * Beginning of the signup period in seconds since unix epoch
     */
    signupStart: number;
    /**
     * End of the signup period in seconds since unix epoch
     */
    signupEnd: number;
    /**
     * Start of the competition period in seconds since unix epoch
     */
    leagueStart: number;
    /**
     * End of the competition period in seconds since unix epoch
     */
    leagueEnd: number;
}

export function EmptyLeague(): League {
    return {
        description: "",
        game: "genericsport",
        leagueEnd: 0,
        leagueId: 0,
        leagueStart: 0,
        name: "",
        publicJoin: false,
        publicView: false,
        signupEnd: 0,
        signupStart: 0
    }
}

export interface LeagueId {
    leagueId: number;
}
export interface LeaguePermissionsCore {
    /**
     * True if this user has administrator priviliges in the active league
     */
    administrator: boolean;
    /**
     * True if this user has permissions to create new teams
     */
    createTeams: boolean;
    /**
     * True if this user has permissions to edit any information in teams
     */
    editTeams: boolean;
    /**
     * True if this user has permissions to reschedule games and report game results
     */
    editGames: boolean;
}

export function DefaultLeaguePermissions(): LeaguePermissionsCore {
    return {
        administrator: false, createTeams: false, editGames: false, editTeams: false
    }
}

export class Markdown {
    markdown: string = "";
}
