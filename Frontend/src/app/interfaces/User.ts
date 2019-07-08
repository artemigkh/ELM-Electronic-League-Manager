import {DefaultLeaguePermissions, LeaguePermissionsCore} from "./League";
import {TeamPermissions} from "./Team";

export interface UserWithPermissions {
    userId: number;
    email: string;
    /**
     * The permissions this user has in the active league
     */
    leaguePermissions: LeaguePermissionsCore;
    /**
     * The permissions this user has in all teams of the league
     */
    teamPermissions: TeamPermissions[];
}

export function EmptyUser(): UserWithPermissions {
    return {
        email: "",
        leaguePermissions: DefaultLeaguePermissions(),
        teamPermissions: [],
        userId: 0
    }
}

export interface User {
    userId: number;
    email: string;
}
export interface UserCreationInformation {
    email: string;
    password: string; // password
}
export interface UserId {
    userId: number;
}
