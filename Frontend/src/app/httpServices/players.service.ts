// import {Injectable} from "@angular/core";
// import {HttpClient} from "@angular/common/http";
// import {httpOptions} from "./http-options";
// import {Observable} from "rxjs/Rx";
// import {LeagueService} from "./leagues.service";
// import {Player} from "../interfaces/Player";
//
// @Injectable()
// export class PlayersService {
//     constructor(private http: HttpClient, private leagueSevice: LeagueService) {}
//
//     public addPlayer(teamId: number, mainRoster: boolean, player: Player): Observable<Object> {
//         let url = "";
//         let body = {
//             teamId: teamId,
//             name: player.name,
//             gameIdentifier: player.gameIdentifier,
//             position: player.position,
//             mainRoster: mainRoster
//         };
//         switch(this.leagueSevice.getGame()) {
//             case 'leagueoflegends': {
//                 url = 'http://localhost:8080/api/league-of-legends/teams/addPlayer';
//                 break;
//             }
//             default: {
//                 url = 'http://localhost:8080/api/teams/addPlayer';
//             }
//         }
//         return this.http.post(url, body, httpOptions)
//     }
//
//     public removePlayer(teamId: number, playerId: number) {
//         return this.http.request('delete', 'http://localhost:8080/api/teams/removePlayer',
//             {
//                 withCredentials: httpOptions.withCredentials,
//                 headers: httpOptions.headers,
//                 body: {
//                     teamId: teamId,
//                     playerId: playerId
//                 }
//             });
//     }
//
//     public updatePlayer(teamId: number, mainRoster: boolean, player: Player): Observable<Object> {
//         let url = "";
//         let body = {
//             teamId: teamId,
//             playerId: player.id,
//             name: player.name,
//             gameIdentifier: player.gameIdentifier,
//             position: player.position,
//             mainRoster: mainRoster
//         };
//         switch(this.leagueSevice.getGame()) {
//             case 'leagueoflegends': {
//                 url = 'http://localhost:8080/api/league-of-legends/teams/updatePlayer';
//                 break;
//             }
//             default: {
//                 url = 'http://localhost:8080/api/teams/updatePlayer';
//             }
//         }
//         return this.http.put(url, body, httpOptions)
//     }
// }
