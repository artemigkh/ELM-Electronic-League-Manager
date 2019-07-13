export interface PlayerStatsEntry {
    id: string;
    name: string;
    teamId: number;
    averageDuration: number;
    goldPerMinute: number;
    csPerMinute: number;
    damagePerMinute: number;
    averageKills: number;
    averageAssists: number;
    averageKda: number;
    averageWards: number;
}

export interface TeamStatsEntry {
    id: number;
    averageDuration: number;
    numberFirstBloods: number;
    numberFirstTurrets: number;
    averageKda: number;
    averageWards: number;
    averageActionScore: number;
    goldPerMinute: number;
    csPerMinute: number;
}

export interface ChampionStatsEntry {
    name: string;
    bans: number;
    picks: number;
    wins: number;
    losses: number;
    winrate: number;
}

/*
1) Fastest Time
2) Most First bloods
3) Most First Turrets
4) Highest avg KDA
5) Highest avg Action Score
6) Most wards
7) Most GPM
8) Most CSPM
 */
