import {BrowserModule} from '@angular/platform-browser';
import {RouterModule} from '@angular/router';
import {NgModule} from '@angular/core';
import {HttpClientModule} from '@angular/common/http';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {MatButtonModule, MatCardModule, MatDividerModule, MatTabsModule} from '@angular/material';
import {MatTableModule} from '@angular/material/table';

import {NavBar} from "./shared/navbar/navbar";
import {AppComponent} from './app.component';
import {HomeComponent} from "./home/home";
import {MatchHistoryComponent} from "./matchHistory/match-history";
import {StandingsComponent} from "./standings/standings";
import {TeamsComponent} from "./teams/teams";
import {UpcomingGamesComponent} from "./upcomingGames/upcoming-games";

import {LeagueService} from './httpServices/leagues.service';

import {ELM_ROUTES} from './routes'

@NgModule({
    declarations: [
        AppComponent,
        StandingsComponent,
        HomeComponent,
        MatchHistoryComponent,
        TeamsComponent,
        MatchHistoryComponent,
        UpcomingGamesComponent,
        NavBar
    ],
    imports: [
        BrowserModule,
        HttpClientModule,
        BrowserAnimationsModule,
        MatTabsModule,
        MatTableModule,
        MatButtonModule,
        MatCardModule,
        MatDividerModule,
        RouterModule.forRoot(ELM_ROUTES)
    ],
    providers: [LeagueService],
    bootstrap: [AppComponent]
})
export class AppModule {
}
