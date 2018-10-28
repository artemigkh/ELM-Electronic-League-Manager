import {HomeComponent} from './home/home'
import {StandingsComponent} from './standings/standings'
import {TeamsComponent} from './teams/teams'
import {MatchHistoryComponent} from './matchHistory/match-history'
import {UpcomingGamesComponent} from './upcomingGames/upcoming-games'

import {Routes} from "@angular/router";
import {ManageComponent} from "./manage/manage";
import {ManageLeagueComponent} from "./manage/league/manage-league";
import {ManageTeamsComponent} from "./manage/teams/manage-teams";
import {ManagePermissionsComponent} from "./manage/permissions/manage-permissions";
import {ManageDatesComponent} from "./manage/dates/manage-dates";
import {ManagePlayersComponent} from "./manage/players/manage-players";
import {ManageGamesComponent} from "./manage/games/manage-games";
import {LoginComponent} from "./login/login";

export const ELM_ROUTES: Routes = [
    {path: '', component: HomeComponent, pathMatch: 'full', data: {}},
    {path: 'standings', component: StandingsComponent, data: {}},
    {path: 'teams', component: TeamsComponent, data: {}},
    {path: 'matchHistory', component: MatchHistoryComponent, data: {}},
    {path: 'upcomingGames', component: UpcomingGamesComponent, data: {}},
    {path: 'login', component: LoginComponent, data: {}},
    {
        path: 'manage',
        component: ManageComponent,
        data: {},
        children: [
            {path: 'league', component: ManageLeagueComponent},
            {path: 'permissions', component: ManagePermissionsComponent},
            {path: 'teams', component: ManageTeamsComponent},
            {path: 'dates', component: ManageDatesComponent},
            {path: 'players', component: ManagePlayersComponent},
            {path: 'games', component: ManageGamesComponent}
        ]
    },
    {path: '**', redirectTo: ''},
];
