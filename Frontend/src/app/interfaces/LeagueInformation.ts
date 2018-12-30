export interface LeagueInformation {
    id: number;
    name: string;
    description: string;
    game: string;
    publicView: boolean;
    publicJoin: boolean;
    signupStart: number;
    signupEnd: number;
    leagueStart: number;
    leagueEnd: number;
}
