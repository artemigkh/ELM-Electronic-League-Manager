import * as moment from "moment";
import {Moment} from "moment";
import {League} from "./League";

export interface AvailabilityId {
    availabilityId: number;
}

export interface Availability {
    availabilityId: number;
    /**
     * The start time of the availability in seconds since unix epoch
     */
    startTime: number;
    /**
     * The end time of the availability in seconds since unix epoch
     */
    endTime: number;
}

export interface AvailabilityCore {
    /**
     * The start time of the availability in seconds since unix epoch
     */
    startTime: number;
    /**
     * The end time of the availability in seconds since unix epoch
     */
    endTime: number;
}

export class SchedulingParameters {
    tournamentType: "roundrobin" | "doubleroundrobin";
    /**
     * Where a round is when each team plays exactly once (or gets a bye)
     */
    roundsPerWeek: number;
    /**
     * How many games can be scheduled that play at the same time
     */
    concurrentGameNum: number;
    /**
     * Duration of a game in minutes
     */
    gameDuration: number;

    constructor() {
        this.tournamentType = "roundrobin";
        this.roundsPerWeek = 1;
        this.concurrentGameNum = 1;
        this.gameDuration = 60;
    }
}

export class WeeklyAvailability {
    availabilityId: number;
    /**
     * The start time of the availability in seconds since unix epoch
     */
    startTime: number;
    /**
     * The end time of the availability in seconds since unix epoch
     */
    endTime: number;
    weekday: "sunday" | "monday" | "tuesday" | "wednesday" | "thursday" | "friday" | "saturday";
    /**
     * timezone as offset in seconds east from UTC
     */
    timezone: number;
    hour: number;
    minute: number;
    /**
     * Duration of a game in minutes
     */
    duration: number;

    constructor(league: League) {
        this.availabilityId = 0;
        this.duration = 60;
        this.endTime = league.leagueEnd;
        this.hour = 12;
        this.minute = 0;
        this.startTime = league.leagueStart;
        this.timezone = 0;
        this.weekday = 'saturday';
    }
}

export function getStartMoment(a: WeeklyAvailability): Moment {
    return moment().milliseconds(0).seconds(0).minute(a.minute).hour(a.hour).utcOffset(a.timezone / 60).local();
}

export function updateFromMoments(a: WeeklyAvailability, start: Moment, end: Moment) {
    a.timezone = start.utcOffset() * 60;
    a.hour = start.hour();
    a.minute = start.minute();
    a.duration = moment.duration(end.diff(start)).asMinutes();
}

export function defaultWeeklyAvailability(league: League): WeeklyAvailability {
    return <WeeklyAvailability>{
        availabilityId: 0,
        duration: 60,
        endTime: league.leagueEnd,
        hour: 12,
        minute: 0,
        startTime: league.leagueStart,
        timezone: 0,
        weekday: 'saturday'
    }
}

//
// export class MomentWeeklyAvailability extends WeeklyAvailability {
//     availabilityId: number;
//     /**
//      * The start time of the availability in seconds since unix epoch
//      */
//     startTime: number;
//     /**
//      * The end time of the availability in seconds since unix epoch
//      */
//     endTime: number;
//     weekday: "sunday" | "monday" | "tuesday" | "wednesday" | "thursday" | "friday" | "saturday";
//     /**
//      * timezone as offset in seconds east from UTC
//      */
//     timezone: number;
//     hour: number;
//     minute: number;
//     /**
//      * Duration of a game in minutes
//      */
//     duration: number;
//
//     start: Moment;
//     end: Moment;
//
//     constructor(a: WeeklyAvailability) {
//         super();
//         this.availabilityId = a.availabilityId;
//         this.startTime = a.startTime;
//         this.endTime = a.endTime;
//         this.weekday = a.weekday;
//         this.timezone = a.timezone;
//         this.hour = a.hour;
//         this.minute = a.minute;
//         this.duration = a.duration;
//         this.start = moment().milliseconds(0).seconds(0).minute(a.minute).hour(a.hour).utcOffset(a.timezone / 60);
//         this.end = this.start.clone().add(this.duration, 'm')
//     }
//
//     updateFromMoments(start: Moment, end: Moment) {
//         this.duration = moment.duration(end.diff(start)).asMinutes();
//         this.timezone = start.utcOffset() * 60;
//         this.hour = start.hour();
//         this.minute = start.minute();
//     }
// }

export class WeeklyAvailabilityCore {
    /**
     * The start time of the availability in seconds since unix epoch
     */
    startTime: number;
    /**
     * The end time of the availability in seconds since unix epoch
     */
    endTime: number;
    weekday: "sunday" | "monday" | "tuesday" | "wednesday" | "thursday" | "friday" | "saturday";
    /**
     * timezone as offset in seconds east from UTC
     */
    timezone: number;
    hour: number;
    minute: number;
    /**
     * Duration of a game in minutes
     */
    duration: number;
}
