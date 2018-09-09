import { HomeComponent } from './home/home'
import { StandingsComponent } from './standings/standings'
import { TeamsComponent } from './teams/teams'
import { MatchHistoryComponent } from './matchHistory/match-history'
import { UpcomingGamesComponent } from './upcomingGames/upcoming-games'

import {Routes} from "@angular/router";

export const ELM_ROUTES: Routes = [
  {path: '', component: HomeComponent, pathMatch: 'full', data: {}},
  {path: 'standings', component: StandingsComponent, data: {}},
  {path: 'teams', component: TeamsComponent, data: {}},
  {path: 'matchHistory', component: MatchHistoryComponent, data: {}},
  {path: 'upcomingGames', component: UpcomingGamesComponent, data: {}},
  {path: '**', redirectTo: ''},
];
