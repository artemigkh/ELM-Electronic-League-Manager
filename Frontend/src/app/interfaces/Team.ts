import {Player} from "./Player";

export class Team {
    teamId: number;
    name: string;
    description?: string;
    tag: string;
    iconSmall: string;
    iconLarge: string;
    wins: number;
    losses: number;
    constructor() {
        this.name = "";
        this.tag = "";
        this.description = "";
    }
}
export class TeamCore {
    name: string = "";
    description: string = "";
    tag: string = "";
}
export interface TeamCoreWithIcon {
    name: string;
    description?: string;
    tag: string;
    icon: string; // binary
}
export interface TeamDisplay {
    teamId: number;
    name: string;
    tag: string;
    iconSmall: string;
    wins: number;
    losses: number;
}
export interface TeamId {
    teamId: number;
}
export interface TeamManager {
    userId: number;
    email: string;
    /**
     * True if this user has administrator priviliges in the team
     */
    administrator: boolean;
    /**
     * True if this user has permissions to change information and players of the team
     */
    information: boolean;
    /**
     * True if this user has permissions to schedule, reschedule, and report games of this team
     */
    games: boolean;
}
export interface TeamPermissions {
    teamId: number;
    name: string;
    tag: string;
    iconSmall: string;
    /**
     * True if this user has administrator priviliges in the team
     */
    administrator: boolean;
    /**
     * True if this user has permissions to change information and players of the team
     */
    information: boolean;
    /**
     * True if this user has permissions to schedule, reschedule, and report games of this team
     */
    games: boolean;
}
export interface TeamPermissionsCore {
    /**
     * True if this user has administrator priviliges in the team
     */
    administrator: boolean;
    /**
     * True if this user has permissions to change information and players of the team
     */
    information: boolean;
    /**
     * True if this user has permissions to schedule, reschedule, and report games of this team
     */
    games: boolean;
}
export interface TeamWithManagers {
    teamId: number;
    name: string;
    tag: string;
    iconSmall: string;
    managers: TeamManager[];
}
export interface TeamWithPlayers extends Team{
    teamId: number;
    name: string;
    description?: string;
    tag: string;
    iconSmall: string;
    iconLarge: string;
    wins: number;
    losses: number;
    players: Player[];
}

export interface TeamWithRosters {
    teamId: number;
    name: string;
    description?: string;
    tag: string;
    iconSmall: string;
    iconLarge: string;
    wins: number;
    losses: number;
    mainRoster: Player[];
    substituteRoster: Player[];
}
