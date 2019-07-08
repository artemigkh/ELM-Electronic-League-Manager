// import {Component, ViewEncapsulation} from "@angular/core";
// import {StatsPageInterface} from "./stats-page";
// import {ActivatedRoute} from "@angular/router";
// import {LeagueService} from "../../httpServices/leagues.service";
// import {TeamsService} from "../../httpServices/teams.service";
// import {GamesService} from "../../httpServices/games.service";
// import {
//     ChampionStatsEntry,
//     PlayerStatsEntry,
//     TeamStatsEntry
// } from "../../interfaces/league-of-legends";
// import {LeagueOfLegendsStatsService} from "../../httpServices/league-of-legends-stats.service";
// import {forkJoin} from "rxjs/index";
// import {Team} from "../../interfaces/Team";
// import * as moment from "moment";
//
// class PlayerEntry {
//     name: string;
//     value: string;
// }
//
// class TeamEntry {
//     id: number;
//     name: string;
//     iconSmall: string;
//     value: string;
// }
//
// class AugmentedTeamStatsEntry implements TeamStatsEntry {
//     id: number;
//     name: string;
//     tag: string;
//     iconSmall: string;
//     averageDuration: number;
//     numberFirstBloods: number;
//     numberFirstTurrets: number;
//     averageKda: number;
//     averageWards: number;
//     averageActionScore: number;
//     goldPerMinute: number;
//     csPerMinute: number;
// }
//
// @Component({
//     templateUrl: './league-of-legends-stats-page.html',
//     styleUrls: ['./league-of-legends-stats-page.scss'],
//     encapsulation: ViewEncapsulation.None
// })
// export class LeagueOfLegendsStatsPage implements StatsPageInterface {
//     playerStats: PlayerStatsEntry[] = [];
//     teamStats: AugmentedTeamStatsEntry[] = [];
//     championStats: ChampionStatsEntry[] = [];
//     constructor(private teamsService: TeamsService, private statsService: LeagueOfLegendsStatsService) {
//         this.statsService.getPlayerStats().subscribe(
//             (next: PlayerStatsEntry[]) => {
//                 this.playerStats = next;
//                 console.log(this.playerStats);
//             }, error => {
//                 console.log(error);
//             }
//         );
//
//         this.statsService.getChampionStats().subscribe(
//             (next: ChampionStatsEntry[]) => {
//                 this.championStats = next;
//                 console.log(this.championStats);
//             }, error => {
//                 console.log(error);
//             }
//         );
//
//         this.statsService.getTeamStats().subscribe(
//             (next: TeamStatsEntry[]) => {
//                 let augmentedTeamStats: AugmentedTeamStatsEntry[] = [];
//                 forkJoin(next.map((team: TeamStatsEntry) => {
//                     return this.teamsService.getTeamInformation(team.id);
//                 })).subscribe(
//                     (teamsInfo: Team[]) => {
//                         teamsInfo.forEach((teamInfo: Team) => {
//                             next.forEach((team: TeamStatsEntry) => {
//                                 if(teamInfo.id == team.id) {
//                                     augmentedTeamStats.push({
//                                         id: teamInfo.id,
//                                         name: teamInfo.name,
//                                         tag: teamInfo.tag,
//                                         iconSmall: teamInfo.iconSmall,
//                                         averageDuration: team.averageDuration,
//                                         numberFirstBloods: team.numberFirstBloods,
//                                         numberFirstTurrets: team.numberFirstTurrets,
//                                         averageKda: team.averageKda,
//                                         averageWards: team.averageWards,
//                                         averageActionScore: team.averageActionScore,
//                                         goldPerMinute: team.goldPerMinute,
//                                         csPerMinute: team.csPerMinute
//                                     });
//                                 }
//                             });
//                         });
//
//                         this.teamStats = augmentedTeamStats;
//                     },
//                     error2 => {console.log(error2)}
//                 )
//             }, error => {
//                 console.log(error);
//             }
//         );
//     }
//
//     getPlayerStats(stat: string, desc = true): PlayerEntry[] {
//         let comp = desc ? 1 : -1;
//         let sortedPlayerEntries: PlayerEntry[] = [];
//         this.playerStats.sort((a,b) => (a[stat] > b[stat]) ? -1 * comp : ((b[stat] > a[stat]) ? comp : 0)).
//         slice(0,3).forEach((playerStats: PlayerStatsEntry) => {
//             let value: string;
//             if(Number.isInteger(playerStats[stat])) {
//                 value = playerStats[stat].toFixed(0);
//             } else {
//                 value = playerStats[stat].toFixed(2);
//             }
//             sortedPlayerEntries.push({
//                 name: playerStats.name,
//                 value: value
//             });
//         });
//
//         return sortedPlayerEntries;
//     }
//
//     getChampionStats(stat: string): ChampionStatsEntry[] {
//         return this.championStats.
//             sort((a,b) => (a[stat] > b[stat]) ? -1 : ((b[stat] > a[stat]) ? 1 : 0)).
//             filter(a => {
//                 if (stat == 'winrate') {
//                     return 1;
//                 } else {
//                     return a[stat] > 0;
//                 }
//         });
//     }
//
//     formatWinrate(winrate: number): string {
//         winrate *= 100;
//         let sn: string;
//         if(Number.isInteger(winrate)) {
//             sn = winrate.toFixed(0);
//         } else {
//             sn = winrate.toFixed(2);
//         }
//         return sn + "%";
//     }
//
//
//     getTeamStats(stat: string, desc = true): TeamEntry[] {
//         let comp = desc ? 1 : -1;
//         let sortedTeamEntries: TeamEntry[] = [];
//         this.teamStats.sort((a,b) => (a[stat] > b[stat]) ? -1 * comp : ((b[stat] > a[stat]) ? comp : 0)).
//         slice(0,3).forEach((teamStats: AugmentedTeamStatsEntry) => {
//             let value: string;
//             if(stat == "averageDuration") {
//                 value = moment.utc(teamStats[stat]*1000).format('mm:ss');
//             }
//             else if(Number.isInteger(teamStats[stat])) {
//                 value = teamStats[stat].toFixed(0);
//             } else {
//                 value = teamStats[stat].toFixed(2);
//             }
//             sortedTeamEntries.push({
//                 id: teamStats.id,
//                 name: teamStats.name,
//                 iconSmall: teamStats.iconSmall,
//                 value: value
//             });
//         });
//
//         return sortedTeamEntries;
//     }
// }
